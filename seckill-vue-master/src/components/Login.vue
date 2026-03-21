<template>
  <div class="auth-container">
    <!-- 背景动画 -->
    <div class="bg-animation">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
      <div class="bg-shape shape-4"></div>
    </div>

    <!-- 左侧装饰区域 -->
    <div class="auth-left">
      <div class="left-content">
        <div class="brand-logo">
          <div class="logo-icon">⚡</div>
          <h2 class="logo-text">秒杀商城</h2>
        </div>
        <p class="brand-tagline">限时抢购，每日爆款</p>
        
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
              </svg>
            </div>
            <div class="feature-text">
              <h4>闪电秒杀</h4>
              <p>限时特价，手慢无</p>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
            </div>
            <div class="feature-text">
              <h4>正品保障</h4>
              <p>品质有保证</p>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="1" y="3" width="15" height="13"/>
                <polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/>
                <circle cx="5.5" cy="18.5" r="2.5"/>
                <circle cx="18.5" cy="18.5" r="2.5"/>
              </svg>
            </div>
            <div class="feature-text">
              <h4>极速配送</h4>
              <p>次日达服务</p>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
            </div>
            <div class="feature-text">
              <h4>安全支付</h4>
              <p>多重安全保障</p>
            </div>
          </div>
        </div>

        <div class="stats-section">
          <div class="stat-item">
            <span class="stat-number">1000+</span>
            <span class="stat-label">爆款商品</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">50万+</span>
            <span class="stat-label">活跃用户</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">99.9%</span>
            <span class="stat-label">好评率</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧登录/注册区域 -->
    <div class="auth-right">
      <div class="auth-card">
        <!-- 顶部装饰条 -->
        <div class="card-top-decoration"></div>
        
        <!-- 标题区域 -->
        <div class="auth-header">
          <div class="header-icon">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
          </div>
          <h1 class="auth-title">{{ isLoginMode ? '欢迎回来' : '创建账号' }}</h1>
          <p class="auth-subtitle">{{ isLoginMode ? '登录您的账号继续购物' : '开启您的秒杀之旅' }}</p>
        </div>

        <!-- 模式切换标签 -->
        <div class="mode-tabs">
          <button 
            :class="['tab-btn', { active: isLoginMode }]" 
            @click="switchToLogin"
          >
            登录
          </button>
          <button 
            :class="['tab-btn', { active: !isLoginMode }]" 
            @click="switchToRegister"
          >
            注册
          </button>
          <div class="tab-indicator" :class="{ 'move-right': !isLoginMode }"></div>
        </div>

        <!-- 管理员切换 -->
        <div class="admin-switch">
          <label class="switch-label">
            <input type="checkbox" v-model="isAdminMode" @change="switchAdminMode">
            <span class="switch-slider"></span>
            <span class="switch-text">{{ isAdminMode ? '管理员登录' : '普通用户登录' }}</span>
          </label>
        </div>

        <!-- 表单区域 -->
        <form @submit.prevent="handleSubmit" class="auth-form">
          <!-- 注册模式 - 用户名 -->
          <transition name="slide-fade">
            <div v-if="!isLoginMode && !isAdminMode" class="form-group">
              <label class="form-label">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
                用户名
              </label>
              <div class="input-wrapper" :class="{ 'has-error': errors.username, 'has-value': formData.username }">
                <input
                  type="text"
                  v-model="formData.username"
                  class="form-input"
                  placeholder="请输入用户名"
                  required
                >
                <div class="input-focus-border"></div>
              </div>
              <transition name="error-fade">
                <p v-if="errors.username" class="error-message">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                  </svg>
                  {{ errors.username }}
                </p>
              </transition>
            </div>
          </transition>

          <!-- 管理员模式 - 管理员账号 -->
          <transition name="slide-fade">
            <div v-if="isAdminMode" class="form-group">
              <label class="form-label">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
                管理员账号
              </label>
              <div class="input-wrapper" :class="{ 'has-error': errors.username, 'has-value': formData.username }">
                <input
                  type="text"
                  v-model="formData.username"
                  class="form-input"
                  placeholder="请输入管理员账号"
                  required
                >
                <div class="input-focus-border"></div>
              </div>
              <transition name="error-fade">
                <p v-if="errors.username" class="error-message">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                  </svg>
                  {{ errors.username }}
                </p>
              </transition>
            </div>
          </transition>

          <!-- 普通用户模式 - 手机号 -->
          <transition name="slide-fade">
            <div v-if="!isAdminMode" class="form-group">
              <label class="form-label">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
                  <line x1="12" y1="18" x2="12.01" y2="18"/>
                </svg>
                手机号
              </label>
              <div class="input-wrapper" :class="{ 'has-error': errors.phone, 'has-value': formData.phone }">
                <input
                  type="tel"
                  v-model="formData.phone"
                  class="form-input"
                  placeholder="请输入手机号"
                  required
                >
                <div class="input-focus-border"></div>
              </div>
              <transition name="error-fade">
                <p v-if="errors.phone" class="error-message">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                  </svg>
                  {{ errors.phone }}
                </p>
              </transition>
            </div>
          </transition>

          <!-- 密码 -->
          <div class="form-group">
            <label class="form-label">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              密码
            </label>
            <div class="input-wrapper" :class="{ 'has-error': errors.password, 'has-value': formData.password }">
              <input
                :type="showPassword ? 'text' : 'password'"
                v-model="formData.password"
                class="form-input"
                placeholder="请输入密码"
                required
              >
              <button type="button" class="password-toggle" @click="showPassword = !showPassword">
                <svg v-if="!showPassword" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
              <div class="input-focus-border"></div>
            </div>
            <transition name="error-fade">
              <p v-if="errors.password" class="error-message">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="12" y1="8" x2="12" y2="12"/>
                  <line x1="12" y1="16" x2="12.01" y2="16"/>
                </svg>
                {{ errors.password }}
              </p>
            </transition>
          </div>

          <!-- 注册模式 - 确认密码 -->
          <transition name="slide-fade">
            <div v-if="!isLoginMode" class="form-group">
              <label class="form-label">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                确认密码
              </label>
              <div class="input-wrapper" :class="{ 'has-error': errors.confirmPassword, 'has-value': formData.confirmPassword }">
                <input
                  :type="showConfirmPassword ? 'text' : 'password'"
                  v-model="formData.confirmPassword"
                  class="form-input"
                  placeholder="请再次输入密码"
                  required
                >
                <button type="button" class="password-toggle" @click="showConfirmPassword = !showConfirmPassword">
                  <svg v-if="!showConfirmPassword" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                </button>
                <div class="input-focus-border"></div>
              </div>
              <transition name="error-fade">
                <p v-if="errors.confirmPassword" class="error-message">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                  </svg>
                  {{ errors.confirmPassword }}
                </p>
              </transition>
            </div>
          </transition>

          <!-- 记住我 & 忘记密码 -->
          <div v-if="isLoginMode" class="form-options">
            <label class="checkbox-wrapper">
              <input type="checkbox" v-model="formData.remember">
              <span class="checkbox-custom"></span>
              <span class="checkbox-label">记住我</span>
            </label>
            <a href="#" class="forgot-link">忘记密码？</a>
          </div>

          <!-- 提交按钮 -->
          <button type="submit" class="submit-btn" :disabled="isLoading">
            <span v-if="!isLoading" class="btn-content">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
                <polyline points="10 17 15 12 10 7"/>
                <line x1="15" y1="12" x2="3" y2="12"/>
              </svg>
              {{ isLoginMode ? '登录' : '注册' }}
            </span>
            <span v-else class="btn-loading">
              <span class="spinner"></span>
              处理中...
            </span>
          </button>

          <!-- 分割线 -->
          <div class="divider">
            <span class="divider-text">或使用以下方式</span>
          </div>

          <!-- 第三方登录 -->
          <div class="social-login">
            <button type="button" class="social-btn wechat">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.03-.406-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z"/>
              </svg>
              <span>微信</span>
            </button>
            <button type="button" class="social-btn qq">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12.003 2c-2.265 0-6.29 1.364-6.29 7.325v1.195S3.55 14.96 3.55 17.474c0 .665.17 1.025.281 1.025.114 0 .902-.484 1.748-2.072 0 0-.18 2.197 1.904 3.967 0 0-1.77.495-1.77 1.182 0 .686 4.078.43 6.29.43 2.212 0 6.29.256 6.29-.43 0-.687-1.77-1.182-1.77-1.182 2.085-1.77 1.904-3.967 1.904-3.967.846 1.588 1.634 2.072 1.748 2.072.111 0 .281-.36.281-1.025 0-2.514-2.164-6.954-2.164-6.954V9.325C18.292 3.364 14.268 2 12.003 2z"/>
              </svg>
              <span>QQ</span>
            </button>
            <button type="button" class="social-btn alipay">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                <path d="M20.422 18.627c-.408-.14-4.356-1.848-5.934-2.565.39-.642.75-1.34 1.056-2.085h-3.54v-1.2h4.32v-.96h-4.32V9.9h-1.92v.024H6.546v.96h4.536v1.2H6.066v.96h8.31c-.24.54-.51 1.056-.81 1.536-2.97-.912-6.078-1.392-7.272-.768-1.392.744-1.464 2.184-.6 3.072 1.344 1.392 4.296 1.248 6.648-.192.816-.504 1.608-1.176 2.352-1.992 2.064 1.008 6.072 2.808 6.072 2.808.312.12.672.264.672.672 0 .456-.456.744-.912.744H1.2c-.66 0-1.2-.54-1.2-1.2V4.8c0-.66.54-1.2 1.2-1.2h21.6c.66 0 1.2.54 1.2 1.2v12.6c0 .66-.54 1.2-1.2 1.2h-2.378v.027z"/>
              </svg>
              <span>支付宝</span>
            </button>
          </div>
        </form>

        <!-- 底部链接 -->
        <div class="auth-footer">
          <p v-if="isLoginMode">
            还没有账号？
            <a href="#" @click.prevent="switchToRegister">立即注册</a>
          </p>
          <p v-else>
            已有账号？
            <a href="#" @click.prevent="switchToLogin">立即登录</a>
          </p>
        </div>
      </div>
    </div>

    <!-- 提示消息 -->
    <transition name="toast-fade">
      <div v-if="toast.show" class="toast-container" :class="toast.type">
        <div class="toast-icon">
          <svg v-if="toast.type === 'success'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
          <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
        </div>
        <div class="toast-content">
          <p class="toast-title">{{ toast.type === 'success' ? '成功' : '错误' }}</p>
          <p class="toast-message">{{ toast.message }}</p>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import axios from "axios";
