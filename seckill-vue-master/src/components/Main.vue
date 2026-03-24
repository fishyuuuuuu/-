<template>
  <div class="seckill-container">
    <!-- 导航栏 -->
    <Navbar />

    <!-- 图片轮播区域 -->
    <div class="carousel-section">
      <div class="carousel-container">
        <div class="carousel-wrapper">
          <div 
            class="carousel-track" 
            :style="{ transform: `translateX(-${currentSlide * 100}%)` }"
          >
            <div 
              v-for="(slide, index) in carouselSlides" 
              :key="index" 
              class="carousel-slide"
            >
              <img :src="slide.image" :alt="slide.title" class="slide-image">
              <div class="slide-content">
                <h2 class="slide-title">{{ slide.title }}</h2>
                <p class="slide-subtitle">{{ slide.subtitle }}</p>
                <button class="slide-btn" @click="handleSlideClick(slide)">
                  {{ slide.btnText }}
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 轮播控制按钮 -->
        <button class="carousel-btn prev" @click="prevSlide">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
        </button>
        <button class="carousel-btn next" @click="nextSlide">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>
        
        <!-- 轮播指示器 -->
        <div class="carousel-indicators">
          <span 
            v-for="(slide, index) in carouselSlides" 
            :key="index"
            :class="['indicator', { active: currentSlide === index }]"
            @click="goToSlide(index)"
          ></span>
        </div>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-section">
      <div class="container">
        <div class="search-wrapper">
          <div class="search-box">
            <svg class="search-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.35-4.35"></path>
            </svg>
            <input 
              type="text" 
              v-model="searchKeyword" 
              placeholder="搜索商品、品牌、分类..." 
              class="search-input"
              @keyup.enter="handleSearch"
            >
            <button class="search-btn" @click="handleSearch">搜索</button>
          </div>
          <div class="hot-keywords">
            <span class="hot-label">热门搜索:</span>
            <span 
              v-for="keyword in hotKeywords" 
              :key="keyword" 
              class="keyword-tag"
              @click="searchByKeyword(keyword)"
            >
              {{ keyword }}
            </span>
          </div>
        </div>
      </div>
    </div>

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

// 轮播相关状态
const currentSlide = ref(0);
const carouselInterval = ref(null);
const carouselSlides = ref([
  {
    image: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=1400&h=400&fit=crop',
    title: 'iPhone 17 Pro Max 限时特惠',
    subtitle: 'A19 Pro芯片，钛金属设计，最高立省2000元',
    btnText: '立即抢购',
    productId: 1
  },
  {
    image: 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=1400&h=400&fit=crop',
    title: 'MacBook Pro M4 新品首发',
    subtitle: 'M4 Pro芯片，性能怪兽，限时秒杀价',
    btnText: '查看详情',
    productId: 11
  },
  {
    image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=1400&h=400&fit=crop',
    title: 'Sony WH-1000XM6',
    subtitle: '业界顶级降噪，40小时续航，限时特价',
    btnText: '立即购买',
    productId: 6
  },
  {
    image: 'https://images.unsplash.com/photo-1592107761705-30a1bbc641e7?w=1400&h=400&fit=crop',
    title: 'Nintendo Switch 2',
    subtitle: '8英寸OLED屏幕，4K输出，游戏新体验',
    btnText: '抢购中',
    productId: 7
  }
]);

// 搜索相关状态
const searchKeyword = ref('');
const hotKeywords = ref(['iPhone 17', '手机', '笔记本', '耳机', '鞋子', '护肤', '运动', '食品']);

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
  { id: 6, name: '美妆护肤', icon: '💄' },
  { id: 7, name: '运动户外', icon: '⚽' },
  { id: 8, name: '食品饮料', icon: '🍎' }
]);

// 品类关键词映射（用于搜索）
const categoryKeywords = {
  '手机': [2],
  '手机数码': [2],
  '数码': [2],
  '智能手机': [2],
  '电话': [2],
  '电脑': [3],
  '电脑办公': [3],
  '笔记本': [3],
  '笔记本电脑': [3],
  '平板': [3],
  '平板电脑': [3],
  '家电': [4],
  '家用电器': [4],
  '电器': [4],
  '家居': [4],
  '服装': [5],
  '服装鞋包': [5],
  '衣服': [5],
  '鞋子': [5],
  '鞋': [5],
  '运动鞋': [5],
  '休闲鞋': [5],
  '包包': [5],
  '包': [5],
  '美妆': [6],
  '美妆护肤': [6],
  '护肤': [6],
  '化妆品': [6],
  '彩妆': [6],
  '运动': [7],
  '运动户外': [7],
  '户外': [7],
  '健身': [7],
  '食品': [8],
  '食品饮料': [8],
  '饮料': [8],
  '零食': [8],
  '耳机': [2],
  '音响': [2, 4],
  '键盘': [3],
  '鼠标': [3],
  '显示器': [3],
  '显卡': [3],
  '游戏': [2, 7],
  '游戏机': [2],
  '吸尘器': [4],
  '扫地机': [4],
  '空调': [4],
  '冰箱': [4],
  '洗衣机': [4],
  '电视': [4],
  '手表': [2],
  '智能手表': [2],
  '相机': [2],
  '摄影': [2]
};

