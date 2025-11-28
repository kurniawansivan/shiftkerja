<script setup>
import { onMounted, ref } from 'vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

const map = ref(null);

onMounted(() => {
  // 1. Initialize the map centered on a default location (e.g., Bali)
  // [Lat, Lng], Zoom Level
  map.value = L.map('mapContainer').setView([-8.409518, 115.188919], 13);

  // 2. Add the "Tile Layer" (The skin of the map)
  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map.value);

  // 3. Add a test marker to prove it works
  L.marker([-8.409518, 115.188919]).addTo(map.value)
    .bindPopup('ShiftKerja HQ')
    .openPopup();
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