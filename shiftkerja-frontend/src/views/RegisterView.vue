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
  <div class="h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-blue-100">
    <div class="bg-white p-8 rounded-lg shadow-xl w-full max-w-md">
      <h1 class="text-3xl font-bold mb-2 text-blue-900">Create Account</h1>
      <p class="text-gray-600 mb-6">Join ShiftKerja as a Worker or Business</p>
      
      <!-- Error Message -->
      <div v-if="errorMsg" class="bg-red-100 text-red-700 p-3 mb-4 rounded-lg text-sm">
        {{ errorMsg }}
      </div>

      <!-- Success Message -->
      <div v-if="successMsg" class="bg-green-100 text-green-700 p-3 mb-4 rounded-lg text-sm">
        {{ successMsg }}
      </div>

      <!-- Form -->
      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
          <input 
            v-model="fullName" 
            type="text" 
            placeholder="John Doe" 
            required
            class="w-full border border-gray-300 p-3 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input 
            v-model="email" 
            type="email" 
            placeholder="john@example.com" 
            required
            class="w-full border border-gray-300 p-3 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input 
            v-model="password" 
            type="password" 
            placeholder="Min. 6 characters" 
            required
            minlength="6"
            class="w-full border border-gray-300 p-3 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">I am a:</label>
          <div class="flex gap-4">
            <label class="flex items-center cursor-pointer">
              <input 
                type="radio" 
                v-model="role" 
                value="worker" 
                class="mr-2 w-4 h-4 text-blue-600"
              />
              <span class="text-gray-700">Worker (Looking for shifts)</span>
            </label>
            
            <label class="flex items-center cursor-pointer">
              <input 
                type="radio" 
                v-model="role" 
                value="business" 
                class="mr-2 w-4 h-4 text-blue-600"
              />
              <span class="text-gray-700">Business (Posting shifts)</span>
            </label>
          </div>
        </div>
        
        <button 
          type="submit"
          class="w-full bg-blue-900 text-white p-3 rounded-lg hover:bg-blue-800 transition font-semibold"
        >
          Register
        </button>
      </form>
      
      <p class="mt-6 text-sm text-center text-gray-600">
        Already have an account? 
        <span @click="goToLogin" class="font-bold text-blue-900 cursor-pointer hover:underline">
          Login here
        </span>
      </p>
    </div>
  </div>
</template>
