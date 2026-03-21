<template>
  <div class="profile-page">
    <!-- 顶部导航 -->
    <Navbar />

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <div class="header-content">
          <div class="header-title">
            <h1>个人中心</h1>
            <p class="header-subtitle">管理您的账户信息和订单</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 用户信息卡片 -->
    <div class="container">
      <div class="user-card">
        <div class="user-info-section">
          <div class="user-avatar">
            <img :src="user.avatar" :alt="user.name" class="avatar-img">
            <div class="avatar-edit" @click="showEditAvatar = true">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
            </div>
          </div>
          <div class="user-details">
            <div class="user-name-row">
              <h2 class="user-name">{{ user.name }}</h2>
              <span class="user-level" :class="user.levelClass">{{ user.level }}</span>
            </div>
            <p class="user-phone">{{ user.phone }}</p>
            <p class="user-email">{{ user.email }}</p>
            <div class="user-tags">
              <span class="user-tag">实名认证</span>
              <span class="user-tag">手机绑定</span>
            </div>
          </div>
          <div class="user-actions">
            <button class="btn btn-outline" @click="showEditProfile = true">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
              编辑资料
            </button>
          </div>
        </div>

        <!-- 用户数据统计 -->
        <div class="user-stats">
          <div class="stat-card" @click="$router.push('/order')">
            <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
                <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
              </svg>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ user.ordersCount }}</span>
              <span class="stat-label">全部订单</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
              </svg>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ user.favoritesCount }}</span>
              <span class="stat-label">我的收藏</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
                <line x1="1" y1="10" x2="23" y2="10"/>
              </svg>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ user.couponsCount }}</span>
              <span class="stat-label">优惠券</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="8" r="7"/>
                <polyline points="8.21 13.89 7 23 12 20 17 23 15.79 13.88"/>
              </svg>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ user.points }}</span>
              <span class="stat-label">积分</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 订单快捷入口 -->
      <div class="section-card">
        <div class="section-header">
          <h3 class="section-title">
            <span class="title-icon">📦</span>
            我的订单
          </h3>
          <button class="view-all-btn" @click="$router.push('/order')">
            查看全部
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 18l6-6-6-6"/>
            </svg>
          </button>
        </div>
        <div class="order-shortcuts">
          <div class="shortcut-item" @click="goToOrders('pending')">
            <div class="shortcut-icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
            </div>
            <span class="shortcut-label">待支付</span>
            <span v-if="user.pendingOrders > 0" class="shortcut-badge">{{ user.pendingOrders }}</span>
          </div>
          <div class="shortcut-item" @click="goToOrders('paid')">
            <div class="shortcut-icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
              </svg>
            </div>
            <span class="shortcut-label">待发货</span>
            <span v-if="user.paidOrders > 0" class="shortcut-badge">{{ user.paidOrders }}</span>
          </div>
          <div class="shortcut-item" @click="goToOrders('shipped')">
            <div class="shortcut-icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="1" y="3" width="15" height="13"/>
                <polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/>
                <circle cx="5.5" cy="18.5" r="2.5"/>
                <circle cx="18.5" cy="18.5" r="2.5"/>
              </svg>
            </div>
            <span class="shortcut-label">待收货</span>
            <span v-if="user.shippedOrders > 0" class="shortcut-badge">{{ user.shippedOrders }}</span>
          </div>
          <div class="shortcut-item" @click="goToOrders('completed')">
            <div class="shortcut-icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
            </div>
            <span class="shortcut-label">已完成</span>
          </div>
          <div class="shortcut-item">
            <div class="shortcut-icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
            </div>
            <span class="shortcut-label">售后/退款</span>
          </div>
        </div>
      </div>

      <!-- 功能菜单 -->
      <div class="menu-grid">
        <!-- 账户管理 -->
        <div class="section-card">
          <div class="section-header">
            <h3 class="section-title">
              <span class="title-icon">👤</span>
              账户管理
            </h3>
          </div>
          <div class="menu-list">
            <div class="menu-item" @click="showEditProfile = true">
              <div class="menu-icon" style="background: #e3f2fd; color: #1976d2;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
              </div>
              <span class="menu-text">个人信息</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item" @click="showAddressModal = true">
              <div class="menu-icon" style="background: #f3e5f5; color: #7b1fa2;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                  <circle cx="12" cy="10" r="3"/>
                </svg>
              </div>
              <span class="menu-text">收货地址</span>
              <span class="menu-badge">{{ user.addressCount }}个</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #e8f5e9; color: #388e3c;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
                  <line x1="1" y1="10" x2="23" y2="10"/>
                </svg>
              </div>
              <span class="menu-text">我的钱包</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #fff3e0; color: #f57c00;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="12" y1="16" x2="12" y2="12"/>
                  <line x1="12" y1="8" x2="12.01" y2="8"/>
                </svg>
              </div>
              <span class="menu-text">账户安全</span>
              <span class="menu-status safe">安全</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
          </div>
        </div>

        <!-- 我的服务 -->
        <div class="section-card">
          <div class="section-header">
            <h3 class="section-title">
              <span class="title-icon">🎯</span>
              我的服务
            </h3>
          </div>
          <div class="menu-list">
            <div class="menu-item">
              <div class="menu-icon" style="background: #fce4ec; color: #c2185b;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                </svg>
              </div>
              <span class="menu-text">我的收藏</span>
              <span class="menu-badge">{{ user.favoritesCount }}件</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #e0f2f1; color: #00796b;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="12" y1="8" x2="12" y2="12"/>
                  <line x1="12" y1="16" x2="12.01" y2="16"/>
                </svg>
              </div>
              <span class="menu-text">浏览记录</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #e8eaf6; color: #3f51b5;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>
              <span class="menu-text">我的评价</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #fbe9e7; color: #d84315;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                </svg>
              </div>
              <span class="menu-text">秒杀提醒</span>
              <span class="menu-badge">{{ user.reminders }}个</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
          </div>
        </div>

        <!-- 客户服务 -->
        <div class="section-card">
          <div class="section-header">
            <h3 class="section-title">
              <span class="title-icon">🎧</span>
              客户服务
            </h3>
          </div>
          <div class="menu-list">
            <div class="menu-item">
              <div class="menu-icon" style="background: #e1f5fe; color: #0288d1;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
                  <line x1="12" y1="17" x2="12.01" y2="17"/>
                </svg>
              </div>
              <span class="menu-text">帮助中心</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #f1f8e9; color: #689f38;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
                </svg>
              </div>
              <span class="menu-text">在线客服</span>
              <span class="menu-status online">在线</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #efebe9; color: #5d4037;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14 2 14 8 20 8"/>
                  <line x1="16" y1="13" x2="8" y2="13"/>
                  <line x1="16" y1="17" x2="8" y2="17"/>
                  <polyline points="10 9 9 9 8 9"/>
                </svg>
              </div>
              <span class="menu-text">意见反馈</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
            <div class="menu-item">
              <div class="menu-icon" style="background: #fff9c4; color: #f9a825;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="12" y1="16" x2="12" y2="12"/>
                  <line x1="12" y1="8" x2="12.01" y2="8"/>
                </svg>
              </div>
              <span class="menu-text">关于我们</span>
              <svg class="menu-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- 退出登录 -->
      <div class="logout-section">
        <button class="btn btn-logout" @click="handleLogout">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
          退出登录
        </button>
      </div>
    </div>

    <!-- 编辑资料弹窗 -->
    <div class="modal" v-if="showEditProfile" @click.self="showEditProfile = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>编辑资料</h3>
          <button class="close-btn" @click="showEditProfile = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>昵称</label>
            <input type="text" v-model="editForm.name" placeholder="请输入昵称">
          </div>
          <div class="form-group">
            <label>手机号</label>
            <input type="tel" v-model="editForm.phone" placeholder="请输入手机号" disabled>
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input type="email" v-model="editForm.email" placeholder="请输入邮箱">
          </div>
          <div class="form-group">
            <label>性别</label>
            <div class="radio-group">
              <label class="radio-label">
                <input type="radio" v-model="editForm.gender" value="male">
                <span>男</span>
              </label>
              <label class="radio-label">
                <input type="radio" v-model="editForm.gender" value="female">
                <span>女</span>
              </label>
              <label class="radio-label">
                <input type="radio" v-model="editForm.gender" value="secret">
                <span>保密</span>
              </label>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showEditProfile = false">取消</button>
          <button class="btn btn-primary" @click="saveProfile" :disabled="isSaving">
            <span v-if="!isSaving">保存</span>
            <span v-else class="btn-spinner-white"></span>
          </button>
        </div>
      </div>
    </div>

    <!-- 收货地址弹窗 -->
    <div class="modal" v-if="showAddressModal" @click.self="showAddressModal = false">
      <div class="modal-content address-modal">
        <div class="modal-header">
          <h3>收货地址</h3>
          <button class="close-btn" @click="showAddressModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="address-list">
            <div v-for="address in addresses" :key="address.id" class="address-card" :class="{ default: address.isDefault }">
              <div class="address-info">
                <div class="address-header">
                  <span class="address-name">{{ address.name }}</span>
                  <span class="address-phone">{{ address.phone }}</span>
                  <span v-if="address.isDefault" class="default-badge">默认</span>
                </div>
                <p class="address-detail">{{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}</p>
              </div>
              <div class="address-actions">
                <button class="action-btn" @click="editAddress(address)">编辑</button>
                <button class="action-btn delete" @click="deleteAddress(address)">删除</button>
              </div>
            </div>
          </div>
          <button class="btn btn-primary btn-add" @click="showAddAddress = true">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            添加新地址
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
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from './Navbar.vue';

const router = useRouter();

// 用户信息
const user = reactive({
  name: '张三',
  phone: '138****8888',
  email: 'zhangsan@example.com',
  level: 'VIP会员',
  levelClass: 'vip',
  avatar: 'https://placehold.co/120x120/667eea/ffffff?text=ZS',
  ordersCount: 12,
  pendingOrders: 2,
  paidOrders: 1,
  shippedOrders: 3,
  favoritesCount: 8,
  couponsCount: 5,
  points: 2580,
  addressCount: 3,
  reminders: 4
});

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
  },
  {
    id: 3,
    name: '李四',
    phone: '139****9999',
    province: '广东省',
    city: '深圳市',
    district: '南山区',
    detail: '科技园南区腾讯大厦',
    isDefault: false
  }
]);

