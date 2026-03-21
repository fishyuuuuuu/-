<template>
  <div class="upcoming-page">
    <!-- 顶部导航 -->
    <Navbar />

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <div class="header-content">
          <div class="header-title">
            <h1>即将开始</h1>
            <p class="header-subtitle">精彩秒杀即将开启，提前锁定心仪好物</p>
          </div>
          <div class="header-stats">
            <div class="stat-item">
              <span class="stat-number">{{ timeSlots.length }}</span>
              <span class="stat-label">今日场次</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">{{ upcomingProductsCount }}</span>
              <span class="stat-label">待秒杀商品</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 场次选择栏 -->
    <div class="slots-bar">
      <div class="container">
        <div class="slots-content">
          <div class="slots-timeline">
            <div
              v-for="(slot, index) in timeSlots"
              :key="slot.id"
              :class="['slot-item', {
                active: activeSlot === slot.id,
                current: slot.status === 'active',
                ended: slot.status === 'ended'
              }]"
              @click="activeSlot = slot.id"
            >
              <div class="slot-time">{{ slot.time }}</div>
              <div class="slot-status-badge" :class="slot.status">
                {{ getSlotStatusText(slot.status) }}
              </div>
              <div v-if="index < timeSlots.length - 1" class="slot-connector"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 倒计时横幅 -->
    <div class="countdown-banner" v-if="currentSlot">
      <div class="container">
        <div class="countdown-content">
          <div class="countdown-info">
            <span class="countdown-label">{{ currentSlot.status === 'active' ? '本场倒计时' : '距离开始' }}</span>
            <div class="countdown-timer">
              <div class="time-unit">
                <span class="time-value">{{ countdown.hours }}</span>
                <span class="time-label">时</span>
              </div>
              <span class="time-separator">:</span>
              <div class="time-unit">
                <span class="time-value">{{ countdown.minutes }}</span>
                <span class="time-label">分</span>
              </div>
              <span class="time-separator">:</span>
              <div class="time-unit">
                <span class="time-value">{{ countdown.seconds }}</span>
                <span class="time-label">秒</span>
              </div>
            </div>
          </div>
          <div class="countdown-actions">
            <button class="btn btn-primary" @click="scrollToProducts">
              <span class="btn-icon">👇</span>
              查看商品
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 商品列表 -->
    <div class="container" id="products-section">
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
            :class="{ 'reminded': isReminded(product.id) }"
          >
            <!-- 商品图片 -->
            <div class="product-image">
              <img :src="product.image" :alt="product.name" loading="lazy">
              <div class="product-overlay">
                <div class="countdown-small">
                  <span class="countdown-icon">⏰</span>
                  <span>{{ formatCountdown(product.countdown) }}</span>
                </div>
              </div>
              <div class="remind-badge" v-if="isReminded(product.id)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                已设置提醒
              </div>
            </div>

            <!-- 商品信息 -->
            <div class="product-info">
              <h3 class="product-name">{{ product.name }}</h3>
              <p class="product-desc">{{ product.description }}</p>

              <!-- 价格 -->
              <div class="product-price">
                <div class="price-seckill">
                  <span class="price-symbol">¥</span>
                  <span class="price-value">{{ formatPrice(product.seckillPrice) }}</span>
                </div>
                <div class="price-original">
                  <span class="label">原价</span>
                  <span class="value">¥{{ formatPrice(product.price) }}</span>
                </div>
                <div class="discount-badge">
                  省 ¥{{ formatPrice(product.price - product.seckillPrice) }}
                </div>
              </div>

              <!-- 库存进度 -->
              <div class="stock-section">
                <div class="stock-header">
                  <span class="stock-label">限量 {{ product.stock }} 件</span>
                  <span class="stock-ratio">{{ Math.round((product.sold || 0) / product.stock * 100) }}% 已预约</span>
                </div>
                <div class="stock-progress-bar">
                  <div class="stock-progress" :style="{ width: ((product.sold || 0) / product.stock * 100) + '%' }"></div>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="product-actions">
                <button
                  :class="['btn', isReminded(product.id) ? 'btn-secondary' : 'btn-primary', 'btn-lg']"
                  @click="toggleRemind(product)"
                  :disabled="currentSlot?.status === 'ended'"
                >
                  <span class="btn-icon">{{ isReminded(product.id) ? '🔔' : '🔕' }}</span>
                  {{ isReminded(product.id) ? '已提醒' : '提醒我' }}
                </button>
                <button
                  class="btn btn-seckill btn-lg"
                  v-if="currentSlot?.status === 'active'"
                  @click="goToSeckill(product)"
                >
                  <span class="btn-icon">⚡</span>
                  立即秒杀
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="empty-state">
          <div class="empty-icon">📅</div>
          <h3 class="empty-title">该场次暂无商品</h3>
          <p class="empty-desc">请选择其他场次查看</p>
          <button class="btn btn-primary" @click="activeSlot = timeSlots.find(s => s.status === 'active' || s.status === 'upcoming')?.id">
            查看其他场次
          </button>
        </div>
      </div>
    </div>

    <!-- 提醒成功提示 -->
    <transition name="fade">
      <div class="toast" v-if="showToast">
        <div class="toast-content">
          <span class="toast-icon">✓</span>
          <span class="toast-message">{{ toastMessage }}</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from './Navbar.vue';

