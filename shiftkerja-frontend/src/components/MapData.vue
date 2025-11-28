<script setup>
import { onMounted, ref } from 'vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
// 👇 1. Import the Store
import { useAuthStore } from '@/stores/auth';

const map = ref(null);
// 👇 2. Initialize the Store
const authStore = useAuthStore();

onMounted(async () => {
  // 1. Initialize the map
  map.value = L.map('mapContainer').setView([-8.6478, 115.1385], 13);

  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; OpenStreetMap'
  }).addTo(map.value);

  // 2. FETCH DATA from the Backend 🚀
  try {
    const response = await fetch('http://localhost:8080/shifts?lat=-8.6478&lng=115.1385&rad=10', {
      headers: {
        // 👇 Now this works because authStore is defined
        'Authorization': `Bearer ${authStore.token}`
      }
    });

    // Handle "Unauthorized" (Token expired or missing)
    if (response.status === 401) {
      console.log("Unauthorized! Redirecting to login...");
      authStore.logout(); // Kick user out
      return;
    }

    const shifts = await response.json();

    // Check if shifts is actually an array before looping (safety check)
    if (Array.isArray(shifts)) {
      shifts.forEach(shift => {
        L.marker([shift.lat, shift.lng])
          .addTo(map.value)
          .bindPopup(`<b>${shift.title}</b><br>Pay: Rp ${shift.pay_rate}`);
      });
    }

  } catch (error) {
    console.error("Error connecting to backend:", error);
  }
});
</script>

<template>
  <div id="mapContainer" class="h-full w-full z-0"></div>
</template>

<style>
#mapContainer {
  z-index: 1;
}
</style>