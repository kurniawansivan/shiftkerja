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
  <div class="min-h-screen flex items-center justify-center relative overflow-hidden bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900 p-4">
    <!-- Animated Background Elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-1/4 -left-10 w-72 h-72 bg-blue-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
      <div class="absolute top-1/3 -right-10 w-72 h-72 bg-indigo-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
      <div class="absolute -bottom-10 left-1/2 w-72 h-72 bg-purple-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-4000"></div>
    </div>

    <!-- Login Card -->
    <div class="relative w-full max-w-md">
      <!-- Glassmorphism Card -->
      <div class="bg-white/10 backdrop-blur-2xl p-8 sm:p-10 rounded-3xl shadow-2xl border border-white/20 transform hover:scale-[1.02] transition-all duration-300">
        <!-- Logo & Branding -->
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl mb-4 shadow-lg shadow-blue-500/50">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <h1 class="text-4xl font-bold text-white mb-2">
            Welcome Back
          </h1>
          <p class="text-blue-200">Sign in to continue to ShiftKerja</p>
        </div>
        
        <!-- Error Message -->
        <div v-if="errorMsg" class="bg-red-500/20 border border-red-500/50 text-red-100 p-4 mb-6 rounded-2xl text-sm flex items-center gap-3 backdrop-blur-xl animate-shake">
          <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{{ errorMsg }}</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Email Input -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-blue-100">Email Address</label>
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-blue-300 group-focus-within:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
                </svg>
              </div>
              <input 
                v-model="email" 
                type="email" 
                placeholder="you@example.com"
                required
                class="w-full pl-12 pr-4 py-3.5 bg-white/10 border border-white/20 text-white placeholder-blue-300 rounded-2xl focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all backdrop-blur-xl"
              />
            </div>
          </div>

          <!-- Password Input -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-blue-100">Password</label>
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-blue-300 group-focus-within:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input 
                v-model="password" 
                type="password" 
                placeholder="Enter your password"
                required
                class="w-full pl-12 pr-4 py-3.5 bg-white/10 border border-white/20 text-white placeholder-blue-300 rounded-2xl focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all backdrop-blur-xl"
              />
            </div>
          </div>
          
          <!-- Submit Button -->
          <button 
            type="submit"
            class="w-full bg-gradient-to-r from-blue-500 to-indigo-600 text-white py-4 rounded-2xl hover:from-blue-600 hover:to-indigo-700 transition-all duration-300 font-bold shadow-xl shadow-blue-500/30 transform hover:scale-[1.02] active:scale-[0.98] flex items-center justify-center gap-2 group"
          >
            <span>Sign In</span>
            <svg class="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </button>
        </form>
        
        <!-- Divider -->
        <div class="relative my-8">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-white/20"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-4 bg-transparent text-blue-200">New to ShiftKerja?</span>
          </div>
        </div>

        <!-- Register Link -->
        <router-link 
          to="/register" 
          class="block w-full text-center px-6 py-4 bg-white/5 hover:bg-white/10 border border-white/20 text-white font-semibold rounded-2xl transition-all duration-300 transform hover:scale-[1.02]"
        >
          Create Account
        </router-link>
      </div>

      <!-- Additional Info -->
      <p class="text-center text-blue-200 text-sm mt-6 opacity-75">
        Secure login powered by ShiftKerja
      </p>
    </div>
  </div>
</template>

<style scoped>
@keyframes blob {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
  20%, 40%, 60%, 80% { transform: translateX(5px); }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}

.animate-shake {
  animation: shake 0.5s;
}
</style>