// 扩展的商品数据（全面更新和扩充）
const extendedProducts = [
  // ==================== 手机数码类 (category: 2) ====================
  {
    id: 1,
    name: 'iPhone 17 Pro Max',
    price: 9999.99,
    originalPrice: 11999.99,
    stock: 100,
    description: 'A19 Pro芯片，钛金属设计，4800万像素三摄，AI智能摄影',
    image: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Apple',
    tags: ['手机', '智能手机', '苹果', 'iPhone', '旗舰'],
    specs: {
      '屏幕': '6.9英寸 OLED',
      '处理器': 'A19 Pro',
      '存储': '256GB',
      '颜色': '原色钛金属'
    }
  },
  {
    id: 2,
    name: 'iPhone 17',
    price: 6999.99,
    originalPrice: 7999.99,
    stock: 150,
    description: 'A18芯片，全新设计语言，超视网膜显示屏',
    image: 'https://images.unsplash.com/photo-1591337676887-a217a6970a8a?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Apple',
    tags: ['手机', '智能手机', '苹果', 'iPhone'],
    specs: {
      '屏幕': '6.1英寸 OLED',
      '处理器': 'A18',
      '存储': '128GB',
      '颜色': '星光色'
    }
  },
  {
    id: 3,
    name: 'Samsung Galaxy S25 Ultra',
    price: 8999.99,
    originalPrice: 9999.99,
    stock: 80,
    description: '骁龙8 Gen4，2亿像素主摄，S Pen手写笔，AI智能助手',
    image: 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Samsung',
    tags: ['手机', '智能手机', '三星', 'Galaxy', '旗舰'],
    specs: {
      '屏幕': '6.8英寸 AMOLED',
      '处理器': '骁龙8 Gen4',
      '存储': '512GB',
      '相机': '2亿像素'
    }
  },
  {
    id: 4,
    name: 'Xiaomi 15 Pro',
    price: 4999.99,
    originalPrice: 5999.99,
    stock: 120,
    description: '骁龙8 Gen4，徕卡影像，120W快充，IP68防水',
    image: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Xiaomi',
    tags: ['手机', '智能手机', '小米', '徕卡'],
    specs: {
      '屏幕': '6.73英寸 AMOLED',
      '处理器': '骁龙8 Gen4',
      '存储': '256GB',
      '快充': '120W'
    }
  },
  {
    id: 5,
    name: 'AirPods Pro 3',
    price: 1999.99,
    originalPrice: 2499.99,
    stock: 200,
    description: 'H3芯片，主动降噪2.0，空间音频，自适应通透模式',
    image: 'https://images.unsplash.com/photo-1600294037681-c80b4cb5b434?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Apple',
    tags: ['耳机', '无线耳机', '降噪耳机', '苹果'],
    specs: {
      '芯片': 'H3',
      '降噪': '主动降噪2.0',
      '续航': '36小时',
      '防水': 'IPX5'
    }
  },
  {
    id: 6,
    name: 'Sony WH-1000XM6',
    price: 2999.99,
    originalPrice: 3499.99,
    stock: 80,
    description: '业界顶级降噪，40小时续航，Hi-Res认证，舒适佩戴',
    image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Sony',
    tags: ['耳机', '头戴式耳机', '降噪耳机', '索尼'],
    specs: {
      '降噪': '顶级主动降噪',
      '续航': '40小时',
      '连接': '蓝牙5.3',
      '重量': '248g'
    }
  },
  {
    id: 7,
    name: 'Nintendo Switch 2',
    price: 2999.99,
    originalPrice: 3499.99,
    stock: 60,
    description: '8英寸OLED屏幕，4K输出，向下兼容，增强Joy-Con',
    image: 'https://images.unsplash.com/photo-1592107761705-30a1bbc641e7?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Nintendo',
    tags: ['游戏机', '游戏', '掌机', '任天堂'],
    specs: {
      '屏幕': '8英寸 OLED',
      '存储': '128GB',
      '输出': '4K',
      '续航': '12小时'
    }
  },
  {
    id: 8,
    name: 'Sony PlayStation 5 Pro',
    price: 4999.99,
    originalPrice: 5499.99,
    stock: 40,
    description: '8K游戏支持，光追加速，2TB SSD，DualSense Edge手柄',
    image: 'https://images.unsplash.com/photo-1606144042614-b2417e99c4e3?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Sony',
    tags: ['游戏机', '游戏', '主机', 'PS5', 'PlayStation'],
    specs: {
      '存储': '2TB SSD',
      '输出': '8K',
      '光追': '硬件加速',
      '手柄': 'DualSense Edge'
    }
  },
  {
    id: 9,
    name: 'Canon EOS R6 Mark III',
    price: 18999.99,
    originalPrice: 21999.99,
    stock: 30,
    description: '4500万像素，8K视频，机身防抖8档，双卡槽',
    image: 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Canon',
    tags: ['相机', '单反', '微单', '佳能', '摄影'],
    specs: {
      '像素': '4500万',
      '视频': '8K 60fps',
      '防抖': '8档',
      '对焦': '全像素双核'
    }
  },
  {
    id: 10,
    name: 'Apple Watch Ultra 3',
    price: 6999.99,
    originalPrice: 7999.99,
    stock: 70,
    description: '钛金属表壳，蓝宝石玻璃，100米防水，双频GPS',
    image: 'https://images.unsplash.com/photo-1434493789847-2f02dc6ca35d?w=400&h=400&fit=crop',
    category: 2,
    brand: 'Apple',
    tags: ['手表', '智能手表', '苹果', '运动'],
    specs: {
      '表壳': '钛金属',
      '防水': '100米',
      '续航': '72小时',
      'GPS': '双频L1+L5'
    }
  },

  // ==================== 电脑办公类 (category: 3) ====================
  {
    id: 11,
    name: 'MacBook Pro 16 M4 Pro',
    price: 19999.99,
    originalPrice: 22999.99,
    stock: 50,
    description: 'M4 Pro芯片，Liquid视网膜XDR显示屏，24小时续航',
    image: 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400&h=400&fit=crop',
    category: 3,
    brand: 'Apple',
    tags: ['笔记本', '笔记本电脑', '苹果', 'MacBook', '电脑'],
    specs: {
      '屏幕': '16.2英寸 XDR',
      '处理器': 'M4 Pro',
      '内存': '32GB',
      '存储': '1TB SSD'
    }
  },
  {
    id: 12,
    name: 'MacBook Air 15 M3',
    price: 10999.99,
    originalPrice: 12999.99,
    stock: 80,
    description: 'M3芯片，15.3英寸Liquid视网膜屏，18小时续航',
    image: 'https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=400&h=400&fit=crop',
    category: 3,
    brand: 'Apple',
    tags: ['笔记本', '笔记本电脑', '苹果', 'MacBook', '电脑'],
    specs: {
      '屏幕': '15.3英寸',
      '处理器': 'M3',
      '内存': '16GB',
      '存储': '512GB SSD'
    }
  },
  {
    id: 13,
    name: 'Dell XPS 15',
    price: 13999.99,
    originalPrice: 16999.99,
    stock: 45,
    description: 'Intel酷睿Ultra 9，RTX 4070，4K OLED触控屏',
    image: 'https://images.unsplash.com/photo-1593642632559-0c6d3fc62b89?w=400&h=400&fit=crop',
    category: 3,
    brand: 'Dell',
    tags: ['笔记本', '笔记本电脑', '戴尔', '电脑', '游戏本'],
    specs: {
      '屏幕': '15.6英寸 4K OLED',
      '处理器': 'Intel Ultra 9',
      '显卡': 'RTX 4070',
      '存储': '1TB SSD'
    }
  },
  {
    id: 14,
    name: 'iPad Pro 13 M4',
    price: 9999.99,
    originalPrice: 11999.99,
    stock: 60,
    description: 'M4芯片，Ultra Retina XDR显示屏，Apple Pencil Pro',
    image: 'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=400&h=400&fit=crop',
    category: 3,
    brand: 'Apple',
    tags: ['平板', '平板电脑', 'iPad', '苹果'],
    specs: {
      '屏幕': '13英寸 Ultra Retina',
      '处理器': 'M4',
      '存储': '256GB',
      '网络': 'WiFi + 5G'
    }
  },
  {
    id: 15,
    name: 'ASUS ROG 游戏显示器 32',
    price: 5999.99,
    originalPrice: 6999.99,
    stock: 35,
    description: '32英寸4K 240Hz OLED，0.03ms响应，G-Sync',
    image: 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=400&h=400&fit=crop',
    category: 3,
    brand: 'ASUS',
    tags: ['显示器', '游戏显示器', '电脑配件', '华硕'],
    specs: {
      '尺寸': '32英寸',
      '分辨率': '4K',
      '刷新率': '240Hz',
      '面板': 'OLED'
    }
  },
  {
    id: 16,
    name: 'NVIDIA RTX 5090',
    price: 14999.99,
    originalPrice: 16999.99,
    stock: 20,
    description: '32GB GDDR7，DLSS 4.0，光追3.0，AI加速',
    image: 'https://images.unsplash.com/photo-1591488320449-011701bb6704?w=400&h=400&fit=crop',
    category: 3,
    brand: 'NVIDIA',
    tags: ['显卡', '电脑配件', '游戏', 'GPU'],
    specs: {
      '显存': '32GB GDDR7',
      '架构': 'Blackwell',
      'DLSS': '4.0',
      '接口': 'PCIe 5.0'
    }
  },
  {
    id: 17,
    name: 'Logitech MX Master 4',
    price: 899.99,
    originalPrice: 1099.99,
    stock: 150,
    description: '人体工学设计，8000 DPI，多设备切换，静音点击',
    image: 'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=400&h=400&fit=crop',
    category: 3,
    brand: 'Logitech',
    tags: ['鼠标', '电脑配件', '办公', '罗技'],
    specs: {
      'DPI': '8000',
      '连接': '蓝牙/2.4G',
      '续航': '70天',
      '按键': '8个'
    }
  },
  {
    id: 18,
    name: 'Keychron Q1 Pro',
    price: 1299.99,
    originalPrice: 1599.99,
    stock: 100,
    description: 'QMK/VIA支持，Gasket结构，热插拔轴体，RGB背光',
    image: 'https://images.unsplash.com/photo-1595225476474-87563907a212?w=400&h=400&fit=crop',
    category: 3,
    brand: 'Keychron',
    tags: ['键盘', '机械键盘', '电脑配件', '办公'],
    specs: {
      '轴体': '热插拔',
      '连接': '蓝牙/有线',
      '背光': 'RGB',
      '结构': 'Gasket'
    }
  },

  // ==================== 家用电器类 (category: 4) ====================
  {
    id: 19,
    name: 'Dyson V16 Detect',
    price: 5499.99,
    originalPrice: 6499.99,
    stock: 50,
    description: '激光探测灰尘，智能吸力调节，70分钟续航',
    image: 'https://images.unsplash.com/photo-1558317374-067fb5f30001?w=400&h=400&fit=crop',
    category: 4,
    brand: 'Dyson',
    tags: ['吸尘器', '家电', '戴森', '清洁'],
    specs: {
      '吸力': '262AW',
      '续航': '70分钟',
      '重量': '2.9kg',
      '尘杯': '0.76L'
    }
  },
  {
    id: 20,
    name: 'Roborock S8 Pro Ultra',
    price: 5999.99,
    originalPrice: 6999.99,
    stock: 40,
    description: 'AI避障，自动集尘，自动洗拖布，60°C热水洗',
    image: 'https://images.unsplash.com/photo-1589894404892-7310b9b71f7b?w=400&h=400&fit=crop',
    category: 4,
    brand: 'Roborock',
    tags: ['扫地机', '扫地机器人', '家电', '清洁'],
    specs: {
      '吸力': '6000Pa',
      '导航': 'LDS+AI',
      '续航': '180分钟',
      '水箱': '250ml'
    }
  },
  {
    id: 21,
    name: 'DJI Mini 5 Pro',
    price: 6999.99,
    originalPrice: 7999.99,
    stock: 35,
    description: '4K/120fps，48MP照片，47分钟续航，全向避障',
    image: 'https://images.unsplash.com/photo-1473968512647-3e447244af8f?w=400&h=400&fit=crop',
    category: 4,
    brand: 'DJI',
    tags: ['无人机', '航拍', '大疆', '摄影'],
    specs: {
      '重量': '249g',
      '视频': '4K/120fps',
      '续航': '47分钟',
      '避障': '全向'
    }
  },
  {
    id: 22,
    name: 'Midea 智能空调 3匹',
    price: 7999.99,
    originalPrice: 9999.99,
    stock: 25,
    description: '新一级能效，AI智能控温，自清洁，静音设计',
    image: 'https://images.unsplash.com/photo-1631567091046-3b31a31d1f76?w=400&h=400&fit=crop',
    category: 4,
    brand: 'Midea',
    tags: ['空调', '家电', '美的', '制冷'],
    specs: {
      '功率': '3匹',
      '能效': '新一级',
      '噪音': '18dB',
      '功能': 'AI控温'
    }
  },
  {
    id: 23,
    name: 'Samsung 65寸 Neo QLED电视',
    price: 12999.99,
    originalPrice: 15999.99,
    stock: 20,
    description: '8K分辨率，Mini LED背光，AI画质增强，120Hz',
    image: 'https://images.unsplash.com/photo-1593359677879-a4bb92f829d1?w=400&h=400&fit=crop',
    category: 4,
    brand: 'Samsung',
    tags: ['电视', '家电', '三星', '智能电视'],
    specs: {
      '尺寸': '65英寸',
      '分辨率': '8K',
      '刷新率': '120Hz',
      '背光': 'Mini LED'
    }
  },

  // ==================== 服装鞋包类 (category: 5) ====================
  {
    id: 24,
    name: 'Nike Air Max 2026',
    price: 1299.99,
    originalPrice: 1699.99,
    stock: 200,
    description: '全新Air缓震科技，Flyknit鞋面，透气舒适',
    image: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=400&h=400&fit=crop',
    category: 5,
    brand: 'Nike',
    tags: ['鞋子', '运动鞋', '耐克', '休闲鞋'],
    specs: {
      '鞋面': 'Flyknit',
      '中底': 'Air Max',
      '闭合': '系带',
      '适用': '跑步/休闲'
    }
  },
  {
    id: 25,
    name: 'Adidas Ultraboost 24',
    price: 1499.99,
    originalPrice: 1899.99,
    stock: 180,
    description: 'Boost中底，Primeknit鞋面，Continental橡胶大底',
    image: 'https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=400&h=400&fit=crop',
    category: 5,
    brand: 'Adidas',
    tags: ['鞋子', '运动鞋', '阿迪达斯', '跑步鞋'],
    specs: {
      '鞋面': 'Primeknit',
      '中底': 'Boost',
      '大底': 'Continental',
      '重量': '310g'
    }
  },
  {
    id: 26,
    name: 'New Balance 990v6',
    price: 1599.99,
    originalPrice: 1999.99,
    stock: 120,
    description: '美产经典，ENCAP中底，猪巴革鞋面，复古设计',
    image: 'https://images.unsplash.com/photo-1539185441755-769473a23570?w=400&h=400&fit=crop',
    category: 5,
    brand: 'New Balance',
    tags: ['鞋子', '运动鞋', '休闲鞋', '复古'],
    specs: {
      '鞋面': '猪巴革/网布',
      '中底': 'ENCAP',
      '产地': '美国',
      '风格': '复古'
    }
  },
  {
    id: 27,
    name: 'Louis Vuitton Neverfull MM',
    price: 15999.99,
    originalPrice: 18999.99,
    stock: 15,
    description: '经典老花帆布，植鞣牛皮手柄，大容量设计',
    image: 'https://images.unsplash.com/photo-1548036328-c9fa89d128fa?w=400&h=400&fit=crop',
    category: 5,
    brand: 'Louis Vuitton',
    tags: ['包包', '手提包', '奢侈品', 'LV'],
    specs: {
      '材质': '老花帆布',
      '手柄': '植鞣牛皮',
      '尺寸': '31x28x14cm',
      '内衬': '条纹帆布'
    }
  },
  {
    id: 28,
    name: 'Canada Goose 远征派克大衣',
    price: 9999.99,
    originalPrice: 12999.99,
    stock: 30,
    description: '625蓬白鸭绒，防风防水，-30°C保暖',
    image: 'https://images.unsplash.com/photo-1539533018447-63fcce2678e3?w=400&h=400&fit=crop',
    category: 5,
    brand: 'Canada Goose',
    tags: ['衣服', '羽绒服', '大衣', '加拿大鹅'],
    specs: {
      '填充': '625蓬白鸭绒',
      '面料': 'Arctic Tech',
      '温度': '-30°C',
      '毛领': '郊狼毛'
    }
  },
  {
    id: 29,
    name: 'Uniqlo U系列羽绒服',
    price: 799.99,
    originalPrice: 999.99,
    stock: 300,
    description: '轻薄保暖，便携收纳，90%白鸭绒',
    image: 'https://images.unsplash.com/photo-1591047139829-d91aecb6caea?w=400&h=400&fit=crop',
    category: 5,
    brand: 'Uniqlo',
    tags: ['衣服', '羽绒服', '优衣库', '休闲'],
    specs: {
      '填充': '90%白鸭绒',
      '面料': '尼龙',
      '特点': '便携收纳',
      '风格': '简约'
    }
  },

  // ==================== 美妆护肤类 (category: 6) ====================
  {
    id: 30,
    name: 'SK-II 神仙水 330ml',
    price: 1799.99,
    originalPrice: 2199.99,
    stock: 100,
    description: 'Pitera精华，改善肤质，提亮肤色，深层保湿',
    image: 'https://images.unsplash.com/photo-1571781926291-c477ebfd024b?w=400&h=400&fit=crop',
    category: 6,
    brand: 'SK-II',
    tags: ['护肤', '精华', '化妆水', '化妆品'],
    specs: {
      '容量': '330ml',
      '成分': 'Pitera',
      '功效': '保湿/提亮',
      '适用': '所有肤质'
    }
  },
  {
    id: 31,
    name: 'La Mer 海蓝之谜面霜',
    price: 3299.99,
    originalPrice: 3899.99,
    stock: 50,
    description: '深海神奇活性精萃，修护肌肤，奢华滋养',
    image: 'https://images.unsplash.com/photo-1608248597279-f99d06bf9dd2?w=400&h=400&fit=crop',
    category: 6,
    brand: 'La Mer',
    tags: ['护肤', '面霜', '奢侈品', '化妆品'],
    specs: {
      '容量': '60ml',
      '成分': '深海精萃',
      '功效': '修护/滋养',
      '适用': '干性/熟龄肌'
    }
  },
  {
    id: 32,
    name: 'Estée Lauder 小棕瓶精华',
    price: 899.99,
    originalPrice: 1099.99,
    stock: 150,
    description: '二裂酵母，夜间修护，抗初老，淡化细纹',
    image: 'https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=400&h=400&fit=crop',
    category: 6,
    brand: 'Estée Lauder',
    tags: ['护肤', '精华', '抗衰老', '化妆品'],
    specs: {
      '容量': '50ml',
      '成分': '二裂酵母',
      '功效': '修护/抗老',
      '适用': '所有肤质'
    }
  },
  {
    id: 33,
    name: 'Dior 迪奥999口红',
    price: 399.99,
    originalPrice: 499.99,
    stock: 200,
    description: '经典正红色，丝绒质地，持久显色',
    image: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=400&h=400&fit=crop',
    category: 6,
    brand: 'Dior',
    tags: ['彩妆', '口红', '迪奥', '化妆品'],
    specs: {
      '色号': '999',
      '质地': '丝绒',
      '显色': '高显色',
      '持久': '持久'
    }
  },
  {
    id: 34,
    name: 'CHANEL 香奈儿5号香水',
    price: 1299.99,
    originalPrice: 1599.99,
    stock: 80,
    description: '经典花香调，优雅迷人，持久留香',
    image: 'https://images.unsplash.com/photo-1541643600914-78b084683601?w=400&h=400&fit=crop',
    category: 6,
    brand: 'CHANEL',
    tags: ['香水', '香奈儿', '奢侈品', '化妆品'],
    specs: {
      '容量': '100ml',
      '香调': '花香调',
      '浓度': 'EDP',
      '留香': '8-10小时'
    }
  },

  // ==================== 运动户外类 (category: 7) ====================
  {
    id: 35,
    name: 'Lululemon Align瑜伽裤',
    price: 899.99,
    originalPrice: 1099.99,
    stock: 150,
    description: 'Nulu面料，裸感体验，高腰设计，透气速干',
    image: 'https://images.unsplash.com/photo-1506629082955-511b1aa562c8?w=400&h=400&fit=crop',
    category: 7,
    brand: 'Lululemon',
    tags: ['运动', '瑜伽', '健身', '瑜伽裤'],
    specs: {
      '面料': 'Nulu',
      '特点': '裸感',
      '腰型': '高腰',
      '功能': '透气速干'
    }
  },
  {
    id: 36,
    name: 'Garmin Fenix 8',
    price: 6999.99,
    originalPrice: 7999.99,
    stock: 40,
    description: '专业户外手表，GPS导航，心率血氧，太阳能充电',
    image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400&h=400&fit=crop',
    category: 7,
    brand: 'Garmin',
    tags: ['运动', '手表', '户外', '健身'],
    specs: {
      '防水': '100米',
      '续航': '37天',
      'GPS': '多频多星',
      '功能': '心率/血氧/导航'
    }
  },
  {
    id: 37,
    name: 'The North Face 冲锋衣',
    price: 2499.99,
    originalPrice: 2999.99,
    stock: 80,
    description: 'Gore-Tex面料，防风防水，透气舒适，专业户外',
    image: 'https://images.unsplash.com/photo-1544923246-77307dd628b8?w=400&h=400&fit=crop',
    category: 7,
    brand: 'The North Face',
    tags: ['运动', '户外', '冲锋衣', '衣服'],
    specs: {
      '面料': 'Gore-Tex',
      '防水': '10000mm',
      '透气': '透气膜',
      '适用': '登山/徒步'
    }
  },
  {
    id: 38,
    name: 'Yeti 保温杯 36oz',
    price: 399.99,
    originalPrice: 499.99,
    stock: 200,
    description: '18/8不锈钢，真空保温，防漏设计，户外必备',
    image: 'https://images.unsplash.com/photo-1602143407151-7111542de6e8?w=400&h=400&fit=crop',
    category: 7,
    brand: 'Yeti',
    tags: ['运动', '户外', '保温杯', '水杯'],
    specs: {
      '容量': '36oz/1L',
      '材质': '18/8不锈钢',
      '保温': '24小时',
      '特点': '防漏'
    }
  },

  // ==================== 食品饮料类 (category: 8) ====================
  {
    id: 39,
    name: '星巴克精选咖啡豆 250g',
    price: 128.99,
    originalPrice: 158.99,
    stock: 500,
    description: '阿拉比卡咖啡豆，中度烘焙，香气浓郁',
    image: 'https://images.unsplash.com/photo-1559056199-641a0ac8b55e?w=400&h=400&fit=crop',
    category: 8,
    brand: 'Starbucks',
    tags: ['食品', '咖啡', '饮料', '零食'],
    specs: {
      '重量': '250g',
      '产地': '拉丁美洲',
      '烘焙': '中度',
      '风味': '坚果/巧克力'
    }
  },
  {
    id: 40,
    name: '三只松鼠坚果礼盒',
    price: 199.99,
    originalPrice: 259.99,
    stock: 300,
    description: '精选坚果组合，每日坚果，营养健康',
    image: 'https://images.unsplash.com/photo-1536816579748-4ecb3f03d72a?w=400&h=400&fit=crop',
    category: 8,
    brand: '三只松鼠',
    tags: ['食品', '坚果', '零食', '礼盒'],
    specs: {
      '规格': '30包/盒',
      '内容': '核桃/腰果/杏仁等',
      '保质期': '12个月',
      '特点': '独立包装'
    }
  },
  {
    id: 41,
    name: '农夫山泉矿泉水 24瓶',
    price: 49.99,
    originalPrice: 59.99,
    stock: 1000,
    description: '天然矿泉水，深层水源，矿物质丰富',
    image: 'https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=400&h=400&fit=crop',
    category: 8,
    brand: '农夫山泉',
    tags: ['饮料', '矿泉水', '水', '食品'],
    specs: {
      '规格': '550ml x 24瓶',
      '水源': '深层地下水',
      'pH值': '7.3±0.5',
      '特点': '天然矿物质'
    }
  },
  {
    id: 42,
    name: '茅台飞天53度 500ml',
    price: 2999.99,
    originalPrice: 3499.99,
    stock: 50,
    description: '酱香型白酒，国酒经典，收藏送礼佳品',
    image: 'https://images.unsplash.com/photo-1607622750671-6cd9a99eabd1?w=400&h=400&fit=crop',
    category: 8,
    brand: '茅台',
    tags: ['饮料', '白酒', '酒', '礼品'],
    specs: {
      '度数': '53度',
      '容量': '500ml',
      '香型': '酱香型',
      '产地': '贵州茅台镇'
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

// 搜索过滤后的商品
const filteredProducts = computed(() => {
  if (!searchKeyword.value.trim()) {
    return products.value;
  }
  const keyword = searchKeyword.value.toLowerCase().trim();
  
  // 检查是否匹配品类关键词
  const matchedCategoryIds = categoryKeywords[keyword] || [];
  
  return products.value.filter(product => {
    // 品类匹配
    if (matchedCategoryIds.length > 0 && matchedCategoryIds.includes(product.category)) {
      return true;
    }
    
    // 商品名称匹配
    if (product.name.toLowerCase().includes(keyword)) {
      return true;
    }
    
    // 描述匹配
    if (product.description.toLowerCase().includes(keyword)) {
      return true;
    }
    
    // 品牌匹配
    if (product.brand && product.brand.toLowerCase().includes(keyword)) {
      return true;
    }
    
    // 标签匹配
    if (product.tags && product.tags.some(tag => tag.toLowerCase().includes(keyword))) {
      return true;
    }
    
    // 规格匹配
    if (product.specs && Object.values(product.specs).some(v => 
      String(v).toLowerCase().includes(keyword)
    )) {
      return true;
    }
    
    return false;
  });
});

// 轮播控制方法
const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % carouselSlides.value.length;
};

const prevSlide = () => {
  currentSlide.value = (currentSlide.value - 1 + carouselSlides.value.length) % carouselSlides.value.length;
};

const goToSlide = (index) => {
  currentSlide.value = index;
};

const startCarouselAutoplay = () => {
  stopCarouselAutoplay();
  carouselInterval.value = setInterval(() => {
    nextSlide();
  }, 5000);
};

const stopCarouselAutoplay = () => {
  if (carouselInterval.value) {
    clearInterval(carouselInterval.value);
    carouselInterval.value = null;
  }
};

const handleSlideClick = (slide) => {
  if (slide.productId) {
    const product = products.value.find(p => p.id === slide.productId);
    if (product) {
      showProductDetail(product);
    }
  }
};

// 搜索方法
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    showToastMessage(`搜索: ${searchKeyword.value}`, 'info');
  }
};

