<template>
  <div class="high-concurrency-container">
    <h2 class="page-title">高并发操作可视化</h2>
    
    <!-- 测试配置区 -->
    <div class="test-config">
      <h3>测试配置</h3>
      <div class="config-form">
        <div class="form-group">
          <label for="totalRequests">总请求数：</label>
          <input type="number" id="totalRequests" v-model.number="testConfig.totalRequests" min="1" max="1000">
        </div>
        <div class="form-group">
          <label for="concurrency">并发数：</label>
          <input type="number" id="concurrency" v-model.number="testConfig.concurrency" min="1" max="100">
        </div>
        <div class="form-group">
          <label for="productId">商品ID：</label>
          <input type="number" id="productId" v-model.number="testConfig.productId" min="1">
        </div>
        <button class="start-test-btn" @click="startTest" :disabled="isTesting">
          {{ isTesting ? '测试中...' : '开始测试' }}
        </button>
      </div>
    </div>

    <!-- 实时监控区 -->
    <div class="realtime-monitor">
      <h3>实时监控</h3>
      <div class="monitor-cards">
        <div class="monitor-card">
          <div class="card-title">总请求数</div>
          <div class="card-value">{{ stats.totalRequests }}</div>
        </div>
        <div class="monitor-card">
          <div class="card-title">成功数</div>
          <div class="card-value success">{{ stats.successCount }}</div>
        </div>
        <div class="monitor-card">
          <div class="card-title">失败数</div>
          <div class="card-value error">{{ stats.errorCount }}</div>
        </div>
        <div class="monitor-card">
          <div class="card-title">QPS</div>
          <div class="card-value">{{ stats.qps.toFixed(2) }}</div>
        </div>
        <div class="monitor-card">
          <div class="card-title">平均响应时间</div>
          <div class="card-value">{{ stats.avgResponseTime.toFixed(2) }}ms</div>
        </div>
        <div class="monitor-card">
          <div class="card-title">P99响应时间</div>
          <div class="card-value">{{ stats.p99ResponseTime.toFixed(2) }}ms</div>
        </div>
      </div>
    </div>

    <!-- 响应时间分布图 -->
    <div class="chart-section">
      <h3>响应时间分布</h3>
      <div class="chart-container">
        <canvas id="responseTimeChart"></canvas>
      </div>
    </div>

    <!-- 成功率趋势图 -->
    <div class="chart-section">
      <h3>成功率趋势</h3>
      <div class="chart-container">
        <canvas id="successRateChart"></canvas>
      </div>
    </div>

    <!-- 测试日志 -->
    <div class="test-logs">
      <h3>测试日志</h3>
      <div class="logs-container">
        <div v-for="(log, index) in testLogs" :key="index" class="log-item" :class="{ 'success': log.type === 'success', 'error': log.type === 'error' }">
          <span class="log-time">{{ log.time }}</span>
          <span class="log-content">{{ log.content }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Chart from 'chart.js/auto';

export default {
  name: 'HighConcurrency',
  data() {
    return {
      testConfig: {
        totalRequests: 100,
        concurrency: 10,
        productId: 1
      },
      isTesting: false,
      stats: {
        totalRequests: 0,
        successCount: 0,
        errorCount: 0,
        qps: 0,
        avgResponseTime: 0,
        p99ResponseTime: 0
      },
      testLogs: [],
      responseTimeChart: null,
      successRateChart: null,
      responseTimes: [],
      successRates: [],
      timestamps: []
    };
  },
  mounted() {
    this.initCharts();
  },
  methods: {
    initCharts() {
      // 初始化响应时间分布图
      const responseTimeCtx = document.getElementById('responseTimeChart').getContext('2d');
      this.responseTimeChart = new Chart(responseTimeCtx, {
        type: 'bar',
        data: {
          labels: ['0-50ms', '50-100ms', '100-200ms', '200-500ms', '>500ms'],
          datasets: [{
            label: '请求数',
            data: [0, 0, 0, 0, 0],
            backgroundColor: [
              'rgba(75, 192, 192, 0.6)',
              'rgba(54, 162, 235, 0.6)',
              'rgba(255, 206, 86, 0.6)',
              'rgba(255, 159, 64, 0.6)',
              'rgba(255, 99, 132, 0.6)'
            ],
            borderColor: [
              'rgba(75, 192, 192, 1)',
              'rgba(54, 162, 235, 1)',
              'rgba(255, 206, 86, 1)',
              'rgba(255, 159, 64, 1)',
              'rgba(255, 99, 132, 1)'
            ],
            borderWidth: 1
          }]
        },
        options: {
          responsive: true,
          scales: {
            y: {
              beginAtZero: true,
              title: {
                display: true,
                text: '请求数'
              }
            },
            x: {
              title: {
                display: true,
                text: '响应时间'
              }
            }
          }
        }
      });

      // 初始化成功率趋势图
      const successRateCtx = document.getElementById('successRateChart').getContext('2d');
      this.successRateChart = new Chart(successRateCtx, {
        type: 'line',
        data: {
          labels: [],
          datasets: [{
            label: '成功率',
            data: [],
            borderColor: 'rgba(75, 192, 192, 1)',
            backgroundColor: 'rgba(75, 192, 192, 0.2)',
            tension: 0.1
          }]
        },
        options: {
          responsive: true,
          scales: {
            y: {
              beginAtZero: true,
              max: 100,
              title: {
                display: true,
                text: '成功率 (%)'
              }
            },
            x: {
              title: {
                display: true,
                text: '时间'
              }
            }
          }
        }
      });
    },
    startTest() {
      this.isTesting = true;
      this.testLogs = [];
      this.stats = {
        totalRequests: 0,
        successCount: 0,
        errorCount: 0,
        qps: 0,
        avgResponseTime: 0,
        p99ResponseTime: 0
      };
      this.responseTimes = [];
      this.successRates = [];
      this.timestamps = [];

      const startTime = Date.now();
      let completedRequests = 0;
      let totalResponseTime = 0;

      // 并发测试函数
      const sendRequest = (workerId) => {
        return new Promise((resolve) => {
          const start = Date.now();
          axios.post('/api/product/seckill', {
            productId: this.testConfig.productId,
            captchaId: '',
            captchaStr: ''
          }, {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          }).then(response => {
            const end = Date.now();
            const responseTime = end - start;
            totalResponseTime += responseTime;
            this.responseTimes.push(responseTime);
            this.stats.successCount++;
            this.addLog('success', `Worker ${workerId}: 请求成功，响应时间: ${responseTime}ms`);
          }).catch(error => {
            const end = Date.now();
            const responseTime = end - start;
            totalResponseTime += responseTime;
            this.responseTimes.push(responseTime);
            this.stats.errorCount++;
            this.addLog('error', `Worker ${workerId}: 请求失败，响应时间: ${responseTime}ms, 错误: ${error.message}`);
          }).finally(() => {
            completedRequests++;
            this.stats.totalRequests = completedRequests;
            this.stats.avgResponseTime = totalResponseTime / completedRequests;
            
            // 计算 P99 响应时间
            if (this.responseTimes.length > 0) {
              const sortedTimes = [...this.responseTimes].sort((a, b) => a - b);
              const p99Index = Math.floor(sortedTimes.length * 0.99);
              this.stats.p99ResponseTime = sortedTimes[p99Index] || 0;
            }

            // 计算 QPS
            const elapsedTime = (Date.now() - startTime) / 1000;
            this.stats.qps = completedRequests / elapsedTime;

            // 更新成功率趋势
            if (completedRequests % 10 === 0) {
              const successRate = (this.stats.successCount / completedRequests) * 100;
              this.successRates.push(successRate);
              this.timestamps.push(new Date().toLocaleTimeString());
              this.updateCharts();
            }

            resolve();
          });
        });
      };

      // 启动并发测试
      const workers = this.testConfig.concurrency;
      const requestsPerWorker = Math.ceil(this.testConfig.totalRequests / workers);

      const workerPromises = [];
      for (let i = 0; i < workers; i++) {
        const workerRequests = i < this.testConfig.totalRequests % workers ? requestsPerWorker + 1 : requestsPerWorker;
        for (let j = 0; j < workerRequests; j++) {
          workerPromises.push(sendRequest(i));
        }
      }

      // 等待所有测试完成
      Promise.all(workerPromises).then(() => {
        this.isTesting = false;
        this.addLog('info', `测试完成，总请求数: ${this.stats.totalRequests}, 成功数: ${this.stats.successCount}, 失败数: ${this.stats.errorCount}, 平均响应时间: ${this.stats.avgResponseTime.toFixed(2)}ms, P99响应时间: ${this.stats.p99ResponseTime.toFixed(2)}ms, QPS: ${this.stats.qps.toFixed(2)}`);
        this.updateCharts();
      });
    },
    addLog(type, content) {
      const time = new Date().toLocaleTimeString();
      this.testLogs.push({ type, content, time });
      // 限制日志数量
      if (this.testLogs.length > 100) {
        this.testLogs.shift();
      }
    },
    updateCharts() {
      // 更新响应时间分布图
      const timeRanges = [0, 0, 0, 0, 0];
      this.responseTimes.forEach(time => {
        if (time < 50) {
          timeRanges[0]++;
        } else if (time < 100) {
          timeRanges[1]++;
        } else if (time < 200) {
          timeRanges[2]++;
        } else if (time < 500) {
          timeRanges[3]++;
        } else {
          timeRanges[4]++;
        }
      });

      this.responseTimeChart.data.datasets[0].data = timeRanges;
      this.responseTimeChart.update();

      // 更新成功率趋势图
      this.successRateChart.data.labels = this.timestamps;
      this.successRateChart.data.datasets[0].data = this.successRates;
      this.successRateChart.update();
    }
  }
};
</script>

<style scoped>
.high-concurrency-container {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: 100vh;
}

.page-title {
  font-size: 24px;
  margin-bottom: 30px;
  color: #333;
  text-align: center;
}

.test-config {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.test-config h3 {
  margin-bottom: 15px;
  color: #555;
}

.config-form {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  align-items: center;
}

.form-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.form-group label {
  font-weight: 500;
  color: #666;
}

.form-group input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 120px;
}

.start-test-btn {
  background-color: #4CAF50;
  color: white;
  padding: 8px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.3s;
}

.start-test-btn:hover {
  background-color: #45a049;
}

.start-test-btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.realtime-monitor {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.realtime-monitor h3 {
  margin-bottom: 15px;
  color: #555;
}

.monitor-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 15px;
}

.monitor-card {
  background-color: #f9f9f9;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.card-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.card-value.success {
  color: #4CAF50;
}

.card-value.error {
  color: #f44336;
}

.chart-section {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.chart-section h3 {
  margin-bottom: 15px;
  color: #555;
}

.chart-container {
  height: 300px;
}

.test-logs {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.test-logs h3 {
  margin-bottom: 15px;
  color: #555;
}

.logs-container {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  background-color: #f9f9f9;
}

.log-item {
  padding: 5px 10px;
  margin-bottom: 5px;
  border-radius: 4px;
  font-size: 14px;
}

.log-item.success {
  background-color: rgba(76, 175, 80, 0.1);
  border-left: 3px solid #4CAF50;
}

.log-item.error {
  background-color: rgba(244, 67, 54, 0.1);
  border-left: 3px solid #f44336;
}

.log-time {
  font-weight: bold;
  margin-right: 10px;
  color: #666;
}

.log-content {
  color: #333;
}

@media (max-width: 768px) {
  .config-form {
    flex-direction: column;
    align-items: flex-start;
  }

  .form-group {
    width: 100%;
  }

  .form-group input {
    flex: 1;
  }

  .monitor-cards {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
