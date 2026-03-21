<template>
  <div class="order-page">
    <!-- 顶部导航 -->
    <Navbar />

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <div class="header-content">
          <div class="header-title">
            <h1>我的订单</h1>
            <p class="header-subtitle">查看和管理您的所有订单</p>
          </div>
          <div class="header-stats">
            <div class="stat-item">
              <span class="stat-number">{{ orders.length }}</span>
              <span class="stat-label">全部订单</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">{{ pendingOrders.length }}</span>
              <span class="stat-label">待支付</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">{{ shippedOrders.length }}</span>
              <span class="stat-label">待收货</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 订单标签栏 -->
    <div class="tabs-bar">
      <div class="container">
        <div class="tabs-content">
          <div class="order-tabs">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              :class="['tab-btn', { active: activeTab === tab.key }]"
              @click="activeTab = tab.key"
            >
              <span class="tab-name">{{ tab.name }}</span>
              <span v-if="tab.count > 0" class="tab-badge">{{ tab.count }}</span>
            </button>
          </div>
          <div class="search-box">
            <input
              type="text"
              v-model="searchQuery"
              placeholder="搜索订单号或商品名称"
              class="search-input"
            >
            <button class="search-btn" @click="handleSearch">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <path d="m21 21-4.35-4.35"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 订单列表 -->
    <div class="container">
      <div class="orders-section">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <p>正在加载订单...</p>
        </div>

        <!-- 订单列表 -->
        <div v-else-if="filteredOrders.length > 0" class="orders-list">
          <div
            v-for="order in filteredOrders"
            :key="order.id"
            class="order-card"
          >
            <!-- 订单头部 -->
            <div class="order-header">
              <div class="order-meta">
                <span class="order-id">订单号: {{ order.id }}</span>
                <span class="order-date">{{ formatDate(order.createdAt) }}</span>
              </div>
              <div class="order-status" :class="order.status">
                <span class="status-dot"></span>
                <span class="status-text">{{ order.statusText }}</span>
              </div>
            </div>

            <!-- 订单商品 -->
            <div class="order-body">
              <div class="product-list">
                <div
                  v-for="(item, index) in order.items"
                  :key="index"
                  class="product-item"
                >
                  <img :src="item.image" :alt="item.name" class="product-image">
                  <div class="product-info">
                    <h4 class="product-name">{{ item.name }}</h4>
                    <p class="product-spec">{{ item.spec || '标准版' }}</p>
                    <div class="product-price-qty">
                      <span class="product-price">¥{{ formatPrice(item.price) }}</span>
                      <span class="product-qty">x{{ item.quantity }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 订单底部 -->
            <div class="order-footer">
              <div class="order-total">
                <span class="total-label">实付金额:</span>
                <span class="total-price">
                  <span class="price-symbol">¥</span>
                  <span class="price-value">{{ formatPrice(order.totalPrice) }}</span>
                </span>
              </div>
              <div class="order-actions">
                <button
                  v-if="order.status === 'pending'"
                  class="btn btn-primary"
                  @click="payOrder(order)"
                >
                  立即支付
                </button>
                <button
                  v-if="order.status === 'pending'"
                  class="btn btn-outline"
                  @click="cancelOrder(order)"
                >
                  取消订单
                </button>
                <button
                  v-if="order.status === 'shipped'"
                  class="btn btn-primary"
                  @click="confirmReceive(order)"
                >
                  确认收货
                </button>
                <button
                  v-if="order.status === 'shipped'"
                  class="btn btn-outline"
                  @click="viewLogistics(order)"
                >
                  查看物流
                </button>
                <button
                  v-if="['paid', 'completed'].includes(order.status)"
                  class="btn btn-outline"
                  @click="viewDetail(order)"
                >
                  查看详情
                </button>
                <button
                  v-if="order.status === 'completed'"
                  class="btn btn-secondary"
                  @click="buyAgain(order)"
                >
                  再次购买
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="empty-state">
          <div class="empty-icon">📋</div>
          <h3 class="empty-title">暂无订单</h3>
          <p class="empty-desc">您还没有相关订单，快去选购心仪商品吧</p>
          <button class="btn btn-primary btn-lg" @click="goToShopping">
            <span class="btn-icon">🛒</span>
            去购物
          </button>
        </div>
      </div>
    </div>

    <!-- 支付弹窗 -->
    <div class="modal" v-if="showPayModal" @click.self="showPayModal = false">
      <div class="modal-content pay-modal">
        <div class="modal-header">
          <h3>订单支付</h3>
          <button class="close-btn" @click="showPayModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="pay-amount">
            <span class="pay-label">支付金额</span>
            <span class="pay-price">
              <span class="price-symbol">¥</span>
              <span class="price-value">{{ currentOrder ? formatPrice(currentOrder.totalPrice) : '0' }}</span>
            </span>
          </div>
          <div class="pay-methods">
            <div
              v-for="method in payMethods"
              :key="method.id"
              :class="['pay-method', { selected: selectedPayMethod === method.id }]"
              @click="selectedPayMethod = method.id"
            >
              <div class="method-icon" :style="{ backgroundColor: method.color + '20', color: method.color }">
                {{ method.icon }}
              </div>
              <span class="method-name">{{ method.name }}</span>
              <div class="method-radio">
                <span class="radio-inner" v-if="selectedPayMethod === method.id">✓</span>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary btn-lg" @click="confirmPay" :disabled="isPaying">
            <span v-if="!isPaying">确认支付</span>
            <span v-else class="btn-spinner-white"></span>
          </button>
        </div>
      </div>
    </div>

    <!-- 取消订单确认弹窗 -->
    <div class="modal" v-if="showCancelModal" @click.self="showCancelModal = false">
      <div class="modal-content confirm-modal">
        <div class="modal-header">
          <h3>取消订单</h3>
          <button class="close-btn" @click="showCancelModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="confirm-icon">⚠️</div>
          <p class="confirm-text">确定要取消订单 {{ currentOrder?.id }} 吗？</p>
          <p class="confirm-hint">取消后无法恢复，优惠券将返还至您的账户</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showCancelModal = false">再想想</button>
          <button class="btn btn-error" @click="confirmCancel" :disabled="isCancelling">
            <span v-if="!isCancelling">确认取消</span>
            <span v-else class="btn-spinner-white"></span>
          </button>
        </div>
      </div>
    </div>

    <!-- 提示消息 -->
    <transition name="fade">
      <div class="toast" v-if="showToast">
        <div class="toast-content" :class="toastType">
          <span class="toast-icon">{{ toastIcon }}</span>
          <span class="toast-message">{{ toastMessage }}</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from './Navbar.vue';

const router = useRouter();

// 状态
const loading = ref(false);
const activeTab = ref('all');
const searchQuery = ref('');
const showPayModal = ref(false);
const showCancelModal = ref(false);
const currentOrder = ref(null);
const selectedPayMethod = ref('alipay');
const isPaying = ref(false);
const isCancelling = ref(false);
const showToast = ref(false);
const toastMessage = ref('');
const toastType = ref('success');

// 标签页
const tabs = computed(() => [
  { key: 'all', name: '全部', count: orders.value.length },
  { key: 'pending', name: '待支付', count: pendingOrders.value.length },
  { key: 'paid', name: '待发货', count: paidOrders.value.length },
  { key: 'shipped', name: '待收货', count: shippedOrders.value.length },
  { key: 'completed', name: '已完成', count: completedOrders.value.length }
]);

// 支付方式
const payMethods = ref([
  { id: 'alipay', name: '支付宝', icon: '💳', color: '#1677ff' },
  { id: 'wechat', name: '微信支付', icon: '💬', color: '#07c160' },
  { id: 'card', name: '银行卡', icon: '💰', color: '#f59e0b' }
]);

// 订单数据
const orders = ref([
  {
    id: '20260320001',
    items: [
      {
        name: 'iPhone 15 Pro',
        price: 7999,
        quantity: 1,
        image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPhone%2015%20Pro%20titanium%20metal%20smartphone%20with%20triple%20camera%20system%20on%20white%20background%20professional%20product%20photography&image_size=square',
        spec: '256GB 钛金属'
      }
    ],
    totalPrice: 7999,
    status: 'pending',
    statusText: '待支付',
    createdAt: '2026-03-20 10:30:00'
  },
  {
    id: '20260319002',
    items: [
      {
        name: 'MacBook Pro 14',
        price: 12999,
        quantity: 1,
        image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=MacBook%20Pro%2014%20inch%20laptop%20with%20silver%20aluminum%20body%20open%20on%20white%20background%20professional%20product%20photography&image_size=square',
        spec: 'M3 Pro 18GB+512GB'
      },
      {
        name: 'AirPods Pro 2',
        price: 1899,
        quantity: 1,
        image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=AirPods%20Pro%202%20wireless%20earbuds%20with%20charging%20case%20on%20white%20background%20professional%20product%20photography&image_size=square',
        spec: 'MagSafe充电盒'
      }
    ],
    totalPrice: 14898,
    status: 'paid',
    statusText: '待发货',
    createdAt: '2026-03-19 15:45:00'
  },
  {
    id: '20260318003',
    items: [
      {
        name: 'iPad Pro 12.9',
        price: 8999,
        quantity: 1,
        image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPad%20Pro%2012.9%20inch%20tablet%20with%20Apple%20Pencil%20on%20white%20background%20professional%20product%20photography&image_size=square',
        spec: 'M2 256GB WiFi'
      }
    ],
    totalPrice: 8999,
    status: 'shipped',
    statusText: '待收货',
    createdAt: '2026-03-18 09:20:00'
  },
  {
    id: '20260315004',
    items: [
      {
        name: 'Apple Watch Ultra 2',
        price: 6299,
        quantity: 1,
        image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Apple%20Watch%20Ultra%20titanium%20smartwatch%20on%20white%20background%20professional%20product%20photography&image_size=square',
        spec: '钛金属表壳 海洋表带'
      }
    ],
    totalPrice: 6299,
    status: 'completed',
    statusText: '已完成',
    createdAt: '2026-03-15 14:00:00'
  }
]);

// 计算属性
const pendingOrders = computed(() => orders.value.filter(o => o.status === 'pending'));
const paidOrders = computed(() => orders.value.filter(o => o.status === 'paid'));
const shippedOrders = computed(() => orders.value.filter(o => o.status === 'shipped'));
const completedOrders = computed(() => orders.value.filter(o => o.status === 'completed'));

const filteredOrders = computed(() => {
  let result = orders.value;

  // 按状态筛选
  if (activeTab.value !== 'all') {
    result = result.filter(order => order.status === activeTab.value);
  }

  // 按搜索词筛选
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(order =>
      order.id.toLowerCase().includes(query) ||
      order.items.some(item => item.name.toLowerCase().includes(query))
    );
  }

  // 按时间倒序
  return result.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
});