const searchByKeyword = (keyword) => {
  searchKeyword.value = keyword;
  handleSearch();
};

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

// 排序商品（已移至上方的computed定义，此处保留分类过滤逻辑）
const sortedProducts = computed(() => {
  let result = [...filteredProducts.value];
  
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
    // 使用相对路径，通过Vite代理转发到后端
    const response = await axios.get('/api/captcha', {
      responseType: 'blob',
      headers: {
        'Cache-Control': 'no-cache'
      }
    });
    
    // 从响应头获取验证码ID
    captchaId.value = response.headers['x-captcha-id'] || '';
    
    // 将blob转换为图片URL
    const blob = new Blob([response.data], { type: 'image/png' });
    captchaImageUrl.value = URL.createObjectURL(blob);
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

    // 调用后端秒杀API
    const response = await axios.post('/api/product/seckill', {
      productId: productId,
      captchaId: captchaId.value,
      captchaStr: captchaInput.value
    }, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    // 秒杀成功
    product.stock--;

    showCaptchaModal.value = false;
    successProductName.value = product.name;
    showSuccessModal.value = true;
    showToastMessage('秒杀成功！订单已创建', 'success');

    // 刷新商品列表以获取最新库存
    await fetchProducts();
  } catch (err) {
    console.error('秒杀失败:', err);
    showToastMessage(err.response?.data?.error || err.message || '秒杀失败，请重试', 'error');
    // 刷新验证码
    await fetchCaptcha();
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

  // 启动轮播自动播放
  startCarouselAutoplay();

  const cleanup = () => {
    clearInterval(timer);
    stopCarouselAutoplay();
  };
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

/* 图片轮播区域 */
.carousel-section {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  padding: 0;
}

.carousel-container {
  position: relative;
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  overflow: hidden;
}

.carousel-wrapper {
  width: 100%;
  overflow: hidden;
}

.carousel-track {
  display: flex;
  transition: transform 0.5s ease-in-out;
}

.carousel-slide {
  min-width: 100%;
  position: relative;
  height: 400px;
}

.slide-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.slide-content {
  position: absolute;
  top: 50%;
  left: 60px;
  transform: translateY(-50%);
  color: white;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
  z-index: 2;
}

.slide-title {
  font-size: 42px;
  font-weight: 700;
  margin-bottom: 12px;
  background: linear-gradient(90deg, #fff, #f0f0f0);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.slide-subtitle {
  font-size: 18px;
  margin-bottom: 24px;
  opacity: 0.9;
}

.slide-btn {
  padding: 14px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 30px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.slide-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.carousel-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border: none;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  z-index: 10;
}

.carousel-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-50%) scale(1.1);
}

.carousel-btn.prev {
  left: 20px;
}

.carousel-btn.next {
  right: 20px;
}

.carousel-indicators {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 10px;
  z-index: 10;
}

.indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator.active {
  background: white;
  width: 32px;
  border-radius: 6px;
}

/* 搜索栏区域 */
.search-section {
  background: white;
  padding: 24px 0;
  border-bottom: 1px solid #eee;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.search-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  background: #f5f5f5;
  border-radius: 30px;
  padding: 4px;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.search-box:focus-within {
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.search-icon {
  margin-left: 16px;
  color: #999;
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  padding: 12px 16px;
  font-size: 15px;
  outline: none;
}

.search-input::placeholder {
  color: #999;
}

.search-btn {
  padding: 12px 28px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.search-btn:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.hot-keywords {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.hot-label {
  font-size: 13px;
  color: #999;
}

.keyword-tag {
  padding: 6px 14px;
  background: #f0f0f0;
  border-radius: 16px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  transition: all 0.3s ease;
}

.keyword-tag:hover {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
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
