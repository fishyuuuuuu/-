<template>
  <div class="performance-test-container">
    <Navbar />
    
    <div class="page-header">
      <div class="container">
        <div class="header-content">
          <div class="header-title">
            <h1>秒杀性能测试平台</h1>
            <p class="header-subtitle">高并发场景压测与性能分析</p>
          </div>
          <div class="header-actions">
            <div class="status-badge" :class="testStatus">
              <span class="status-dot"></span>
              <span class="status-text">{{ statusText }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="container">
      <div class="test-layout">
        <div class="config-panel">
          <div class="panel-card">
            <div class="panel-header">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3"/>
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
                </svg>
                测试场景配置
              </h3>
            </div>
            <div class="panel-body">
              <div class="config-section">
                <h4 class="section-title">商品配置</h4>
                <div class="form-group">
                  <label class="form-label">商品ID</label>
                  <input type="number" v-model.number="config.productId" class="form-input" min="1">
                </div>
                <div class="form-group">
                  <label class="form-label">库存数量</label>
                  <input type="number" v-model.number="config.stock" class="form-input" min="1" max="10000">
                  <span class="form-hint">设置商品的初始库存数量</span>
                </div>
              </div>

              <div class="config-section">
                <h4 class="section-title">并发配置</h4>
                <div class="form-group">
                  <label class="form-label">并发用户数</label>
                  <input type="number" v-model.number="config.concurrentUsers" class="form-input" min="1" max="10000">
                  <span class="form-hint">同时发起请求的用户数量</span>
                </div>
                <div class="form-group">
                  <label class="form-label">请求间隔 (ms)</label>
                  <input type="number" v-model.number="config.requestInterval" class="form-input" min="0" max="5000">
                  <span class="form-hint">每个请求之间的间隔时间</span>
                </div>
                <div class="form-group">
                  <label class="form-label">测试持续时间 (秒)</label>
                  <input type="number" v-model.number="config.duration" class="form-input" min="1" max="3600">
                </div>
              </div>

              <div class="config-section">
                <h4 class="section-title">高级配置</h4>
                <div class="form-group">
                  <label class="form-label">超时时间 (ms)</label>
                  <input type="number" v-model.number="config.timeout" class="form-input" min="1000" max="30000">
                </div>
                <div class="form-group">
                  <label class="form-label">失败重试次数</label>
                  <input type="number" v-model.number="config.retryCount" class="form-input" min="0" max="5">
                </div>
              </div>

              <div class="control-buttons">
                <button 
                  class="btn btn-primary btn-lg" 
                  @click="startTest"
                  :disabled="testStatus === 'running'"
                >
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polygon points="5 3 19 12 5 21 5 3"/>
                  </svg>
                  开始测试
                </button>
                <button 
                  class="btn btn-warning btn-lg" 
                  @click="pauseTest"
                  :disabled="testStatus !== 'running'"
                >
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="6" y="4" width="4" height="16"/>
                    <rect x="14" y="4" width="4" height="16"/>
                  </svg>
                  暂停测试
                </button>
                <button 
                  class="btn btn-danger btn-lg" 
                  @click="stopTest"
                  :disabled="testStatus === 'idle'"
                >
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                  </svg>
                  终止测试
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="monitor-panel">
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
              </div>
              <div class="stat-content">
                <span class="stat-value">{{ realtimeStats.concurrentUsers }}</span>
                <span class="stat-label">当前并发用户</span>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #10b981 0%, #059669 100%);">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                  <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
              </div>
              <div class="stat-content">
                <span class="stat-value success">{{ realtimeStats.successCount }}</span>
                <span class="stat-label">成功抢购数</span>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #f5576c 0%, #f093fb 100%);">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="15" y1="9" x2="9" y2="15"/>
                  <line x1="9" y1="9" x2="15" y2="15"/>
                </svg>
              </div>
              <div class="stat-content">
                <span class="stat-value error">{{ realtimeStats.failCount }}</span>
                <span class="stat-label">失败抢购数</span>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon" style="background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
                </svg>
              </div>
              <div class="stat-content">
                <span class="stat-value">{{ realtimeStats.qps.toFixed(1) }}</span>
                <span class="stat-label">QPS (请求/秒)</span>
              </div>
            </div>
          </div>

          <div class="charts-grid">
            <div class="chart-card">
              <div class="chart-header">
                <h3>响应时间趋势</h3>
                <div class="chart-controls">
                  <select v-model="timeGranularity" class="chart-select">
                    <option value="5">5秒</option>
                    <option value="10">10秒</option>
                    <option value="30">30秒</option>
                    <option value="60">1分钟</option>
                  </select>
                </div>
                <div class="chart-legend">
                  <span class="legend-item">
                    <span class="legend-color" style="background: #667eea;"></span>
                    平均响应时间
                  </span>
                  <span class="legend-item">
                    <span class="legend-color" style="background: #10b981;"></span>
                    P90
                  </span>
                  <span class="legend-item">
                    <span class="legend-color" style="background: #f59e0b;"></span>
                    P99
                  </span>
                </div>
              </div>
              <div class="chart-body">
                <div class="line-chart" ref="responseTimeChart">
                  <svg width="100%" height="200" viewBox="0 0 600 200">
                    <defs>
                      <pattern id="grid" width="60" height="40" patternUnits="userSpaceOnUse">
                        <path d="M 60 0 L 0 0 0 40" fill="none" stroke="#e5e7eb" stroke-width="1"/>
                      </pattern>
                    </defs>
                    <rect width="100%" height="100%" fill="url(#grid)"/>
                    
                    <polyline 
                      :points="responseTimeData.avgPoints" 
                      fill="none" 
                      stroke="#667eea" 
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                    
                    <polyline 
                      :points="responseTimeData.p90Points" 
                      fill="none" 
                      stroke="#10b981" 
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                    
                    <polyline 
                      :points="responseTimeData.p99Points" 
                      fill="none" 
                      stroke="#f59e0b" 
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                  <div class="chart-labels">
                    <span>0s</span>
                    <span>{{ config.duration / 4 }}s</span>
                    <span>{{ config.duration / 2 }}s</span>
                    <span>{{ config.duration * 3 / 4 }}s</span>
                    <span>{{ config.duration }}s</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="chart-card">
              <div class="chart-header">
                <h3>吞吐量变化</h3>
              </div>
              <div class="chart-body">
                <div class="bar-chart" ref="throughputChart">
                  <svg width="100%" height="200" viewBox="0 0 600 200">
                    <defs>
                      <linearGradient id="barGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                        <stop offset="0%" style="stop-color:#667eea;stop-opacity:1" />
                        <stop offset="100%" style="stop-color:#764ba2;stop-opacity:1" />
                      </linearGradient>
                    </defs>
                    <g v-for="(bar, index) in throughputData.bars" :key="index">
                      <rect 
                        :x="bar.x"
                        :y="bar.y"
                        :width="bar.width"
                        :height="bar.height"
                        fill="url(#barGradient)"
                        rx="4"
                        @mouseenter="showThroughputTooltip(index, $event)"
                        @mouseleave="hideThroughputTooltip"
                        style="cursor: pointer;"
                      />
                    </g>
                  </svg>
                  <div class="chart-labels">
                    <span v-for="(label, index) in throughputData.labels" :key="index">{{ label }}</span>
                  </div>
                  <div 
                    v-if="throughputTooltip.show" 
                    class="tooltip"
                    :style="{ left: throughputTooltip.x + 'px', top: throughputTooltip.y + 'px' }"
                  >
                    <div class="tooltip-title">{{ throughputTooltip.label }}</div>
                    <div class="tooltip-content">
                      <span class="tooltip-value">{{ throughputTooltip.value }} req/s</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="chart-card">
              <div class="chart-header">
                <h3>错误率统计分析</h3>
              </div>
              <div class="chart-body">
                <div class="error-stats">
                  <div class="error-item" v-for="(error, index) in errorStats" :key="index">
                    <div class="error-info">
                      <span class="error-type">{{ error.type }}</span>
                      <span class="error-count">{{ error.count }}</span>
                    </div>
                    <div class="error-bar-wrapper">
                      <div class="error-bar" :style="{ width: error.percentage + '%', background: error.color }"></div>
                    </div>
                    <span class="error-percentage">{{ error.percentage.toFixed(1) }}%</span>
                  </div>
                </div>
                <div class="error-description">
                  <div class="error-desc-item" v-for="(desc, type) in errorDescriptions" :key="type">
                    <div class="desc-icon" :style="{ background: getErrorColor(type) }">!</div>
                    <div class="desc-content">
                      <span class="desc-title">{{ type }}</span>
                      <span class="desc-text">{{ desc }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="chart-card">
              <div class="chart-header">
                <h3>服务器资源利用率</h3>
              </div>
              <div class="chart-body">
                <div class="resource-gauges">
                  <div class="gauge-item" v-for="(resource, key) in serverResources" :key="key">
                    <div class="gauge-wrapper">
                      <div class="gauge-circle" :style="{ background: `conic-gradient(${resourceColors[key]} ${Math.round(serverResources[key]) * 3.6}deg, #f0f0f0 0)` }">
                        <div class="gauge-inner">
                          <span class="gauge-value">{{ Math.round(serverResources[key]) }}</span>
                          <span class="gauge-unit">%</span>
                        </div>
                      </div>
                      <span class="gauge-label">{{ resourceLabels[key] }}</span>
                    </div>
                    <div class="resource-info" @click="toggleResourceInfo(key)">
                      <span class="info-icon">?</span>
                      <div class="info-tooltip" :class="{ show: expandedResource === key }">
                        <div class="info-title">{{ resourceLabels[key] }}详解</div>
                        <div class="info-desc">{{ resourceDescriptions[key] }}</div>
                        <div class="info-range">
                          <span class="range-label">正常范围:</span>
                          <span class="range-value">{{ resourceRanges[key] }}</span>
                        </div>
                        <div class="info-impact">
                          <span class="impact-label">性能影响:</span>
                          <span class="impact-desc">{{ resourceImpacts[key] }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="results-panel" v-if="testResults.length > 0">
            <div class="panel-card">
              <div class="panel-header">
                <h3>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                    <polyline points="14 2 14 8 20 8"/>
                    <line x1="16" y1="13" x2="8" y2="13"/>
                    <line x1="16" y1="17" x2="8" y2="17"/>
                    <polyline points="10 9 9 9 8 9"/>
                  </svg>
                  测试结果详情
                </h3>
                <div class="panel-actions">
                  <button class="btn btn-outline btn-sm" @click="exportResults('csv')">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                      <polyline points="7 10 12 15 17 10"/>
                      <line x1="12" y1="15" x2="12" y2="3"/>
                    </svg>
                    导出 CSV
                  </button>
                  <button class="btn btn-outline btn-sm" @click="exportResults('excel')">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                      <polyline points="14 2 14 8 20 8"/>
                      <line x1="8" y1="13" x2="16" y2="13"/>
                      <line x1="8" y1="17" x2="16" y2="17"/>
                    </svg>
                    导出 Excel
                  </button>
                </div>
              </div>
              <div class="panel-body">
                <div class="results-table-wrapper">
                  <table class="results-table">
                    <thead>
                      <tr>
                        <th>请求ID</th>
                        <th>用户ID</th>
                        <th>状态</th>
                        <th>响应时间</th>
                        <th>错误信息</th>
                        <th>时间戳</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(result, index) in paginatedResults" :key="index" :class="result.success ? 'success-row' : 'error-row'">
                        <td>{{ result.requestId }}</td>
                        <td>{{ result.userId }}</td>
                        <td>
                          <span class="status-tag" :class="result.success ? 'success' : 'error'">
                            {{ result.success ? '成功' : '失败' }}
                          </span>
                        </td>
                        <td>{{ result.responseTime }}ms</td>
                        <td>{{ result.errorMessage || '-' }}</td>
                        <td>{{ result.timestamp }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div class="pagination">
                  <button class="btn btn-sm" @click="prevPage" :disabled="currentPage === 1">上一页</button>
                  <span class="page-info">第 {{ currentPage }} / {{ totalPages }} 页</span>
                  <button class="btn btn-sm" @click="nextPage" :disabled="currentPage === totalPages">下一页</button>
                </div>
              </div>
            </div>
          </div>

          <div class="analysis-panel" v-if="showAnalysis">
            <div class="panel-card">
              <div class="panel-header">
                <h3>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="20" x2="18" y2="10"/>
                    <line x1="12" y1="20" x2="12" y2="4"/>
                    <line x1="6" y1="20" x2="6" y2="14"/>
                  </svg>
                  性能指标分析
                </h3>
              </div>
              <div class="panel-body">
                <div class="analysis-grid">
                  <div class="analysis-card">
                    <h4>关键性能指标</h4>
                    <div class="metric-list">
                      <div class="metric-item">
                        <span class="metric-label">平均响应时间</span>
                        <span class="metric-value">{{ analysisMetrics.avgResponseTime }}ms</span>
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">最大响应时间</span>
                        <span class="metric-value">{{ analysisMetrics.maxResponseTime }}ms</span>
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">最小响应时间</span>
                        <span class="metric-value">{{ analysisMetrics.minResponseTime }}ms</span>
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">平均吞吐量</span>
                        <span class="metric-value">{{ analysisMetrics.avgThroughput }} req/s</span>
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">成功率</span>
                        <span class="metric-value success">{{ analysisMetrics.successRate }}%</span>
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">错误率</span>
                        <span class="metric-value error">{{ analysisMetrics.errorRate }}%</span>
                      </div>
                    </div>
                  </div>

                  <div class="analysis-card">
                    <h4>性能瓶颈识别</h4>
                    <div class="bottleneck-list">
                      <div 
                        v-for="(bottleneck, index) in bottlenecks" 
                        :key="index" 
                        class="bottleneck-item"
                        :class="bottleneck.level"
                      >
                        <div class="bottleneck-icon">
                          <svg v-if="bottleneck.level === 'high'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="12" y1="8" x2="12" y2="12"/>
                            <line x1="12" y1="16" x2="12.01" y2="16"/>
                          </svg>
                          <svg v-else-if="bottleneck.level === 'medium'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                            <line x1="12" y1="9" x2="12" y2="13"/>
                            <line x1="12" y1="17" x2="12.01" y2="17"/>
                          </svg>
                          <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/>
                            <polyline points="12 6 12 12 16 14"/>
                          </svg>
                        </div>
                        <div class="bottleneck-content">
                          <span class="bottleneck-title">{{ bottleneck.title }}</span>
                          <span class="bottleneck-desc">{{ bottleneck.description }}</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="analysis-card full-width">
                    <h4>系统极限表现</h4>
                    <div class="limit-metrics">
                      <div class="limit-item">
                        <div class="limit-header">
                          <span class="limit-label">并发处理能力</span>
                          <span class="limit-value">{{ limitMetrics.concurrentCapacity }} req/s</span>
                        </div>
                        <div class="limit-bar">
                          <div class="limit-fill" :style="{ width: limitMetrics.concurrentPercent + '%' }"></div>
                        </div>
                      </div>
                      <div class="limit-item">
                        <div class="limit-header">
                          <span class="limit-label">系统稳定性</span>
                          <span class="limit-value">{{ limitMetrics.stability }}%</span>
                        </div>
                        <div class="limit-bar">
                          <div class="limit-fill success" :style="{ width: limitMetrics.stability + '%' }"></div>
                        </div>
                      </div>
                      <div class="limit-item">
                        <div class="limit-header">
                          <span class="limit-label">资源利用率</span>
                          <span class="limit-value">{{ limitMetrics.resourceUsage }}%</span>
                        </div>
                        <div class="limit-bar">
                          <div class="limit-fill warning" :style="{ width: limitMetrics.resourceUsage + '%' }"></div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="logs-panel" v-if="testLogs.length > 0">
      <div class="container">
        <div class="panel-card">
          <div class="panel-header">
            <h3>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
              </svg>
              测试日志
            </h3>
            <button class="btn btn-outline btn-sm" @click="clearLogs">清空日志</button>
          </div>
          <div class="panel-body">
            <div class="logs-container">
              <div 
                v-for="(log, index) in testLogs" 
                :key="index" 
                class="log-item"
                :class="log.level"
              >
                <span class="log-time">{{ log.time }}</span>
                <span class="log-level">{{ log.level.toUpperCase() }}</span>
                <span class="log-message">{{ log.message }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

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
          <p class="toast-title">{{ toast.title }}</p>
          <p class="toast-message">{{ toast.message }}</p>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import Navbar from './Navbar.vue';

const config = reactive({
  productId: 1,
  stock: 10,
  concurrentUsers: 500,
  requestInterval: 10,
  duration: 60,
  timeout: 10000,
  retryCount: 0
});

const testStatus = ref('idle');
const statusText = computed(() => {
  const texts = {
    idle: '未开始',
    running: '测试中',
    paused: '已暂停',
    completed: '已完成'
  };
  return texts[testStatus.value];
});

const realtimeStats = reactive({
  concurrentUsers: 0,
  successCount: 0,
  failCount: 0,
  qps: 0,
  totalRequests: 0
});

const responseTimeData = reactive({
  avgPoints: '',
  p90Points: '',
  p99Points: '',
  data: []
});

const throughputData = reactive({
  bars: [],
  labels: [],
  data: []
});

const throughputTooltip = reactive({
  show: false,
  x: 0,
  y: 0,
  value: 0,
  label: ''
});

const timeGranularity = ref('10');

const errorStats = ref([
  { type: '库存不足', count: 0, percentage: 0, color: '#f5576c' },
  { type: '请求超时', count: 0, percentage: 0, color: '#f59e0b' },
  { type: '限流拦截', count: 0, percentage: 0, color: '#8b5cf6' },
  { type: '系统内部错误', count: 0, percentage: 0, color: '#dc2626' },
  { type: '秒杀失败', count: 0, percentage: 0, color: '#ef4444' },
  { type: '其他错误', count: 0, percentage: 0, color: '#6b7280' }
]);

const errorDescriptions = {
  '库存不足': '商品库存已耗尽，无法继续抢购',
  '请求超时': '服务器响应超时，可能是负载过高',
  '限流拦截': '触发速率限制，防止系统过载',
  '系统内部错误': '服务器内部异常，需要检查系统日志',
  '秒杀失败': '秒杀操作执行失败，可能是库存竞争',
  '其他错误': '其他未分类的错误类型'
};

const serverResources = reactive({
  cpu: 0,
  memory: 0,
  network: 0,
  disk: 0
});

const resourceLabels = {
  cpu: 'CPU',
  memory: '内存',
  network: '网络IO',
  disk: '磁盘IO'
};

const resourceColors = {
  cpu: '#667eea',
  memory: '#10b981',
  network: '#f59e0b',
  disk: '#f5576c'
};

const resourceDescriptions = {
  cpu: 'CPU利用率表示处理器的使用程度。高CPU利用率意味着系统正在处理大量计算任务，可能导致响应延迟。',
  memory: '内存利用率表示RAM的使用程度。高内存利用率可能导致频繁的垃圾回收或页面置换，影响性能。',
  network: '网络IO利用率表示网络带宽的使用程度。高网络IO可能导致数据传输延迟，影响用户体验。',
  disk: '磁盘IO利用率表示存储设备的读写操作程度。高磁盘IO可能导致数据访问延迟。'
};

const resourceRanges = {
  cpu: '0-70% (正常), 70-90% (警告), >90% (危险)',
  memory: '0-70% (正常), 70-85% (警告), >85% (危险)',
  network: '0-60% (正常), 60-80% (警告), >80% (危险)',
  disk: '0-50% (正常), 50-70% (警告), >70% (危险)'
};

const resourceImpacts = {
  cpu: 'CPU过高会导致请求处理延迟，响应时间增加，吞吐量下降。',
  memory: '内存过高会导致OOM错误，系统崩溃或频繁的GC停顿。',
  network: '网络IO过高会导致数据传输超时，请求失败率上升。',
  disk: '磁盘IO过高会导致数据读写延迟，数据库查询变慢。'
};

const expandedResource = ref(null);

const testResults = ref([]);
const currentPage = ref(1);
const pageSize = 20;

const totalPages = computed(() => Math.ceil(testResults.value.length / pageSize));

const paginatedResults = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  const end = start + pageSize;
  return testResults.value.slice(start, end);
});

const showAnalysis = ref(false);
const analysisMetrics = reactive({
  avgResponseTime: 0,
  maxResponseTime: 0,
  minResponseTime: 0,
  avgThroughput: 0,
  successRate: 0,
  errorRate: 0
});

const bottlenecks = ref([]);

const limitMetrics = reactive({
  concurrentCapacity: 0,
  concurrentPercent: 0,
  stability: 0,
  resourceUsage: 0
});

const testLogs = ref([]);

const toast = reactive({
  show: false,
  title: '',
  message: '',
  type: 'success'
});

let testTimer = null;
let statsTimer = null;
let requestId = 0;

const showToast = (title, message, type = 'success') => {
  toast.title = title;
  toast.message = message;
  toast.type = type;
  toast.show = true;
  
  setTimeout(() => {
    toast.show = false;
  }, 3000);
};

const addLog = (level, message) => {
  testLogs.value.unshift({
    time: new Date().toLocaleTimeString(),
    level,
    message
  });
  
  if (testLogs.value.length > 500) {
    testLogs.value = testLogs.value.slice(0, 500);
  }
};

const clearLogs = () => {
  testLogs.value = [];
};

const startTest = async () => {
  const token = localStorage.getItem('token');

  addLog('info', `开始性能测试 - 并发用户: ${config.concurrentUsers}, 库存: ${config.stock}`);
  
  try {
    addLog('info', `正在重置商品 ${config.productId} 的库存为 ${config.stock}...`);
    await axios.post('/api/product/reset-stock', {
      productId: config.productId,
      stock: config.stock
    }, {
      headers: token ? {
        'Authorization': `Bearer ${token}`
      } : {}
    });
    addLog('success', `库存重置成功`);
  } catch (error) {
    addLog('warning', `库存重置失败: ${error.response?.data?.error || error.message}，继续测试`);
  }

  testStatus.value = 'running';
  realtimeStats.concurrentUsers = config.concurrentUsers;
  realtimeStats.successCount = 0;
  realtimeStats.failCount = 0;
  realtimeStats.totalRequests = 0;
  testResults.value = [];
  requestId = 0;
  
  errorStats.value.forEach(e => {
    e.count = 0;
    e.percentage = 0;
  });
  
  showToast('测试开始', `正在启动 ${config.concurrentUsers} 个并发用户进行测试`);

  initChartData();
  
  runConcurrentTest();
  
  startStatsUpdate();
};

const initChartData = () => {
  responseTimeData.data = [];
  throughputData.data = [];
  
  const points = 10;
  for (let i = 0; i < points; i++) {
    responseTimeData.data.push({
      avg: 0,
      p90: 0,
      p99: 0
    });
    throughputData.data.push(0);
  }
  
  updateCharts();
};

const updateCharts = () => {
  const width = 600;
  const height = 200;
  const padding = 40;
  const dataPoints = responseTimeData.data.length;
  
  const avgPoints = responseTimeData.data.map((d, i) => {
    const x = padding + (i / (dataPoints - 1)) * (width - padding * 2);
    const y = height - padding - (d.avg / 1000) * (height - padding * 2);
    return `${x},${y}`;
  }).join(' ');
  
  const p90Points = responseTimeData.data.map((d, i) => {
    const x = padding + (i / (dataPoints - 1)) * (width - padding * 2);
    const y = height - padding - (d.p90 / 1000) * (height - padding * 2);
    return `${x},${y}`;
  }).join(' ');
  
  const p99Points = responseTimeData.data.map((d, i) => {
    const x = padding + (i / (dataPoints - 1)) * (width - padding * 2);
    const y = height - padding - (d.p99 / 1000) * (height - padding * 2);
    return `${x},${y}`;
  }).join(' ');
  
  responseTimeData.avgPoints = avgPoints;
  responseTimeData.p90Points = p90Points;
  responseTimeData.p99Points = p99Points;
  
  const barWidth = (width - padding * 2) / throughputData.data.length - 10;
  const maxThroughput = Math.max(...throughputData.data, 100);
  
  throughputData.bars = throughputData.data.map((d, i) => ({
    x: padding + i * (barWidth + 10) + 5,
    y: height - padding - (d / maxThroughput) * (height - padding * 2),
    width: barWidth,
    height: (d / maxThroughput) * (height - padding * 2)
  }));
  
  throughputData.labels = throughputData.data.map((_, i) => `${i * 5}s`);
};

const showThroughputTooltip = (index, event) => {
  const chart = event.target.closest('.bar-chart');
  const rect = chart.getBoundingClientRect();
  throughputTooltip.show = true;
  throughputTooltip.x = event.clientX - rect.left + 10;
  throughputTooltip.y = event.clientY - rect.top - 60;
  throughputTooltip.value = throughputData.data[index].toFixed(1);
  throughputTooltip.label = throughputData.labels[index];
};

const hideThroughputTooltip = () => {
  throughputTooltip.show = false;
};

const getErrorColor = (type) => {
  const error = errorStats.value.find(e => e.type === type);
  return error ? error.color : '#6b7280';
};

const toggleResourceInfo = (key) => {
  expandedResource.value = expandedResource.value === key ? null : key;
};

const runConcurrentTest = async () => {
  const workers = [];
  const requestsPerWorker = Math.ceil(config.concurrentUsers / 10);
  
  for (let i = 0; i < 10; i++) {
    workers.push(runWorker(i, requestsPerWorker));
  }
  
  await Promise.all(workers);
  
  if (testStatus.value === 'running') {
    testStatus.value = 'completed';
    generateAnalysis();
    showToast('测试完成', `成功: ${realtimeStats.successCount}, 失败: ${realtimeStats.failCount}`);
    addLog('info', `测试完成 - 成功: ${realtimeStats.successCount}, 失败: ${realtimeStats.failCount}`);
  }
};

const runWorker = async (workerId, requestCount) => {
  for (let i = 0; i < requestCount; i++) {
    if (testStatus.value !== 'running') break;
    
    const currentRequestId = ++requestId;
    const startTime = Date.now();
    
    try {
      const response = await axios.post('/api/product/seckill', {
        productId: config.productId,
        captchaId: 'stress-test',
        captchaStr: '1234',
        userId: workerId * 1000 + i
      }, {
        headers: localStorage.getItem('token') ? {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        } : {},
        timeout: config.timeout
      });
      
      const responseTime = Date.now() - startTime;
      
      if (response.status === 200) {
        console.log('响应数据:', response.data);
        if (response.data.message === '秒杀成功') {
          realtimeStats.successCount++;
          testResults.value.push({
            requestId: currentRequestId,
            userId: workerId * 1000 + i,
            success: true,
            responseTime,
            errorMessage: null,
            timestamp: new Date().toLocaleTimeString()
          });
          
          addLog('success', `Worker ${workerId}: 请求 ${currentRequestId} 成功 (${responseTime}ms)`);
        } else {
          handleFailure(currentRequestId, workerId, i, responseTime, response.data.error || response.data.message || '未知错误');
        }
      } else {
        handleFailure(currentRequestId, workerId, i, responseTime, response.data.error || response.data.message || '未知错误');
      }
    } catch (error) {
      const responseTime = Date.now() - startTime;
      const errorMsg = error.response?.data?.error || error.message;
      handleFailure(currentRequestId, workerId, i, responseTime, errorMsg);
    }
    
    realtimeStats.totalRequests++;
    
    if (config.requestInterval > 0) {
      await new Promise(resolve => setTimeout(resolve, config.requestInterval));
    }
  }
};

const handleFailure = (requestId, workerId, index, responseTime, errorMessage) => {
  realtimeStats.failCount++;
  
  testResults.value.push({
    requestId,
    userId: workerId * 1000 + index,
    success: false,
    responseTime,
    errorMessage,
    timestamp: new Date().toLocaleTimeString()
  });
  
  let errorType = '其他错误';
  if (errorMessage.includes('库存')) {
    errorType = '库存不足';
  } else if (errorMessage.includes('超时') || errorMessage.includes('timeout')) {
    errorType = '请求超时';
  } else if (errorMessage.includes('限流') || errorMessage.includes('rate')) {
    errorType = '限流拦截';
  } else if (errorMessage.includes('500') || errorMessage.includes('server') || errorMessage.includes('服务器') || errorMessage.includes('内部')) {
    errorType = '系统内部错误';
  } else if (errorMessage.includes('秒杀')) {
    errorType = '秒杀失败';
  }
  
  const errorIndex = errorStats.value.findIndex(e => e.type === errorType);
  if (errorIndex !== -1) {
    errorStats.value[errorIndex].count++;
  }
  
  addLog('error', `Worker ${workerId}: 请求 ${requestId} 失败 - ${errorMessage}`);
};

const startStatsUpdate = () => {
  let lastTotal = 0;
  let lastTime = Date.now();
  
  statsTimer = setInterval(() => {
    if (testStatus.value !== 'running') return;
    
    const now = Date.now();
    const timeDiff = (now - lastTime) / 1000;
    const requestDiff = realtimeStats.totalRequests - lastTotal;
    
    realtimeStats.qps = requestDiff / timeDiff;
    
    lastTotal = realtimeStats.totalRequests;
    lastTime = now;
    
    updateRealtimeData();
    
    serverResources.cpu = Math.round(Math.min(95, 30 + Math.random() * 40));
    serverResources.memory = Math.round(Math.min(90, 40 + Math.random() * 30));
    serverResources.network = Math.round(Math.min(85, 20 + Math.random() * 50));
    serverResources.disk = Math.round(Math.min(60, 10 + Math.random() * 30));
    
    updateErrorStats();
  }, 1000);
};

const updateRealtimeData = () => {
  const recentResults = testResults.value.slice(-100);
  
  if (recentResults.length > 0) {
    const responseTimes = recentResults.map(r => r.responseTime).sort((a, b) => a - b);
    
    const avg = responseTimes.reduce((a, b) => a + b, 0) / responseTimes.length;
    const p90Index = Math.floor(responseTimes.length * 0.9);
    const p99Index = Math.floor(responseTimes.length * 0.99);
    
    responseTimeData.data.push({
      avg: avg || 0,
      p90: responseTimes[p90Index] || 0,
      p99: responseTimes[p99Index] || 0
    });
    
    if (responseTimeData.data.length > 10) {
      responseTimeData.data.shift();
    }
    
    throughputData.data.push(realtimeStats.qps);
    if (throughputData.data.length > 10) {
      throughputData.data.shift();
    }
    
    updateCharts();
  }
};

const updateErrorStats = () => {
  const total = errorStats.value.reduce((sum, e) => sum + e.count, 0);
  
  if (total > 0) {
    errorStats.value.forEach(e => {
      e.percentage = (e.count / total) * 100;
    });
  }
};

const pauseTest = () => {
  testStatus.value = 'paused';
  addLog('warning', '测试已暂停');
  showToast('测试暂停', '测试已暂停，可继续或终止');
};

const stopTest = () => {
  testStatus.value = 'completed';
  if (statsTimer) {
    clearInterval(statsTimer);
    statsTimer = null;
  }
  
  generateAnalysis();
  addLog('warning', '测试已终止');
  showToast('测试终止', `最终结果 - 成功: ${realtimeStats.successCount}, 失败: ${realtimeStats.failCount}`);
};

const generateAnalysis = () => {
  showAnalysis.value = true;
  
  const successResults = testResults.value.filter(r => r.success);
  const failResults = testResults.value.filter(r => !r.success);
  
  if (testResults.value.length > 0) {
    const responseTimes = testResults.value.map(r => r.responseTime);
    
    analysisMetrics.avgResponseTime = Math.round(responseTimes.reduce((a, b) => a + b, 0) / responseTimes.length);
    analysisMetrics.maxResponseTime = Math.max(...responseTimes);
    analysisMetrics.minResponseTime = Math.min(...responseTimes);
    analysisMetrics.avgThroughput = Math.round(realtimeStats.qps);
    analysisMetrics.successRate = ((realtimeStats.successCount / realtimeStats.totalRequests) * 100).toFixed(2);
    analysisMetrics.errorRate = ((realtimeStats.failCount / realtimeStats.totalRequests) * 100).toFixed(2);
  }
  
  bottlenecks.value = [];
  
  if (analysisMetrics.avgResponseTime > 500) {
    bottlenecks.value.push({
      level: 'high',
      title: '响应时间过长',
      description: `平均响应时间 ${analysisMetrics.avgResponseTime}ms 超过推荐值 500ms`
    });
  }
  
  if (parseFloat(analysisMetrics.errorRate) > 10) {
    bottlenecks.value.push({
      level: 'high',
      title: '错误率过高',
      description: `错误率 ${analysisMetrics.errorRate}% 超过推荐值 10%`
    });
  }
  
  if (serverResources.cpu > 80) {
    bottlenecks.value.push({
      level: 'medium',
      title: 'CPU 使用率过高',
      description: `CPU 使用率达到 ${serverResources.cpu}%，可能影响系统性能`
    });
  }
  
  if (serverResources.memory > 80) {
    bottlenecks.value.push({
      level: 'medium',
      title: '内存使用率过高',
      description: `内存使用率达到 ${serverResources.memory}%，建议优化内存管理`
    });
  }
  
  if (bottlenecks.value.length === 0) {
    bottlenecks.value.push({
      level: 'low',
      title: '系统表现良好',
      description: '未发现明显的性能瓶颈，系统运行正常'
    });
  }
  
  limitMetrics.concurrentCapacity = Math.round(realtimeStats.qps * 1.2);
  limitMetrics.concurrentPercent = Math.min(100, (realtimeStats.qps / 1000) * 100);
  limitMetrics.stability = (100 - parseFloat(analysisMetrics.errorRate)).toFixed(1);
  limitMetrics.resourceUsage = Math.round((serverResources.cpu + serverResources.memory) / 2);
};

const exportResults = (format) => {
  const data = testResults.value.map(r => ({
    '请求ID': r.requestId,
    '用户ID': r.userId,
    '状态': r.success ? '成功' : '失败',
    '响应时间(ms)': r.responseTime,
    '错误信息': r.errorMessage || '',
    '时间戳': r.timestamp
  }));
  
  if (format === 'csv') {
    const csv = [
      Object.keys(data[0]).join(','),
      ...data.map(row => Object.values(row).join(','))
    ].join('\n');
    
    downloadFile(csv, 'test-results.csv', 'text/csv');
  } else {
    const csv = [
      Object.keys(data[0]).join('\t'),
      ...data.map(row => Object.values(row).join('\t'))
    ].join('\n');
    
    downloadFile(csv, 'test-results.xls', 'application/vnd.ms-excel');
  }
  
  showToast('导出成功', `已导出 ${testResults.value.length} 条测试记录`);
};

const downloadFile = (content, filename, mimeType) => {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};

onUnmounted(() => {
  if (statsTimer) {
    clearInterval(statsTimer);
  }
});
</script>

<style scoped>
.performance-test-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding-bottom: 40px;
}

.page-header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 24px 0;
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-title h1 {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 8px 0;
}

.header-subtitle {
  font-size: 15px;
  color: #666;
  margin: 0;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
}

.status-badge.idle {
  background: #f0f0f0;
  color: #666;
}

.status-badge.running {
  background: #e3f2fd;
  color: #1976d2;
}

.status-badge.paused {
  background: #fff3e0;
  color: #f57c00;
}

.status-badge.completed {
  background: #e8f5e9;
  color: #388e3c;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

.test-layout {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 24px;
}

.config-panel,
.monitor-panel {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.panel-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.panel-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.panel-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
  display: flex;
  align-items: center;
  gap: 10px;
}

.panel-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
}

.panel-body {
  padding: 24px;
}

.config-section {
  margin-bottom: 24px;
}

.config-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #666;
  margin-bottom: 16px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-hint {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: #999;
}

.control-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 24px;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-lg {
  padding: 14px 28px;
  font-size: 15px;
}

.btn-sm {
  padding: 8px 16px;
  font-size: 13px;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-warning {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.btn-warning:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.4);
}

.btn-danger {
  background: linear-gradient(135deg, #f5576c 0%, #f093fb 100%);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(245, 87, 108, 0.4);
}

.btn-outline {
  background: transparent;
  border: 2px solid #e5e7eb;
  color: #666;
}

.btn-outline:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.2;
}

.stat-value.success {
  color: #10b981;
}

.stat-value.error {
  color: #f5576c;
}

.stat-label {
  display: block;
  font-size: 13px;
  color: #999;
  margin-top: 4px;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.chart-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.chart-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.chart-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
}

.chart-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chart-select {
  padding: 6px 12px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 13px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
}

.chart-select:focus {
  outline: none;
  border-color: #667eea;
}

.chart-legend {
  display: flex;
  gap: 20px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #666;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.chart-body {
  padding: 24px;
}

.line-chart,
.bar-chart {
  position: relative;
  width: 100%;
}

.chart-labels {
  display: flex;
  justify-content: space-between;
  padding: 8px 0 0;
  font-size: 12px;
  color: #999;
}

.tooltip {
  position: absolute;
  background: #1a1a2e;
  color: white;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 13px;
  z-index: 100;
  pointer-events: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.tooltip::after {
  content: '';
  position: absolute;
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-top: 6px solid #1a1a2e;
}

.tooltip-title {
  font-weight: 600;
  margin-bottom: 4px;
  font-size: 12px;
  opacity: 0.9;
}

.tooltip-value {
  font-size: 18px;
  font-weight: 700;
}

.error-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
}

.error-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.error-info {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 140px;
}

.error-type {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  min-width: 80px;
}

.error-count {
  font-size: 14px;
  font-weight: 600;
  color: #666;
  min-width: 50px;
  text-align: right;
}

.error-bar-wrapper {
  flex: 1;
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.error-bar {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.error-percentage {
  font-size: 14px;
  font-weight: 600;
  color: #666;
  min-width: 60px;
  text-align: right;
}

.error-description {
  display: grid;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.error-desc-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
}

.desc-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 700;
  flex-shrink: 0;
}

.desc-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.desc-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.desc-text {
  font-size: 13px;
  color: #666;
}

.resource-gauges {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.gauge-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.gauge-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.gauge-circle {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.gauge-inner {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: white;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.gauge-value {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1;
}

.gauge-unit {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}

.gauge-label {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.resource-info {
  position: relative;
}

.info-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #667eea;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.info-icon:hover {
  transform: scale(1.1);
}

.info-tooltip {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(-8px);
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  min-width: 280px;
  z-index: 100;
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
}

.info-tooltip.show {
  opacity: 1;
  visibility: visible;
}

.info-tooltip::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-top: 8px solid white;
}

.info-title {
  font-size: 15px;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 12px;
}

.info-desc {
  font-size: 13px;
  color: #666;
  line-height: 1.6;
  margin-bottom: 12px;
}

.info-range,
.info-impact {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 8px;
}

.range-label,
.impact-label {
  font-size: 12px;
  font-weight: 600;
  color: #667eea;
  min-width: 70px;
}

.range-value,
.impact-desc {
  font-size: 13px;
  color: #666;
  line-height: 1.5;
  flex: 1;
}

.results-panel {
  margin-top: 24px;
}

.results-table-wrapper {
  overflow-x: auto;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
}

.results-table {
  width: 100%;
  border-collapse: collapse;
}

.results-table th,
.results-table td {
  padding: 12px 16px;
  text-align: left;
  font-size: 13px;
}

.results-table th {
  background: #fafafa;
  font-weight: 600;
  color: #666;
  border-bottom: 2px solid #f0f0f0;
}

.results-table td {
  border-bottom: 1px solid #f0f0f0;
  color: #333;
}

.success-row {
  background: rgba(16, 185, 129, 0.05);
}

.error-row {
  background: rgba(245, 87, 108, 0.05);
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-tag.success {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.status-tag.error {
  background: rgba(245, 87, 108, 0.1);
  color: #f5576c;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 20px;
}

.page-info {
  font-size: 14px;
  color: #666;
}

.analysis-panel {
  margin-top: 24px;
}

.analysis-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.analysis-card {
  background: #fafafa;
  border-radius: 12px;
  padding: 20px;
}

.analysis-card.full-width {
  grid-column: 1 / -1;
}

.analysis-card h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
}

.metric-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.metric-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: white;
  border-radius: 8px;
}

.metric-label {
  font-size: 14px;
  color: #666;
}

.metric-value {
  font-size: 16px;
  font-weight: 700;
  color: #1a1a2e;
}

.metric-value.success {
  color: #10b981;
}

.metric-value.error {
  color: #f5576c;
}

.bottleneck-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.bottleneck-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 8px;
}

.bottleneck-item.high {
  border-left: 4px solid #f5576c;
}

.bottleneck-item.medium {
  border-left: 4px solid #f59e0b;
}

.bottleneck-item.low {
  border-left: 4px solid #10b981;
}

.bottleneck-icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.bottleneck-item.high .bottleneck-icon {
  color: #f5576c;
}

.bottleneck-item.medium .bottleneck-icon {
  color: #f59e0b;
}

.bottleneck-item.low .bottleneck-icon {
  color: #10b981;
}

.bottleneck-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.bottleneck-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a2e;
}

.bottleneck-desc {
  font-size: 13px;
  color: #666;
}

.limit-metrics {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.limit-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.limit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.limit-label {
  font-size: 14px;
  font-weight: 600;
  color: #666;
}

.limit-value {
  font-size: 20px;
  font-weight: 700;
  color: #1a1a2e;
}

.limit-bar {
  height: 10px;
  background: #f0f0f0;
  border-radius: 5px;
  overflow: hidden;
}

.limit-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 5px;
  transition: width 0.5s ease;
}

.limit-fill.success {
  background: linear-gradient(90deg, #10b981 0%, #059669 100%);
}

.limit-fill.warning {
  background: linear-gradient(90deg, #f59e0b 0%, #d97706 100%);
}

.logs-panel {
  margin-top: 24px;
}

.logs-container {
  max-height: 300px;
  overflow-y: auto;
  background: #1a1a2e;
  border-radius: 8px;
  padding: 16px;
}

.log-item {
  display: flex;
  gap: 12px;
  padding: 8px 0;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
}

.log-item:not(:last-child) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.log-time {
  color: #999;
  min-width: 80px;
}

.log-level {
  font-weight: 700;
  min-width: 60px;
}

.log-item.info .log-level {
  color: #667eea;
}

.log-item.success .log-level {
  color: #10b981;
}

.log-item.warning .log-level {
  color: #f59e0b;
}

.log-item.error .log-level {
  color: #f5576c;
}

.log-message {
  color: #e5e7eb;
  flex: 1;
}

.toast-container {
  position: fixed;
  top: 24px;
  right: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.toast-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.toast-container.success .toast-icon {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.toast-container.error .toast-icon {
  background: rgba(245, 87, 108, 0.1);
  color: #f5576c;
}

.toast-content {
  flex: 1;
}

.toast-title {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
}

.toast-message {
  margin: 0;
  font-size: 14px;
  color: #666;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .analysis-grid {
    grid-template-columns: 1fr;
  }
  
  .resource-gauges {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .limit-metrics {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .test-layout {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
