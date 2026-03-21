<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <!-- 关闭按钮 -->
      <button class="close-btn" @click="$emit('close')">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6L6 18M6 6l12 12"/>
        </svg>
      </button>

      <div class="modal-body">
        <!-- 左侧图片 -->
        <div class="product-gallery">
          <div class="main-image">
            <img :src="product.image" :alt="product.name">
            <div v-if="product.isSeckill" class="seckill-badge">秒杀</div>
          </div>
          <div class="thumbnail-list">
            <div class="thumb active">
              <img :src="product.image" :alt="product.name">
            </div>
            <div v-for="i in 3" :key="i" class="thumb">
              <img :src="getPlaceholderImage(i)" :alt="`${product.name} ${i}`">
            </div>
          </div>
        </div>

        <!-- 右侧信息 -->
        <div class="product-details">
          <div class="product-header">
            <h2 class="product-name">{{ product.name }}</h2>
            <p class="product-desc">{{ product.description }}</p>
          </div>

          <!-- 价格区域 -->
          <div class="price-section">
            <div class="current-price">
              <span class="price-symbol">¥</span>
              <span class="price-value">{{ formatPrice(product.price) }}</span>
            </div>
            <div v-if="product.originalPrice" class="original-price">
              <span class="label">原价</span>
              <span class="value">¥{{ formatPrice(product.originalPrice) }}</span>
            </div>
            <div v-if="product.discount" class="discount-tag">
              省 ¥{{ formatPrice(product.originalPrice - product.price) }}
            </div>
          </div>

          <!-- 库存状态 -->
          <div class="stock-section">
            <div class="stock-header">
              <span class="stock-label">库存状态</span>
              <span class="stock-value" :class="stockStatusClass">
                {{ stockStatusText }}
              </span>
            </div>
            <div class="stock-progress-bar">
              <div class="progress-fill" :style="{ width: stockPercentage + '%' }"></div>
            </div>
            <p class="stock-hint">已售 {{ product.sales || 0 }} 件</p>
          </div>

          <!-- 规格选择 -->
          <div class="specs-section">
            <h4>商品规格</h4>
            <div class="spec-options">
              <button
                v-for="spec in specs"
                :key="spec.id"
                :class="['spec-btn', { active: selectedSpec === spec.id }]"
                @click="selectedSpec = spec.id"
              >
                {{ spec.name }}
              </button>
            </div>
          </div>

          <!-- 数量选择 -->
          <div class="quantity-section">
            <h4>购买数量</h4>
            <div class="quantity-selector">
              <button class="qty-btn" @click="decreaseQty" :disabled="quantity <= 1">-</button>
              <span class="qty-value">{{ quantity }}</span>
              <button class="qty-btn" @click="increaseQty" :disabled="quantity >= maxQuantity">+</button>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="action-buttons">
            <button
              v-if="product.isSeckill && product.stock > 0"
              class="btn btn-seckill btn-lg"
              @click="handleSeckill"
            >
              <span class="btn-icon">⚡</span>
              立即秒杀
            </button>
            <button
              v-else-if="product.stock > 0"
              class="btn btn-primary btn-lg"
              @click="handleAddToCart"
            >
              <span class="btn-icon">🛒</span>
              加入购物车
            </button>
            <button v-else class="btn btn-secondary btn-lg" disabled>
              已售罄
            </button>
            <button class="btn btn-outline btn-lg" @click="$emit('close')">
              再看看
            </button>
          </div>

          <!-- 服务承诺 -->
          <div class="service-tags">
            <span class="service-tag">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              正品保障
            </span>
            <span class="service-tag">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              极速发货
            </span>
            <span class="service-tag">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 14l6-6m-5.5.5h.01m4.99 5h.01M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16l3.5-2 3.5 2 3.5-2 3.5 2z"/>
              </svg>
              7天无理由
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';

const props = defineProps({
  product: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['close', 'seckill', 'addToCart']);

// 状态
const quantity = ref(1);
const selectedSpec = ref('default');

// 规格选项
const specs = ref([
  { id: 'default', name: '标准版' },
  { id: 'pro', name: '专业版' },
  { id: 'max', name: '至尊版' }
]);

// 计算属性
const maxQuantity = computed(() => props.product.stock || 1);

const stockPercentage = computed(() => {
  if (!props.product.totalStock) return 0;
  return (props.product.stock / props.product.totalStock) * 100;
});

const stockStatusClass = computed(() => {
  if (props.product.stock === 0) return 'out-of-stock';
  if (props.product.stock < 10) return 'low-stock';
  return 'in-stock';
});

const stockStatusText = computed(() => {
  if (props.product.stock === 0) return '已售罄';
  if (props.product.stock < 10) return `仅剩 ${props.product.stock} 件`;
  return '库存充足';
});

// 方法
const formatPrice = (price) => {
  return price?.toLocaleString('zh-CN') || '0';
};

const getPlaceholderImage = (index) => {
  // 基于产品名称生成相关的实物图
  const productName = props.product.name || 'product';
  const angles = ['side view', 'back view', 'close up detail', 'in use'];
  const angle = angles[index % angles.length];
  return `https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=${encodeURIComponent(productName + ' ' + angle + ' on white background professional product photography')}&image_size=square`;
};

const increaseQty = () => {
  if (quantity.value < maxQuantity.value) {
    quantity.value++;
  }
};

const decreaseQty = () => {
  if (quantity.value > 1) {
    quantity.value--;
  }
};

const handleSeckill = () => {
  emit('seckill', { ...props.product, quantity: quantity.value, spec: selectedSpec.value });
};

const handleAddToCart = () => {
  emit('addToCart', { ...props.product, quantity: quantity.value, spec: selectedSpec.value });
};
</script>

<style scoped>
@import '../styles/design-system.css';

.modal-overlay {
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
  max-width: 1000px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
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

.close-btn {
  position: absolute;
  top: var(--space-4);
  right: var(--space-4);
  width: 40px;
  height: 40px;
  border: none;
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  z-index: 10;
}

.close-btn:hover {
  background: var(--bg-tertiary);
  transform: rotate(90deg);
}

.modal-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-8);
  padding: var(--space-8);
}

