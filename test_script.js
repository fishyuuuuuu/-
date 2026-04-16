// 测试脚本：模拟不同类型的用户行为
const axios = require('axios');
const fs = require('fs');
const path = require('path');

// 配置
const config = {
  baseURL: 'http://localhost:8081/api',
  productID: 1,
  stock: 10,
  concurrency: 200,
  testDuration: 60000, // 60秒
  outputFile: path.join(__dirname, 'test_results.json'),
  testPassword: 'test123456'
};

// 测试结果
const testResults = {
  startTime: new Date().toISOString(),
  normalUsers: {
    total: 0,
    successful: 0,
    failed: 0,
    errorMessages: []
  },
  systemMetrics: {
    responseTimes: [],
    successRate: 0
  }
};

// 生成随机用户ID
function generateUserID() {
  return Math.floor(Math.random() * 1000000).toString();
}

// 生成随机IP（模拟不同用户）
function generateRandomIP() {
  return `${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}`;
}

// 获取验证码ID（测试模式下验证码固定使用 1234）
async function getCaptchaID(ip) {
  const response = await axios.get(`${config.baseURL}/captcha`, {
    headers: {
      'X-Forwarded-For': ip,
      'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    }
  });

  return response.headers['x-captcha-id'];
}

// 为测试用户自动注册并登录，获取有效 token
async function getAuthToken(username, phone, ip) {
  const captchaID = await getCaptchaID(ip);

  try {
    await axios.post(`${config.baseURL}/user/register`, {
      username,
      phone,
      password: config.testPassword,
      captcha_id: captchaID,
      captcha_code: '1234'
    }, {
      headers: {
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });
  } catch (error) {
    const status = error.response?.status;
    if (status !== 409) {
      throw error;
    }
  }

  const loginCaptchaID = await getCaptchaID(ip);
  const loginResponse = await axios.post(`${config.baseURL}/user/login`, {
    phone,
    password: config.testPassword,
    captcha_id: loginCaptchaID,
    captcha_code: '1234'
  }, {
    headers: {
      'X-Forwarded-For': ip,
      'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    }
  });

  return loginResponse.data.token;
}

// 模拟正常用户行为
async function simulateNormalUser(userID, ip) {
  try {
    const username = `normal_${userID}`;
    const phone = `13${userID.toString().padStart(9, '0').slice(-9)}`;
    const token = await getAuthToken(username, phone, ip);

    // 1. 先浏览商品页面
    await axios.get(`${config.baseURL}/product/list`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });

    // 2. 随机延迟（模拟用户思考时间）
    await new Promise(resolve => setTimeout(resolve, Math.random() * 1000 + 500));

    // 3. 获取验证码
    const captchaID = await getCaptchaID(ip);

    // 4. 执行秒杀（使用测试验证码）
    const startTime = Date.now();
    const seckillResponse = await axios.post(`${config.baseURL}/product/seckill`, {
      productId: config.productID,
      captchaId: captchaID,
      captchaStr: '1234' // 测试模式验证码
    }, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });

    const responseTime = Date.now() - startTime;
    testResults.systemMetrics.responseTimes.push(responseTime);

    if (seckillResponse.data.message === "秒杀成功") {
      testResults.normalUsers.successful++;
      console.log(`正常用户 ${userID} 秒杀成功`);
    } else {
      testResults.normalUsers.failed++;
      testResults.normalUsers.errorMessages.push(seckillResponse.data.message || seckillResponse.data.error);
      console.log(`正常用户 ${userID} 秒杀失败: ${seckillResponse.data.message || seckillResponse.data.error}`);
    }
  } catch (error) {
    testResults.normalUsers.failed++;
    const message = error.response?.data?.error || error.response?.data?.message || error.message;
    testResults.normalUsers.errorMessages.push(message);
    console.log(`正常用户 ${userID} 秒杀失败: ${message}`);
  } finally {
    testResults.normalUsers.total++;
  }
}

// 重置测试数据
async function resetTestData() {
  try {
    // 清理测试数据
    await axios.post(`${config.baseURL}/product/clear-test-data`);
    // 重置库存
    await axios.post(`${config.baseURL}/product/reset-stock`, {
      productId: config.productID,
      stock: config.stock
    });
    console.log('测试数据重置成功');
  } catch (error) {
    console.error('测试数据重置失败:', error.message);
  }
}

// 执行测试
async function runTest() {
  console.log('开始测试...');
  
  // 重置测试数据
  await resetTestData();

  // 并发执行测试
  const promises = [];

  // 全部按正常用户路径执行并发测试
  for (let i = 0; i < config.concurrency; i++) {
    const userID = `normal-${generateUserID()}`;
    const ip = generateRandomIP();
    promises.push(simulateNormalUser(userID, ip));
  }

  // 等待所有测试完成
  await Promise.all(promises);

  // 计算成功率
  const totalRequests = testResults.normalUsers.total;
  const totalSuccessful = testResults.normalUsers.successful;
  testResults.systemMetrics.successRate = totalRequests > 0 ? (totalSuccessful / totalRequests) * 100 : 0;

  // 计算平均响应时间
  if (testResults.systemMetrics.responseTimes.length > 0) {
    const sum = testResults.systemMetrics.responseTimes.reduce((a, b) => a + b, 0);
    testResults.systemMetrics.averageResponseTime = sum / testResults.systemMetrics.responseTimes.length;
  }

  // 结束时间
  testResults.endTime = new Date().toISOString();

  // 保存测试结果
  fs.writeFileSync(config.outputFile, JSON.stringify(testResults, null, 2));
  console.log(`测试完成，结果已保存到 ${config.outputFile}`);

  // 打印测试结果
  console.log('\n测试结果:');
  console.log(`总请求数: ${totalRequests}`);
  console.log(`成功数: ${totalSuccessful}`);
  console.log(`成功率: ${testResults.systemMetrics.successRate.toFixed(2)}%`);
  console.log(`平均响应时间: ${testResults.systemMetrics.averageResponseTime?.toFixed(2) || 0}ms`);
  console.log('\n正常用户:');
  console.log(`总请求: ${testResults.normalUsers.total}`);
  console.log(`成功: ${testResults.normalUsers.successful}`);
  console.log(`失败: ${testResults.normalUsers.failed}`);
}

// 检查是否安装了axios
if (!fs.existsSync(path.join(__dirname, 'node_modules', 'axios'))) {
  console.log('请先安装axios: npm install axios');
  process.exit(1);
}

// 运行测试
runTest().catch(console.error);