import { useRouter } from 'vue-router';

const router = useRouter();

// 状态管理
const isLoginMode = ref(true);
const isAdminMode = ref(false);
const isLoading = ref(false);
const showPassword = ref(false);
const showConfirmPassword = ref(false);

const formData = reactive({
  username: '',
  phone: '',
  password: '',
  confirmPassword: '',
  remember: false
});

const errors = reactive({});

// 提示消息状态
const toast = reactive({
  show: false,
  message: '',
  type: 'success',
  duration: 3000
});

// 显示提示
const showToast = (message, type = 'success', duration = 3000) => {
  toast.message = message;
  toast.type = type;
  toast.show = true;
  toast.duration = duration;

  clearTimeout(window.toastTimer);
  window.toastTimer = setTimeout(() => {
    toast.show = false;
  }, duration);
};

// 切换到登录模式
const switchToLogin = () => {
  if (!isLoginMode.value) {
    isLoginMode.value = true;
    resetForm();
  }
};

// 切换到注册模式
const switchToRegister = () => {
  if (isLoginMode.value) {
    isLoginMode.value = false;
    resetForm();
  }
};

// 切换管理员模式
const switchAdminMode = () => {
  resetForm();
  if (isAdminMode.value) {
    showToast('已切换到管理员登录模式', 'info');
  } else {
    showToast('已切换到普通用户登录模式', 'info');
  }
};

