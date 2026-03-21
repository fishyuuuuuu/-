<template>
  <div class="admin-page">
    <!-- 顶部导航 -->
    <div class="admin-header">
      <div class="header-content">
        <div class="logo-section">
          <div class="logo-icon">🔐</div>
          <div class="logo-text">
            <h1>管理员后台</h1>
            <p>Admin Dashboard</p>
          </div>
        </div>
        <div class="header-actions">
          <div class="admin-info">
            <span class="admin-name">管理员</span>
            <span class="login-time">登录时间: {{ formatLoginTime }}</span>
          </div>
          <button class="logout-btn" @click="handleLogout">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            退出登录
          </button>
        </div>
      </div>
    </div>

    <div class="admin-content">
      <!-- 侧边栏 -->
      <div class="sidebar">
        <nav class="sidebar-nav">
          <div class="nav-section">
            <h3 class="nav-section-title">主要功能</h3>
            <ul class="nav-list">
              <li 
                v-for="item in mainMenuItems" 
                :key="item.id"
                :class="['nav-item', { active: activeMenu === item.id }]"
                @click="activeMenu = item.id"
              >
                <div class="nav-icon" :style="{ background: item.color }">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path :d="item.icon"/>
                  </svg>
                </div>
                <span class="nav-label">{{ item.label }}</span>
                <span v-if="item.badge" class="nav-badge">{{ item.badge }}</span>
              </li>
            </ul>
          </div>

          <div class="nav-section">
            <h3 class="nav-section-title">系统管理</h3>
            <ul class="nav-list">
              <li 
                v-for="item in systemMenuItems" 
                :key="item.id"
                :class="['nav-item', { active: activeMenu === item.id }]"
                @click="activeMenu = item.id"
              >
                <div class="nav-icon" :style="{ background: item.color }">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path :d="item.icon"/>
                  </svg>
                </div>
                <span class="nav-label">{{ item.label }}</span>
              </li>
            </ul>
          </div>
        </nav>
      </div>

      <!-- 主内容区域 -->
      <div class="main-content">
        <!-- 仪表板 -->
        <div v-if="activeMenu === 'dashboard'" class="content-section">
          <h2 class="section-title">数据概览</h2>
          
          <!-- 统计卡片 -->
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
              </div>
              <div class="stat-info">
                <p class="stat-value">1,234</p>
                <p class="stat-label">总用户数</p>
                <p class="stat-trend positive">+12.5%</p>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
                  <line x1="3" y1="6" x2="21" y2="6"/>
                  <path d="M16 10a4 4 0 0 1-8 0"/>
                </svg>
              </div>
              <div class="stat-info">
                <p class="stat-value">5,678</p>
                <p class="stat-label">总订单数</p>
                <p class="stat-trend positive">+8.2%</p>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 0 0-7H17z"/>
                  <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 0 0-7H17z"/>
                  <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 0 0-7H17z"/>
                </svg>
              </div>
              <div class="stat-info">
                <p class="stat-value">¥128,456</p>
                <p class="stat-label">总销售额</p>
                <p class="stat-trend positive">+23.1%</p>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                </svg>
              </div>
              <div class="stat-info">
                <p class="stat-value">89</p>
                <p class="stat-label">商品总数</p>
                <p class="stat-trend positive">+5.3%</p>
              </div>
            </div>
          </div>

          <!-- 图表区域 -->
          <div class="charts-grid">
            <div class="chart-card">
              <h3 class="chart-title">销售趋势</h3>
              <div class="chart-placeholder">
                <div class="chart-bars">
                  <div class="bar" v-for="(height, index) in salesData" :key="index" :style="{ height: height + '%' }"></div>
                </div>
              </div>
            </div>

            <div class="chart-card">
              <h3 class="chart-title">订单状态分布</h3>
              <div class="chart-placeholder">
                <div class="pie-chart">
                  <div class="pie-segment" style="background: #667eea; --percentage: 45;"></div>
                  <div class="pie-segment" style="background: #f093fb; --percentage: 30;"></div>
                  <div class="pie-segment" style="background: #4facfe; --percentage: 15;"></div>
                  <div class="pie-segment" style="background: #43e97b; --percentage: 10;"></div>
                </div>
                <div class="legend">
                  <div class="legend-item">
                    <span class="legend-color" style="background: #667eea;"></span>
                    <span class="legend-label">待支付 (45%)</span>
                  </div>
                  <div class="legend-item">
                    <span class="legend-color" style="background: #f093fb;"></span>
                    <span class="legend-label">已支付 (30%)</span>
                  </div>
                  <div class="legend-item">
                    <span class="legend-color" style="background: #4facfe;"></span>
                    <span class="legend-label">已发货 (15%)</span>
                  </div>
                  <div class="legend-item">
                    <span class="legend-color" style="background: #43e97b;"></span>
                    <span class="legend-label">已完成 (10%)</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 用户管理 -->
        <div v-else-if="activeMenu === 'users'" class="content-section">
          <div class="section-header">
            <h2 class="section-title">用户管理</h2>
            <button class="btn btn-primary" @click="fetchUsers">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21.21 15.89A10 10 0 1 1 8 2.83"/>
                <path d="M22 12A10 10 0 0 0 12 2v10z"/>
              </svg>
              刷新
            </button>
          </div>

          <!-- 搜索框 -->
          <div class="search-bar">
            <div class="search-group">
              <label for="username-search">用户名</label>
              <input 
                type="text" 
                id="username-search" 
                placeholder="输入用户名搜索" 
                v-model="usersFilter.username"
                @keyup.enter="searchUsers"
              >
            </div>
            <div class="search-group">
              <label for="phone-search">手机号</label>
              <input 
                type="text" 
                id="phone-search" 
                placeholder="输入手机号搜索" 
                v-model="usersFilter.phone"
                @keyup.enter="searchUsers"
              >
            </div>
            <button class="btn btn-primary" @click="searchUsers">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              搜索
            </button>
            <button class="btn btn-secondary" @click="resetSearch">
              重置
            </button>
            <button 
              class="btn" 
              :class="pollingEnabled ? 'btn-success' : 'btn-secondary'"
              @click="togglePolling"
              :title="pollingEnabled ? '停止实时同步' : '开始实时同步'"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path v-if="pollingEnabled" d="M18 12H6"/>
                <path v-else d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8" fill="none"/>
                <circle cx="12" cy="12" r="10" fill="none"/>
              </svg>
              {{ pollingEnabled ? '实时同步中...' : '实时同步' }}
            </button>
          </div>

          <!-- 加载状态 -->
          <div v-if="usersLoading" class="loading-overlay">
            <div class="spinner"></div>
            <p>正在加载用户数据...</p>
          </div>

          <!-- 错误信息 -->
          <div v-if="usersError && !usersLoading" class="error-message">
            <p>{{ usersError }}</p>
            <button class="btn btn-secondary" @click="fetchUsers">重试</button>
          </div>

          <div v-if="!usersLoading" class="table-container">
            <table class="data-table" v-if="users.length > 0">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>用户名</th>
                  <th>手机号</th>
                  <th>注册时间</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="user in users" :key="user.id">
                  <td>{{ user.id }}</td>
                  <td>{{ user.username }}</td>
                  <td>{{ user.phone }}</td>
                  <td>{{ user.registerTime }}</td>
                  <td>
                    <span :class="['status-badge', user.status]">{{ user.statusText }}</span>
                  </td>
                  <td>
                    <div class="action-buttons">
                      <button class="btn-icon" title="编辑">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                      </button>
                      <button class="btn-icon danger" title="删除" @click="deleteUser(user.id)">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <polyline points="3 6 5 6 21 6"/>
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>

            <!-- 空状态 -->
            <div v-if="users.length === 0 && !usersError && !usersLoading" class="empty-state">
              <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              <p>暂无用户数据</p>
            </div>
          </div>

          <!-- 分页控件 -->
          <div v-if="usersPagination.totalPages > 1 && !usersLoading && users.length > 0" class="pagination">
            <button 
              class="btn-pagination" 
              :disabled="usersPagination.page === 1" 
              @click="changePage(usersPagination.page - 1)"
            >
              上一页
            </button>
            
            <span class="pagination-info">
              第 {{ usersPagination.page }} 页 / 共 {{ usersPagination.totalPages }} 页
              （总计 {{ usersPagination.total }} 条记录）
            </span>
            
            <button 
              class="btn-pagination" 
              :disabled="usersPagination.page === usersPagination.totalPages" 
              @click="changePage(usersPagination.page + 1)"
            >
              下一页
            </button>
          </div>
        </div>

        <!-- 商品管理 -->
        <div v-else-if="activeMenu === 'products'" class="content-section">
          <div class="section-header">
            <h2 class="section-title">商品管理</h2>
            <button class="btn btn-primary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              添加商品
            </button>
          </div>

          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>商品名称</th>
                  <th>价格</th>
                  <th>库存</th>
                  <th>秒杀</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="product in products" :key="product.id">
                  <td>{{ product.id }}</td>
                  <td>{{ product.name }}</td>
                  <td>¥{{ product.price.toLocaleString() }}</td>
                  <td>{{ product.stock }}</td>
                  <td>
                    <span :class="['status-badge', product.isSeckill ? 'seckill' : 'normal']">
                      {{ product.isSeckill ? '是' : '否' }}
                    </span>
                  </td>
                  <td>
                    <div class="action-buttons">
                      <button class="btn-icon" title="编辑">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                      </button>
                      <button class="btn-icon danger" title="删除">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <polyline points="3 6 5 6 21 6"/>
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 订单管理 -->
        <div v-else-if="activeMenu === 'orders'" class="content-section">
          <div class="section-header">
            <h2 class="section-title">订单管理</h2>
            <div class="filter-group">
              <select class="filter-select">
                <option value="">全部状态</option>
                <option value="pending">待支付</option>
                <option value="paid">已支付</option>
                <option value="shipped">已发货</option>
                <option value="completed">已完成</option>
              </select>
            </div>
          </div>

          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>订单号</th>
                  <th>用户</th>
                  <th>商品</th>
                  <th>金额</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="order in orders" :key="order.id">
                  <td>{{ order.orderNo }}</td>
                  <td>{{ order.user }}</td>
                  <td>{{ order.product }}</td>
                  <td>¥{{ order.amount.toLocaleString() }}</td>
                  <td>
                    <span :class="['status-badge', order.status]">{{ order.statusText }}</span>
                  </td>
                  <td>
                    <div class="action-buttons">
                      <button class="btn-icon" title="查看详情">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                          <circle cx="12" cy="12" r="3"/>
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 系统设置 -->
        <div v-else-if="activeMenu === 'settings'" class="content-section">
          <h2 class="section-title">系统设置</h2>
          
          <div class="settings-container">
            <div class="setting-card">
              <h3 class="setting-title">基本信息</h3>
              <div class="setting-form">
                <div class="form-row">
                  <label class="form-label">系统名称</label>
                  <input type="text" class="form-input" value="秒杀商城">
                </div>
                <div class="form-row">
                  <label class="form-label">管理员邮箱</label>
                  <input type="email" class="form-input" value="admin@example.com">
                </div>
                <div class="form-row">
                  <label class="form-label">联系电话</label>
                  <input type="tel" class="form-input" value="400-123-4567">
                </div>
              </div>
            </div>

            <div class="setting-card">
              <h3 class="setting-title">秒杀设置</h3>
              <div class="setting-form">
                <div class="form-row">
                  <label class="form-label">秒杀开始时间</label>
                  <input type="time" class="form-input" value="10:00">
                </div>
                <div class="form-row">
                  <label class="form-label">秒杀结束时间</label>
                  <input type="time" class="form-input" value="22:00">
                </div>
                <div class="form-row">
                  <label class="form-label">每人限购数量</label>
                  <input type="number" class="form-input" value="1" min="1">
                </div>
              </div>
            </div>

            <div class="setting-card">
              <h3 class="setting-title">安全设置</h3>
              <div class="setting-form">
                <div class="form-row checkbox-row">
                  <label class="checkbox-label">
                    <input type="checkbox" checked>
                    <span>启用验证码</span>
                  </label>
                </div>
                <div class="form-row checkbox-row">
                  <label class="checkbox-label">
                    <input type="checkbox" checked>
                    <span>启用IP限制</span>
                  </label>
                </div>
                <div class="form-row checkbox-row">
                  <label class="checkbox-label">
                    <input type="checkbox" checked>
                    <span>启用登录日志</span>
                  </label>
                </div>
              </div>
            </div>

            <div class="setting-actions">
              <button class="btn btn-primary">保存设置</button>
              <button class="btn btn-outline">重置默认</button>
            </div>
          </div>
        </div>

        <!-- 操作日志 -->
        <div v-else-if="activeMenu === 'logs'" class="content-section">
          <h2 class="section-title">操作日志</h2>
          
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>时间</th>
                  <th>操作人</th>
                  <th>操作类型</th>
                  <th>详情</th>
                  <th>IP地址</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in logs" :key="log.id">
                  <td>{{ log.time }}</td>
                  <td>{{ log.operator }}</td>
                  <td>{{ log.type }}</td>
                  <td>{{ log.detail }}</td>
                  <td>{{ log.ip }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const activeMenu = ref('dashboard');

const mainMenuItems = ref([
  { id: 'dashboard', label: '数据概览', icon: 'M3 3h18v18H3zM3 9h18M9 21V9', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { id: 'users', label: '用户管理', icon: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M12 11a4 4 0 1 0 0-8 0 4 4 0 0 0 0 8zm0 0v9m0-9a3 3 0 0 1 0-6 3 3 0 0 1 0 6zm0 0H9m3 0h3', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { id: 'products', label: '商品管理', icon: 'M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z M3 6h18 M16 10a4 4 0 1 1-8 0', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { id: 'orders', label: '订单管理', icon: 'M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2M9 5a2 2 0 0 0 2 2h2a2 2 0 0 0 2-2M9 5a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' }
]);

const systemMenuItems = ref([
  { id: 'settings', label: '系统设置', icon: 'M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm0 0v9m0-9a3 3 0 0 1 0-6 3 3 0 0 1 0 6zm0 0H9m3 0h3 M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 0 0-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 0 0-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 0 0-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 0 0-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 0 0 1.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 1 1-6 0 3 3 0 0 1 1 6 0', color: 'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)' },
  { id: 'logs', label: '操作日志', icon: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z M14 2v6 M10 18H6 M14 18h-4', color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)' }
]);

const salesData = ref([65, 78, 90, 81, 56, 55, 40, 70, 85, 95, 88, 72]);

// 用户管理相关状态
const users = ref([]);
const usersLoading = ref(false);
const usersError = ref('');
const usersPagination = ref({
  page: 1,
  pageSize: 10,
  total: 0,
  totalPages: 0
});
const usersFilter = ref({
  username: '',
  phone: ''
});

// 商品数据（暂时保持模拟数据）
const products = ref([
  { id: 1, name: 'iPhone 15 Pro', price: 7999, stock: 15, isSeckill: true },
  { id: 2, name: 'MacBook Pro 14', price: 12999, stock: 8, isSeckill: true },
  { id: 3, name: 'AirPods Pro 2', price: 1899, stock: 45, isSeckill: true },
  { id: 4, name: 'iPad Pro 12.9', price: 8999, stock: 0, isSeckill: false },
  { id: 5, name: 'Apple Watch Ultra 2', price: 6299, stock: 23, isSeckill: false }
]);

// 订单数据（暂时保持模拟数据）
const orders = ref([
  { id: 1, orderNo: '202403200001', user: '张三', product: 'iPhone 15 Pro', amount: 7999, status: 'pending', statusText: '待支付' },
  { id: 2, orderNo: '202403200002', user: '李四', product: 'MacBook Pro 14', amount: 12999, status: 'paid', statusText: '已支付' },
  { id: 3, orderNo: '202403200003', user: '王五', product: 'AirPods Pro 2', amount: 1899, status: 'shipped', statusText: '已发货' },
  { id: 4, orderNo: '202403200004', user: '赵六', product: 'iPad Pro 12.9', amount: 8999, status: 'completed', statusText: '已完成' }
]);

// 日志数据（暂时保持模拟数据）
const logs = ref([
  { id: 1, time: '2024-03-20 10:30:25', operator: 'admin', type: '登录', detail: '管理员登录系统', ip: '192.168.1.100' },
  { id: 2, time: '2024-03-20 11:15:30', operator: 'admin', type: '商品管理', detail: '添加商品 iPhone 15 Pro', ip: '192.168.1.100' },
  { id: 3, time: '2024-03-20 12:45:18', operator: 'admin', type: '用户管理', detail: '禁用用户 王五', ip: '192.168.1.100' },
  { id: 4, time: '2024-03-20 14:20:45', operator: 'admin', type: '订单管理', detail: '处理订单 202403200002', ip: '192.168.1.100' }
]);

const formatLoginTime = computed(() => {
  const adminUser = JSON.parse(localStorage.getItem('admin_user') || '{}');
  if (adminUser.loginTime) {
    const date = new Date(adminUser.loginTime);
    return date.toLocaleString('zh-CN');
  }
  return '未知';
});

const handleLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    localStorage.removeItem('admin_token');
    localStorage.removeItem('admin_user');
    router.push('/login');
  }
};

// 轮询相关状态
const usersPollingInterval = ref(null);
const pollingEnabled = ref(false);
const POLLING_INTERVAL = 10000; // 10秒

// 开始轮询用户数据
const startUsersPolling = () => {
  if (usersPollingInterval.value) {
    clearInterval(usersPollingInterval.value);
  }
  
  pollingEnabled.value = true;
  usersPollingInterval.value = setInterval(() => {
    if (activeMenu.value === 'users' && !usersLoading.value) {
      fetchUsers();
    }
  }, POLLING_INTERVAL);
};

// 停止轮询用户数据
const stopUsersPolling = () => {
  pollingEnabled.value = false;
  if (usersPollingInterval.value) {
    clearInterval(usersPollingInterval.value);
    usersPollingInterval.value = null;
  }
};

// 切换轮询状态
const togglePolling = () => {
  if (pollingEnabled.value) {
    stopUsersPolling();
  } else {
    startUsersPolling();
  }
};

// 获取用户列表
const fetchUsers = async () => {
  if (usersLoading.value) return;
  
  usersLoading.value = true;
  usersError.value = '';
  
  try {
    // 构建查询参数
    const params = {
      page: usersPagination.value.page,
      pageSize: usersPagination.value.pageSize
    };
    
    if (usersFilter.value.username) {
      params.username = usersFilter.value.username;
    }
    
    if (usersFilter.value.phone) {
      params.phone = usersFilter.value.phone;
    }
    
    // 调用API
    const response = await axios.get('/api/user/list', {
      params
    });
    
    // 处理响应数据
    users.value = response.data.users.map(user => ({
      id: user.id,
      username: user.username,
      phone: user.phone,
      registerTime: new Date(user.created_at).toLocaleDateString('zh-CN'),
      status: 'active', // 这里可以根据实际数据设置状态
      statusText: '正常' // 这里可以根据实际数据设置状态文本
    }));
    
    // 更新分页信息
    usersPagination.value.total = response.data.total;
    usersPagination.value.totalPages = response.data.totalPages || 
      Math.ceil(response.data.total / usersPagination.value.pageSize);
    
  } catch (error) {
    console.error('获取用户列表失败:', error);
    usersError.value = error.response?.data?.error || error.message || '获取用户列表失败';
    
    // 如果是未授权错误，跳转到登录页面
    if (error.response?.status === 401) {
      localStorage.removeItem('admin_token');
      localStorage.removeItem('admin_user');
      router.push('/login');
    }
  } finally {
    usersLoading.value = false;
  }
};

// 搜索用户
const searchUsers = () => {
  usersPagination.value.page = 1;
  fetchUsers();
};

// 重置搜索
const resetSearch = () => {
  usersFilter.value.username = '';
  usersFilter.value.phone = '';
  usersPagination.value.page = 1;
  fetchUsers();
};

// 切换页面
const changePage = (page) => {
  if (page < 1 || page > usersPagination.value.totalPages) return;
  usersPagination.value.page = page;
  fetchUsers();
};

// 当切换到用户管理页面时自动加载数据
watch(() => activeMenu.value, (newValue) => {
  if (newValue === 'users' && users.value.length === 0) {
    fetchUsers();
  }
});

// 删除用户
const deleteUser = async (id) => {
  if (!confirm('确定要删除这个用户吗？')) return;
  
  try {
    await axios.delete(`/api/user/delete/${id}`);
    
    // 删除成功后刷新用户列表
    fetchUsers();
    
    // 显示成功消息（可以添加toast组件）
    alert('用户删除成功');
  } catch (error) {
    console.error('删除用户失败:', error);
    alert(error.response?.data?.error || error.message || '删除用户失败');
  }
};

onMounted(() => {
  // 不再检查token，直接加载数据
  if (activeMenu.value === 'users') {
    fetchUsers();
  }
});

// 组件卸载时清理轮询
onUnmounted(() => {
  stopUsersPolling();
});
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  background: #f5f7fa;
}