// 弹窗状态
const showEditProfile = ref(false);
const showAddressModal = ref(false);
const showAddAddress = ref(false);
const isSaving = ref(false);

// 编辑表单
const editForm = reactive({
  name: user.name,
  phone: user.phone,
  email: user.email,
  gender: 'male'
});

// 提示消息
const showToast = ref(false);
const toastMessage = ref('');
const toastType = ref('success');
const toastIcon = ref('✓');

// 显示提示
const showToastMessage = (message, type = 'success') => {
  toastMessage.value = message;
  toastType.value = type;
  toastIcon.value = type === 'success' ? '✓' : type === 'error' ? '✗' : '⚠';
  showToast.value = true;
  setTimeout(() => {
    showToast.value = false;
  }, 3000);
};

// 跳转到订单页面
const goToOrders = (status) => {
  router.push({
    path: '/order',
    query: { status }
  });
};

// 保存资料
const saveProfile = async () => {
  isSaving.value = true;
  await new Promise(resolve => setTimeout(resolve, 1000));
  user.name = editForm.name;
  user.email = editForm.email;
  
  // 保存用户信息到 localStorage
  localStorage.setItem('user', JSON.stringify(user));
  
  isSaving.value = false;
  showEditProfile.value = false;
  showToastMessage('资料保存成功');
};