const router = useRouter();

// 状态
const loading = ref(false);
const activeSlot = ref(1);
const showToast = ref(false);
const toastMessage = ref('');
const remindedProducts = ref(JSON.parse(localStorage.getItem('remindedProducts') || '[]'));

// 倒计时
const countdown = ref({ hours: '00', minutes: '00', seconds: '00' });
let countdownTimer = null;

// 场次数据
const timeSlots = ref([
  { id: 1, time: '10:00', status: 'ended' },
  { id: 2, time: '14:00', status: 'active' },
  { id: 3, time: '18:00', status: 'upcoming' },
  { id: 4, time: '20:00', status: 'upcoming' },
  { id: 5, time: '22:00', status: 'upcoming' }
]);

// 商品数据
const products = ref([
  {
    id: 1,
    name: 'iPhone 15 Pro',
    description: '钛金属设计，A17 Pro芯片，专业级摄像系统',
    price: 8999,
    seckillPrice: 7999,
    stock: 100,
    sold: 67,
    timeSlotId: 2,
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPhone%2015%20Pro%20titanium%20metal%20smartphone%20with%20triple%20camera%20system%20on%20white%20background%20professional%20product%20photography&image_size=square',
    countdown: { hours: 0, minutes: 45, seconds: 30 }
  },
  {
    id: 2,
    name: 'MacBook Pro 14',
    description: 'M3芯片，Liquid视网膜XDR显示屏',
    price: 14999,
    seckillPrice: 12999,
    stock: 50,
    sold: 32,
    timeSlotId: 2,
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=MacBook%20Pro%2014%20inch%20laptop%20with%20silver%20aluminum%20body%20open%20on%20white%20background%20professional%20product%20photography&image_size=square',
    countdown: { hours: 0, minutes: 45, seconds: 30 }
  },
  {
    id: 3,
    name: 'AirPods Pro 2',
    description: '主动降噪，空间音频，MagSafe充电盒',
    price: 2299,
    seckillPrice: 1899,
    stock: 200,
    sold: 156,
    timeSlotId: 2,
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=AirPods%20Pro%202%20wireless%20earbuds%20with%20charging%20case%20on%20white%20background%20professional%20product%20photography&image_size=square',
    countdown: { hours: 0, minutes: 45, seconds: 30 }
  },
  {
    id: 4,
    name: 'iPad Pro 12.9',
    description: 'M2芯片，12.9英寸Liquid视网膜XDR屏',
    price: 9999,
    seckillPrice: 8999,
    stock: 80,
    sold: 45,
    timeSlotId: 3,
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPad%20Pro%2012.9%20inch%20tablet%20with%20Apple%20Pencil%20on%20white%20background%20professional%20product%20photography&image_size=square',
    countdown: { hours: 4, minutes: 45, seconds: 30 }
  },
  {
    id: 5,
    name: 'Apple Watch Ultra 2',
    description: '钛金属表壳，双频GPS，100米防水',
    price: 6999,
    seckillPrice: 6299,
    stock: 60,
    sold: 28,
    timeSlotId: 3,
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Apple%20Watch%20Ultra%20titanium%20smartwatch%20on%20white%20background%20professional%20product%20photography&image_size=square',
    countdown: { hours: 4, minutes: 45, seconds: 30 }
  },
  {
    id: 6,
    name: 'Sony WH-1000XM5',
    description: '业界领先的降噪耳机，30小时续航',
    price: 2999,
    seckillPrice: 2499,
    stock: 150,
    sold: 89,
    timeSlotId: 4,
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Sony%20WH-1000XM5%20black%20over-ear%20headphones%20on%20white%20background%20professional%20product%20photography&image_size=square',
    countdown: { hours: 6, minutes: 45, seconds: 30 }
  }
]);

// 计算属性
const currentSlot = computed(() => {
  return timeSlots.value.find(slot => slot.id === activeSlot.value);
});

