// 测试脚本：验证行为分析系统
const axios = require('axios');

async function testBehaviorAnalysis() {
  const baseURL = 'http://localhost:8081/api';
  const productID = 1;
  const userID = 1000000;
  const ip = '192.168.1.100';

  console.log('测试1：直接秒杀（不浏览页面）');
  try {
    // 直接获取验证码
    const captchaResponse = await axios.get(`${baseURL}/captcha`, {
      headers: {
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });

    const captchaID = captchaResponse.headers['x-captcha-id'];
    console.log('验证码ID:', captchaID);

    // 直接执行秒杀
    const seckillResponse = await axios.post(`${baseURL}/product/seckill`, {
      productId: productID,
      userId: userID,
      captchaId: captchaID,
      captchaStr: '1234'
    }, {
      headers: {
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });

    console.log('秒杀结果:', seckillResponse.data);
  } catch (error) {
    console.log('秒杀失败:', error.response?.data || error.message);
  }

  console.log('\n测试2：先浏览页面，再秒杀');
  try {
    // 先浏览页面
    await axios.get(`${baseURL}/product/list`, {
      headers: {
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });
    console.log('浏览页面成功');

    // 等待1秒
    await new Promise(resolve => setTimeout(resolve, 1000));

    // 获取验证码
    const captchaResponse = await axios.get(`${baseURL}/captcha`, {
      headers: {
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });

    const captchaID = captchaResponse.headers['x-captcha-id'];
    console.log('验证码ID:', captchaID);

    // 执行秒杀
    const seckillResponse = await axios.post(`${baseURL}/product/seckill`, {
      productId: productID,
      userId: userID + 1,
      captchaId: captchaID,
      captchaStr: '1234'
    }, {
      headers: {
        'X-Forwarded-For': ip,
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
      }
    });

    console.log('秒杀结果:', seckillResponse.data);
  } catch (error) {
    console.log('秒杀失败:', error.response?.data || error.message);
  }

  console.log('\n测试3：多次快速秒杀（模拟脚本）');
  for (let i = 0; i < 5; i++) {
    try {
      // 获取验证码
      const captchaResponse = await axios.get(`${baseURL}/captcha`, {
        headers: {
          'X-Forwarded-For': ip,
          'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
        }
      });

      const captchaID = captchaResponse.headers['x-captcha-id'];
      console.log(`验证码ID (${i}):`, captchaID);

      // 执行秒杀
      const seckillResponse = await axios.post(`${baseURL}/product/seckill`, {
        productId: productID,
        userId: userID + 2 + i,
        captchaId: captchaID,
        captchaStr: '1234'
      }, {
        headers: {
          'X-Forwarded-For': ip,
          'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
        }
      });

      console.log(`秒杀结果 (${i}):`, seckillResponse.data);
    } catch (error) {
      console.log(`秒杀失败 (${i}):`, error.response?.data || error.message);
    }
  }
}

testBehaviorAnalysis().catch(console.error);
