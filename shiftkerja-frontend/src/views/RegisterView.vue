<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const email = ref('');
const password = ref('');
const fullName = ref('');
const role = ref('worker');
const errorMsg = ref('');
const successMsg = ref('');
const authStore = useAuthStore();
const router = useRouter();

const handleRegister = async () => {
  errorMsg.value = '';
  successMsg.value = '';
  
  try {
    const res = await fetch('http://localhost:8080/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: email.value,
        password: password.value,
        full_name: fullName.value,
        role: role.value
      })
    });

    if (!res.ok) {
      const data = await res.text();
      throw new Error(data || 'Registration failed');
    }

    successMsg.value = 'Registration successful! Redirecting to login...';
    
    // Auto-login after 1.5 seconds
    setTimeout(async () => {
      const loginSuccess = await authStore.login(email.value, password.value);
      if (loginSuccess) {
        router.push('/');
      }
    }, 1500);

  } catch (err) {
    errorMsg.value = err.message || 'Registration failed';
  }
};

const goToLogin = () => {
  router.push('/login');
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-indigo-50 to-slate-100 p-4">
    <div class="bg-white/80 backdrop-blur-lg p-8 sm:p-10 rounded-2xl shadow-2xl border border-white w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent mb-2">
          Join ShiftKerja
        </h1>
        <p class="text-slate-600 text-sm">Create your account to get started</p>
      </div>
      
      <!-- Error Message -->
      <div v-if="errorMsg" class="bg-red-50 border border-red-200 text-red-700 p-3 mb-4 rounded-xl text-sm flex items-center gap-2">
        <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ errorMsg }}
      </div>

      <!-- Success Message -->
      <div v-if="successMsg" class="bg-green-50 border border-green-200 text-green-700 p-3 mb-4 rounded-xl text-sm flex items-center gap-2">
        <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ successMsg }}
      </div>

      <!-- Form -->
      <form @submit.prevent="handleRegister" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Full Name</label>
          <input 
            v-model="fullName" 
            type="text" 
            placeholder="John Doe" 
            required
            class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          />
        </div>

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
            placeholder="Minimum 6 characters" 
            required
            minlength="6"
            class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-3">I am a:</label>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <label 
              :class="role === 'worker' ? 'border-blue-500 bg-blue-50' : 'border-slate-200 bg-white'"
              class="flex items-center p-4 border-2 rounded-xl cursor-pointer transition-all hover:border-blue-300"
            >
              <input 
                type="radio" 
                v-model="role" 
                value="worker" 
                class="w-4 h-4 text-blue-600 focus:ring-blue-500"
              />
              <div class="ml-3">
                <p class="font-semibold text-slate-800">Worker</p>
                <p class="text-xs text-slate-500">Find shifts</p>
              </div>
            </label>
            
            <label 
              :class="role === 'business' ? 'border-blue-500 bg-blue-50' : 'border-slate-200 bg-white'"
              class="flex items-center p-4 border-2 rounded-xl cursor-pointer transition-all hover:border-blue-300"
            >
              <input 
                type="radio" 
                v-model="role" 
                value="business" 
                class="w-4 h-4 text-blue-600 focus:ring-blue-500"
              />
              <div class="ml-3">
                <p class="font-semibold text-slate-800">Business</p>
                <p class="text-xs text-slate-500">Post shifts</p>
              </div>
            </label>
          </div>
        </div>
        
        <button 
          type="submit"
          class="w-full bg-gradient-to-r from-blue-600 to-indigo-600 text-white py-3 rounded-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-200 font-semibold shadow-lg shadow-blue-500/30"
        >
          Create Account
        </button>
      </form>
      
      <p class="mt-6 text-center text-sm text-slate-600">
        Already have an account? 
        <span @click="goToLogin" class="font-semibold text-blue-600 cursor-pointer hover:text-blue-700 transition-colors">
          Sign in
        </span>
      </p>
    </div>
  </div>
</template>
