<template>
  <div class="seckill-container">
    <!-- 导航栏 -->
    <Navbar />

    <!-- 秒杀活动横幅 -->
    <div class="seckill-banner">
      <div class="banner-bg-animation"></div>
      <div class="container">
        <div class="banner-content">
          <h2 class="banner-title">
            <span class="title-highlight">限时秒杀</span>
            <span class="title-pulse">🔥</span>
          </h2>
          <p class="banner-subtitle">每日爆款，限时抢购，手慢无！</p>
          <div class="countdown-section">
            <span class="countdown-label">距离结束:</span>
            <div class="countdown" v-if="countdown">
              <div class="time-box">
                <span class="time-value">{{ countdown.hours }}</span>
                <span class="time-label">时</span>
              </div>
              <span class="colon">:</span>
              <div class="time-box">
                <span class="time-value">{{ countdown.minutes }}</span>
                <span class="time-label">分</span>
              </div>
              <span class="colon">:</span>
              <div class="time-box">
                <span class="time-value">{{ countdown.seconds }}</span>
                <span class="time-label">秒</span>
              </div>
            </div>
          </div>
          <div class="banner-stats">
            <div class="stat-item">
              <span class="stat-number">{{ products.length }}</span>
              <span class="stat-label">爆款商品</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">{{ totalStock }}</span>
              <span class="stat-label">剩余库存</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">10K+</span>
              <span class="stat-label">活跃用户</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 商品分类 -->
    <div class="category-section">
      <div class="container">
        <div class="categories">
          <div 
            v-for="category in categories" 
            :key="category.id"
            :class="{ active: activeCategory === category.id }"
            class="category-item"
            @click="activeCategory = category.id"
          >
            <span class="category-icon">{{ category.icon }}</span>
            <span class="category-name">{{ category.name }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 商品列表区域 -->
    <div class="container">
      <main class="product-list">
        <div class="section-header">
          <div class="section-title-wrapper">
            <h2 class="section-title">🔥 秒杀商品</h2>
            <span class="section-subtitle">限时特惠，抢购从速</span>
          </div>
          <div class="sort-options">
            <span class="sort-label">排序:</span>
            <select class="sort-select" v-model="sortBy" @change="sortProducts">
              <option value="default">默认推荐</option>
              <option value="price-asc">价格从低到高</option>
              <option value="price-desc">价格从高到低</option>
              <option value="stock">库存优先</option>
              <option value="discount">折扣力度</option>
            </select>
          </div>
        </div>

        <!-- 加载状态 -->
        <div class="loading" v-if="loading">
          <div class="spinner"></div>
          <p>正在加载商品...</p>
        </div>

        <!-- 错误状态 -->
        <div class="error-message" v-if="error">
          <i class="icon-error">⚠️</i>
          <p>{{ error }}</p>
          <button @click="fetchProducts" class="retry-btn">重试</button>
        </div>

        <!-- 商品列表 -->
        <div class="products-grid" v-if="!loading && !error && sortedProducts.length">
          <div 
            class="product-card" 
            v-for="product in sortedProducts" 
            :key="product.id"
            :class="{ 'sold-out': product.stock === 0, 'low-stock': product.stock > 0 && product.stock <= 10 }"
          >
            <div class="product-badge" v-if="product.stock <= 10 && product.stock > 0">
              仅剩{{ product.stock }}件
            </div>
            <div class="product-badge sold-out-badge" v-if="product.stock === 0">
              已抢完
            </div>
            <div class="product-image" @click="showProductDetail(product)">
              <div class="image-container">
                <img 
                  :src="product.image" 
                  :alt="product.name" 
                  loading="lazy"
                  @error="handleImageError(product)"
                  @load="handleImageLoad(product)"
                >
                <div v-if="imageLoading[product.id]" class="image-loading">
                  <div class="loading-spinner"></div>
                  <span>加载中...</span>
                </div>
                <div v-else-if="imageError[product.id]" class="image-error">
                  <span>📷</span>
                  <span>图片加载失败</span>
                  <button @click="retryLoadImage(product)" class="retry-btn">重试</button>
                </div>
              </div>
              <div class="seckill-tag">-{{ calculateDiscount(product) }}%</div>
              <div class="product-overlay">
                <button class="view-detail-btn">查看详情</button>
              </div>
            </div>
            <div class="product-info">
              <h3 class="product-name">{{ product.name }}</h3>
              <p class="product-desc">{{ product.description }}</p>
              <div class="product-price">
                <span class="current-price">¥{{ product.price.toFixed(2) }}</span>
                <span class="original-price">¥{{ product.originalPrice.toFixed(2) }}</span>
              </div>
              <div class="product-stock">
                <div class="stock-progress">
                  <div class="stock-bar" :style="{ width: getStockPercentage(product) + '%' }"></div>
                </div>
                <span class="stock-text" :class="getStockClass(product)">
                  {{ getStockText(product) }}
                </span>
              </div>
              <div class="product-actions">
                <button
                  class="seckill-btn"
                  @click.stop="handleSeckill(product.id)"
                  :disabled="product.stock === 0 || seckilling[product.id]"
                  :class="{ 'seckilling': seckilling[product.id] }"
                >
                  <span v-if="!seckilling[product.id] && product.stock > 0">立即秒杀</span>
                  <span v-if="seckilling[product.id]">
                    <span class="btn-spinner"></span>秒杀中...
                  </span>
                  <span v-if="product.stock === 0">已售罄</span>
                </button>
                <button 
                  class="cart-btn" 
                  @click="addToCart(product)"
                  :disabled="product.stock === 0"
                  :class="{ 'added': isInCart(product.id) }"
                >
                  <span v-if="!isInCart(product.id)">🛒</span>
                  <span v-else>✓</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div class="empty-state" v-if="!loading && !error && products.length === 0">
          <i class="icon-empty">🛍️</i>
          <p>当前暂无秒杀商品</p>
          <button @click="fetchProducts" class="retry-btn">刷新</button>
        </div>
      </main>
    </div>

    <!-- 底部信息 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-section">
            <h3>关于我们</h3>
            <p>秒杀商城是一个专注于限时抢购的电商平台，为用户提供优质的商品和服务。</p>
            <div class="footer-features">
              <span class="feature-tag">正品保障</span>
              <span class="feature-tag">极速发货</span>
              <span class="feature-tag">售后无忧</span>
            </div>
          </div>
          <div class="footer-section">
            <h3>联系我们</h3>
            <p>客服热线: 400-123-4567</p>
            <p>邮箱: service@seckill.com</p>
            <p>工作时间: 9:00-21:00</p>
          </div>
          <div class="footer-section">
            <h3>关注我们</h3>
            <div class="social-links">
              <a href="#" class="social-link">微信</a>
              <a href="#" class="social-link">微博</a>
              <a href="#" class="social-link">抖音</a>
              <a href="#" class="social-link">小红书</a>
            </div>
          </div>
        </div>
        <div class="footer-bottom">
          <p>&copy; 2026 秒杀商城 版权所有 | 京ICP备12345678号</p>
        </div>
      </div>
    </footer>

    <!-- 商品详情弹窗 -->
    <div class="modal product-detail-modal" v-if="showDetailModal" @click.self="closeDetailModal">
      <div class="modal-content detail-content">
        <button class="close-btn" @click="closeDetailModal">&times;</button>
        <div class="detail-layout" v-if="selectedProduct">
          <div class="detail-image-section">
            <img :src="selectedProduct.image" :alt="selectedProduct.name" class="detail-image">
            <div class="image-thumbnails">
              <div class="thumb active">
                <img :src="selectedProduct.image" :alt="selectedProduct.name">
              </div>
              <div class="thumb" v-for="i in 3" :key="i">
                <img :src="getProductImage(i)" :alt="selectedProduct.name">
              </div>
            </div>
          </div>
          <div class="detail-info-section">
            <h2 class="detail-title">{{ selectedProduct.name }}</h2>
            <p class="detail-desc">{{ selectedProduct.description }}</p>
            <div class="detail-price-section">
              <div class="detail-price">
                <span class="price-label">秒杀价</span>
                <span class="price-value">¥{{ selectedProduct.price.toFixed(2) }}</span>
              </div>
              <div class="detail-original-price">
                <span class="original-label">原价</span>
                <span class="original-value">¥{{ selectedProduct.originalPrice.toFixed(2) }}</span>
              </div>
              <div class="detail-discount">
                <span class="discount-tag">-{{ calculateDiscount(selectedProduct) }}%</span>
                <span class="savings">省 ¥{{ (selectedProduct.originalPrice - selectedProduct.price).toFixed(2) }}</span>
              </div>
            </div>
            <div class="detail-stock">
              <span class="stock-label">库存状态:</span>
              <span class="stock-value" :class="getStockClass(selectedProduct)">
                {{ getStockText(selectedProduct) }}
              </span>
              <div class="stock-progress detail-progress">
                <div class="stock-bar" :style="{ width: getStockPercentage(selectedProduct) + '%' }"></div>
              </div>
            </div>
            <div class="detail-specs" v-if="selectedProduct.specs">
              <h4>商品规格</h4>
              <div class="spec-list">
                <div class="spec-item" v-for="(value, key) in selectedProduct.specs" :key="key">
                  <span class="spec-key">{{ key }}:</span>
                  <span class="spec-value">{{ value }}</span>
                </div>
              </div>
            </div>
            <div class="detail-actions">
              <button 
                class="detail-seckill-btn"
                @click="handleSeckillFromDetail"
                :disabled="selectedProduct.stock === 0 || seckilling[selectedProduct.id]"
              >
                <span v-if="!seckilling[selectedProduct.id] && selectedProduct.stock > 0">立即秒杀</span>
                <span v-if="seckilling[selectedProduct.id]">
                  <span class="btn-spinner"></span>秒杀中...
                </span>
                <span v-if="selectedProduct.stock === 0">已售罄</span>
              </button>
              <button 
                class="detail-cart-btn"
                @click="addToCart(selectedProduct)"
                :disabled="selectedProduct.stock === 0"
                :class="{ 'added': isInCart(selectedProduct.id) }"
              >
                <span v-if="!isInCart(selectedProduct.id)">加入购物车</span>
                <span v-else>已在购物车</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 购物车侧边栏 -->
    <div class="cart-sidebar" :class="{ 'open': showCartSidebar }">
      <div class="cart-header">
        <h3>🛒 购物车</h3>
        <button class="close-cart-btn" @click="showCartSidebar = false">&times;</button>
      </div>
      <div class="cart-content">
        <div class="cart-empty" v-if="cart.length === 0">
          <span class="empty-icon">🛒</span>
          <p>购物车是空的</p>
          <button class="continue-shopping-btn" @click="showCartSidebar = false">继续购物</button>
        </div>
        <div class="cart-items" v-else>
          <div class="cart-item" v-for="item in cart" :key="item.id">
            <img :src="item.image" :alt="item.name" class="cart-item-image">
            <div class="cart-item-info">
              <h4 class="cart-item-name">{{ item.name }}</h4>
              <p class="cart-item-price">¥{{ item.price.toFixed(2) }}</p>
            </div>
            <div class="cart-item-actions">
              <button class="remove-btn" @click="removeFromCart(item.id)">&times;</button>
            </div>
          </div>
        </div>
      </div>
      <div class="cart-footer" v-if="cart.length > 0">
        <div class="cart-total">
          <span class="total-label">合计:</span>
          <span class="total-value">¥{{ cartTotal.toFixed(2) }}</span>
        </div>
        <button class="checkout-btn" @click="checkout">去结算</button>
      </div>
    </div>

    <!-- 购物车按钮 -->
    <button class="floating-cart-btn" @click="showCartSidebar = true" :class="{ 'has-items': cart.length > 0 }">
      <span class="cart-icon">🛒</span>
      <span class="cart-count" v-if="cart.length > 0">{{ cart.length }}</span>
    </button>

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
          <button class="view-order-btn" @click="viewOrder">查看订单</button>
          <button class="continue-btn" @click="showSuccessModal = false">继续购物</button>
        </div>
      </div>
    </div>

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
          <button class="cancel-btn" @click="showCaptchaModal = false">取消</button>
          <button class="confirm-btn" @click="confirmCaptchaAndSeckill">确认秒杀</button>
        </div>
      </div>
    </div>

    <!-- 结算弹窗 -->
    <div class="modal checkout-modal" v-if="showCheckoutModal" @click.self="closeCheckoutModal">
      <div class="modal-content checkout-content">
        <div class="checkout-header">
          <h3>🛒 确认订单</h3>
          <button class="close-btn" @click="closeCheckoutModal">&times;</button>
        </div>
        
        <!-- 步骤指示器 -->
        <div class="checkout-steps">
          <div class="step" :class="{ active: checkoutStep >= 1, current: checkoutStep === 1 }">
            <span class="step-number">1</span>
            <span class="step-label">确认订单</span>
          </div>
          <div class="step-line" :class="{ active: checkoutStep >= 2 }"></div>
          <div class="step" :class="{ active: checkoutStep >= 2, current: checkoutStep === 2 }">
            <span class="step-number">2</span>
            <span class="step-label">选择地址</span>
          </div>
          <div class="step-line" :class="{ active: checkoutStep >= 3 }"></div>
          <div class="step" :class="{ active: checkoutStep >= 3, current: checkoutStep === 3 }">
            <span class="step-number">3</span>
            <span class="step-label">支付方式</span>
          </div>
        </div>

        <!-- 步骤1: 确认订单 -->
        <div class="checkout-body" v-if="checkoutStep === 1">
          <div class="order-items">
            <h4>商品清单</h4>
            <div class="order-item" v-for="item in cart" :key="item.id">
              <img :src="item.image" :alt="item.name" class="order-item-img">
              <div class="order-item-info">
                <h5>{{ item.name }}</h5>
                <p class="order-item-price">¥{{ item.price.toFixed(2) }}</p>
              </div>
            </div>
          </div>
          <div class="order-summary">
            <div class="summary-row">
              <span>商品总额</span>
              <span>¥{{ cartTotal.toFixed(2) }}</span>
            </div>
            <div class="summary-row">
              <span>运费</span>
              <span class="free">免运费</span>
            </div>
            <div class="summary-row total">
              <span>应付总额</span>
              <span class="total-price">¥{{ cartTotal.toFixed(2) }}</span>
            </div>
          </div>
        </div>

        <!-- 步骤2: 选择地址 -->
        <div class="checkout-body" v-if="checkoutStep === 2">
          <div class="address-list">
            <h4>选择收货地址</h4>
            <div 
              v-for="addr in addresses" 
              :key="addr.id"
              class="address-card"
              :class="{ selected: selectedAddress?.id === addr.id }"
              @click="selectAddress(addr)"
            >
              <div class="address-radio">
                <span class="radio-inner" v-if="selectedAddress?.id === addr.id">✓</span>
              </div>
              <div class="address-info">
                <div class="address-header">
                  <span class="address-name">{{ addr.name }}</span>
                  <span class="address-phone">{{ addr.phone }}</span>
                  <span class="default-tag" v-if="addr.isDefault">默认</span>
                </div>
                <p class="address-detail">{{ addr.province }}{{ addr.city }}{{ addr.district }}{{ addr.detail }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 步骤3: 支付方式 -->
        <div class="checkout-body" v-if="checkoutStep === 3">
          <div class="payment-methods">
            <h4>选择支付方式</h4>
            <div 
              v-for="method in paymentMethods" 
              :key="method.id"
              class="payment-card"
              :class="{ selected: selectedPayment === method.id }"
              @click="selectPayment(method.id)"
            >
              <div class="payment-icon" :style="{ backgroundColor: method.color + '20', color: method.color }">
                {{ method.icon }}
              </div>
              <span class="payment-name">{{ method.name }}</span>
              <div class="payment-radio">
                <span class="radio-inner" v-if="selectedPayment === method.id">✓</span>
              </div>
            </div>
          </div>
          <div class="payment-summary">
            <div class="final-total">
              <span>实付金额:</span>
              <span class="final-price">¥{{ cartTotal.toFixed(2) }}</span>
            </div>
          </div>
        </div>

        <!-- 步骤4: 支付成功 -->
        <div class="checkout-body success-body" v-if="checkoutStep === 4">
          <div class="success-animation-large">
            <div class="success-circle-large">
              <span class="success-icon-large">✓</span>
            </div>
          </div>
          <h3 class="success-title">支付成功!</h3>
          <p class="success-desc">您的订单已提交，我们将尽快为您发货</p>
          <div class="success-actions-large">
            <button class="view-orders-btn" @click="viewMyOrders">查看订单</button>
            <button class="continue-shop-btn" @click="continueShopping">继续购物</button>
          </div>
        </div>

        <!-- 底部操作按钮 -->
        <div class="checkout-footer" v-if="checkoutStep < 4">
          <button class="prev-btn" v-if="checkoutStep > 1" @click="prevStep">上一步</button>
          <button class="next-btn" v-if="checkoutStep < 3" @click="nextStep">下一步</button>
          <button class="submit-btn" v-if="checkoutStep === 3" @click="submitOrder" :disabled="isProcessing">
            <span v-if="!isProcessing">确认支付</span>
            <span v-else class="btn-spinner-white"></span>
          </button>
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
import { ref, onMounted, onUnmounted, computed, reactive } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import Navbar from './Navbar.vue';

const router = useRouter();

// 状态管理
const products = ref([]);
const loading = ref(false);
const error = ref('');
const seckilling = reactive({});
const showSuccessModal = ref(false);
const successProductName = ref('');
const imageLoading = ref({});
const imageError = ref({});
const countdown = ref(null);
const userId = ref(localStorage.getItem('user_id'));
const activeCategory = ref(1);
const sortBy = ref('default');

// 验证码相关状态
const showCaptchaModal = ref(false);
const captchaImageUrl = ref('');
const captchaId = ref('');
const captchaInput = ref('');
const pendingProductForCaptcha = ref(null);
const isSubmitting = ref(false);

// 商品详情弹窗
const showDetailModal = ref(false);
const selectedProduct = ref(null);

// 购物车
const showCartSidebar = ref(false);
const cart = ref(JSON.parse(localStorage.getItem('cart') || '[]'));

// 结算弹窗
const showCheckoutModal = ref(false);
const checkoutStep = ref(1);
const selectedAddress = ref(null);
const selectedPayment = ref('alipay');
const isProcessing = ref(false);

// 地址列表
const addresses = ref([
  {
    id: 1,
    name: '张三',
    phone: '138****8888',
    province: '北京市',
    city: '北京市',
    district: '朝阳区',
    detail: '建国路88号SOHO现代城A座1201',
    isDefault: true
  },
  {
    id: 2,
    name: '张三',
    phone: '138****8888',
    province: '上海市',
    city: '上海市',
    district: '浦东新区',
    detail: '陆家嘴环路1000号恒生银行大厦',
    isDefault: false
  }
]);

// 支付方式
const paymentMethods = ref([
  { id: 'alipay', name: '支付宝', icon: '💳', color: '#1677FF' },
  { id: 'wechat', name: '微信支付', icon: '💬', color: '#07C160' },
  { id: 'card', name: '银行卡', icon: '💰', color: '#FF6B6B' }
]);

// 提示消息
const showToast = ref(false);
const toastMessage = ref('');
const toastType = ref('success');
const toastIcon = ref('✓');

// 商品分类
const categories = ref([
  { id: 1, name: '全部', icon: '📦' },
  { id: 2, name: '手机数码', icon: '📱' },
  { id: 3, name: '电脑办公', icon: '💻' },
  { id: 4, name: '家用电器', icon: '🏠' },
  { id: 5, name: '服装鞋包', icon: '👕' },
  { id: 6, name: '美妆护肤', icon: '💄' }
]);

// 扩展的商品数据
const extendedProducts = [
  {
    id: 1,
    name: 'iPhone 15 Pro',
    price: 7999.99,
    originalPrice: 9999.99,
    stock: 100,
    description: '钛金属设计，A17 Pro芯片，4800万像素主摄',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPhone%2015%20Pro%20titanium%20metal%20smartphone%20with%20triple%20camera%20system%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 2,
    specs: {
      '屏幕': '6.1英寸 OLED',
      '处理器': 'A17 Pro',
      '存储': '256GB',
      '颜色': '原色钛金属'
    }
  },
  {
    id: 2,
    name: 'MacBook Pro 14',
    price: 12999.99,
    originalPrice: 14999.99,
    stock: 50,
    description: 'M3芯片，Liquid视网膜XDR显示屏，18小时续航',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=MacBook%20Pro%2014%20inch%20laptop%20with%20silver%20aluminum%20body%20open%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 3,
    specs: {
      '屏幕': '14.2英寸 XDR',
      '处理器': 'M3',
      '内存': '16GB',
      '存储': '512GB SSD'
    }
  },
  {
    id: 3,
    name: 'AirPods Pro 2',
    price: 1899.99,
    originalPrice: 2299.99,
    stock: 200,
    description: '主动降噪，空间音频，自适应通透模式',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=AirPods%20Pro%202%20wireless%20earbuds%20with%20charging%20case%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 2,
    specs: {
      '芯片': 'H2',
      '降噪': '主动降噪',
      '续航': '30小时',
      '防水': 'IPX4'
    }
  },
  {
    id: 4,
    name: 'iPad Pro 12.9',
    price: 8999.99,
    originalPrice: 10999.99,
    stock: 50,
    description: 'M2芯片，Liquid视网膜XDR显示屏，支持Apple Pencil',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=iPad%20Pro%2012.9%20inch%20tablet%20with%20Apple%20Pencil%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 3,
    specs: {
      '屏幕': '12.9英寸 XDR',
      '处理器': 'M2',
      '存储': '256GB',
      '网络': 'WiFi + 5G'
    }
  },
  {
    id: 5,
    name: 'Sony WH-1000XM5',
    price: 2499.99,
    originalPrice: 2999.99,
    stock: 80,
    description: '业界领先的降噪，30小时续航，舒适佩戴',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Sony%20WH-1000XM5%20black%20over-ear%20headphones%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 2,
    specs: {
      '降噪': '主动降噪',
      '续航': '30小时',
      '连接': '蓝牙5.2',
      '重量': '250g'
    }
  },
  {
    id: 6,
    name: 'Nintendo Switch OLED',
    price: 2199.99,
    originalPrice: 2599.99,
    stock: 60,
    description: '7英寸OLED屏幕，64GB存储，增强音频',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Nintendo%20Switch%20OLED%20white%20handheld%20gaming%20console%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 2,
    specs: {
      '屏幕': '7英寸 OLED',
      '存储': '64GB',
      '续航': '9小时',
      '颜色': '白色'
    }
  },
  {
    id: 7,
    name: 'Dyson V15 Detect',
    price: 4599.99,
    originalPrice: 5499.99,
    stock: 40,
    description: '激光探测，智能吸力调节，60分钟续航',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Dyson%20V15%20Detect%20cordless%20vacuum%20cleaner%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 4,
    specs: {
      '吸力': '230AW',
      '续航': '60分钟',
      '重量': '3kg',
      '尘杯': '0.77L'
    }
  },
  {
    id: 8,
    name: 'Nike Air Max 270',
    price: 899.99,
    originalPrice: 1299.99,
    stock: 150,
    description: '大气垫缓震，透气网面，时尚设计',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Nike%20Air%20Max%20270%20sneakers%20with%20visible%20air%20cushion%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 5,
    specs: {
      '鞋面': '网眼布',
      '中底': 'Air Max',
      '闭合': '系带',
      '适用': '跑步/休闲'
    }
  },
  {
    id: 9,
    name: 'SK-II 神仙水',
    price: 1299.99,
    originalPrice: 1699.99,
    stock: 70,
    description: 'Pitera精华，改善肤质，提亮肤色',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=SK-II%20Facial%20Treatment%20Essence%20bottle%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 6,
    specs: {
      '容量': '230ml',
      '成分': 'Pitera',
      '功效': '保湿/提亮',
      '适用': '所有肤质'
    }
  },
  {
    id: 10,
    name: 'Samsung Galaxy S24',
    price: 5999.99,
    originalPrice: 6999.99,
    stock: 90,
    description: 'AI智能手机，超视觉影像，钛金属边框',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Samsung%20Galaxy%20S24%20smartphone%20with%20titanium%20frame%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 2,
    specs: {
      '屏幕': '6.2英寸 AMOLED',
      '处理器': '骁龙8 Gen3',
      '存储': '256GB',
      '相机': '5000万像素'
    }
  },
  {
    id: 11,
    name: 'Lego 兰博基尼',
    price: 2499.99,
    originalPrice: 3299.99,
    stock: 30,
    description: '1:8比例，3696颗粒，V12发动机',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Lego%20Lamborghini%20model%20car%20kit%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 4,
    specs: {
      '颗粒': '3696',
      '比例': '1:8',
      '尺寸': '60x25x13cm',
      '年龄': '18+'
    }
  },
  {
    id: 12,
    name: 'Adidas Yeezy Boost',
    price: 1899.99,
    originalPrice: 2499.99,
    stock: 45,
    description: 'Boost缓震，Primeknit鞋面，潮流设计',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Adidas%20Yeezy%20Boost%20sneakers%20with%20Primeknit%20upper%20on%20white%20background%20professional%20product%20photography&image_size=square',
    category: 5,
    specs: {
      '鞋面': 'Primeknit',
      '中底': 'Boost',
      '闭合': '系带',
      '风格': '潮流'
    }
  }
];

// 计算总库存
const totalStock = computed(() => {
  return products.value.reduce((sum, p) => sum + p.stock, 0);
});

// 购物车总价
const cartTotal = computed(() => {
  return cart.value.reduce((sum, item) => sum + item.price, 0);
});

// 格式化数字为两位数
const formatNumber = (num) => {
  return num.toString().padStart(2, '0');
};

// 计算倒计时
const calculateCountdown = () => {
  const now = new Date();
  const endTime = new Date();
  endTime.setHours(23, 59, 59, 999);

  const diff = endTime - now;
  if (diff <= 0) {
    return { hours: '00', minutes: '00', seconds: '00' };
  }

  const hours = Math.floor(diff / (1000 * 60 * 60));
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((diff % (1000 * 60)) / 1000);

  return {
    hours: formatNumber(hours),
    minutes: formatNumber(minutes),
    seconds: formatNumber(seconds)
  };
};

// 计算折扣
const calculateDiscount = (product) => {
  return Math.round((1 - product.price / product.originalPrice) * 100);
};

// 获取库存百分比
const getStockPercentage = (product) => {
  const maxStock = 200;
  return Math.min(100, (product.stock / maxStock) * 100);
};

// 获取库存文本
const getStockText = (product) => {
  if (product.stock === 0) return '已抢完';
  if (product.stock <= 10) return `仅剩${product.stock}件`;
  if (product.stock <= 50) return '库存紧张';
  return '库存充足';
};

// 获取库存样式类
const getStockClass = (product) => {
  if (product.stock === 0) return 'out-of-stock';
  if (product.stock <= 10) return 'low-stock';
  if (product.stock <= 50) return 'medium-stock';
  return 'in-stock';
};

// 获取商品列表
const fetchProducts = async () => {
  loading.value = true;
  error.value = '';

  try {
    products.value = extendedProducts;
  } catch (err) {
    console.error('获取商品失败:', err);
    error.value = '获取商品列表失败，请稍后重试';
  } finally {
    loading.value = false;
  }
};

// 排序商品
const sortedProducts = computed(() => {
  let result = [...products.value];
  
  if (activeCategory.value !== 1) {
    result = result.filter(p => p.category === activeCategory.value);
  }
  
  switch (sortBy.value) {
    case 'price-asc':
      result.sort((a, b) => a.price - b.price);
      break;
    case 'price-desc':
      result.sort((a, b) => b.price - a.price);
      break;
    case 'stock':
      result.sort((a, b) => b.stock - a.stock);
      break;
    case 'discount':
      result.sort((a, b) => calculateDiscount(b) - calculateDiscount(a));
      break;
    default:
      break;
  }
  
  return result;
});

// 显示商品详情
const showProductDetail = (product) => {
  selectedProduct.value = product;
  showDetailModal.value = true;
  document.body.style.overflow = 'hidden';
};

// 关闭详情弹窗
const closeDetailModal = () => {
  showDetailModal.value = false;
  selectedProduct.value = null;
  document.body.style.overflow = '';
};

// 从详情页秒杀
const handleSeckillFromDetail = () => {
  if (selectedProduct.value) {
    // 先保存商品信息
    const product = selectedProduct.value;
    if (product && product.id) {
      closeDetailModal();
      handleSeckill(product.id);
    } else {
      showToastMessage('商品信息错误', 'error');
    }
  }
};

// 处理秒杀
const handleSeckill = async (productId) => {
  try {
    const product = products.value.find(p => p.id === productId);
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

// 添加到购物车
const addToCart = (product) => {
  if (product.stock === 0) {
    showToastMessage('商品已售罄', 'error');
    return;
  }
  
  const existingItem = cart.value.find(item => item.id === product.id);
  if (existingItem) {
    showToastMessage('商品已在购物车中', 'warning');
    return;
  }
  
  cart.value.push({
    id: product.id,
    name: product.name,
    price: product.price,
    image: product.image
  });
  
  localStorage.setItem('cart', JSON.stringify(cart.value));
  showToastMessage('已添加到购物车', 'success');
};

// 从购物车移除
const removeFromCart = (productId) => {
  cart.value = cart.value.filter(item => item.id !== productId);
  localStorage.setItem('cart', JSON.stringify(cart.value));
  showToastMessage('已从购物车移除', 'info');
};

// 检查是否在购物车中
const isInCart = (productId) => {
  return cart.value.some(item => item.id === productId);
};

// 结算
const checkout = () => {
  if (cart.value.length === 0) {
    showToastMessage('购物车是空的', 'warning');
    return;
  }
  selectedAddress.value = addresses.value.find(a => a.isDefault) || addresses.value[0];
  checkoutStep.value = 1;
  showCheckoutModal.value = true;
  showCartSidebar.value = false;
};

// 关闭结算弹窗
const closeCheckoutModal = () => {
  showCheckoutModal.value = false;
  checkoutStep.value = 1;
  isProcessing.value = false;
};

// 选择地址
const selectAddress = (address) => {
  selectedAddress.value = address;
};

// 选择支付方式
const selectPayment = (methodId) => {
  selectedPayment.value = methodId;
};

// 下一步
const nextStep = () => {
  if (checkoutStep.value < 3) {
    checkoutStep.value++;
  }
};

// 上一步
const prevStep = () => {
  if (checkoutStep.value > 1) {
    checkoutStep.value--;
  }
};

// 提交订单
const submitOrder = async () => {
  isProcessing.value = true;
  
  await new Promise(resolve => setTimeout(resolve, 2000));
  
  cart.value = [];
  localStorage.setItem('cart', JSON.stringify(cart.value));
  
  checkoutStep.value = 4;
  isProcessing.value = false;
  showToastMessage('订单提交成功！', 'success');
};

// 查看我的订单
const viewMyOrders = () => {
  closeCheckoutModal();
  router.push('/order');
};

// 继续购物
const continueShopping = () => {
  closeCheckoutModal();
};

// 查看订单
const viewOrder = () => {
  showSuccessModal.value = false;
  router.push('/order');
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

// 获取商品图片
const getProductImage = (productId) => {
  const product = extendedProducts.find(p => p.id === productId);
  if (product) {
    return product.image;
  }
  
  // 为详情页生成不同角度的实物图
  if (selectedProduct.value) {
    const angles = ['side view', 'back view', 'close up detail', 'in use'];
    const index = productId % angles.length;
    const angle = angles[index];
    return `https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=${encodeURIComponent(selectedProduct.value.name + ' ' + angle + ' on white background professional product photography')}&image_size=square`;
  }
  
  return `https://placehold.co/400x400/667eea/ffffff?text=Product+${productId}`;
};

// 图片加载处理
const handleImageError = (product) => {
  imageError.value[product.id] = true;
  imageLoading.value[product.id] = false;
};

const handleImageLoad = (product) => {
  imageError.value[product.id] = false;
  imageLoading.value[product.id] = false;
};

const retryLoadImage = (product) => {
  imageError.value[product.id] = false;
  imageLoading.value[product.id] = true;
  // 重新设置src以触发重新加载
  const originalSrc = product.image;
  product.image = '';
  setTimeout(() => {
    product.image = originalSrc;
  }, 100);
};

// 初始化
onMounted(async () => {
  if (!localStorage.getItem('token')) {
    router.push('/login');
    return;
  }
  
  await fetchProducts();

  countdown.value = calculateCountdown();
  const timer = setInterval(() => {
    countdown.value = calculateCountdown();
  }, 1000);

  const cleanup = () => clearInterval(timer);
  window.addEventListener('beforeunload', cleanup);

  // 在组件卸载时清理
  onUnmounted(() => {
    cleanup();
    window.removeEventListener('beforeunload', cleanup);
  });
});
</script>

<style scoped>
/* 全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.seckill-container {
  font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
  color: #333;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  min-height: 100vh;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

/* 秒杀活动横幅 */
.seckill-banner {
  position: relative;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  color: white;
  padding: 80px 0 100px;
  margin-bottom: 0;
  overflow: hidden;
}

.banner-bg-animation {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 80%, rgba(255,255,255,0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255,255,255,0.1) 0%, transparent 50%);
  animation: bgPulse 4s ease-in-out infinite;
}

@keyframes bgPulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

.banner-content {
  position: relative;
  z-index: 1;
  text-align: center;
}

.banner-title {
  font-size: 56px;
  font-weight: 800;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.title-highlight {
  background: linear-gradient(90deg, #fff, #ffd700, #fff);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: shine 3s linear infinite;
}

@keyframes shine {
  to { background-position: 200% center; }
}

.title-pulse {
  animation: pulse 1.5s ease-in-out infinite;
  display: inline-block;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.2); }
}

.banner-subtitle {
  font-size: 22px;
  margin-bottom: 40px;
  opacity: 0.95;
  font-weight: 300;
}

.countdown-section {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-bottom: 50px;
}

.countdown-label {
  font-size: 18px;
  font-weight: 500;
  opacity: 0.9;
}

.countdown {
  display: flex;
  align-items: center;
  gap: 12px;
}

.time-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 12px 20px;
  min-width: 70px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.time-value {
  font-size: 36px;
  font-weight: 700;
  line-height: 1;
}

.time-label {
  font-size: 12px;
  opacity: 0.8;
  margin-top: 4px;
}

.colon {
  font-size: 32px;
  font-weight: 700;
  opacity: 0.8;
}

.banner-stats {
  display: flex;
  justify-content: center;
  gap: 60px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.stat-number {
  font-size: 36px;
  font-weight: 700;
  background: linear-gradient(90deg, #ffd700, #ffed4a);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

/* 商品分类 */
.category-section {
  background: white;
  padding: 24px 0;
  margin-bottom: 40px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}

.categories {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  cursor: pointer;
  border-radius: 30px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
  background: #f8f9fa;
  border: 2px solid transparent;
}

.category-item:hover {
  background: #e9ecef;
  transform: translateY(-2px);
}

.category-item.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: #667eea;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.category-icon {
  font-size: 20px;
}

/* 商品列表样式 */
.product-list {
  margin-bottom: 60px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
}

.section-title-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-title {
  font-size: 32px;
  font-weight: 700;
  color: #333;
  margin: 0;
}

.section-subtitle {
  font-size: 16px;
  color: #666;
}

.sort-options {
  display: flex;
  align-items: center;
  gap: 12px;
}

.sort-label {
  font-weight: 500;
  color: #666;
}

.sort-select {
  padding: 10px 16px;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  transition: all 0.3s;
}

.sort-select:hover {
  border-color: #667eea;
}

.sort-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* 加载状态 */
.loading {
  text-align: center;
  padding: 100px 0;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.spinner {
  width: 60px;
  height: 60px;
  margin: 0 auto 24px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 错误状态 */
.error-message {
  text-align: center;
  padding: 100px 0;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  color: #e63946;
}

.icon-error {
  font-size: 60px;
  margin-bottom: 20px;
  display: inline-block;
}

.retry-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 32px;
  border-radius: 8px;
  cursor: pointer;
  margin-top: 20px;
  transition: all 0.3s;
  font-size: 16px;
  font-weight: 500;
}

.retry-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

/* 商品网格 */
.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 30px;
}

.product-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  position: relative;
}

.product-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

.product-card.sold-out {
  opacity: 0.7;
}

.product-card.sold-out .product-image img {
  filter: grayscale(100%);
}

.product-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a5a 100%);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  z-index: 10;
  animation: badgePulse 2s ease-in-out infinite;
}

@keyframes badgePulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.sold-out-badge {
  background: linear-gradient(135deg, #868e96 0%, #495057 100%);
  animation: none;
}

.product-image {
  position: relative;
  height: 240px;
  overflow: hidden;
  cursor: pointer;
}

.image-container {
  position: relative;
  width: 100%;
  height: 100%;
  background: #f8f9fa;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.product-card:hover .product-image img {
  transform: scale(1.1);
}

.image-loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1;
}

.image-error {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1;
  text-align: center;
  padding: 20px;
}

.image-error span:first-child {
  font-size: 32px;
  margin-bottom: 10px;
}

.image-error span:last-child {
  margin-bottom: 10px;
  color: #666;
}

.image-error .retry-btn {
  background: #1890ff;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background 0.3s ease;
  margin-top: 0;
  padding: 8px 16px;
  font-size: 14px;
}

.image-error .retry-btn:hover {
  background: #40a9ff;
  transform: none;
  box-shadow: none;
}

.product-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.product-image:hover .product-overlay {
  opacity: 1;
}

.view-detail-btn {
  background: white;
  color: #333;
  border: none;
  padding: 12px 24px;
  border-radius: 25px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.view-detail-btn:hover {
  background: #667eea;
  color: white;
}

.seckill-tag {
  position: absolute;
  top: 12px;
  left: 12px;
  background: linear-gradient(135deg, #ffd700 0%, #ffed4a 100%);
  color: #333;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 700;
  z-index: 10;
}

.product-info {
  padding: 24px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.product-name {
  font-size: 18px;
  margin-bottom: 8px;
  color: #333;
  font-weight: 600;
  line-height: 1.4;
}

.product-desc {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-price {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.current-price {
  font-size: 28px;
  color: #e63946;
  font-weight: 700;
}

.original-price {
  font-size: 14px;
  color: #999;
  text-decoration: line-through;
}

.product-stock {
  margin-bottom: 20px;
}

.stock-progress {
  height: 8px;
  background: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.stock-bar {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transition: width 0.5s ease;
  border-radius: 4px;
}

.stock-text {
  font-size: 13px;
  font-weight: 500;
}

.stock-text.out-of-stock {
  color: #868e96;
}

.stock-text.low-stock {
  color: #e63946;
}

.stock-text.medium-stock {
  color: #f59e0b;
}

.stock-text.in-stock {
  color: #10b981;
}

.product-actions {
  display: flex;
  gap: 12px;
  margin-top: auto;
}

.seckill-btn {
  flex: 1;
  background: linear-gradient(135deg, #e63946 0%, #c1121f 100%);
  color: white;
  border: none;
  padding: 14px 0;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.seckill-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(230, 57, 70, 0.4);
}

.seckill-btn:disabled {
  background: linear-gradient(135deg, #868e96 0%, #495057 100%);
  cursor: not-allowed;
}

.seckill-btn.seckilling {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.btn-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  display: inline-block;
}

.cart-btn {
  width: 50px;
  height: 50px;
  border: 2px solid #e9ecef;
  background: white;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cart-btn:hover:not(:disabled) {
  border-color: #667eea;
  background: #f8f9ff;
}

.cart-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.cart-btn.added {
  background: #10b981;
  border-color: #10b981;
  color: white;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 120px 0;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  color: #999;
}

.icon-empty {
  font-size: 80px;
  margin-bottom: 20px;
  display: inline-block;
}

/* 底部信息 */
.footer {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: white;
  padding: 60px 0 30px;
  margin-top: 80px;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 40px;
  margin-bottom: 40px;
}

.footer-section h3 {
  font-size: 20px;
  margin-bottom: 20px;
  color: #fff;
  font-weight: 600;
}

.footer-section p {
  line-height: 1.8;
  color: #adb5bd;
  margin-bottom: 16px;
}

.footer-features {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 16px;
}

.feature-tag {
  background: rgba(102, 126, 234, 0.2);
  color: #a5b4fc;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
}

.social-links {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.social-link {
  color: #adb5bd;
  text-decoration: none;
  padding: 10px 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  transition: all 0.3s;
}

.social-link:hover {
  background: rgba(102, 126, 234, 0.3);
  color: white;
  border-color: #667eea;
}

.footer-bottom {
  text-align: center;
  padding-top: 30px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  color: #6c757d;
}

/* 弹窗通用样式 */
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
  padding: 20px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: white;
  border-radius: 20px;
  max-width: 900px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
  animation: slideUp 0.3s ease;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
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
  top: 20px;
  right: 20px;
  width: 40px;
  height: 40px;
  border: none;
  background: #f8f9fa;
  border-radius: 50%;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  z-index: 10;
}

.close-btn:hover {
  background: #e9ecef;
  transform: rotate(90deg);
}

/* 商品详情弹窗 */
.product-detail-modal .modal-content {
  max-width: 1000px;
}

.detail-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  padding: 40px;
}

.detail-image-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-image {
  width: 100%;
  border-radius: 16px;
  object-fit: cover;
  aspect-ratio: 1;
}

.image-thumbnails {
  display: flex;
  gap: 12px;
}

.thumb {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.thumb.active {
  border-color: #667eea;
}

.thumb:hover {
  border-color: #667eea;
}

.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.detail-info-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-title {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  line-height: 1.3;
}

.detail-desc {
  font-size: 16px;
  color: #666;
  line-height: 1.6;
}

.detail-price-section {
  background: linear-gradient(135deg, #fff5f5 0%, #fff 100%);
  padding: 24px;
  border-radius: 16px;
  border: 1px solid #ffe0e0;
}

.detail-price {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 12px;
}

.price-label {
  font-size: 14px;
  color: #e63946;
  font-weight: 500;
}

.price-value {
  font-size: 36px;
  font-weight: 700;
  color: #e63946;
}

.detail-original-price {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.original-label {
  font-size: 14px;
  color: #999;
}

.original-value {
  font-size: 16px;
  color: #999;
  text-decoration: line-through;
}

.detail-discount {
  display: flex;
  align-items: center;
  gap: 12px;
}

.discount-tag {
  background: linear-gradient(135deg, #e63946 0%, #c1121f 100%);
  color: white;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
}

.savings {
  color: #e63946;
  font-weight: 500;
}

.detail-stock {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stock-label {
  font-weight: 500;
  color: #666;
}

.stock-value {
  font-weight: 600;
}

.detail-progress {
  height: 10px;
}

.detail-specs h4 {
  font-size: 18px;
  margin-bottom: 16px;
  color: #333;
}

.spec-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.spec-item {
  display: flex;
  gap: 8px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.spec-key {
  color: #666;
  font-weight: 500;
}

.spec-value {
  color: #333;
  font-weight: 600;
}

.detail-actions {
  display: flex;
  gap: 16px;
  margin-top: auto;
}

.detail-seckill-btn {
  flex: 1;
  background: linear-gradient(135deg, #e63946 0%, #c1121f 100%);
  color: white;
  border: none;
  padding: 16px 32px;
  border-radius: 12px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.detail-seckill-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(230, 57, 70, 0.4);
}

.detail-seckill-btn:disabled {
  background: linear-gradient(135deg, #868e96 0%, #495057 100%);
  cursor: not-allowed;
}

.detail-cart-btn {
  padding: 16px 32px;
  border: 2px solid #667eea;
  background: white;
  color: #667eea;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.detail-cart-btn:hover:not(:disabled) {
  background: #667eea;
  color: white;
}

.detail-cart-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.detail-cart-btn.added {
  background: #10b981;
  border-color: #10b981;
  color: white;
}

/* 购物车侧边栏 */
.cart-sidebar {
  position: fixed;
  top: 0;
  right: -420px;
  width: 400px;
  height: 100vh;
  background: white;
  box-shadow: -10px 0 40px rgba(0, 0, 0, 0.15);
  z-index: 1001;
  display: flex;
  flex-direction: column;
  transition: right 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.cart-sidebar.open {
  right: 0;
}

.cart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #e9ecef;
}

.cart-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.close-cart-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f8f9fa;
  border-radius: 50%;
  font-size: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.close-cart-btn:hover {
  background: #e9ecef;
}

.cart-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.cart-empty {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty-icon {
  font-size: 60px;
  margin-bottom: 16px;
  display: block;
}

.continue-shopping-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  margin-top: 20px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.continue-shopping-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.cart-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  transition: all 0.3s;
}

.cart-item:hover {
  background: #e9ecef;
}

.cart-item-image {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
}

.cart-item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
}

.cart-item-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  line-height: 1.4;
}

.cart-item-price {
  font-size: 18px;
  font-weight: 700;
  color: #e63946;
}

.cart-item-actions {
  display: flex;
  align-items: center;
}

.remove-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: #e63946;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 18px;
  transition: all 0.3s;
}

.remove-btn:hover {
  background: #c1121f;
  transform: scale(1.1);
}

.cart-footer {
  padding: 24px;
  border-top: 1px solid #e9ecef;
  background: white;
}

.cart-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.total-label {
  font-size: 16px;
  color: #666;
}

.total-value {
  font-size: 28px;
  font-weight: 700;
  color: #e63946;
}

.checkout-btn {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 16px;
  border-radius: 12px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.checkout-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

/* 悬浮购物车按钮 */
.floating-cart-btn {
  position: fixed;
  bottom: 30px;
  right: 30px;
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  transition: all 0.3s;
}

.floating-cart-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 12px 35px rgba(102, 126, 234, 0.5);
}

.floating-cart-btn.has-items {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.cart-icon {
  font-size: 28px;
}

.cart-count {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #e63946;
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 秒杀成功弹窗 */
.success-modal .modal-content {
  max-width: 450px;
  text-align: center;
  padding: 50px 40px;
}

.success-animation {
  position: relative;
  margin-bottom: 30px;
}

.success-circle {
  width: 100px;
  height: 100px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  animation: scaleIn 0.5s ease;
}

@keyframes scaleIn {
  0% { transform: scale(0); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.success-icon {
  font-size: 50px;
  color: white;
  font-weight: 700;
}

.confetti {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
  pointer-events: none;
}

.confetti span {
  position: absolute;
  width: 10px;
  height: 10px;
  background: #ffd700;
  border-radius: 50%;
  animation: confetti-fall 1s ease-out forwards;
}

.confetti-1 { transform: rotate(0deg) translateY(-80px); animation-delay: 0s; background: #ff6b6b; }
.confetti-2 { transform: rotate(30deg) translateY(-80px); animation-delay: 0.1s; background: #4ecdc4; }
.confetti-3 { transform: rotate(60deg) translateY(-80px); animation-delay: 0.2s; background: #45b7d1; }
.confetti-4 { transform: rotate(90deg) translateY(-80px); animation-delay: 0.3s; background: #f9ca24; }
.confetti-5 { transform: rotate(120deg) translateY(-80px); animation-delay: 0.4s; background: #6c5ce7; }
.confetti-6 { transform: rotate(150deg) translateY(-80px); animation-delay: 0.5s; background: #a29bfe; }
.confetti-7 { transform: rotate(180deg) translateY(-80px); animation-delay: 0.1s; background: #fd79a8; }
.confetti-8 { transform: rotate(210deg) translateY(-80px); animation-delay: 0.2s; background: #e17055; }
.confetti-9 { transform: rotate(240deg) translateY(-80px); animation-delay: 0.3s; background: #00b894; }
.confetti-10 { transform: rotate(270deg) translateY(-80px); animation-delay: 0.4s; background: #0984e3; }
.confetti-11 { transform: rotate(300deg) translateY(-80px); animation-delay: 0.5s; background: #fdcb6e; }
.confetti-12 { transform: rotate(330deg) translateY(-80px); animation-delay: 0.6s; background: #e84393; }

@keyframes confetti-fall {
  0% {
    opacity: 1;
    transform: rotate(var(--rotation, 0deg)) translateY(-80px) scale(1);
  }
  100% {
    opacity: 0;
    transform: rotate(var(--rotation, 0deg)) translateY(20px) scale(0);
  }
}

.success-content h3 {
  font-size: 28px;
  color: #333;
  margin-bottom: 12px;
}

.success-content p {
  color: #666;
  font-size: 16px;
  margin-bottom: 30px;
}

.success-product {
  color: #e63946;
  font-weight: 600;
}

.success-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.view-order-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 32px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.view-order-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.continue-btn {
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
  padding: 14px 32px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.continue-btn:hover {
  background: #f8f9ff;
}

/* 验证码弹窗 */
.captcha-modal-wrapper .modal-content {
  max-width: 420px;
  padding: 40px;
  text-align: center;
}

.captcha-content h3 {
  font-size: 24px;
  margin-bottom: 8px;
  color: #333;
}

.captcha-content p {
  color: #666;
  margin-bottom: 24px;
}

.captcha-row {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.captcha-image-wrapper {
  cursor: pointer;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid #e9ecef;
  transition: all 0.3s;
  position: relative;
}

.captcha-image-wrapper:hover {
  border-color: #667eea;
}

.captcha-image {
  width: 100%;
  height: 80px;
  object-fit: contain;
  background: #f8f9fa;
}

.captcha-refresh {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.captcha-input {
  width: 100%;
  padding: 16px;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  font-size: 18px;
  text-align: center;
  letter-spacing: 8px;
  transition: all 0.3s;
}

.captcha-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.captcha-actions {
  display: flex;
  gap: 16px;
}

.cancel-btn {
  flex: 1;
  background: #f8f9fa;
  color: #666;
  border: none;
  padding: 14px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.cancel-btn:hover {
  background: #e9ecef;
}

.confirm-btn {
  flex: 1;
  background: linear-gradient(135deg, #e63946 0%, #c1121f 100%);
  color: white;
  border: none;
  padding: 14px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(230, 57, 70, 0.4);
}

/* 结算弹窗 */
.checkout-modal .modal-content {
  max-width: 600px;
  max-height: 85vh;
}

.checkout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 30px;
  border-bottom: 1px solid #e9ecef;
}

.checkout-header h3 {
  font-size: 22px;
  font-weight: 600;
  color: #333;
}

/* 步骤指示器 */
.checkout-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 30px;
  background: #f8f9fa;
  gap: 8px;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  opacity: 0.5;
  transition: all 0.3s;
}

.step.active {
  opacity: 1;
}

.step.current .step-number {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.step-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #e9ecef;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s;
}

.step-label {
  font-size: 12px;
  font-weight: 500;
  color: #666;
}

.step-line {
  flex: 1;
  height: 2px;
  background: #e9ecef;
  max-width: 60px;
  transition: all 0.3s;
}

.step-line.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* 结算内容区 */
.checkout-body {
  padding: 30px;
  min-height: 300px;
}

/* 步骤1: 确认订单 */
.order-items h4 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.order-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  margin-bottom: 12px;
}

.order-item-img {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
}

.order-item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
}

.order-item-info h5 {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.order-item-price {
  font-size: 18px;
  font-weight: 700;
  color: #e63946;
}

.order-summary {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #e9ecef;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 15px;
  color: #666;
}

.summary-row.total {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e9ecef;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.total-price {
  font-size: 28px;
  font-weight: 700;
  color: #e63946;
}

.free {
  color: #10b981;
  font-weight: 500;
}

/* 步骤2: 选择地址 */
.address-list h4 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.address-card {
  display: flex;
  gap: 16px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
  margin-bottom: 12px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.address-card:hover {
  background: #e9ecef;
}

.address-card.selected {
  border-color: #667eea;
  background: #f8f9ff;
}

.address-radio {
  width: 24px;
  height: 24px;
  border: 2px solid #d1d5db;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.address-card.selected .address-radio {
  border-color: #667eea;
  background: #667eea;
}

.radio-inner {
  color: white;
  font-size: 12px;
  font-weight: 700;
}

.address-info {
  flex: 1;
}

.address-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.address-name {
  font-weight: 600;
  color: #333;
}

.address-phone {
  color: #666;
  font-size: 14px;
}

.default-tag {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.address-detail {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

/* 步骤3: 支付方式 */
.payment-methods h4 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.payment-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
  margin-bottom: 12px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.payment-card:hover {
  background: #e9ecef;
}

.payment-card.selected {
  border-color: #667eea;
  background: #f8f9ff;
}

.payment-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.payment-name {
  flex: 1;
  font-weight: 500;
  color: #333;
}

.payment-radio {
  width: 24px;
  height: 24px;
  border: 2px solid #d1d5db;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.payment-card.selected .payment-radio {
  border-color: #667eea;
  background: #667eea;
}

.payment-summary {
  margin-top: 30px;
  padding: 20px;
  background: linear-gradient(135deg, #fff5f5 0%, #fff 100%);
  border-radius: 12px;
  border: 1px solid #ffe0e0;
}

.final-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.final-total span:first-child {
  font-size: 16px;
  color: #666;
}

.final-price {
  font-size: 32px;
  font-weight: 700;
  color: #e63946;
}

/* 步骤4: 支付成功 */
.success-body {
  text-align: center;
  padding: 40px 30px;
}

.success-animation-large {
  margin-bottom: 30px;
}

.success-circle-large {
  width: 120px;
  height: 120px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  animation: scaleIn 0.5s ease;
  box-shadow: 0 10px 30px rgba(16, 185, 129, 0.3);
}

.success-icon-large {
  font-size: 60px;
  color: white;
  font-weight: 700;
}

.success-title {
  font-size: 28px;
  color: #333;
  margin-bottom: 12px;
}

.success-desc {
  color: #666;
  margin-bottom: 30px;
}

.success-actions-large {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.view-orders-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 32px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.view-orders-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.continue-shop-btn {
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
  padding: 14px 32px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.continue-shop-btn:hover {
  background: #f8f9ff;
}

/* 结算底部按钮 */
.checkout-footer {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  padding: 20px 30px;
  border-top: 1px solid #e9ecef;
  background: white;
}

.prev-btn {
  background: #f8f9fa;
  color: #666;
  border: none;
  padding: 14px 28px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.prev-btn:hover {
  background: #e9ecef;
}

.next-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 32px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.next-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.submit-btn {
  background: linear-gradient(135deg, #e63946 0%, #c1121f 100%);
  color: white;
  border: none;
  padding: 14px 48px;
  border-radius: 10px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(230, 57, 70, 0.4);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-spinner-white {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  display: inline-block;
}

/* 提示消息 */
.toast {
  position: fixed;
  top: 100px;
  left: 50%;
  transform: translateX(-50%) translateY(-100px);
  background: white;
  padding: 16px 28px;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  gap: 12px;
  z-index: 2000;
  opacity: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast.show {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

.toast.success {
  border-left: 4px solid #10b981;
}

.toast.error {
  border-left: 4px solid #e63946;
}

.toast.warning {
  border-left: 4px solid #f59e0b;
}

.toast.info {
  border-left: 4px solid #667eea;
}

.toast-icon {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
}

.toast.success .toast-icon {
  background: #d1fae5;
  color: #10b981;
}

.toast.error .toast-icon {
  background: #fee2e2;
  color: #e63946;
}

.toast.warning .toast-icon {
  background: #fef3c7;
  color: #f59e0b;
}

.toast.info .toast-icon {
  background: #e0e7ff;
  color: #667eea;
}

.toast-message {
  font-weight: 500;
  color: #333;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .banner-title {
    font-size: 36px;
    flex-direction: column;
    gap: 8px;
  }
  
  .banner-subtitle {
    font-size: 16px;
  }
  
  .countdown-section {
    flex-direction: column;
    gap: 12px;
  }
  
  .time-box {
    padding: 8px 12px;
    min-width: 50px;
  }
  
  .time-value {
    font-size: 24px;
  }
  
  .banner-stats {
    gap: 30px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .products-grid {
    grid-template-columns: 1fr;
  }
  
  .detail-layout {
    grid-template-columns: 1fr;
    padding: 20px;
  }
  
  .cart-sidebar {
    width: 100%;
    right: -100%;
  }
  
  .section-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .checkout-modal .modal-content {
    max-width: 100%;
    max-height: 100vh;
    border-radius: 0;
  }
  
  .checkout-steps {
    padding: 16px;
  }
  
  .step-label {
    font-size: 10px;
  }
  
  .step-line {
    max-width: 30px;
  }
  
  .modal {
    padding: 0;
  }
  
  .modal-content {
    max-width: 100%;
    max-height: 100vh;
    border-radius: 0;
  }
  
  .floating-cart-btn {
    width: 56px;
    height: 56px;
    bottom: 20px;
    right: 20px;
  }
  
  .cart-icon {
    font-size: 24px;
  }
}

@media (max-width: 480px) {
  .container {
    padding: 0 16px;
  }
  
  .seckill-banner {
    padding: 60px 0 80px;
  }
  
  .category-item {
    padding: 10px 16px;
    font-size: 14px;
  }
  
  .product-info {
    padding: 16px;
  }
  
  .current-price {
    font-size: 22px;
  }
  
  .seckill-btn {
    padding: 12px 0;
    font-size: 14px;
  }
  
  .detail-actions {
    flex-direction: column;
  }
  
  .detail-seckill-btn,
  .detail-cart-btn {
    width: 100%;
  }
  
  .checkout-footer {
    flex-direction: column;
    gap: 12px;
  }
  
  .prev-btn,
  .next-btn,
  .submit-btn {
    width: 100%;
  }
  
  .success-actions,
  .success-actions-large {
    flex-direction: column;
  }
  
  .view-order-btn,
  .continue-btn,
  .view-orders-btn,
  .continue-shop-btn {
    width: 100%;
  }
}
</style>