// 重置表单
const resetForm = () => {
  Object.keys(formData).forEach(key => {
    formData[key] = key === 'remember' ? false : '';
  });
  Object.keys(errors).forEach(key => {
    errors[key] = '';
  });
};

// 表单验证
const validateForm = () => {
  let isValid = true;
  const newErrors = {};

  if (!isLoginMode.value) {
    if (!formData.username) {
      newErrors.username = '请输入用户名';
      isValid = false;
    } else if (formData.username.length < 2 || formData.username.length > 20) {
      newErrors.username = '用户名长度必须在2-20个字符之间';
      isValid = false;
    }
  }

  if (isAdminMode.value) {
    if (!formData.username) {
      newErrors.username = '请输入管理员账号';
      isValid = false;
    } else if (formData.username.length < 2 || formData.username.length > 20) {
      newErrors.username = '管理员账号长度必须在2-20个字符之间';
      isValid = false;
    }
  } else {
    if (!formData.phone) {
      newErrors.phone = '请输入手机号';
      isValid = false;
    } else if (!/^1[3-9]\d{9}$/.test(formData.phone)) {
      newErrors.phone = '请输入有效的手机号';
      isValid = false;
    }
  }

  if (!formData.password) {
    newErrors.password = '请输入密码';
    isValid = false;
  } else if (formData.password.length < 6) {
    newErrors.password = '密码长度至少为6位';
    isValid = false;
  }

  if (!isLoginMode.value && formData.password !== formData.confirmPassword) {
    newErrors.confirmPassword = '两次输入的密码不一致';
    isValid = false;
  }

  Object.assign(errors, newErrors);
  return isValid;
};