const toastIcon = computed(() => {
  const icons = {
    success: '✓',
    error: '✗',
    warning: '⚠️'
  };
  return icons[toastType.value] || '✓';
});

// 方法
const formatPrice = (price) => {
  return price?.toLocaleString('zh-CN') || '0';
};

const formatDate = (dateStr) => {
  const date = new Date(dateStr);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
};

const handleSearch = () => {
  // 搜索逻辑已在计算属性中处理
};

const payOrder = (order) => {
  currentOrder.value = order;
  showPayModal.value = true;
};

const confirmPay = () => {
  isPaying.value = true;
  setTimeout(() => {
    isPaying.value = false;
    showPayModal.value = false;
    if (currentOrder.value) {
      currentOrder.value.status = 'paid';
      currentOrder.value.statusText = '待发货';
    }
    showToastMessage('支付成功！', 'success');
  }, 1500);
};

const cancelOrder = (order) => {
  currentOrder.value = order;
  showCancelModal.value = true;
};

const confirmCancel = () => {
  isCancelling.value = true;
  setTimeout(() => {
    isCancelling.value = false;
    showCancelModal.value = false;
    if (currentOrder.value) {
      const index = orders.value.findIndex(o => o.id === currentOrder.value.id);
      if (index > -1) {
        orders.value.splice(index, 1);
      }
    }
    showToastMessage('订单已取消', 'success');
  }, 1000);
};