// 编辑地址
const editAddress = (address) => {
  showToastMessage('编辑地址功能开发中', 'info');
};

// 删除地址
const deleteAddress = (address) => {
  addresses.value = addresses.value.filter(a => a.id !== address.id);
  showToastMessage('地址已删除');
};

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user_id');
  showToastMessage('已退出登录');
  setTimeout(() => {
    router.push('/login');
  }, 1000);
};

// 初始化
onMounted(() => {
  if (!localStorage.getItem('token')) {
    router.push('/login');
  }
  
  // 从 localStorage 获取用户信息
  const savedUser = localStorage.getItem('user');
  if (savedUser) {
    try {
      const userData = JSON.parse(savedUser);
      Object.assign(user, userData);
    } catch (error) {
      console.error('解析用户信息失败:', error);
    }
  }
});
</script>

<style scoped>
/* 页面头部 */
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 60px 0 100px;
  margin-bottom: -60px;
}

.header-content {
  text-align: center;
}

.header-title h1 {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 12px;
}

.header-subtitle {
  font-size: 16px;
  opacity: 0.9;
}

/* 用户卡片 */
.user-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  padding: 32px;
  margin-bottom: 24px;
}

.user-info-section {
  display: flex;
  align-items: center;
  gap: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f0f0f0;
}

.user-avatar {
  position: relative;
  flex-shrink: 0;
}

.avatar-img {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #f0f0f0;
}

.avatar-edit {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 32px;
  height: 32px;
  background: #667eea;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 3px solid white;
}

.avatar-edit:hover {
  background: #764ba2;
  transform: scale(1.1);
}

.user-details {
  flex: 1;
}

.user-name-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.user-name {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a2e;
}

