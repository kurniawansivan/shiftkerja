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
  <div class="h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded shadow-md w-96">
      <h1 class="text-2xl font-bold mb-6 text-blue-900">ShiftKerja Login</h1>
      
      <div v-if="errorMsg" class="bg-red-100 text-red-700 p-2 mb-4 rounded">
        {{ errorMsg }}
      </div>

      <input v-model="email" type="email" placeholder="Email" class="w-full border p-2 mb-4 rounded" />
      <input v-model="password" type="password" placeholder="Password" class="w-full border p-2 mb-4 rounded" />
      
      <button @click="handleLogin" class="w-full bg-blue-900 text-white p-2 rounded hover:bg-blue-800 transition">
        Login
      </button>
      
      <p class="mt-4 text-sm text-gray-600">
        Don't have an account? 
        <router-link to="/register" class="font-bold text-blue-900 hover:underline">
          Register
        </router-link>
      </p>
    </div>
  </div>
</template>