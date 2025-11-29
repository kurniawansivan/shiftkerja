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
  <div class="min-h-screen flex items-center justify-center relative overflow-hidden bg-gradient-to-br from-slate-900 via-purple-900 to-indigo-900 p-4">
    <!-- Animated Background Elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-1/4 -left-10 w-72 h-72 bg-purple-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
      <div class="absolute top-1/3 -right-10 w-72 h-72 bg-indigo-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
      <div class="absolute -bottom-10 left-1/2 w-72 h-72 bg-pink-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-4000"></div>
    </div>

    <!-- Register Card -->
    <div class="relative w-full max-w-md">
      <!-- Glassmorphism Card -->
      <div class="bg-white/10 backdrop-blur-2xl p-8 sm:p-10 rounded-3xl shadow-2xl border border-white/20 transform hover:scale-[1.01] transition-all duration-300">
      <!-- Logo & Branding -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-purple-500 to-indigo-600 rounded-2xl mb-4 shadow-lg shadow-purple-500/50">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
          </svg>
        </div>
        <h1 class="text-4xl font-bold text-white mb-2">
          Join ShiftKerja
        </h1>
        <p class="text-purple-200">Create your account to get started</p>
      </div>
      
      <!-- Error Message -->
      <div v-if="errorMsg" class="bg-red-500/20 border border-red-500/50 text-red-100 p-4 mb-6 rounded-2xl text-sm flex items-center gap-3 backdrop-blur-xl animate-shake">
        <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{{ errorMsg }}</span>
      </div>

      <!-- Success Message -->
      <div v-if="successMsg" class="bg-green-500/20 border border-green-500/50 text-green-100 p-4 mb-6 rounded-2xl text-sm flex items-center gap-3 backdrop-blur-xl animate-pulse">
        <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{{ successMsg }}</span>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleRegister" class="space-y-5">
        <!-- Full Name Input -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-purple-100">Full Name</label>
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <svg class="w-5 h-5 text-purple-300 group-focus-within:text-purple-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <input 
              v-model="fullName" 
              type="text" 
              placeholder="John Doe" 
              required
              class="w-full pl-12 pr-4 py-3.5 bg-white/10 border border-white/20 text-white placeholder-purple-300 rounded-2xl focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all backdrop-blur-xl"
            />
          </div>
        </div>

        <!-- Email Input -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-purple-100">Email Address</label>
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <svg class="w-5 h-5 text-purple-300 group-focus-within:text-purple-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
              </svg>
            </div>
            <input 
              v-model="email" 
              type="email" 
              placeholder="you@example.com" 
              required
              class="w-full pl-12 pr-4 py-3.5 bg-white/10 border border-white/20 text-white placeholder-purple-300 rounded-2xl focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all backdrop-blur-xl"
            />
          </div>
        </div>

        <!-- Password Input -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-purple-100">Password</label>
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <svg class="w-5 h-5 text-purple-300 group-focus-within:text-purple-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <input 
              v-model="password" 
              type="password" 
              placeholder="Minimum 6 characters" 
              required
              minlength="6"
              class="w-full pl-12 pr-4 py-3.5 bg-white/10 border border-white/20 text-white placeholder-purple-300 rounded-2xl focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all backdrop-blur-xl"
            />
          </div>
        </div>

        <!-- Role Selection -->
        <div class="space-y-3">
          <label class="block text-sm font-medium text-purple-100">Choose Your Role</label>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <label 
              :class="role === 'worker' ? 'border-purple-400 bg-purple-500/20 shadow-lg shadow-purple-500/30' : 'border-white/20 bg-white/5'"
              class="relative flex items-center p-5 border-2 rounded-2xl cursor-pointer transition-all hover:border-purple-400 hover:bg-white/10 group"
            >
              <input 
                type="radio" 
                v-model="role" 
                value="worker" 
                class="sr-only"
              />
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <div :class="role === 'worker' ? 'bg-purple-500' : 'bg-white/10'" class="w-10 h-10 rounded-xl flex items-center justify-center transition-colors">
                    <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div>
                    <p class="font-bold text-white">Worker</p>
                    <p class="text-xs text-purple-200">Find & apply for shifts</p>
                  </div>
                </div>
              </div>
              <div v-if="role === 'worker'" class="absolute top-3 right-3">
                <svg class="w-6 h-6 text-purple-400" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
              </div>
            </label>
            
            <label 
              :class="role === 'business' ? 'border-purple-400 bg-purple-500/20 shadow-lg shadow-purple-500/30' : 'border-white/20 bg-white/5'"
              class="relative flex items-center p-5 border-2 rounded-2xl cursor-pointer transition-all hover:border-purple-400 hover:bg-white/10 group"
            >
              <input 
                type="radio" 
                v-model="role" 
                value="business" 
                class="sr-only"
              />
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <div :class="role === 'business' ? 'bg-purple-500' : 'bg-white/10'" class="w-10 h-10 rounded-xl flex items-center justify-center transition-colors">
                    <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                    </svg>
                  </div>
                  <div>
                    <p class="font-bold text-white">Business</p>
                    <p class="text-xs text-purple-200">Post & manage shifts</p>
                  </div>
                </div>
              </div>
              <div v-if="role === 'business'" class="absolute top-3 right-3">
                <svg class="w-6 h-6 text-purple-400" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
              </div>
            </label>
          </div>
        </div>
        
        <!-- Submit Button -->
        <button 
          type="submit"
          :disabled="successMsg"
          class="w-full bg-gradient-to-r from-purple-500 to-indigo-600 text-white py-4 rounded-2xl hover:from-purple-600 hover:to-indigo-700 transition-all duration-300 font-bold shadow-xl shadow-purple-500/30 transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 group"
        >
          <span>{{ successMsg ? 'Redirecting...' : 'Create Account' }}</span>
          <svg v-if="!successMsg" class="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
          <svg v-else class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </button>
      </form>
      
      <!-- Divider -->
      <div class="relative my-6">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-white/20"></div>
        </div>
        <div class="relative flex justify-center text-sm">
          <span class="px-4 bg-transparent text-purple-200">Already have an account?</span>
        </div>
      </div>

      <!-- Login Link -->
      <button
        @click="goToLogin"
        type="button"
        class="w-full text-center px-6 py-4 bg-white/5 hover:bg-white/10 border border-white/20 text-white font-semibold rounded-2xl transition-all duration-300 transform hover:scale-[1.02]"
      >
        Sign In Instead
      </button>
    </div>

    <!-- Additional Info -->
    <p class="text-center text-purple-200 text-sm mt-6 opacity-75">
      Join thousands of workers and businesses on ShiftKerja
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