/* 商品图片区域 */
.product-gallery {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.main-image {
  position: relative;
  aspect-ratio: 1;
  border-radius: var(--radius-xl);
  overflow: hidden;
  background: var(--bg-secondary);
}

.main-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.seckill-badge {
  position: absolute;
  top: var(--space-4);
  left: var(--space-4);
  background: var(--color-seckill);
  color: white;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-lg);
  font-size: var(--text-sm);
  font-weight: var(--font-bold);
}

.thumbnail-list {
  display: flex;
  gap: var(--space-3);
}

.thumb {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-lg);
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all var(--transition-fast);
}

.thumb:hover {
  border-color: var(--border-medium);
}

.thumb.active {
  border-color: var(--color-primary);
}

.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 商品详情区域 */
.product-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.product-header {
  border-bottom: 1px solid var(--border-light);
  padding-bottom: var(--space-4);
}

.product-name {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
  line-height: var(--leading-tight);
}

.product-desc {
  font-size: var(--text-base);
  color: var(--text-secondary);
  line-height: var(--leading-relaxed);
}

/* 价格区域 */
.price-section {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  flex-wrap: wrap;
  padding: var(--space-4);
  background: linear-gradient(135deg, var(--color-error-50) 0%, var(--bg-primary) 100%);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-error-100);
}

.current-price {
  color: var(--color-seckill);
}

.price-symbol {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
}

.price-value {
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
}

.original-price {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.original-price .label {
  font-size: var(--text-xs);
  color: var(--text-muted);
}

.original-price .value {
  font-size: var(--text-base);
  color: var(--text-muted);
  text-decoration: line-through;
}

.discount-tag {
  background: var(--color-seckill);
  color: white;
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
}

/* 库存区域 */
.stock-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.stock-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stock-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

.stock-value {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
}

.stock-value.in-stock {
  color: var(--color-success);
}

.stock-value.low-stock {
  color: var(--color-warning);
}

.stock-value.out-of-stock {
  color: var(--color-error);
}

.stock-progress-bar {
  height: 8px;
  background: var(--border-light);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-success) 0%, var(--color-warning) 50%, var(--color-error) 100%);
  border-radius: var(--radius-full);
  transition: width var(--transition-base);
}

.stock-hint {
  font-size: var(--text-sm);
  color: var(--text-muted);
}

/* 规格选择 */
.specs-section h4,
.quantity-section h4 {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-3);
}

.spec-options {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.spec-btn {
  padding: var(--space-2) var(--space-4);
  border: 1px solid var(--border-medium);
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  font-size: var(--text-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.spec-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.spec-btn.active {
  background: var(--color-primary-50);
  border-color: var(--color-primary);
  color: var(--color-primary);
  font-weight: var(--font-semibold);
}

/* 数量选择 */
.quantity-selector {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-lg);
  width: fit-content;
}

.qty-btn {
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  font-size: var(--text-lg);
  color: var(--text-primary);
  cursor: pointer;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
}

.qty-btn:hover:not(:disabled) {
  background: var(--bg-secondary);
}

.qty-btn:disabled {
  color: var(--text-muted);
  cursor: not-allowed;
}

.qty-value {
  width: 50px;
  text-align: center;
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: var(--space-3);
  margin-top: var(--space-2);
}

.action-buttons .btn {
  flex: 1;
}

.btn-icon {
  margin-right: var(--space-2);
}

/* 服务标签 */
.service-tags {
  display: flex;
  gap: var(--space-4);
  padding-top: var(--space-4);
  border-top: 1px solid var(--border-light);
}

.service-tag {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.service-tag svg {
  color: var(--color-success);
}

/* 响应式 */
@media (max-width: 768px) {
  .modal-content {
    max-height: 100vh;
    border-radius: 0;
  }

  .modal-body {
    grid-template-columns: 1fr;
    gap: var(--space-6);
    padding: var(--space-6);
    padding-top: var(--space-12);
  }

  .product-gallery {
    order: -1;
  }

  .main-image {
    max-height: 300px;
  }

  .thumbnail-list {
    justify-content: center;
  }

  .price-section {
    flex-direction: column;
    align-items: flex-start;
  }

  .action-buttons {
    flex-direction: column;
  }

  .service-tags {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