const confirmReceive = (order) => {
  order.status = 'completed';
  order.statusText = '已完成';
  showToastMessage('确认收货成功！', 'success');
};

const viewLogistics = (order) => {
  showToastMessage('物流信息功能开发中...', 'warning');
};

const viewDetail = (order) => {
  showToastMessage('订单详情功能开发中...', 'warning');
};

const buyAgain = (order) => {
  showToastMessage('已加入购物车', 'success');
};

const goToShopping = () => {
  router.push('/category');
};

const showToastMessage = (message, type = 'success') => {
  toastMessage.value = message;
  toastType.value = type;
  showToast.value = true;
  setTimeout(() => {
    showToast.value = false;
  }, 3000);
};

// 生命周期
onMounted(() => {
  loading.value = true;
  setTimeout(() => {
    loading.value = false;
  }, 500);
});
</script>

<style scoped>
@import '../styles/design-system.css';

.order-page {
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

/* 标签栏 */
.tabs-bar {
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-light);
  padding: var(--space-4) 0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow-sm);
}

.tabs-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

.order-tabs {
  display: flex;
  gap: var(--space-2);
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

.tab-badge {
  background: var(--color-error);
  color: white;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  min-width: 18px;
  text-align: center;
}

.tab-btn.active .tab-badge {
  background: rgba(255, 255, 255, 0.3);
}

/* 搜索框 */
.search-box {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.search-input {
  padding: var(--space-2) var(--space-4);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-lg);
  font-size: var(--text-sm);
  width: 250px;
  transition: all var(--transition-fast);
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-100);
}

