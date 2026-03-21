import { createRouter, createWebHistory } from 'vue-router'
import MainView from './components/Main.vue'
import LoginView from './components/Login.vue'
import StressTestView from './components/StressTest.vue'
import SecurityVisualizationView from './components/SecurityVisualization.vue'
import OrderView from './components/Order.vue'
import CategoryView from './components/Category.vue'
import ProfileView from './components/Profile.vue'
import UpcomingView from './components/Upcoming.vue'
import AdminView from './components/Admin.vue'

const routes = [
    { path: '/', redirect: '/login' },
    { path: '/main', component: MainView },
    { path: '/login', component: LoginView },
    { path: '/stress-test', component: StressTestView },
    { path: '/security', component: SecurityVisualizationView },
    { path: '/order', component: OrderView },
    { path: '/category', component: CategoryView },
    { path: '/profile', component: ProfileView },
    { path: '/upcoming', component: UpcomingView },
    { path: '/seckill/:id', redirect: '/main' },
    { 
        path: '/admin', 
        component: AdminView,
        meta: { requiresAdmin: true }
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, from, next) => {
    const adminToken = localStorage.getItem('admin_token');
    
    if (to.meta.requiresAdmin) {
        if (!adminToken) {
            next('/login');
        } else {
            next();
        }
    } else {
        next();
    }
})

export default router