<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const email = ref('');
const password = ref('');
const errorMsg = ref('');
const authStore = useAuthStore();
const router = useRouter();

const handleLogin = async () => {
  const success = await authStore.login(email.value, password.value);
  if (success) {
    router.push('/'); // Redirect to Map
  } else {
    errorMsg.value = "Invalid email or password";
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-indigo-50 to-slate-100 p-4">
    <div class="bg-white/80 backdrop-blur-lg p-8 sm:p-10 rounded-2xl shadow-2xl border border-white w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent mb-2">
          ShiftKerja
        </h1>
        <p class="text-slate-600 text-sm">Sign in to your account</p>
      </div>
      
      <div v-if="errorMsg" class="bg-red-50 border border-red-200 text-red-700 p-3 mb-6 rounded-xl text-sm flex items-center gap-2">
        <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ errorMsg }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Email</label>
          <input 
            v-model="email" 
            type="email" 
            placeholder="you@example.com"
            required
            class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Password</label>
          <input 
            v-model="password" 
            type="password" 
            placeholder="••••••••"
            required
            class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          />
        </div>
        
        <button 
          type="submit"
          class="w-full bg-gradient-to-r from-blue-600 to-indigo-600 text-white py-3 rounded-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-200 font-semibold shadow-lg shadow-blue-500/30"
        >
          Sign In
        </button>
      </form>
      
      <p class="mt-6 text-center text-sm text-slate-600">
        Don't have an account? 
        <router-link to="/register" class="font-semibold text-blue-600 hover:text-blue-700 transition-colors">
          Create one
        </router-link>
      </p>
    </div>
  </div>
</template>