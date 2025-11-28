import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '');
  const userRole = ref(localStorage.getItem('role') || '');
  const router = useRouter();

  // Function to Login
  const login = async (email, password) => {
    try {
      const res = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      if (!res.ok) throw new Error("Login failed");

      const data = await res.json();
      
      // 1. Save to State
      token.value = data.token;
      userRole.value = data.role;

      // 2. Persist to LocalStorage (so refresh doesn't kill session)
      localStorage.setItem('token', data.token);
      localStorage.setItem('role', data.role);

      return true; // Success
    } catch (err) {
      console.error(err);
      return false; // Failed
    }
  };

  // Function to Logout
  const logout = () => {
    token.value = '';
    userRole.value = '';
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    router.push('/login');
  };

  return { token, userRole, login, logout };
});