const filteredProducts = computed(() => {
  return products.value.filter(product => product.timeSlotId === activeSlot.value);
});

const upcomingProductsCount = computed(() => {
  return products.value.filter(p => {
    const slot = timeSlots.value.find(s => s.id === p.timeSlotId);
    return slot && slot.status !== 'ended';
  }).length;
});

// 方法
const getSlotStatusText = (status) => {
  const statusMap = {
    'active': '进行中',
    'upcoming': '即将开始',
    'ended': '已结束'
  };
  return statusMap[status] || status;
};

const formatPrice = (price) => {
  return price?.toLocaleString('zh-CN') || '0';
};

const formatCountdown = (cd) => {
  if (!cd) return '00:00:00';
  const h = String(cd.hours).padStart(2, '0');
  const m = String(cd.minutes).padStart(2, '0');
  const s = String(cd.seconds).padStart(2, '0');
  return `${h}:${m}:${s}`;
};

const isReminded = (productId) => {
  return remindedProducts.value.includes(productId);
};

const toggleRemind = (product) => {
  if (currentSlot.value?.status === 'ended') return;

  const index = remindedProducts.value.indexOf(product.id);
  if (index > -1) {
    remindedProducts.value.splice(index, 1);
    showToastMessage('已取消提醒');
  } else {
    remindedProducts.value.push(product.id);
    showToastMessage(`已设置提醒，${product.name} 开始时会通知您`);
  }
  localStorage.setItem('remindedProducts', JSON.stringify(remindedProducts.value));
};

const showToastMessage = (message) => {
  toastMessage.value = message;
  showToast.value = true;
  setTimeout(() => {
    showToast.value = false;
  }, 3000);
};

const goToSeckill = (product) => {
  router.push(`/seckill/${product.id}`);
};

const scrollToProducts = () => {
  document.getElementById('products-section')?.scrollIntoView({ behavior: 'smooth' });
};

// 倒计时更新
const updateCountdown = () => {
  const now = new Date();
  const currentHour = now.getHours();
  const currentMinute = now.getMinutes();
  const currentSecond = now.getSeconds();

  // 找到当前选中的场次时间
  const slot = timeSlots.value.find(s => s.id === activeSlot.value);
  if (!slot) return;

  const [targetHour, targetMinute] = slot.time.split(':').map(Number);

  let targetTime = new Date();
  targetTime.setHours(targetHour, targetMinute, 0, 0);

  if (slot.status === 'ended') {
    countdown.value = { hours: '00', minutes: '00', seconds: '00' };
    return;
  }

  if (slot.status === 'active') {
    // 进行中场次，倒计时到结束（假设每场2小时）
    targetTime = new Date(targetTime.getTime() + 2 * 60 * 60 * 1000);
  }

  let diff = targetTime - now;
  if (diff < 0) diff = 0;

  const hours = Math.floor(diff / (1000 * 60 * 60));
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((diff % (1000 * 60)) / 1000);

  countdown.value = {
    hours: String(hours).padStart(2, '0'),
    minutes: String(minutes).padStart(2, '0'),
    seconds: String(seconds).padStart(2, '0')
  };
};

// 生命周期
onMounted(() => {
  loading.value = true;
  setTimeout(() => {
    loading.value = false;
  }, 500);

  // 启动倒计时
  countdownTimer = setInterval(updateCountdown, 1000);
  updateCountdown();

  // 默认选中进行中的场次
  const activeSlotItem = timeSlots.value.find(s => s.status === 'active');
  if (activeSlotItem) {
    activeSlot.value = activeSlotItem.id;
  }
});

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer);
  }
});
</script>

<style scoped>
@import '../styles/design-system.css';

.upcoming-page {
  min-height: 100vh;
  background: var(--bg-secondary);
}

/* 页面头部 */
.page-header {
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-dark) 100%);
  color: white;
  padding: var(--space-12) 0 var(--space-8);
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

.stat-item {
  text-align: center;
}

.stat-number {
  display: block;
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
}

.stat-label {
  font-size: var(--text-sm);
  opacity: 0.8;
}

/* 场次选择栏 */
.slots-bar {
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-light);
  padding: var(--space-6) 0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow-sm);
}

.slots-timeline {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
}

.slot-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-4) var(--space-6);
  border-radius: var(--radius-xl);
  cursor: pointer;
  transition: all var(--transition-base);
  position: relative;
  min-width: 100px;
}

.slot-item:hover:not(.ended) {
  background: var(--bg-secondary);
}

.slot-item.active {
  background: var(--color-primary);
  color: white;
  box-shadow: var(--shadow-primary);
}

.slot-item.ended {
  opacity: 0.5;
  cursor: not-allowed;
}