.search-btn {
  padding: var(--space-2) var(--space-3);
  border: none;
  background: var(--color-primary);
  color: white;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.search-btn:hover {
  background: var(--color-primary-dark);
}

/* 订单区域 */
.orders-section {
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

/* 订单列表 */
.orders-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* 订单卡片 */
.order-card {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  overflow: hidden;
  transition: all var(--transition-base);
}

.order-card:hover {
  box-shadow: var(--shadow-md);
}

/* 订单头部 */
.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) var(--space-6);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-light);
}

.order-meta {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.order-id {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

.order-date {
  font-size: var(--text-sm);
  color: var(--text-muted);
}

.order-status {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
}

.order-status.pending {
  color: var(--color-warning);
}

.order-status.pending .status-dot {
  background: var(--color-warning);
  animation: pulse 2s infinite;
}

.order-status.paid {
  color: var(--color-primary);
}

.order-status.paid .status-dot {
  background: var(--color-primary);
}

.order-status.shipped {
  color: var(--color-info);
}

.order-status.shipped .status-dot {
  background: var(--color-info);
}

.order-status.completed {
  color: var(--color-success);
}

.order-status.completed .status-dot {
  background: var(--color-success);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* 订单商品 */
.order-body {
  padding: var(--space-6);
}

.product-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.product-item {
  display: flex;
  gap: var(--space-4);
}

.product-image {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-lg);
  object-fit: cover;
  background: var(--bg-secondary);
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.product-name {
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-1);
}

.product-spec {
  font-size: var(--text-sm);
  color: var(--text-tertiary);
  margin-bottom: var(--space-2);
}

.product-price-qty {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.product-price {
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--color-seckill);
}

.product-qty {
  font-size: var(--text-sm);
  color: var(--text-muted);
}

/* 订单底部 */
.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) var(--space-6);
  border-top: 1px solid var(--border-light);
  background: var(--bg-secondary);
}

.order-total {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.total-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.total-price {
  color: var(--color-seckill);
}

.total-price .price-symbol {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
}

.total-price .price-value {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
}

.order-actions {
  display: flex;
  gap: var(--space-3);
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

/* 弹窗 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-4);
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: var(--bg-primary);
  border-radius: var(--radius-2xl);
  max-width: 500px;
  width: 100%;
  animation: slideUp 0.3s ease;
  box-shadow: var(--shadow-2xl);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-light);
}

.modal-header h3 {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  color: var(--text-primary);
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  font-size: var(--text-xl);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.close-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.modal-body {
  padding: var(--space-6);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-6);
  border-top: 1px solid var(--border-light);
}

/* 支付弹窗 */
.pay-amount {
  text-align: center;
  margin-bottom: var(--space-6);
  padding-bottom: var(--space-6);
  border-bottom: 1px solid var(--border-light);
}

.pay-label {
  display: block;
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-bottom: var(--space-2);
}

.pay-price {
  color: var(--color-seckill);
}

.pay-price .price-symbol {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
}

.pay-price .price-value {
  font-size: var(--text-4xl);
  font-weight: var(--font-bold);
}

.pay-methods {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.pay-method {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-4);
  border: 2px solid var(--border-light);
  border-radius: var(--radius-xl);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.pay-method:hover {
  border-color: var(--border-medium);
}

.pay-method.selected {
  border-color: var(--color-primary);
  background: var(--color-primary-50);
}

.method-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-2xl);
}

.method-name {
  flex: 1;
  font-size: var(--text-base);
  font-weight: var(--font-medium);
  color: var(--text-primary);
}

.method-radio {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border-medium);
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.pay-method.selected .method-radio {
  border-color: var(--color-primary);
  background: var(--color-primary);
}

.radio-inner {
  color: white;
  font-size: var(--text-sm);
  font-weight: var(--font-bold);
}

/* 确认弹窗 */
.confirm-modal .modal-body {
  text-align: center;
}

.confirm-icon {
  font-size: var(--text-5xl);
  margin-bottom: var(--space-4);
}

.confirm-text {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
}

.confirm-hint {
  font-size: var(--text-sm);
  color: var(--text-tertiary);
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
  padding: var(--space-4) var(--space-6);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  color: white;
}

.toast-content.success {
  background: var(--color-success);
}

.toast-content.error {
  background: var(--color-error);
}

.toast-content.warning {
  background: var(--color-warning);
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

  .tabs-content {
    flex-direction: column;
    align-items: stretch;
  }

  .order-tabs {
    overflow-x: auto;
    flex-wrap: nowrap;
    padding-bottom: var(--space-2);
    -webkit-overflow-scrolling: touch;
  }

  .search-box {
    width: 100%;
  }

  .search-input {
    flex: 1;
    width: auto;
  }

  .order-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-2);
  }

  .order-footer {
    flex-direction: column;
    gap: var(--space-4);
  }

  .order-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .modal-content {
    margin: var(--space-4);
  }
}
</style>
