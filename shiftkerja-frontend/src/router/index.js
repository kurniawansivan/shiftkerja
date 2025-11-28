import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import MapView from '../components/MapData.vue'
import BusinessDashboard from '../views/BusinessDashboard.vue'
import WorkerDashboard from '../views/WorkerDashboard.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/',
      name: 'home',
      component: MapView,
      meta: { requiresAuth: true }
    },
    {
      path: '/business/dashboard',
      name: 'business-dashboard',
      component: BusinessDashboard,
      meta: { requiresAuth: true, role: 'business' }
    },
    {
      path: '/worker/dashboard',
      name: 'worker-dashboard',
      component: WorkerDashboard,
      meta: { requiresAuth: true, role: 'worker' }
    }
  ]
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  const role = localStorage.getItem('role');
  
  if (to.meta.requiresAuth && !token) {
    next('/login');
  } else if (to.meta.role && to.meta.role !== role) {
    // Role-based protection
    next('/');
  } else {
    next();
  }
});

export default router