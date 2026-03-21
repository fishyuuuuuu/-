<template>
  <div class="category-page">
    <!-- 顶部导航 -->
    <Navbar />

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <div class="header-content">
          <div class="header-title">
            <h1>全部商品</h1>
            <p class="header-subtitle">探索精选好物，发现品质生活</p>
          </div>
          <div class="header-stats">
            <div class="stat-item">
              <span class="stat-number">{{ products.length }}</span>
              <span class="stat-label">商品总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">{{ seckillCount }}</span>
              <span class="stat-label">秒杀商品</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="container">
        <div class="filter-content">
          <!-- 分类标签 -->
          <div class="category-tabs">
            <button
              v-for="category in categories"
              :key="category.id"
              :class="['tab-btn', { active: activeCategory === category.id }]"
              @click="activeCategory = category.id"
            >
              <span class="tab-icon">{{ category.icon }}</span>
              <span class="tab-text">{{ category.name }}</span>
              <span v-if="category.count" class="tab-count">{{ category.count }}</span>
            </button>
          </div>

          <!-- 排序和筛选 -->
          <div class="filter-options">
            <div class="sort-select">
              <select v-model="sortBy" @change="handleSort">
                <option value="default">综合排序</option>
                <option value="price-asc">价格从低到高</option>
                <option value="price-desc">价格从高到低</option>
                <option value="sales">销量优先</option>
                <option value="newest">最新上架</option>
              </select>
            </div>
            <div class="filter-tags">
              <label class="filter-tag" :class="{ active: showSeckillOnly }">
                <input type="checkbox" v-model="showSeckillOnly" hidden>
                <span class="tag-dot"></span>
                <span>仅看秒杀</span>
              </label>
              <label class="filter-tag" :class="{ active: showInStockOnly }">
                <input type="checkbox" v-model="showInStockOnly" hidden>
                <span class="tag-dot"></span>
                <span>仅看有货</span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 商品列表 -->
    <div class="container">
      <div class="products-section">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <p>正在加载商品...</p>
        </div>

        <!-- 商品网格 -->
        <div v-else-if="filteredProducts.length > 0" class="products-grid">
          <div
            v-for="product in filteredProducts"
            :key="product.id"
            class="product-card"
            :class="{ 'seckill': product.isSeckill, 'out-of-stock': product.stock === 0 }"
            @click="goToDetail(product)"
          >
            <!-- 商品图片 -->
            <div class="product-image">
              <img :src="product.image" :alt="product.name" loading="lazy">
              <div class="product-badges">
                <span v-if="product.isSeckill" class="badge seckill">秒杀</span>
                <span v-if="product.stock === 0" class="badge soldout">售罄</span>
                <span v-else-if="product.stock < 10" class="badge low-stock">仅剩{{ product.stock }}件</span>
              </div>
              <div v-if="product.isSeckill" class="discount-badge">
                -{{ product.discount }}%
              </div>
            </div>

            <!-- 商品信息 -->
            <div class="product-info">
              <h3 class="product-name">{{ product.name }}</h3>
              <p class="product-desc">{{ product.description }}</p>

              <!-- 价格 -->
              <div class="product-price">
                <div class="price-main">
                  <span class="price-symbol">¥</span>
                  <span class="price-value">{{ formatPrice(product.price) }}</span>
                </div>
                <div v-if="product.originalPrice" class="price-original">
                  ¥{{ formatPrice(product.originalPrice) }}
                </div>
              </div>

              <!-- 销量和库存 -->
              <div class="product-meta">
                <span class="sales">已售 {{ product.sales || 0 }}</span>
                <div class="stock-bar" v-if="product.stock > 0">
                  <div class="stock-progress" :style="{ width: (product.stock / product.totalStock * 100) + '%' }"></div>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="product-actions">
                <button
                  class="btn btn-seckill btn-sm"
                  v-if="product.isSeckill && product.stock > 0"
                  @click.stop="handleSeckill(product)"
                >
                  立即秒杀
                </button>
                <button
                  class="btn btn-primary btn-sm"
                  v-else-if="product.stock > 0"
                  @click.stop="addToCart(product)"
                >
                  加入购物车
                </button>
                <button class="btn btn-secondary btn-sm" disabled v-else>
                  已售罄
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="empty-state">
          <div class="empty-icon">📦</div>
          <h3 class="empty-title">暂无商品</h3>
          <p class="empty-desc">换个筛选条件试试吧</p>
          <button class="btn btn-primary" @click="resetFilters">重置筛选</button>
        </div>
      </div>
    </div>

    <!-- 商品详情弹窗 -->
    <ProductDetailModal
      v-if="selectedProduct"
      :product="selectedProduct"
      @close="selectedProduct = null"
      @seckill="handleSeckill"
      @addToCart="addToCart"
    />

    <!-- 验证码弹窗 -->
    <div class="modal captcha-modal-wrapper" v-if="showCaptchaModal" @click.self="showCaptchaModal = false">
      <div class="modal-content captcha-content">
        <h3>🔒 安全验证</h3>
        <p>请输入图片中的验证码完成秒杀</p>
        <div class="captcha-row">
          <div class="captcha-image-wrapper" @click="fetchCaptcha">
            <img v-if="captchaImageUrl" :src="captchaImageUrl" alt="验证码" class="captcha-image">
            <span class="captcha-refresh">🔄 点击刷新</span>
          </div>
          <input
            v-model="captchaInput"
            type="text"
            maxlength="4"
            placeholder="请输入验证码"
            class="captcha-input"
          >
        </div>
        <div class="captcha-actions">
          <button class="btn btn-outline" @click="showCaptchaModal = false">取消</button>
          <button class="btn btn-primary" @click="confirmCaptchaAndSeckill" :disabled="isSubmitting">
            <span v-if="!isSubmitting">确认秒杀</span>
            <span v-else class="btn-spinner-white"></span>
          </button>
        </div>
      </div>
    </div>

    <!-- 秒杀成功提示 -->
    <div class="modal success-modal" v-if="showSuccessModal" @click.self="showSuccessModal = false">
      <div class="modal-content success-content">
        <div class="success-animation">
          <div class="success-circle">
            <span class="success-icon">✓</span>
          </div>
          <div class="confetti">
            <span v-for="i in 12" :key="i" :class="'confetti-' + i"></span>
          </div>
        </div>
        <h3>秒杀成功!</h3>
        <p>恭喜您成功秒杀 <span class="success-product">{{ successProductName }}</span></p>
        <div class="success-actions">
          <button class="btn btn-primary" @click="showSuccessModal = false">继续购物</button>
        </div>
      </div>
    </div>

    <!-- 提示消息 -->
    <div class="toast" :class="{ 'show': showToast, [toastType]: true }">
      <span class="toast-icon">{{ toastIcon }}</span>
      <span class="toast-message">{{ toastMessage }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from './Navbar.vue';
import ProductDetailModal from './ProductDetailModal.vue';

const router = useRouter();

// 状态
const loading = ref(false);
const activeCategory = ref('all');
const sortBy = ref('default');
const showSeckillOnly = ref(false);
const showInStockOnly = ref(false);
const selectedProduct = ref(null);

// 验证码相关状态
const showCaptchaModal = ref(false);
const captchaImageUrl = ref('');
const captchaId = ref('');
const captchaInput = ref('');
const pendingProductForCaptcha = ref(null);
const isSubmitting = ref(false);

// 秒杀状态
const seckilling = reactive({});
const showSuccessModal = ref(false);
const successProductName = ref('');

// 提示消息
const showToast = ref(false);
const toastMessage = ref('');
const toastType = ref('success');
const toastIcon = ref('✓');

// 分类数据
const categories = ref([
  { id: 'all', name: '全部', icon: '📦', count: 0 },
  { id: 'phone', name: '手机数码', icon: '📱', count: 12 },
  { id: 'computer', name: '电脑办公', icon: '💻', count: 8 },
  { id: 'audio', name: '影音娱乐', icon: '🎧', count: 6 },
  { id: 'home', name: '智能家居', icon: '🏠', count: 10 },
  { id: 'wearable', name: '智能穿戴', icon: '⌚', count: 5 },
]);

// 商品数据
const products = ref([
  {
    id: 1,
    name: 'iPhone 15 Pro',
    description: '钛金属设计，A17 Pro芯片，专业级摄像系统',
    price: 7999,
    originalPrice: 8999,
    stock: 15,
    totalStock: 100,
    sales: 2341,
    isSeckill: true,
    discount: 11,
    category: 'phone',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPhone%2015%20Pro%20titanium%20metal%20smartphone%20with%20triple%20camera%20system%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
  {
    id: 2,
    name: 'MacBook Pro 14',
    description: 'M3芯片，Liquid视网膜XDR显示屏',
    price: 12999,
    originalPrice: 14999,
    stock: 8,
    totalStock: 50,
    sales: 567,
    isSeckill: true,
    discount: 13,
    category: 'computer',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=MacBook%20Pro%2014%20inch%20laptop%20with%20silver%20aluminum%20body%20open%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
  {
    id: 3,
    name: 'AirPods Pro 2',
    description: '主动降噪，空间音频，MagSafe充电盒',
    price: 1899,
    originalPrice: 2299,
    stock: 45,
    totalStock: 200,
    sales: 8923,
    isSeckill: true,
    discount: 17,
    category: 'audio',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=AirPods%20Pro%202%20wireless%20earbuds%20with%20charging%20case%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
  {
    id: 4,
    name: 'iPad Pro 12.9',
    description: 'M2芯片，12.9英寸Liquid视网膜XDR屏',
    price: 8999,
    originalPrice: 9999,
    stock: 0,
    totalStock: 50,
    sales: 1234,
    isSeckill: true,
    discount: 10,
    category: 'phone',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPad%20Pro%2012.9%20inch%20tablet%20with%20Apple%20Pencil%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
  {
    id: 5,
    name: 'Apple Watch Ultra 2',
    description: '钛金属表壳，双频GPS，100米防水',
    price: 6299,
    originalPrice: 6999,
    stock: 23,
    totalStock: 80,
    sales: 1890,
    isSeckill: true,
    discount: 10,
    category: 'wearable',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Apple%20Watch%20Ultra%20titanium%20smartwatch%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
  {
    id: 6,
    name: 'Sony WH-1000XM5',
    description: '业界领先的降噪耳机，30小时续航',
    price: 2499,
    originalPrice: 2999,
    stock: 67,
    totalStock: 150,
    sales: 3456,
    isSeckill: false,
    category: 'audio',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Sony%20WH-1000XM5%20black%20over-ear%20headphones%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
  {
    id: 7,
    name: '小米14 Pro',
    description: '徕卡光学镜头，骁龙8 Gen3',
    price: 4999,
    originalPrice: 5499,
    stock: 89,
    totalStock: 200,
    sales: 5678,
    isSeckill: true,
    discount: 9,
    category: 'phone',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Xiaomi%2014%20Pro%20smartphone%20with%20Leica%20camera%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
  {
    id: 8,
    name: 'Dell XPS 15',
    description: '英特尔酷睿i7，OLED触控屏',
    price: 15999,
    originalPrice: 17999,
    stock: 12,
    totalStock: 30,
    sales: 234,
    isSeckill: false,
    category: 'computer',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Dell%20XPS%2015%20laptop%20with%20OLED%20touch%20screen%20open%20on%20white%20background%20professional%20product%20photography&image_size=square'
  },
]);

// 计算属性
const seckillCount = computed(() => products.value.filter(p => p.isSeckill).length);

const filteredProducts = computed(() => {
  let result = products.value;

  // 分类筛选
  if (activeCategory.value !== 'all') {
    result = result.filter(p => p.category === activeCategory.value);
  }

  // 仅看秒杀
  if (showSeckillOnly.value) {
    result = result.filter(p => p.isSeckill);
  }

  // 仅看有货
  if (showInStockOnly.value) {
    result = result.filter(p => p.stock > 0);
  }

  // 排序
  switch (sortBy.value) {
    case 'price-asc':
      result = [...result].sort((a, b) => a.price - b.price);
      break;
    case 'price-desc':
      result = [...result].sort((a, b) => b.price - a.price);
      break;
    case 'sales':
      result = [...result].sort((a, b) => b.sales - a.sales);
      break;
    case 'newest':
      result = [...result].sort((a, b) => b.id - a.id);
      break;
  }

  return result;
});

// 方法
const formatPrice = (price) => {
  return price.toLocaleString('zh-CN');
};

const handleSort = () => {
  // 排序逻辑已在计算属性中处理
};

const resetFilters = () => {
  activeCategory.value = 'all';
  sortBy.value = 'default';
  showSeckillOnly.value = false;
  showInStockOnly.value = false;
};

const goToDetail = (product) => {
  selectedProduct.value = product;
};

// 显示提示消息
const showToastMessage = (message, type = 'success') => {
  toastMessage.value = message;
  toastType.value = type;
  toastIcon.value = type === 'success' ? '✓' : type === 'error' ? '✗' : type === 'warning' ? '⚠' : 'ℹ';
  showToast.value = true;
  
  setTimeout(() => {
    showToast.value = false;
  }, 3000);
};

// 获取验证码
const fetchCaptcha = async () => {
  try {
    // 模拟验证码请求
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 生成随机验证码图片 URL
    const randomCode = Math.random().toString(36).substring(2, 6).toUpperCase();
    captchaId.value = Date.now().toString();
    captchaImageUrl.value = `https://placehold.co/200x60/667eea/ffffff?text=${randomCode}`;
  } catch (e) {
    console.error('获取验证码失败', e);
    showToastMessage('获取验证码失败', 'error');
    showCaptchaModal.value = false;
    pendingProductForCaptcha.value = null;
  }
};

// 提交秒杀请求
const submitSeckill = async (productId, product) => {
  if (seckilling[productId]) return;
  
  seckilling[productId] = true;
  
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('未登录');
    }
    
    if (!product) {
      throw new Error('商品信息错误');
    }
    
    // 模拟秒杀请求，避免真实请求导致的服务器错误
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // 模拟秒杀成功
    product.stock--;
    
    showCaptchaModal.value = false;
    successProductName.value = product.name;
    showSuccessModal.value = true;
    showToastMessage('秒杀成功！', 'success');
  } catch (err) {
    console.error('秒杀失败:', err);
    showToastMessage(err.response?.data?.error || err.message || '秒杀失败，请重试', 'error');
  } finally {
    seckilling[productId] = false;
    showCaptchaModal.value = false;
    pendingProductForCaptcha.value = null;
  }
};

// 验证码确认按钮
const confirmCaptchaAndSeckill = async () => {
  if (isSubmitting.value) return;
  
  try {
    isSubmitting.value = true;
    
    if (!captchaInput.value) {
      showToastMessage('请输入验证码', 'warning');
      return;
    }
    if (!pendingProductForCaptcha.value) {
      showCaptchaModal.value = false;
      return;
    }
    
    const product = pendingProductForCaptcha.value;
    if (seckilling[product.id]) {
      showToastMessage('秒杀请求已提交，请稍候', 'warning');
      return;
    }
    
    await submitSeckill(product.id, product);
  } catch (error) {
    console.error('验证码确认失败:', error);
    showToastMessage('系统错误，请稍后重试', 'error');
    showCaptchaModal.value = false;
    pendingProductForCaptcha.value = null;
  } finally {
    isSubmitting.value = false;
  }
};

// 处理秒杀
const handleSeckill = async (product) => {
  try {
    if (!product || product.stock === 0) {
      showToastMessage('商品信息错误或已售罄', 'error');
      return;
    }

    const token = localStorage.getItem('token');
    if (!token) {
      showToastMessage('请先登录', 'error');
      router.push('/login');
      return;
    }

    pendingProductForCaptcha.value = product;
    captchaInput.value = '';
    await fetchCaptcha();
    showCaptchaModal.value = true;
  } catch (error) {
    console.error('秒杀处理失败:', error);
    showToastMessage('系统错误，请稍后重试', 'error');
  }
};

const addToCart = (product) => {
  // 添加到购物车逻辑
  const cart = JSON.parse(localStorage.getItem('cart') || '[]');
  const existingItem = cart.find(item => item.id === product.id);

  if (existingItem) {
    existingItem.quantity += 1;
  } else {
    cart.push({
      id: product.id,
      name: product.name,
      price: product.price,
      image: product.image,
      quantity: 1
    });
  }

  localStorage.setItem('cart', JSON.stringify(cart));
  showToastMessage('已添加到购物车', 'success');
};

onMounted(() => {
  loading.value = true;
  setTimeout(() => {
    loading.value = false;
  }, 500);
});
</script>

<style scoped>
@import '../styles/design-system.css';

.category-page {
  min-height: 100vh;
  background: var(--bg-secondary);
}

/* 页面头部 */
.page-header {
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-dark) 100%);
  color: white;
  padding: var(--space-12) 0 var(--space-8);
  margin-bottom: var(--space-6);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.header-title h1 {
  font-size: var(--text-4xl);
  font-weight: var(--font-bold);
  margin-bottom: var(--space-2);
}

.header-subtitle {
  font-size: var(--text-lg);
  opacity: 0.9;
}

.header-stats {
  display: flex;
  gap: var(--space-8);
}

.header-stats .stat-item {
  text-align: center;
}

.header-stats .stat-number {
  display: block;
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
}

.header-stats .stat-label {
  font-size: var(--text-sm);
  opacity: 0.8;
}

/* 筛选栏 */
.filter-bar {
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-light);
  padding: var(--space-4) 0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow-sm);
}

.filter-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

/* 分类标签 */
.category-tabs {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.tab-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.tab-btn.active {
  background: var(--color-primary);
  color: white;
}

.tab-icon {
  font-size: var(--text-base);
}

.tab-count {
  background: var(--color-primary-100);
  color: var(--color-primary);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
}

.tab-btn.active .tab-count {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

/* 筛选选项 */
.filter-options {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.sort-select select {
  padding: var(--space-2) var(--space-8) var(--space-2) var(--space-4);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-lg);
  background: var(--bg-primary);
  font-size: var(--text-sm);
  color: var(--text-primary);
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 8px center;
}

.filter-tags {
  display: flex;
  gap: var(--space-2);
}

.filter-tag {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-full);
  font-size: var(--text-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.filter-tag:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.filter-tag.active {
  background: var(--color-primary-50);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.tag-dot {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
  background: var(--border-dark);
}

.filter-tag.active .tag-dot {
  background: var(--color-primary);
}

/* 商品区域 */
.products-section {
  padding: var(--space-8) 0;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 3px solid var(--border-light);
  border-top-color: var(--color-primary);
  border-radius: var(--radius-full);
  animation: spin 1s linear infinite;
  margin-bottom: var(--space-4);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 商品网格 */
.products-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-6);
}

@media (min-width: 640px) {
  .products-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .products-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 1280px) {
  .products-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

/* 商品卡片 */
.product-card {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  transition: all var(--transition-base);
  cursor: pointer;
}

.product-card:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-4px);
}

.product-card.seckill {
  border-color: var(--color-seckill-light);
}

.product-card.out-of-stock {
  opacity: 0.7;
}

.product-card.out-of-stock .product-image img {
  filter: grayscale(100%);
}

/* 商品图片 */
.product-image {
  position: relative;
  aspect-ratio: 1;
  overflow: hidden;
  background: var(--bg-secondary);
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-slow);
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.product-badges {
  position: absolute;
  top: var(--space-3);
  left: var(--space-3);
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.badge {
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
}

.badge.seckill {
  background: var(--color-seckill);
  color: white;
}

.badge.soldout {
  background: var(--color-gray-600);
  color: white;
}

.badge.low-stock {
  background: var(--color-warning);
  color: white;
}

.discount-badge {
  position: absolute;
  top: var(--space-3);
  right: var(--space-3);
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
}

/* 商品信息 */
.product-info {
  padding: var(--space-4);
}

.product-name {
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-1);
  line-height: var(--leading-snug);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-desc {
  font-size: var(--text-sm);
  color: var(--text-tertiary);
  margin-bottom: var(--space-3);
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 价格 */
.product-price {
  display: flex;
  align-items: baseline;
  gap: var(--space-2);
  margin-bottom: var(--space-3);
}

.price-main {
  color: var(--color-seckill);
}

.price-symbol {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
}

.price-value {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
}

.price-original {
  font-size: var(--text-sm);
  color: var(--text-muted);
  text-decoration: line-through;
}

/* 销量和库存 */
.product-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}

.sales {
  font-size: var(--text-xs);
  color: var(--text-muted);
}

.stock-bar {
  width: 60px;
  height: 4px;
  background: var(--border-light);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.stock-progress {
  height: 100%;
  background: linear-gradient(90deg, var(--color-success) 0%, var(--color-warning) 50%, var(--color-error) 100%);
  border-radius: var(--radius-full);
  transition: width var(--transition-base);
}

/* 操作按钮 */
.product-actions {
  display: flex;
  gap: var(--space-2);
}

.product-actions .btn {
  flex: 1;
}

/* 空状态 */
.empty-state {
  padding: var(--space-16);
}

.empty-icon {
  font-size: var(--text-5xl);
  margin-bottom: var(--space-4);
}

.empty-title {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
}

.empty-desc {
  font-size: var(--text-base);
  color: var(--text-tertiary);
  margin-bottom: var(--space-6);
}

/* 响应式 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-6);
  }

  .filter-content {
    flex-direction: column;
    align-items: stretch;
  }

  .category-tabs {
    overflow-x: auto;
    flex-wrap: nowrap;
    padding-bottom: var(--space-2);
    -webkit-overflow-scrolling: touch;
  }

  .filter-options {
    justify-content: space-between;
  }
}

/* 弹窗样式 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow: hidden;
  animation: modalSlideIn 0.3s ease;
  padding: 32px;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 验证码弹窗 */
.captcha-content {
  text-align: center;
}

.captcha-content h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #1a1a2e;
}

.captcha-content p {
  color: #666;
  margin-bottom: 24px;
}

.captcha-row {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  align-items: center;
}

.captcha-image-wrapper {
  flex-shrink: 0;
  cursor: pointer;
  position: relative;
}

.captcha-image {
  width: 200px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
}

.captcha-refresh {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.captcha-image-wrapper:hover .captcha-refresh {
  opacity: 1;
}

.captcha-input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 16px;
  text-align: center;
}

.captcha-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.captcha-actions {
  display: flex;
  gap: 12px;
}

.captcha-actions .btn {
  flex: 1;
}

/* 成功弹窗 */
.success-content {
  text-align: center;
  padding: 40px 32px;
}

.success-animation {
  position: relative;
  margin-bottom: 24px;
}

.success-circle {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  position: relative;
  z-index: 1;
}

.success-icon {
  font-size: 40px;
  color: white;
  font-weight: bold;
}

.confetti {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;
}

.confetti span {
  position: absolute;
  width: 10px;
  height: 10px;
  background: var(--color-primary);
  border-radius: 50%;
  animation: confetti 2s ease-out forwards;
}

@keyframes confetti {
  0% {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
  }
  100% {
    transform: translateY(100px) rotate(360deg);
    opacity: 0;
  }
}

.confetti-1 { top: 10%; left: 30%; animation-delay: 0s; }
.confetti-2 { top: 20%; right: 25%; animation-delay: 0.2s; }
.confetti-3 { top: 30%; left: 20%; animation-delay: 0.4s; }
.confetti-4 { top: 40%; right: 30%; animation-delay: 0.6s; }
.confetti-5 { top: 50%; left: 25%; animation-delay: 0.8s; }
.confetti-6 { top: 60%; right: 20%; animation-delay: 1s; }
.confetti-7 { top: 70%; left: 30%; animation-delay: 1.2s; }
.confetti-8 { top: 80%; right: 25%; animation-delay: 1.4s; }
.confetti-9 { top: 90%; left: 20%; animation-delay: 1.6s; }
.confetti-10 { top: 100%; right: 30%; animation-delay: 1.8s; }

.success-content h3 {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #1a1a2e;
}

.success-content p {
  color: #666;
  margin-bottom: 24px;
  font-size: 16px;
}

.success-product {
  font-weight: 600;
  color: var(--color-primary);
}

.success-actions {
  display: flex;
  justify-content: center;
}

/* 提示消息 */
.toast {
  position: fixed;
  top: 100px;
  left: 50%;
  transform: translateX(-50%);
  background: #e8f5e9;
  color: #2e7d32;
  padding: 12px 24px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  z-index: 2000;
  opacity: 0;
  transition: all 0.3s ease;
  pointer-events: none;
}

.toast.show {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

.toast.error {
  background: #ffebee;
  color: #c62828;
}

.toast.warning {
  background: #fff3e0;
  color: #ef6c00;
}

.toast.info {
  background: #e3f2fd;
  color: #1565c0;
}

.toast-icon {
  font-size: 16px;
  font-weight: bold;
}

.toast-message {
  font-size: 14px;
  font-weight: 500;
}

/* 按钮样式 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-outline {
  background: white;
  border: 1px solid #e0e0e0;
  color: #666;
}

.btn-outline:hover {
  border-color: #667eea;
  color: #667eea;
}

.btn-spinner-white {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