// 处理表单提交
const handleSubmit = async () => {
  if (!validateForm()) {
    showToast('请完善表单信息后提交', 'error');
    return;
  }

  isLoading.value = true;

  try {
    if (isAdminMode.value) {
      if (formData.username === 'admin' && formData.password === 'admin123') {
        showToast('管理员登录成功！即将跳转后台', 'success');
        localStorage.setItem('admin_token', 'admin_token_' + Date.now());
        localStorage.setItem('admin_user', JSON.stringify({
          username: 'admin',
          role: 'admin',
          loginTime: new Date().toISOString()
        }));
        setTimeout(() => {
          router.push('/admin');
        }, 1500);
      } else {
        showToast('管理员账号或密码错误', 'error');
      }
    } else if (isLoginMode.value) {
      const response = await axios.post('/api/user/login', {
        phone: formData.phone,
        password: formData.password
      });
      
      showToast('登录成功！即将跳转首页', 'success');
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user_id', response.data.user_id);
      
      // 保存用户信息到 localStorage
      const userInfo = {
        name: formData.phone,
        phone: formData.phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2'),
        email: `${formData.phone}@example.com`,
        level: '普通会员',
        levelClass: '',
        avatar: `https://placehold.co/120x120/667eea/ffffff?text=${formData.phone.slice(-2)}`,
        ordersCount: 0,
        pendingOrders: 0,
        paidOrders: 0,
        shippedOrders: 0,
        favoritesCount: 0,
        couponsCount: 0,
        points: 0,
        addressCount: 0,
        reminders: 0
      };
      localStorage.setItem('user', JSON.stringify(userInfo));
      
      setTimeout(() => {
        router.push('/main');
      }, 1500);
    } else {
      const response = await axios.post('/api/user/register', {
        username: formData.username,
        phone: formData.phone,
        password: formData.password
      });
      
      showToast('注册成功！请登录', 'success');
      setTimeout(() => {
        switchToLogin();
      }, 1500);
    }
  } catch (error) {
    const errorMsg = error.response?.data?.error || '网络异常，请稍后重试';
    showToast(errorMsg, 'error');
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
/* 容器样式 */
.auth-container {
  min-height: 100vh;
  display: flex;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

/* 背景动画 */
.bg-animation {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 20s infinite;
}

.shape-1 {
  width: 400px;
  height: 400px;
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.shape-2 {
  width: 300px;
  height: 300px;
  top: 50%;
  right: -150px;
  animation-delay: -5s;
}

.shape-3 {
  width: 200px;
  height: 200px;
  bottom: -100px;
  left: 30%;
  animation-delay: -10s;
}

.shape-4 {
  width: 150px;
  height: 150px;
  top: 20%;
  left: 40%;
  animation-delay: -15s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-30px) rotate(180deg);
  }
}

/* 左侧装饰区域 */
.auth-left {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  position: relative;
  z-index: 1;
}

.left-content {
  color: white;
  max-width: 500px;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.logo-icon {
  width: 60px;
  height: 60px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  backdrop-filter: blur(10px);
}

.logo-text {
  font-size: 36px;
  font-weight: 700;
  margin: 0;
}

.brand-tagline {
  font-size: 18px;
  opacity: 0.9;
  margin-bottom: 48px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 48px;
}

.feature-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.feature-card:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.feature-icon-wrapper {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.feature-text h4 {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 4px 0;
}

.feature-text p {
  font-size: 13px;
  opacity: 0.8;
  margin: 0;
}

.stats-section {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-item {
  text-align: center;
  flex: 1;
}

.stat-number {
  display: block;
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  opacity: 0.8;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
}

/* 右侧登录区域 */
.auth-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  position: relative;
  z-index: 1;
}

.auth-card {
  background: white;
  border-radius: 24px;
  padding: 48px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  position: relative;
  overflow: hidden;
}

.card-top-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2, #f093fb);
}

/* 标题区域 */
.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.header-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
}

.auth-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 8px 0;
}

.auth-subtitle {
  font-size: 15px;
  color: #666;
  margin: 0;
}

/* 模式切换标签 */
.mode-tabs {
  display: flex;
  position: relative;
  background: #f5f5f5;
  border-radius: 12px;
  padding: 4px;
  margin-bottom: 32px;
}

.tab-btn {
  flex: 1;
  padding: 12px;
  border: none;
  background: transparent;
  font-size: 15px;
  font-weight: 500;
  color: #666;
  cursor: pointer;
  transition: color 0.3s ease;
  position: relative;
  z-index: 1;
}

.tab-btn.active {
  color: #667eea;
}

.tab-indicator {
  position: absolute;
  top: 4px;
  left: 4px;
  width: calc(50% - 4px);
  height: calc(100% - 8px);
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: transform 0.3s ease;
}

.tab-indicator.move-right {
  transform: translateX(100%);
}

/* 管理员切换 */
.admin-switch {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.switch-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
}

.switch-label input[type="checkbox"] {
  position: relative;
  width: 50px;
  height: 26px;
  appearance: none;
  background: #e0e0e0;
  border-radius: 13px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.switch-label input[type="checkbox"]:checked {
  background: #667eea;
}

.switch-slider {
  position: absolute;
  top: 3px;
  left: 3px;
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  transition: transform 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.switch-label input[type="checkbox"]:checked + .switch-slider {
  transform: translateX(24px);
}

.switch-text {
  font-size: 14px;
  font-weight: 500;
  color: #666;
}

.switch-label input[type="checkbox"]:checked ~ .switch-text {
  color: #667eea;
}

/* 表单样式 */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.form-label svg {
  color: #667eea;
}

.input-wrapper {
  position: relative;
}

.form-input {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid #e0e0e0;
  border-radius: 12px;
  font-size: 15px;
  transition: all 0.3s ease;
  background: #fafafa;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.input-wrapper.has-error .form-input {
  border-color: #f5576c;
  background: #fff5f5;
}

.input-wrapper.has-value .form-input {
  background: white;
}

.input-focus-border {
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  transition: all 0.3s ease;
  transform: translateX(-50%);
}

.form-input:focus ~ .input-focus-border {
  width: 100%;
}

.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.3s ease;
}

.password-toggle:hover {
  color: #667eea;
}

/* 错误提示 */
.error-message {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #f5576c;
  font-size: 13px;
  margin: 0;
  padding: 8px 12px;
  background: #fff5f5;
  border-radius: 8px;
}

/* 表单选项 */
.form-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.checkbox-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.checkbox-wrapper input {
  display: none;
}

.checkbox-custom {
  width: 20px;
  height: 20px;
  border: 2px solid #d0d0d0;
  border-radius: 6px;
  position: relative;
  transition: all 0.3s ease;
}

.checkbox-wrapper input:checked + .checkbox-custom {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-color: #667eea;
}

.checkbox-wrapper input:checked + .checkbox-custom::after {
  content: '';
  position: absolute;
  left: 6px;
  top: 2px;
  width: 5px;
  height: 10px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.checkbox-label {
  font-size: 14px;
  color: #666;
}

.forgot-link {
  font-size: 14px;
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
}

.forgot-link:hover {
  color: #764ba2;
}

/* 提交按钮 */
.submit-btn {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
  position: relative;
  overflow: hidden;
}

.submit-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.submit-btn:hover:not(:disabled)::before {
  left: 100%;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.4);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.spinner {
  width: 20px;
  height: 20px;
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

/* 分割线 */
.divider {
  display: flex;
  align-items: center;
  margin: 8px 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e0e0e0;
}

.divider-text {
  padding: 0 16px;
  font-size: 13px;
  color: #999;
}

/* 第三方登录 */
.social-login {
  display: flex;
  gap: 12px;
}

.social-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 12px;
  border: 2px solid #e0e0e0;
  border-radius: 12px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
}

.social-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.social-btn span {
  font-size: 12px;
  font-weight: 500;
}

.social-btn.wechat {
  color: #07C160;
  border-color: #07C160;
}

.social-btn.wechat:hover {
  background: #f0fff4;
}

.social-btn.qq {
  color: #12B7F5;
  border-color: #12B7F5;
}

.social-btn.qq:hover {
  background: #f0faff;
}

.social-btn.alipay {
  color: #1677FF;
  border-color: #1677FF;
}

.social-btn.alipay:hover {
  background: #f0f7ff;
}

/* 底部链接 */
.auth-footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

.auth-footer p {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.auth-footer a {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.3s ease;
}

.auth-footer a:hover {
  color: #764ba2;
}

/* 提示消息 */
.toast-container {
  position: fixed;
  top: 24px;
  right: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  z-index: 1000;
  max-width: 400px;
}

.toast-container.success {
  border-left: 4px solid #10b981;
}

.toast-container.success .toast-icon {
  color: #10b981;
}

.toast-container.error {
  border-left: 4px solid #f5576c;
}

.toast-container.error .toast-icon {
  color: #f5576c;
}

.toast-icon {
  flex-shrink: 0;
}

.toast-content {
  flex: 1;
}

.toast-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.toast-message {
  font-size: 13px;
  color: #666;
  margin: 0;
}

/* 动画 */
.slide-fade-enter-active {
  transition: all 0.3s ease;
}

.slide-fade-leave-active {
  transition: all 0.2s ease;
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.error-fade-enter-active,
.error-fade-leave-active {
  transition: all 0.2s ease;
}

.error-fade-enter-from,
.error-fade-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.3s ease;
}

.toast-fade-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-fade-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .auth-container {
    flex-direction: column;
  }

  .auth-left {
    width: 100%;
    padding: 40px;
  }

  .left-content {
    max-width: 100%;
  }

  .brand-logo {
    justify-content: center;
  }

  .brand-tagline {
    text-align: center;
  }

  .features-grid {
    grid-template-columns: repeat(4, 1fr);
  }

  .feature-text h4 {
    font-size: 14px;
  }

  .feature-text p {
    display: none;
  }

  .stats-section {
    justify-content: center;
  }

  .auth-right {
    width: 100%;
    padding: 24px;
  }

  .auth-card {
    max-width: 100%;
    padding: 32px;
  }
}

@media (max-width: 768px) {
  .auth-left {
    padding: 32px 24px;
  }

  .logo-text {
    font-size: 28px;
  }

  .brand-tagline {
    font-size: 16px;
    margin-bottom: 32px;
  }

  .features-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    margin-bottom: 32px;
  }

  .feature-card {
    padding: 16px;
  }

  .feature-icon-wrapper {
    width: 40px;
    height: 40px;
  }

  .feature-icon-wrapper svg {
    width: 22px;
    height: 22px;
  }

  .feature-text p {
    display: block;
  }

  .stats-section {
    padding: 16px;
    gap: 16px;
  }

  .stat-number {
    font-size: 22px;
  }

  .stat-label {
    font-size: 12px;
  }

  .auth-card {
    padding: 24px;
  }

  .header-icon {
    width: 60px;
    height: 60px;
  }

  .auth-title {
    font-size: 24px;
  }

  .social-login {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .auth-left {
    padding: 24px 16px;
  }

  .brand-logo {
    flex-direction: column;
    gap: 12px;
  }

  .logo-icon {
    width: 50px;
    height: 50px;
    font-size: 26px;
  }

  .logo-text {
    font-size: 24px;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .feature-card {
    padding: 12px;
  }

  .stats-section {
    flex-wrap: wrap;
    gap: 12px;
  }

  .stat-item {
    flex: none;
    width: calc(50% - 6px);
  }

  .stat-divider {
    display: none;
  }

  .auth-card {
    padding: 20px;
    border-radius: 20px;
  }

  .form-input {
    padding: 12px 14px;
    font-size: 14px;
  }

  .submit-btn {
    padding: 14px;
    font-size: 15px;
  }

  .toast-container {
    left: 16px;
    right: 16px;
    top: 16px;
    max-width: none;
  }
}
</style>
