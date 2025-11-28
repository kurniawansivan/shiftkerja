<script setup>
import { onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { useAuthStore } from '@/stores/auth';
import { useSocketStore } from '@/stores/socket';

const map = ref(null);
const authStore = useAuthStore();
const socketStore = useSocketStore();
const router = useRouter();
const markers = ref([]);
const selectedShift = ref(null);
const showShiftModal = ref(false);

const applyForShift = async (shiftId) => {
  try {
    const res = await fetch('http://localhost:8080/shifts/apply', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({ shift_id: shiftId })
    });

    if (res.ok) {
      alert('✅ Application submitted successfully!');
      showShiftModal.value = false;
    } else {
      const error = await res.text();
      alert('❌ ' + error);
    }
  } catch (error) {
    console.error('Error applying:', error);
    alert('Failed to apply for shift');
  }
};

const viewShiftDetails = (shift) => {
  selectedShift.value = shift;
  showShiftModal.value = true;
};

const goToDashboard = () => {
  if (authStore.userRole === 'business') {
    router.push('/business/dashboard');
  } else if (authStore.userRole === 'worker') {
    router.push('/worker/dashboard');
  }
};

const logout = () => {
  authStore.logout();
};

// Send test location (for demo purposes)
const sendTestLocation = () => {
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
        const marker = L.marker([shift.lat, shift.lng])
          .addTo(map.value);
        
        // Create custom popup with button
        const popupContent = `
          <div style="min-width: 200px;">
            <h3 style="font-weight: bold; margin-bottom: 8px;">${shift.title}</h3>
            <p style="color: #16a34a; font-weight: 600; margin-bottom: 8px;">
              💰 Rp ${shift.pay_rate?.toLocaleString() || shift.pay_rate}
            </p>
            <p style="font-size: 12px; color: #666; margin-bottom: 8px;">
              ${shift.description || 'No description'}
            </p>
            <button 
              onclick="window.viewShift(${shift.id})"
              style="
                width: 100%;
                background: #1e3a8a;
                color: white;
                padding: 8px;
                border: none;
                border-radius: 6px;
                cursor: pointer;
                font-weight: 600;
              "
            >
              ${authStore.userRole === 'worker' ? '📋 Apply Now' : '👁️ View Details'}
            </button>
          </div>
        `;
        
        marker.bindPopup(popupContent);
        
        // Store for later reference
        markers.value.push({ marker, shift });
      });
    }

  } catch (error) {
    console.error("Error connecting to backend:", error);
  }
  
  // Make viewShift available globally for popup button
  window.viewShift = (shiftId) => {
    const shiftData = markers.value.find(m => m.shift.id === shiftId);
    if (shiftData) {
      viewShiftDetails(shiftData.shift);
    }
  };
});
</script>

<template>
  <div class="relative h-full w-full">
    <!-- Top Navigation Bar -->
    <div class="absolute top-4 left-4 right-4 z-[9999] flex justify-between items-center gap-4">
      <div class="bg-white px-4 py-2 rounded-lg shadow-lg">
        <p class="text-sm text-gray-600">
          Logged in as: <span class="font-bold text-blue-900">{{ authStore.userRole }}</span>
        </p>
      </div>
      
      <div class="flex gap-2">
        <button 
          @click="goToDashboard" 
          class="bg-blue-600 text-white px-4 py-2 rounded-lg shadow-lg hover:bg-blue-700 font-semibold transition"
        >
          📊 Dashboard
        </button>
        
        <button 
          @click="sendTestLocation" 
          class="bg-orange-600 text-white px-4 py-2 rounded-lg shadow-lg hover:bg-orange-700 font-semibold transition"
        >
          📍 Test WS
        </button>
        
        <button 
          @click="logout" 
          class="bg-red-600 text-white px-4 py-2 rounded-lg shadow-lg hover:bg-red-700 font-semibold transition"
        >
          Logout
        </button>
      </div>
    </div>

    <!-- Map Container -->
    <div id="mapContainer" class="h-full w-full z-0"></div>

    <!-- Shift Details Modal -->
    <div 
      v-if="showShiftModal && selectedShift"
      class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[10000]"
      @click.self="showShiftModal = false"
    >
      <div class="bg-white rounded-lg shadow-2xl p-6 max-w-md w-full mx-4">
        <div class="flex justify-between items-start mb-4">
          <h2 class="text-2xl font-bold text-gray-800">{{ selectedShift.title }}</h2>
          <button 
            @click="showShiftModal = false"
            class="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>
        
        <div class="space-y-3 mb-6">
          <div class="flex items-start">
            <span class="text-2xl mr-2">💰</span>
            <div>
              <p class="text-sm text-gray-600">Pay Rate</p>
              <p class="text-xl font-bold text-green-600">
                Rp {{ selectedShift.pay_rate?.toLocaleString() }}
              </p>
            </div>
          </div>
          
          <div class="flex items-start">
            <span class="text-2xl mr-2">📍</span>
            <div>
              <p class="text-sm text-gray-600">Location</p>
              <p class="text-gray-800">{{ selectedShift.lat }}, {{ selectedShift.lng }}</p>
            </div>
          </div>
          
          <div v-if="selectedShift.description" class="flex items-start">
            <span class="text-2xl mr-2">📝</span>
            <div>
              <p class="text-sm text-gray-600">Description</p>
              <p class="text-gray-800">{{ selectedShift.description }}</p>
            </div>
          </div>
          
          <div class="flex items-start">
            <span class="text-2xl mr-2">🏷️</span>
            <div>
              <p class="text-sm text-gray-600">Status</p>
              <span 
                :class="{
                  'bg-green-100 text-green-800': selectedShift.status === 'OPEN',
                  'bg-gray-100 text-gray-800': selectedShift.status !== 'OPEN'
                }"
                class="px-3 py-1 rounded-full text-sm font-semibold"
              >
                {{ selectedShift.status }}
              </span>
            </div>
          </div>
        </div>
        
        <div class="flex gap-2">
          <button 
            v-if="authStore.userRole === 'worker' && selectedShift.status === 'OPEN'"
            @click="applyForShift(selectedShift.id)"
            class="flex-1 bg-blue-900 text-white px-4 py-3 rounded-lg hover:bg-blue-800 transition font-semibold"
          >
            📋 Apply for This Shift
          </button>
          
          <button 
            @click="showShiftModal = false"
            class="flex-1 bg-gray-200 text-gray-800 px-4 py-3 rounded-lg hover:bg-gray-300 transition font-semibold"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
#mapContainer {
  z-index: 1;
}
</style>