<script setup>
import { onMounted, ref, watch } from 'vue'; // 👈 Added 'watch'
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { useAuthStore } from '@/stores/auth';
import { useSocketStore } from '@/stores/socket';

const map = ref(null);
const authStore = useAuthStore();
const socketStore = useSocketStore();

// --- SEND LOGIC (Simulating a Worker) ---
const sendTestLocation = () => {
  // Generate random small movement around Canggu
  const randomLat = -8.6478 + (Math.random() * 0.01 - 0.005);
  const randomLng = 115.1385 + (Math.random() * 0.01 - 0.005);

  const payload = {
    id: "shift_001",
    title: "Moving Barista",
    lat: randomLat,
    lng: randomLng,
    pay_rate: 75000
  };

  console.log("📤 Sending Location:", payload);
  socketStore.sendMessage(payload);
};

onMounted(async () => {
  // 1. Initialize the map
  map.value = L.map('mapContainer').setView([-8.6478, 115.1385], 13);

  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; OpenStreetMap'
  }).addTo(map.value);

  // 2. CONNECT WEBSOCKET
  socketStore.connect();

  // 3. WATCH FOR LIVE UPDATES (The Receiver)
  watch(() => socketStore.messages, (newMessages) => {
    if (newMessages.length === 0) return;

    // Get the latest message
    const lastMsgJson = newMessages[newMessages.length - 1];
    
    try {
      // Parse the JSON string
      const data = JSON.parse(lastMsgJson);
      console.log("📍 MAP UPDATE:", data);

      // Add a marker to the map (Live Update!)
      L.marker([data.lat, data.lng])
        .addTo(map.value)
        .bindPopup(`<b>${data.title}</b><br>LIVE UPDATE!`)
        .openPopup();
        
    } catch (e) {
      console.error("Error parsing WS message", e);
    }
  }, { deep: true });

  // 4. FETCH INITIAL DATA from Backend
  try {
    const response = await fetch('http://localhost:8080/shifts?lat=-8.6478&lng=115.1385&rad=10', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });

    if (response.status === 401) {
      console.log("Unauthorized! Redirecting...");
      authStore.logout();
      return;
    }

    const shifts = await response.json();
    
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
  <div class="relative h-full w-full">
    <button 
      @click="sendTestLocation" 
      class="absolute top-4 right-4 z-[9999] bg-red-600 text-white px-4 py-2 rounded shadow-lg hover:bg-red-700 font-bold"
    >
      📍 Move Me!
    </button>

    <div id="mapContainer" class="h-full w-full z-0"></div>
  </div>
</template>

<style>
#mapContainer {
  z-index: 1;
}
</style>