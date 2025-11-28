<script setup>
import { onMounted, ref } from 'vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

const map = ref(null);

onMounted(async () => {
  // 1. Initialize the map
  map.value = L.map('mapContainer').setView([-8.6478, 115.1385], 13); // Centered on Canggu (near our seed data)

  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; OpenStreetMap'
  }).addTo(map.value);

  // 2. FETCH DATA from the Backend 🚀
  try {
    // We request shifts near the map center
    const response = await fetch('http://localhost:8080/shifts?lat=-8.6478&lng=115.1385&rad=10');
    
    // Parse the JSON text into a JavaScript Array
    const shifts = await response.json();
    console.log("Shifts received:", shifts);

    // 3. RENDER THE PINS
    // We loop through the array and add a marker for every shift found
    shifts.forEach(shift => {
      L.marker([shift.lat, shift.lng])
        .addTo(map.value)
        .bindPopup(`<b>${shift.title}</b><br>Pay: Rp ${shift.pay_rate}`);
    });

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