/* 顶部导航 */
.admin-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0 32px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo-icon {
  font-size: 32px;
}

.logo-text h1 {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
}

.logo-text p {
  font-size: 12px;
  margin: 0;
  opacity: 0.8;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.admin-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.admin-name {
  font-weight: 600;
  font-size: 14px;
}

.login-time {
  font-size: 12px;
  opacity: 0.8;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

/* 主内容区域 */
.admin-content {
  display: flex;
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px 32px;
  gap: 24px;
}

/* 侧边栏 */
.sidebar {
  width: 240px;
  flex-shrink: 0;
}

.sidebar-nav {
  position: sticky;
  top: 88px;
}

.nav-section {
  margin-bottom: 24px;
}

.nav-section-title {
  font-size: 12px;
  font-weight: 600;
  color: #9ca3af;
  text-transform: uppercase;
  margin-bottom: 12px;
  padding: 0 12px;
}

.nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 4px;
  position: relative;
}

.nav-item:hover {
  background: rgba(102, 126, 234, 0.1);
}

.nav-item.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.nav-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.nav-label {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

.nav-badge {
  background: #ef4444;
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 600;
}

/* 主内容区域 */
.main-content {
  flex: 1;
  min-width: 0;
}

.content-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.stat-trend {
  font-size: 12px;
  font-weight: 600;
}

.stat-trend.positive {
  color: #10b981;
}

.stat-trend.negative {
  color: #ef4444;
}

/* 图表区域 */
.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.chart-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.chart-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.chart-placeholder {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  height: 100%;
  width: 100%;
  justify-content: space-around;
}

.bar {
  width: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.3s ease;
}

.pie-chart {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  background: conic-gradient(
    #667eea 0% 45%,
    #f093fb 45% 75%,
    #4facfe 75% 90%,
    #43e97b 90% 100%
  );
  position: relative;
}

.legend {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-left: 20px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

/* 表格样式 */
.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  background: #f9fafb;
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  color: #374151;
  font-size: 14px;
  border-bottom: 2px solid #e5e7eb;
}

.data-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #e5e7eb;
  color: #6b7280;
  font-size: 14px;
}

.data-table tr:hover td {
  background: #f9fafb;
}

/* 状态徽章 */
.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.inactive {
  background: #fee2e2;
  color: #991b1b;
}

.status-badge.pending {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.paid {
  background: #dbeafe;
  color: #1e40af;
}

.status-badge.shipped {
  background: #e0e7ff;
  color: #3730a3;
}

.status-badge.completed {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.seckill {
  background: #fee2e2;
  color: #991b1b;
}

.status-badge.normal {
  background: #d1fae5;
  color: #065f46;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-icon {
  width: 32px;
  height: 32px;
  border: none;
  background: #f3f4f6;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  transition: all 0.3s ease;
}

.btn-icon:hover {
  background: #e5e7eb;
  color: #374151;
}

.btn-icon.danger:hover {
  background: #fee2e2;
  color: #991b1b;
}

/* 设置表单 */
.settings-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.setting-card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
}

.setting-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.setting-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-input {
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.checkbox-row {
  flex-direction: row;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #374151;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.setting-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* 通用按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-secondary {
  background: #6b7280;
  color: white;
}

.btn-secondary:hover {
  background: #4b5563;
  transform: translateY(-1px);
}

.btn-success {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.btn-success:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.btn-outline {
  background: white;
  border: 1px solid #d1d5db;
  color: #374151;
}

.btn-outline:hover {
  background: #f9fafb;
  border-color: #667eea;
  color: #667eea;
}

/* 筛选器 */
.filter-group {
  display: flex;
  gap: 12px;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

/* 响应式 */
@media (max-width: 1024px) {
  .admin-content {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
  }
  
  .sidebar-nav {
    position: static;
  }
  
  .nav-list {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 8px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
  }
}

/* 搜索框样式 */
.search-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  align-items: flex-end;
  flex-wrap: wrap;
}

.search-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 200px;
}

.search-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.search-group input {
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.3s ease;
}

.search-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* 加载状态样式 */
.loading-overlay {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 16px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-overlay p {
  color: #6b7280;
  font-size: 14px;
}

/* 错误信息样式 */
.error-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  gap: 16px;
  text-align: center;
}

.error-message p {
  color: #ef4444;
  font-size: 14px;
  max-width: 400px;
}

/* 空状态样式 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 16px;
  color: #9ca3af;
}

.empty-state svg {
  stroke-width: 1.5;
}

.empty-state p {
  font-size: 14px;
}

/* 分页控件样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 24px;
  margin-top: 24px;
  padding: 16px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.btn-pagination {
  padding: 8px 16px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-pagination:hover:not(:disabled) {
  background: #5a67d8;
}

.btn-pagination:disabled {
  background: #d1d5db;
  cursor: not-allowed;
}

.pagination-info {
  font-size: 14px;
  color: #6b7280;
}
</style>