.user-level {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.user-level.vip {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.user-phone,
.user-email {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.user-tags {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.user-tag {
  padding: 4px 10px;
  background: #e8f5e9;
  color: #2e7d32;
  border-radius: 4px;
  font-size: 12px;
}

.user-actions {
  flex-shrink: 0;
}

/* 用户统计 */
.user-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-top: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.stat-card:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

/* 区块卡片 */
.section-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
}

.title-icon {
  font-size: 20px;
}

.view-all-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  background: transparent;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  color: #666;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.view-all-btn:hover {
  border-color: #667eea;
  color: #667eea;
}

/* 订单快捷入口 */
.order-shortcuts {
  display: flex;
  justify-content: space-around;
  padding: 24px;
}

.shortcut-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  position: relative;
  padding: 16px 24px;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.shortcut-item:hover {
  background: #f8f9fa;
}

.shortcut-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, #667eea20 0%, #764ba220 100%);
  color: #667eea;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.shortcut-item:hover .shortcut-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  transform: scale(1.05);
}

.shortcut-label {
  font-size: 14px;
  color: #666;
}

.shortcut-badge {
  position: absolute;
  top: 12px;
  right: 16px;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  background: #f5576c;
  color: white;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 菜单网格 */
.menu-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.menu-list {
  padding: 8px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 24px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.menu-item:hover {
  background: #f8f9fa;
}

.menu-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.menu-text {
  flex: 1;
  font-size: 15px;
  color: #333;
}

.menu-badge {
  padding: 2px 10px;
  background: #f0f0f0;
  color: #666;
  border-radius: 10px;
  font-size: 12px;
}

.menu-status {
  padding: 2px 10px;
  border-radius: 10px;
  font-size: 12px;
}

.menu-status.safe {
  background: #e8f5e9;
  color: #2e7d32;
}

.menu-status.online {
  background: #e3f2fd;
  color: #1976d2;
}

.menu-arrow {
  color: #999;
}

/* 退出登录 */
.logout-section {
  margin-top: 32px;
  margin-bottom: 48px;
  text-align: center;
}

.btn-logout {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 14px 48px;
  background: white;
  border: 1px solid #f5576c;
  border-radius: 12px;
  color: #f5576c;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-logout:hover {
  background: #f5576c;
  color: white;
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
}

.modal-content.address-modal {
  max-width: 600px;
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

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: #f0f0f0;
  border-radius: 8px;
  font-size: 20px;
  color: #666;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: #e0e0e0;
  color: #333;
}

.modal-body {
  padding: 24px;
  max-height: 60vh;
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid #f0f0f0;
}

.modal-footer .btn {
  flex: 1;
}

/* 表单样式 */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  font-size: 15px;
  transition: all 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group input:disabled {
  background: #f5f5f5;
  color: #999;
}

.radio-group {
  display: flex;
  gap: 20px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.radio-label input {
  width: auto;
}

/* 地址卡片 */
.address-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.address-card {
  padding: 16px;
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.address-card.default {
  border-color: #667eea;
  background: #f8f9ff;
}

.address-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.address-name {
  font-weight: 600;
  color: #1a1a2e;
}

.address-phone {
  color: #666;
  font-size: 14px;
}

.default-badge {
  padding: 2px 8px;
  background: #667eea;
  color: white;
  border-radius: 4px;
  font-size: 12px;
}

.address-detail {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

.address-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.action-btn {
  padding: 6px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  background: white;
  color: #666;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn:hover {
  border-color: #667eea;
  color: #667eea;
}

.action-btn.delete:hover {
  border-color: #f5576c;
  color: #f5576c;
}

.btn-add {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
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
  to {
    transform: rotate(360deg);
  }
}

/* 提示消息 */
.toast {
  position: fixed;
  top: 100px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2000;
}

.toast-content {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  animation: toastSlideIn 0.3s ease;
}

.toast-content.success {
  background: #e8f5e9;
  color: #2e7d32;
}

.toast-content.error {
  background: #ffebee;
  color: #c62828;
}

.toast-content.info {
  background: #e3f2fd;
  color: #1565c0;
}

@keyframes toastSlideIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

/* 动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 容器 */
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .menu-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .user-stats {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .page-header {
    padding: 40px 0 80px;
  }
  
  .header-title h1 {
    font-size: 28px;
  }
  
  .user-card {
    padding: 20px;
  }
  
  .user-info-section {
    flex-direction: column;
    text-align: center;
  }
  
  .user-name-row {
    justify-content: center;
  }
  
  .user-tags {
    justify-content: center;
  }
  
  .user-stats {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
  
  .stat-card {
    padding: 16px;
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }
  
  .order-shortcuts {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .shortcut-item {
    flex: 1;
    min-width: 70px;
    padding: 12px 8px;
  }
  
  .menu-grid {
    grid-template-columns: 1fr;
  }
  
  .modal-content {
    margin: 20px;
  }
}

@media (max-width: 480px) {
  .container {
    padding: 0 16px;
  }
  
  .user-stats {
    grid-template-columns: 1fr;
  }
  
  .shortcut-icon {
    width: 48px;
    height: 48px;
  }
  
  .shortcut-icon svg {
    width: 24px;
    height: 24px;
  }
}
</style>
