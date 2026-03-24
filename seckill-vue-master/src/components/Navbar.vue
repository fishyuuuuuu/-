<template>
  <nav class="navbar">
    <div class="container">
      <div class="navbar-brand">
        <h1 class="logo">秒杀商城</h1>
      </div>
      <div class="navbar-nav">
        <router-link to="/main" class="nav-item" active-class="active">首页</router-link>
        <router-link to="/category" class="nav-item" active-class="active">全部商品</router-link>
        <router-link to="/upcoming" class="nav-item" active-class="active">即将开始</router-link>
        <router-link to="/order" class="nav-item" active-class="active">我的订单</router-link>
        <router-link to="/profile" class="nav-item" active-class="active">个人中心</router-link>
      </div>
      <div class="navbar-user">
        <span class="user-info" v-if="userId">欢迎，用户{{ userId }}</span>
        <button class="logout-btn" @click="logout" v-if="userId">退出登录</button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const userId = ref(localStorage.getItem('user_id'));

const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user_id');
  userId.value = null;
  alert('已退出登录');
  router.push('/login');
};
</script>

<style scoped>
.navbar {
  background-color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
}

.navbar-brand .logo {
  font-size: 22px;
  font-weight: bold;
  color: #e63946;
  margin: 0;
}

.navbar-nav {
  display: flex;
  gap: 24px;
}

.nav-item {
  text-decoration: none;
  color: #333;
  font-weight: 500;
  transition: color 0.3s;
  position: relative;
}

.nav-item:hover {
  color: #e63946;
}

.nav-item.active {
  color: #e63946;
}

.nav-item.active::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: #e63946;
}

.navbar-user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  font-weight: 500;
  color: #333;
}

.logout-btn {
  background-color: #f5f5f5;
  color: #333;
  border: 1px solid #ddd;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.logout-btn:hover {
  background-color: #e63946;
  color: white;
  border-color: #e63946;
}

@media (max-width: 768px) {
  .container {
    flex-direction: column;
    align-items: flex-start;
    height: auto;
    padding: 10px 16px;
    gap: 8px;
  }

  .navbar-nav {
    flex-wrap: wrap;
    gap: 12px;
  }
}
</style>
