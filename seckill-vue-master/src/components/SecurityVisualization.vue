<template>
  <div class="security-visualization-container" :class="{ embedded }">
    <Navbar v-if="!embedded" />
    
    <div v-if="!embedded" class="page-header">
      <div class="container">
        <div class="header-content">
          <div class="header-left">
            <div class="breadcrumb">
              <span class="breadcrumb-item" @click="goBack" style="cursor: pointer;">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="15 18 9 12 15 6"></polyline>
                </svg>
                返回
              </span>
              <span class="breadcrumb-separator">/</span>
              <span class="breadcrumb-item active">安全可视化</span>
            </div>
            <h1>系统安全监控中心</h1>
            <p class="header-subtitle">实时监控安全状态、限流机制、权限控制与审计日志</p>
          </div>
          <div class="header-right">
            <div class="security-status-badge" :class="securityStatus">
              <span class="status-dot"></span>
              <span class="status-text">{{ securityStatusText }}</span>
            </div>
            <button class="btn btn-outline" @click="refreshSecurityData">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 4 23 10 17 10"></polyline>
                <polyline points="1 20 1 14 7 14"></polyline>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
              </svg>
              刷新数据
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="container">
      <div class="security-overview-grid">
        <div class="security-card" v-for="(item, index) in securityOverview" :key="index" :class="{ 'active': item.enabled }">
          <div class="card-icon" :style="{ background: item.gradient }">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <component :is="item.icon" />
            </svg>
          </div>
          <div class="card-content">
            <div class="card-title">{{ item.title }}</div>
            <div class="card-status">
              <span class="status-indicator" :class="item.enabled ? 'on' : 'off'"></span>
              {{ item.enabled ? '已启用' : '已禁用' }}
            </div>
            <div class="card-metric">{{ item.metric }}</div>
          </div>
        </div>
      </div>

      <div class="main-content">
        <div class="left-panel">
          <div class="panel-card">
            <div class="panel-header">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="3" width="7" height="7"></rect>
                  <rect x="14" y="3" width="7" height="7"></rect>
                  <rect x="14" y="14" width="7" height="7"></rect>
                  <rect x="3" y="14" width="7" height="7"></rect>
                </svg>
                限流监控
              </h3>
              <div class="panel-controls">
                <select v-model="rateLimitTimeRange" class="chart-select">
                  <option value="1m">最近1分钟</option>
                  <option value="5m">最近5分钟</option>
                  <option value="15m">最近15分钟</option>
                  <option value="1h">最近1小时</option>
                </select>
              </div>
            </div>
            <div class="panel-body">
              <div class="chart-wrapper">
                <svg width="560" height="280" class="rate-limit-chart">
                  <defs>
                    <linearGradient id="requestGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" style="stop-color:#667eea;stop-opacity:0.3" />
                      <stop offset="100%" style="stop-color:#667eea;stop-opacity:0" />
                    </linearGradient>
                    <linearGradient id="limitGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" style="stop-color:#f5576c;stop-opacity:0.3" />
                      <stop offset="100%" style="stop-color:#f5576c;stop-opacity:0" />
                    </linearGradient>
                  </defs>
                  
                  <g class="grid-lines">
                    <line v-for="i in 5" :key="'grid-' + i"
                          :x1="50" :y1="40 + i * 44"
                          :x2="530" :y2="40 + i * 44"
                          stroke="#e5e7eb" stroke-width="1" />
                  </g>
                  
                  <g class="y-axis-labels">
                    <text v-for="(value, i) in rateLimitYLabels" :key="'y-' + i"
                          x="45" :y="45 + i * 44"
                          text-anchor="end" class="chart-label">
                      {{ value }}
                    </text>
                  </g>
                  
                  <path class="chart-area" :d="requestAreaPath" fill="url(#requestGradient)" />
                  <path class="chart-line" :d="requestLinePath" fill="none" stroke="#667eea" stroke-width="2.5" />
                  <g class="chart-points">
                    <circle v-for="(point, i) in requestPoints" :key="'rp-' + i"
                            :cx="point.x" :cy="point.y" r="4"
                            fill="#667eea"
                            @mouseenter="showRateLimitTooltip('request', i, $event)"
                            @mouseleave="hideRateLimitTooltip"
                            style="cursor: pointer;" />
                  </g>
                  
                  <path class="chart-area" :d="limitAreaPath" fill="url(#limitGradient)" />
                  <path class="chart-line" :d="limitLinePath" fill="none" stroke="#f5576c" stroke-width="2.5" />
                  <g class="chart-points">
                    <circle v-for="(point, i) in limitPoints" :key="'lp-' + i"
                            :cx="point.x" :cy="point.y" r="4"
                            fill="#f5576c"
                            @mouseenter="showRateLimitTooltip('limit', i, $event)"
                            @mouseleave="hideRateLimitTooltip"
                            style="cursor: pointer;" />
                  </g>
                  
                  <g class="x-axis-labels">
                    <text v-for="(label, i) in rateLimitLabels" :key="'x-' + i"
                          :x="70 + i * 51.6" y="270"
                          text-anchor="middle" class="chart-label">
                      {{ label }}
                    </text>
                  </g>
                </svg>
                
                <div class="chart-legend">
                  <div class="legend-item">
                    <span class="legend-color" style="background: #667eea;"></span>
                    <span class="legend-text">总请求数</span>
                  </div>
                  <div class="legend-item">
                    <span class="legend-color" style="background: #f5576c;"></span>
                    <span class="legend-text">限流数</span>
                  </div>
                </div>
                
                <div v-if="rateLimitTooltip.show" class="chart-tooltip" :style="{ left: rateLimitTooltip.x + 'px', top: rateLimitTooltip.y + 'px' }">
                  <div class="tooltip-title">{{ rateLimitTooltip.label }}</div>
                  <div class="tooltip-content">
                    <span class="tooltip-value">{{ rateLimitTooltip.value }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="panel-card">
            <div class="panel-header">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                  <polyline points="22 4 12 14.01 9 11.01"></polyline>
                </svg>
                安全事件实时监控
              </h3>
              <div class="panel-actions">
                <button class="btn btn-outline btn-sm" @click="clearSecurityEvents">
                  清空
                </button>
              </div>
            </div>
            <div class="panel-body">
              <div class="events-container">
                <div v-for="(event, index) in securityEvents" :key="index" class="event-item" :class="event.level">
                  <div class="event-icon" :style="{ background: event.color }">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <component :is="event.icon" />
                    </svg>
                  </div>
                  <div class="event-content">
                    <div class="event-title">{{ event.title }}</div>
                    <div class="event-desc">{{ event.description }}</div>
                  </div>
                  <div class="event-time">{{ event.time }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="right-panel">
          <div class="panel-card">
            <div class="panel-header">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
                </svg>
                安全威胁分布
              </h3>
            </div>
            <div class="panel-body">
              <div class="threat-stats">
                <div class="threat-item" v-for="(threat, index) in threatDistribution" :key="index">
                  <div class="threat-info">
                    <span class="threat-name">{{ threat.name }}</span>
                    <span class="threat-count">{{ threat.count }}</span>
                  </div>
                  <div class="threat-bar-wrapper">
                    <div class="threat-bar" :style="{ width: threat.percentage + '%', background: threat.color }"></div>
                  </div>
                  <span class="threat-percentage">{{ threat.percentage.toFixed(1) }}%</span>
                </div>
              </div>
            </div>
          </div>

          <div class="panel-card">
            <div class="panel-header">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polygon points="12 2 2 7 12 12 22 7 12 2"></polygon>
                  <polyline points="2 17 12 22 22 17"></polyline>
                  <polyline points="2 12 12 17 22 12"></polyline>
                </svg>
                权限控制列表
              </h3>
            </div>
            <div class="panel-body">
              <div class="permission-list">
                <div class="permission-item" v-for="(perm, index) in permissions" :key="index">
                  <div class="permission-left">
                    <div class="permission-icon" :style="{ background: perm.gradient }">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <component :is="perm.icon" />
                      </svg>
                    </div>
                    <div class="permission-info">
                      <div class="permission-name">{{ perm.name }}</div>
                      <div class="permission-code">{{ perm.code }}</div>
                    </div>
                  </div>
                  <div class="permission-right">
                    <span class="permission-badge" :class="perm.enabled ? 'active' : 'inactive'">
                      {{ perm.enabled ? '已授权' : '未授权' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="panel-card">
            <div class="panel-header">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
                安全测试工具
              </h3>
            </div>
            <div class="panel-body">
              <div class="security-tests">
                <div class="test-card" v-for="(test, index) in securityTests" :key="index">
                  <div class="test-header">
                    <div class="test-icon" :style="{ background: test.gradient }">
                      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <component :is="test.icon" />
                      </svg>
                    </div>
                    <div class="test-info">
                      <div class="test-name">{{ test.name }}</div>
                      <div class="test-status" :class="test.status">{{ test.statusText }}</div>
                    </div>
                  </div>
                  <p class="test-desc">{{ test.description }}</p>
                  <button class="btn btn-primary btn-sm" @click="runSecurityTest(index)" :disabled="test.status === 'running'">
                    {{ test.status === 'running' ? '测试中...' : '开始测试' }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from './Navbar.vue';

const { embedded = false } = defineProps({
  embedded: {
    type: Boolean,
    default: false
  }
});

const router = useRouter();

const securityStatus = ref('secure');
const rateLimitTimeRange = ref('5m');

const securityStatusText = computed(() => {
  const statusMap = {
    'secure': '系统安全',
    'warning': '存在警告',
    'danger': '发现威胁'
  };
  return statusMap[securityStatus.value] || '系统安全';
});

const securityOverview = ref([
  {
    title: 'JWT认证',
    enabled: true,
    metric: '令牌有效',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    icon: 'shieldIcon'
  },
  {
    title: '限流机制',
    enabled: true,
    metric: '正常运作',
    gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
    icon: 'clockIcon'
  },
  {
    title: '验证码防护',
    enabled: true,
    metric: '已启用',
    gradient: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
    icon: 'captchaIcon'
  },
  {
    title: 'RBAC权限',
    enabled: true,
    metric: '权限正常',
    gradient: 'linear-gradient(135deg, #ec4899 0%, #be185d 100%)',
    icon: 'userIcon'
  }
]);

const rateLimitData = ref([
  { requests: 45, limits: 2 },
  { requests: 52, limits: 3 },
  { requests: 68, limits: 5 },
  { requests: 75, limits: 8 },
  { requests: 82, limits: 6 },
  { requests: 78, limits: 4 },
  { requests: 85, limits: 10 },
  { requests: 92, limits: 12 },
  { requests: 76, limits: 5 },
  { requests: 68, limits: 3 }
]);

const rateLimitLabels = ref(['0s', '5s', '10s', '15s', '20s', '25s', '30s', '35s', '40s', '45s']);
const rateLimitYLabels = ref(['100', '75', '50', '25', '0']);

const requestLinePath = computed(() => {
  const maxVal = 100;
  const points = rateLimitData.value.map((data, i) => {
    const x = 70 + i * 51.6;
    const y = 260 - (data.requests / maxVal) * 220;
    return `${x},${y}`;
  });
  return 'M ' + points.join(' L ');
});

const requestAreaPath = computed(() => {
  const maxVal = 100;
  const points = rateLimitData.value.map((data, i) => {
    const x = 70 + i * 51.6;
    const y = 260 - (data.requests / maxVal) * 220;
    return `${x},${y}`;
  });
  return 'M 70,260 L ' + points.join(' L ') + ' 530,260 Z';
});

const requestPoints = computed(() => {
  const maxVal = 100;
  return rateLimitData.value.map((data, i) => ({
    x: 70 + i * 51.6,
    y: 260 - (data.requests / maxVal) * 220
  }));
});

const limitLinePath = computed(() => {
  const maxVal = 100;
  const points = rateLimitData.value.map((data, i) => {
    const x = 70 + i * 51.6;
    const y = 260 - (data.limits / maxVal) * 220;
    return `${x},${y}`;
  });
  return 'M ' + points.join(' L ');
});

const limitAreaPath = computed(() => {
  const maxVal = 100;
  const points = rateLimitData.value.map((data, i) => {
    const x = 70 + i * 51.6;
    const y = 260 - (data.limits / maxVal) * 220;
    return `${x},${y}`;
  });
  return 'M 70,260 L ' + points.join(' L ') + ' 530,260 Z';
});

const limitPoints = computed(() => {
  const maxVal = 100;
  return rateLimitData.value.map((data, i) => ({
    x: 70 + i * 51.6,
    y: 260 - (data.limits / maxVal) * 220
  }));
});

const rateLimitTooltip = reactive({
  show: false,
  x: 0,
  y: 0,
  value: 0,
  label: ''
});

const securityEvents = ref([
  {
    level: 'info',
    title: '用户登录成功',
    description: '用户 admin 成功登录系统',
    time: '10:32:15',
    color: '#667eea',
    icon: 'loginIcon'
  },
  {
    level: 'warning',
    title: '限流触发',
    description: 'IP 192.168.1.100 请求频率超限',
    time: '10:31:42',
    color: '#f59e0b',
    icon: 'warningIcon'
  },
  {
    level: 'success',
    title: '权限验证通过',
    description: '用户权限检查成功',
    time: '10:30:58',
    color: '#10b981',
    icon: 'checkIcon'
  },
  {
    level: 'danger',
    title: '验证码错误',
    description: '多次验证码输入错误',
    time: '10:29:33',
    color: '#f5576c',
    icon: 'errorIcon'
  }
]);

const threatDistribution = ref([
  { name: '暴力破解尝试', count: 12, percentage: 35.3, color: '#f5576c' },
  { name: 'SQL注入尝试', count: 8, percentage: 23.5, color: '#f59e0b' },
  { name: 'XSS攻击尝试', count: 6, percentage: 17.6, color: '#ec4899' },
  { name: 'CSRF攻击尝试', count: 4, percentage: 11.8, color: '#8b5cf6' },
  { name: '其他异常', count: 4, percentage: 11.8, color: '#6b7280' }
]);

const permissions = ref([
  {
    name: '查看商品',
    code: 'product:view',
    enabled: true,
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    icon: 'eyeIcon'
  },
  {
    name: '创建商品',
    code: 'product:create',
    enabled: false,
    gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
    icon: 'plusIcon'
  },
  {
    name: '参与秒杀',
    code: 'product:seckill',
    enabled: true,
    gradient: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
    icon: 'flashIcon'
  },
  {
    name: '查看订单',
    code: 'order:view',
    enabled: true,
    gradient: 'linear-gradient(135deg, #ec4899 0%, #be185d 100%)',
    icon: 'orderIcon'
  },
  {
    name: '系统管理',
    code: 'admin:system',
    enabled: false,
    gradient: 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)',
    icon: 'settingsIcon'
  }
]);

const securityTests = ref([
  {
    name: 'JWT认证测试',
    description: '验证JWT令牌的生成、解析和过期机制',
    status: 'idle',
    statusText: '待测试',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    icon: 'shieldIcon'
  },
  {
    name: '限流机制测试',
    description: '测试系统的限流策略和防护能力',
    status: 'idle',
    statusText: '待测试',
    gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
    icon: 'clockIcon'
  },
  {
    name: '验证码验证测试',
    description: '测试验证码的生成和验证功能',
    status: 'idle',
    statusText: '待测试',
    gradient: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
    icon: 'captchaIcon'
  }
]);

const shieldIcon = '<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>';
const clockIcon = '<circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline>';
const captchaIcon = '<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path>';
const userIcon = '<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle>';
const loginIcon = '<path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path><polyline points="10 17 15 12 10 7"></polyline><line x1="15" y1="12" x2="3" y2="12"></line>';
const warningIcon = '<path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line>';
const checkIcon = '<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline>';
const errorIcon = '<circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line>';
const eyeIcon = '<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle>';
const plusIcon = '<line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line>';
const flashIcon = '<polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>';
const orderIcon = '<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline>';
const settingsIcon = '<circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>';

let monitoringInterval = null;

const goBack = () => {
  router.back();
};

const refreshSecurityData = () => {
  const now = new Date().toLocaleTimeString();
  securityEvents.value.unshift({
    level: 'info',
    title: '数据刷新',
    description: '安全数据已手动刷新',
    time: now,
    color: '#667eea',
    icon: 'checkIcon'
  });
  
  if (securityEvents.value.length > 10) {
    securityEvents.value.pop();
  }
};

const showRateLimitTooltip = (type, index, event) => {
  const chart = event.target.closest('.chart-wrapper');
  const rect = chart.getBoundingClientRect();
  rateLimitTooltip.show = true;
  rateLimitTooltip.x = event.clientX - rect.left + 10;
  rateLimitTooltip.y = event.clientY - rect.top - 60;
  
  if (type === 'request') {
    rateLimitTooltip.value = rateLimitData.value[index].requests + ' 次';
    rateLimitTooltip.label = rateLimitLabels.value[index] + ' - 总请求';
  } else {
    rateLimitTooltip.value = rateLimitData.value[index].limits + ' 次';
    rateLimitTooltip.label = rateLimitLabels.value[index] + ' - 限流';
  }
};

const hideRateLimitTooltip = () => {
  rateLimitTooltip.show = false;
};

const clearSecurityEvents = () => {
  securityEvents.value = [];
};

const runSecurityTest = (index) => {
  const test = securityTests.value[index];
  test.status = 'running';
  test.statusText = '测试中...';
  
  setTimeout(() => {
    test.status = 'completed';
    test.statusText = '测试完成';
    
    const now = new Date().toLocaleTimeString();
    securityEvents.value.unshift({
      level: 'success',
      title: test.name + ' 完成',
      description: '安全测试执行成功',
      time: now,
      color: '#10b981',
      icon: 'checkIcon'
    });
    
    setTimeout(() => {
      test.status = 'idle';
      test.statusText = '待测试';
    }, 3000);
  }, 2000);
};

const updateRateLimitData = () => {
  rateLimitData.value.shift();
  rateLimitData.value.push({
    requests: Math.floor(Math.random() * 60) + 40,
    limits: Math.floor(Math.random() * 15)
  });
  
  const now = new Date();
  const seconds = now.getSeconds();
  rateLimitLabels.value.shift();
  rateLimitLabels.value.push(seconds + 's');
};

onMounted(() => {
  monitoringInterval = setInterval(() => {
    updateRateLimitData();
    
    if (Math.random() > 0.7) {
      const eventTypes = [
        { level: 'info', title: '请求成功', desc: 'API请求正常处理', color: '#667eea', icon: 'checkIcon' },
        { level: 'warning', title: '限流预警', desc: '接近限流阈值', color: '#f59e0b', icon: 'warningIcon' },
        { level: 'success', title: '验证通过', desc: '安全验证成功', color: '#10b981', icon: 'checkIcon' }
      ];
      const randomEvent = eventTypes[Math.floor(Math.random() * eventTypes.length)];
      const now = new Date().toLocaleTimeString();
      
      securityEvents.value.unshift({
        level: randomEvent.level,
        title: randomEvent.title,
        description: randomEvent.desc,
        time: now,
        color: randomEvent.color,
        icon: randomEvent.icon
      });
      
      if (securityEvents.value.length > 10) {
        securityEvents.value.pop();
      }
    }
  }, 3000);
});

onUnmounted(() => {
  if (monitoringInterval) {
    clearInterval(monitoringInterval);
  }
});
</script>

<style scoped>
.security-visualization-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea10 0%, #764ba210 100%);
}

.security-visualization-container.embedded {
  min-height: auto;
  background: transparent;
}

.page-header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 24px 0;
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.header-left {
  flex: 1;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 14px;
  color: #6b7280;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 4px;
  transition: color 0.2s;
}

.breadcrumb-item:hover {
  color: #667eea;
}

.breadcrumb-item.active {
  color: #374151;
  font-weight: 500;
}

.breadcrumb-separator {
  color: #d1d5db;
}

.header-left h1 {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.header-subtitle {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.security-status-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 24px;
}

.security-status-badge.warning {
  background: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
}

.security-status-badge.danger {
  background: rgba(245, 87, 108, 0.1);
  border-color: rgba(245, 87, 108, 0.2);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #10b981;
  animation: pulse 2s infinite;
}

.security-status-badge.warning .status-dot {
  background: #f59e0b;
}

.security-status-badge.danger .status-dot {
  background: #f5576c;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-text {
  font-size: 14px;
  font-weight: 500;
  color: #10b981;
}

.security-status-badge.warning .status-text {
  color: #f59e0b;
}

.security-status-badge.danger .status-text {
  color: #f5576c;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.security-overview-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.security-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.security-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.security-card.active {
  border-color: rgba(102, 126, 234, 0.2);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.02) 0%, rgba(118, 75, 162, 0.02) 100%);
}

.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.card-content {
  flex: 1;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.card-status {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.on {
  background: #10b981;
}

.status-indicator.off {
  background: #6b7280;
}

.card-metric {
  font-size: 13px;
  color: #667eea;
  font-weight: 500;
}

.main-content {
  display: grid;
  grid-template-columns: 1fr 420px;
  gap: 24px;
}

.panel-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
  overflow: hidden;
}

.panel-header {
  padding: 16px 20px;
  border-bottom: 1px solid #f3f4f6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.panel-body {
  padding: 20px;
}

.chart-wrapper {
  position: relative;
}

.rate-limit-chart {
  width: 100%;
}

.chart-label {
  font-size: 11px;
  fill: #9ca3af;
}

.chart-line {
  stroke-linecap: round;
  stroke-linejoin: round;
}

.chart-points circle {
  transition: r 0.2s;
}

.chart-points circle:hover {
  r: 6;
}

.chart-legend {
  display: flex;
  gap: 20px;
  margin-top: 16px;
  justify-content: center;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6b7280;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.chart-tooltip {
  position: absolute;
  background: white;
  padding: 10px 14px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  pointer-events: none;
  z-index: 10;
}

.tooltip-title {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
}

.tooltip-value {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.events-container {
  max-height: 320px;
  overflow-y: auto;
}

.event-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  margin-bottom: 8px;
  border-left: 3px solid transparent;
}

.event-item:last-child {
  margin-bottom: 0;
}

.event-item.info {
  border-left-color: #667eea;
}

.event-item.warning {
  border-left-color: #f59e0b;
}

.event-item.success {
  border-left-color: #10b981;
}

.event-item.danger {
  border-left-color: #f5576c;
}

.event-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.event-content {
  flex: 1;
  min-width: 0;
}

.event-title {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 2px;
}

.event-desc {
  font-size: 13px;
  color: #6b7280;
}

.event-time {
  font-size: 12px;
  color: #9ca3af;
}

.threat-stats {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.threat-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.threat-info {
  width: 120px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.threat-name {
  font-size: 13px;
  color: #374151;
}

.threat-count {
  font-size: 13px;
  font-weight: 600;
  color: #1f2937;
}

.threat-bar-wrapper {
  flex: 1;
  height: 8px;
  background: #f3f4f6;
  border-radius: 4px;
  overflow: hidden;
}

.threat-bar {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s;
}

.threat-percentage {
  width: 45px;
  text-align: right;
  font-size: 13px;
  color: #6b7280;
  font-weight: 500;
}

.permission-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.permission-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.permission-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.permission-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.permission-name {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 2px;
}

.permission-code {
  font-size: 12px;
  color: #9ca3af;
  font-family: monospace;
}

.permission-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.permission-badge.active {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.permission-badge.inactive {
  background: #f3f4f6;
  color: #9ca3af;
}

.security-tests {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.test-card {
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.test-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.test-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.test-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.test-status {
  font-size: 12px;
}

.test-status.idle {
  color: #9ca3af;
}

.test-status.running {
  color: #f59e0b;
}

.test-status.completed {
  color: #10b981;
}

.test-desc {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 12px;
  line-height: 1.5;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-outline {
  background: white;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-outline:hover {
  background: #f9fafb;
  border-color: #9ca3af;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.panel-controls {
  display: flex;
  gap: 8px;
}

.panel-actions {
  display: flex;
  gap: 8px;
}

.chart-select {
  padding: 6px 10px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 13px;
  color: #374151;
  background: white;
  cursor: pointer;
}

.chart-select:focus {
  outline: none;
  border-color: #667eea;
}

.events-container::-webkit-scrollbar {
  width: 6px;
}

.events-container::-webkit-scrollbar-track {
  background: #f3f4f6;
  border-radius: 3px;
}

.events-container::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 3px;
}

.events-container::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

@media (max-width: 1024px) {
  .main-content {
    grid-template-columns: 1fr;
  }
  
  .security-overview-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .security-overview-grid {
    grid-template-columns: 1fr;
  }
  
  .threat-info {
    width: 100px;
  }
}
</style>