.slot-time {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
}

.slot-status-badge {
  font-size: var(--text-xs);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  font-weight: var(--font-medium);
}

.slot-status-badge.active {
  background: var(--color-success);
  color: white;
}

.slot-status-badge.upcoming {
  background: var(--color-warning-50);
  color: var(--color-warning);
}

.slot-status-badge.ended {
  background: var(--color-gray-200);
  color: var(--color-gray-500);
}

.slot-item.active .slot-status-badge {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.slot-connector {
  position: absolute;
  right: -20px;
  top: 50%;
  width: 16px;
  height: 2px;
  background: var(--border-medium);
}

/* 倒计时横幅 */
.countdown-banner {
  background: linear-gradient(135deg, var(--color-warning-50) 0%, var(--bg-primary) 100%);
  border-bottom: 1px solid var(--color-warning-100);
  padding: var(--space-6) 0;
}

.countdown-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.countdown-info {
  display: flex;
  align-items: center;
  gap: var(--space-6);
}

.countdown-label {
  font-size: var(--text-base);
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

.countdown-timer {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.time-unit {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-1);
}

.time-value {
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--color-seckill);
  background: white;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  min-width: 60px;
  text-align: center;
}

.time-label {
  font-size: var(--text-xs);
  color: var(--text-muted);
}

.time-separator {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-seckill);
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

/* 商品卡片 */
.product-card {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  transition: all var(--transition-base);
}

.product-card:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-4px);
}

.product-card.reminded {
  border-color: var(--color-primary-light);
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

.product-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7) 0%, transparent 50%);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: var(--space-4);
  opacity: 0;
  transition: opacity var(--transition-base);
}

.product-card:hover .product-overlay {
  opacity: 1;
}

.countdown-small {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  background: rgba(255, 255, 255, 0.95);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-full);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-seckill);
}

.countdown-icon {
  font-size: var(--text-base);
}

.remind-badge {
  position: absolute;
  top: var(--space-3);
  right: var(--space-3);
  background: var(--color-success);
  color: white;
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  display: flex;
  align-items: center;
  gap: var(--space-1);
}

/* 商品信息 */
.product-info {
  padding: var(--space-5);
}

.product-name {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
  line-height: var(--leading-snug);
}

.product-desc {
  font-size: var(--text-sm);
  color: var(--text-tertiary);
  margin-bottom: var(--space-4);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 价格 */
.product-price {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
  flex-wrap: wrap;
}

.price-seckill {
  color: var(--color-seckill);
}

.price-seckill .price-symbol {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
}

.price-seckill .price-value {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
}

.price-original {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.price-original .label {
  font-size: var(--text-xs);
  color: var(--text-muted);
}

.price-original .value {
  font-size: var(--text-sm);
  color: var(--text-muted);
  text-decoration: line-through;
}

.discount-badge {
  background: var(--color-seckill);
  color: white;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
}

/* 库存 */
.stock-section {
  margin-bottom: var(--space-4);
}

.stock-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--space-2);
  font-size: var(--text-sm);
}

.stock-label {
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

.stock-ratio {
  color: var(--color-primary);
  font-weight: var(--font-semibold);
}

.stock-progress-bar {
  height: 6px;
  background: var(--border-light);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.stock-progress {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary) 0%, var(--color-primary-light) 100%);
  border-radius: var(--radius-full);
  transition: width var(--transition-base);
}

/* 操作按钮 */
.product-actions {
  display: flex;
  gap: var(--space-3);
}

.product-actions .btn {
  flex: 1;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: var(--space-16);
}

.empty-icon {
  font-size: var(--text-6xl);
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

/* Toast提示 */
.toast {
  position: fixed;
  top: var(--space-6);
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.toast-content {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  background: var(--color-success);
  color: white;
  padding: var(--space-4) var(--space-6);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
}

.toast-icon {
  font-size: var(--text-xl);
}

.toast-message {
  font-size: var(--text-base);
  font-weight: var(--font-medium);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-6);
  }

  .slots-timeline {
    overflow-x: auto;
    justify-content: flex-start;
    padding-bottom: var(--space-2);
    -webkit-overflow-scrolling: touch;
  }

  .slot-item {
    min-width: 80px;
    padding: var(--space-3) var(--space-4);
  }

  .slot-connector {
    display: none;
  }

  .countdown-content {
    flex-direction: column;
    gap: var(--space-4);
    text-align: center;
  }

  .countdown-info {
    flex-direction: column;
    gap: var(--space-4);
  }

  .time-value {
    font-size: var(--text-2xl);
    min-width: 50px;
  }

  .product-actions {
    flex-direction: column;
  }
}
</style>
