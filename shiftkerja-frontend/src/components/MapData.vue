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
    console.log('Applying for shift:', shiftId);
    console.log('Token:', authStore.token ? 'Present' : 'Missing');
    
    const res = await fetch('http://localhost:8080/shifts/apply', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({ shift_id: shiftId })
    });

    console.log('Response status:', res.status);
    
    if (res.ok) {
      const result = await res.json();
      console.log('Success:', result);
      alert('✅ Application submitted successfully!');
      showShiftModal.value = false;
    } else {
      const error = await res.text();
      console.error('Error response:', error);
      alert('❌ Failed: ' + error);
    }
  } catch (error) {
    console.error('Network error applying:', error);
    alert('❌ Network error: ' + error.message);
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

const searchQuery = ref('');
const searchResults = ref([]);
const isSearching = ref(false);
const userMarker = ref(null);
const currentPosition = ref({ lat: -8.6478, lng: 115.1385 });
const manualPinMode = ref(false);
const tempMarker = ref(null);
const selectedCoords = ref(null);

// Geocoding search using Nominatim (OpenStreetMap)
const searchLocation = async () => {
  if (!searchQuery.value || searchQuery.value.length < 3) {
    searchResults.value = [];
    return;
  }

  isSearching.value = true;
  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(searchQuery.value)}&limit=5&countrycodes=id`
    );
    searchResults.value = await response.json();
  } catch (error) {
    console.error('Search error:', error);
    searchResults.value = [];
  } finally {
    isSearching.value = false;
  }
};

// Select a location from search results
const selectLocation = (result) => {
  const lat = parseFloat(result.lat);
  const lng = parseFloat(result.lon);
  
  map.value.setView([lat, lng], 15);
  searchQuery.value = result.display_name;
  searchResults.value = [];
  currentPosition.value = { lat, lng };
  
  // Fetch shifts near this location
  fetchShiftsNearby(lat, lng);
};

// Get user's current location
const getUserLocation = () => {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const lat = position.coords.latitude;
        const lng = position.coords.longitude;
        
        currentPosition.value = { lat, lng };
        map.value.setView([lat, lng], 14);
        
        // Add user location marker
        if (userMarker.value) {
          map.value.removeLayer(userMarker.value);
        }
        
        const userIcon = L.divIcon({
          className: 'user-location-marker',
          html: '<div style="background: #3B82F6; width: 16px; height: 16px; border-radius: 50%; border: 3px solid white; box-shadow: 0 0 10px rgba(59, 130, 246, 0.5);"></div>',
          iconSize: [16, 16],
          iconAnchor: [8, 8]
        });
        
        userMarker.value = L.marker([lat, lng], { icon: userIcon })
          .addTo(map.value)
          .bindPopup('You are here');
        
        // Fetch nearby shifts
        fetchShiftsNearby(lat, lng);
      },
      (error) => {
        console.error('Error getting location:', error);
        alert('Unable to get your location. Using default location.');
      }
    );
  } else {
    alert('Geolocation is not supported by your browser.');
  }
};

// Fetch shifts near a location
const fetchShiftsNearby = async (lat, lng, radius = 10) => {
  try {
    const response = await fetch(
      `http://localhost:8080/shifts?lat=${lat}&lng=${lng}&rad=${radius}`,
      {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      }
    );

    if (response.ok) {
      const shifts = await response.json();
      
      // Clear existing shift markers (but keep user marker)
      markers.value.forEach(m => map.value.removeLayer(m.marker));
      markers.value = [];
      
      // Add new markers
      shifts.forEach(shift => {
        const shiftIcon = L.divIcon({
          className: 'shift-marker',
          html: `<div style="background: #10B981; color: white; padding: 8px 12px; border-radius: 20px; font-weight: 600; font-size: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); white-space: nowrap;">Rp ${(shift.pay_rate / 1000).toFixed(0)}k</div>`,
          iconSize: [80, 32],
          iconAnchor: [40, 16]
        });
        
        const marker = L.marker([shift.lat, shift.lng], { icon: shiftIcon })
          .addTo(map.value)
          .on('click', () => viewShiftDetails(shift));
        
        markers.value.push({ marker, shift });
      });
    }
  } catch (error) {
    console.error('Error fetching shifts:', error);
  }
};

// Toggle manual pin mode
const togglePinMode = () => {
  manualPinMode.value = !manualPinMode.value;
  
  if (manualPinMode.value) {
    map.value.getContainer().style.cursor = 'crosshair';
    alert('📍 Click anywhere on the map to explore that location and find nearby shifts!');
  } else {
    map.value.getContainer().style.cursor = '';
    if (tempMarker.value) {
      map.value.removeLayer(tempMarker.value);
      tempMarker.value = null;
    }
  }
};

// Handle map click for manual pinning
const handleMapClick = (e) => {
  if (manualPinMode.value) {
    const lat = e.latlng.lat;
    const lng = e.latlng.lng;
    
    if (tempMarker.value) {
      map.value.removeLayer(tempMarker.value);
    }
    
    const pinIcon = L.divIcon({
      className: 'temp-pin-marker',
      html: '<div style="background: #EF4444; width: 24px; height: 24px; border-radius: 50% 50% 50% 0; transform: rotate(-45deg); border: 3px solid white; box-shadow: 0 4px 8px rgba(0,0,0,0.3);"><div style="width: 8px; height: 8px; background: white; border-radius: 50%; position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%) rotate(45deg);"></div></div>',
      iconSize: [24, 24],
      iconAnchor: [12, 24]
    });
    
    tempMarker.value = L.marker([lat, lng], { icon: pinIcon })
      .addTo(map.value)
      .bindPopup(`
        <div style="text-align: center; font-family: sans-serif;">
          <p style="font-weight: 600; margin: 0 0 8px 0; color: #1e293b;">📍 Selected Location</p>
          <p style="font-size: 11px; color: #64748b; margin: 0 0 8px 0; font-family: monospace;">${lat.toFixed(6)}, ${lng.toFixed(6)}</p>
          <p style="font-size: 12px; color: #10b981; margin: 0;">Loading nearby shifts...</p>
        </div>
      `)
      .openPopup();
    
    selectedCoords.value = { lat, lng };
    currentPosition.value = { lat, lng };
    fetchShiftsNearby(lat, lng);
    
    manualPinMode.value = false;
    map.value.getContainer().style.cursor = '';
  }
};

onMounted(async () => {
  // 1. Initialize the map
  map.value = L.map('mapContainer').setView([-8.6478, 115.1385], 13);

  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; OpenStreetMap'
  }).addTo(map.value);
  
  // Enable map click for manual pinning
  map.value.on('click', handleMapClick);

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
      console.log("🔔 Live Event:", data);

      if (data.type === 'shift_created') {
        // New shift posted - add animated live marker
        const liveIcon = L.divIcon({
          className: 'shift-marker-live',
          html: `<div style="background: linear-gradient(135deg, #F59E0B, #EF4444); color: white; padding: 8px 12px; border-radius: 20px; font-weight: 700; font-size: 12px; box-shadow: 0 4px 12px rgba(245, 158, 11, 0.4); white-space: nowrap;">🔥 NEW: Rp ${(data.pay_rate / 1000).toFixed(0)}k</div>`,
          iconSize: [120, 32],
          iconAnchor: [60, 16]
        });
        
        L.marker([data.lat, data.lng], { icon: liveIcon })
          .addTo(map.value)
          .bindPopup(`
            <div style="text-align: center; font-family: sans-serif;">
              <p style="font-weight: 700; margin: 0 0 4px 0; color: #F59E0B;">🔥 LIVE: New Shift Posted!</p>
              <p style="font-weight: 600; margin: 0 0 8px 0; color: #1e293b;">${data.title}</p>
              <p style="font-size: 14px; color: #10b981; font-weight: 600; margin: 0;">Rp ${data.pay_rate?.toLocaleString()}</p>
            </div>
          `)
          .openPopup();
          
        // Refresh shifts after 2 seconds
        setTimeout(() => {
          fetchShiftsNearby(currentPosition.value.lat, currentPosition.value.lng);
        }, 2000);
        
      } else if (data.type === 'shift_applied') {
        console.log("✅ Someone applied to shift:", data.shift_id);
      }
        
    } catch (e) {
      console.error("Error parsing WS message", e);
    }
  }, { deep: true });

  // 4. FETCH INITIAL DATA and enable location services
  fetchShiftsNearby(currentPosition.value.lat, currentPosition.value.lng);
  
  // Try to get user location automatically
  setTimeout(() => {
    getUserLocation();
  }, 1000);
  
  // Watch for WebSocket updates
  watch(() => socketStore.messages, (newMessages) => {
    if (newMessages.length === 0) return;

    const lastMsgJson = newMessages[newMessages.length - 1];
    
    try {
      const data = JSON.parse(lastMsgJson);
      console.log("📍 MAP UPDATE:", data);

      const shiftIcon = L.divIcon({
        className: 'shift-marker-live',
        html: `<div style="background: #F59E0B; color: white; padding: 8px 12px; border-radius: 20px; font-weight: 600; font-size: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); white-space: nowrap; animation: pulse 2s infinite;">LIVE: Rp ${(data.pay_rate / 1000).toFixed(0)}k</div>`,
        iconSize: [100, 32],
        iconAnchor: [50, 16]
      });
      
      L.marker([data.lat, data.lng], { icon: shiftIcon })
        .addTo(map.value)
        .on('click', () => {
          // Fetch fresh data
          fetchShiftsNearby(currentPosition.value.lat, currentPosition.value.lng);
        });
        
    } catch (e) {
      console.error("Error parsing WS message", e);
    }
  }, { deep: true });
});
</script>

<template>
  <div class="relative h-full w-full">
    <!-- Modern Top Navigation Bar -->
    <div class="absolute top-4 left-4 right-4 z-[9999] flex flex-col gap-3">
      <!-- Top Row: User info and Actions -->
      <div class="flex flex-col sm:flex-row justify-between items-stretch sm:items-center gap-3">
        <div class="bg-white/95 backdrop-blur-lg px-4 py-3 rounded-xl shadow-lg border border-slate-200">
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <p class="text-sm text-slate-600">
              <span class="font-semibold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent capitalize">{{ authStore.userRole }}</span>
            </p>
          </div>
        </div>
        
        <div class="flex gap-2">
          <button 
            @click="getUserLocation" 
            class="bg-white/95 backdrop-blur-lg text-blue-600 px-4 py-3 rounded-xl shadow-lg hover:bg-white font-medium transition-all duration-200 flex items-center justify-center gap-2"
            title="Use my current GPS location"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span class="hidden sm:inline">My Location</span>
          </button>
          
          <button 
            @click="goToDashboard" 
            class="flex-1 sm:flex-none bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-4 py-3 rounded-xl shadow-lg hover:from-blue-700 hover:to-indigo-700 font-medium transition-all duration-200 flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <span class="hidden sm:inline">Dashboard</span>
          </button>
          
          <button 
            @click="logout" 
            class="bg-red-50 text-red-600 px-4 py-3 rounded-xl shadow-lg hover:bg-red-100 font-medium transition-all duration-200"
          >
            Logout
          </button>
        </div>
      </div>

      <!-- Enhanced Search Bar -->
      <div class="relative">
        <div class="bg-white/95 backdrop-blur-lg rounded-xl shadow-lg border border-slate-200 overflow-hidden hover:border-blue-300 transition-colors">
          <div class="flex items-center gap-2 px-4">
            <svg class="w-5 h-5 text-slate-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              v-model="searchQuery"
              @input="searchLocation"
              @keydown.escape="searchResults = []; searchQuery = ''"
              type="text"
              placeholder="🔍 Search: Canggu, Denpasar, Seminyak, Ubud..."
              class="flex-1 py-3 bg-transparent border-none outline-none text-slate-700 placeholder-slate-400"
            />
            <button 
              v-if="searchQuery && !isSearching"
              @click="searchResults = []; searchQuery = ''"
              class="text-slate-400 hover:text-slate-600 transition-colors p-1"
              title="Clear search"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            <svg v-if="isSearching" class="w-5 h-5 text-blue-600 animate-spin flex-shrink-0" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
        </div>

        <!-- Search Results Dropdown -->
        <div v-if="searchResults.length > 0" class="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-2xl border border-slate-200 overflow-hidden max-h-64 overflow-y-auto z-[10000]">
          <button
            v-for="(result, index) in searchResults"
            :key="index"
            @click="selectLocation(result)"
            class="w-full px-4 py-3 text-left hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-b-0"
          >
            <div class="flex items-start gap-3">
              <svg class="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <div class="flex-1">
                <p class="text-sm font-medium text-slate-800">{{ result.name || result.display_name.split(',')[0] }}</p>
                <p class="text-xs text-slate-500 mt-0.5">{{ result.display_name }}</p>
              </div>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Map Container -->
    <div id="mapContainer" class="h-full w-full z-0"></div>

    <!-- Modern Shift Details Modal -->
    <div 
      v-if="showShiftModal && selectedShift"
      class="absolute inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[10000] p-4"
      @click.self="showShiftModal = false"
    >
      <div class="bg-white rounded-2xl shadow-2xl p-6 max-w-lg w-full mx-4 transform transition-all">
        <div class="flex justify-between items-start mb-6">
          <div>
            <h2 class="text-2xl font-bold text-slate-800 mb-1">{{ selectedShift.title }}</h2>
            <p v-if="selectedShift.description" class="text-sm text-slate-600">{{ selectedShift.description }}</p>
          </div>
          <button 
            @click="showShiftModal = false"
            class="text-slate-400 hover:text-slate-600 transition-colors p-1"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="space-y-4 mb-6">
          <div class="bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200 rounded-xl p-4">
            <div class="flex items-center gap-3">
              <div class="bg-green-600 p-3 rounded-lg">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <p class="text-sm text-green-700 font-medium">Pay Rate</p>
                <p class="text-2xl font-bold text-green-600">
                  Rp {{ selectedShift.pay_rate?.toLocaleString() }}
                </p>
              </div>
            </div>
          </div>
          
          <div class="bg-slate-50 border border-slate-200 rounded-xl p-4">
            <div class="flex items-start gap-3">
              <svg class="w-5 h-5 text-slate-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <div>
                <p class="text-sm text-slate-600 font-medium">Location</p>
                <p class="text-slate-800 font-mono text-sm">{{ selectedShift.lat?.toFixed(6) }}, {{ selectedShift.lng?.toFixed(6) }}</p>
              </div>
            </div>
          </div>
          
          <div class="flex items-center gap-2">
            <span 
              :class="{
                'bg-green-100 text-green-700 border-green-200': selectedShift.status === 'OPEN',
                'bg-slate-100 text-slate-700 border-slate-200': selectedShift.status !== 'OPEN'
              }"
              class="px-4 py-2 rounded-lg text-sm font-semibold border inline-flex items-center gap-2"
            >
              <div 
                :class="{
                  'bg-green-500': selectedShift.status === 'OPEN',
                  'bg-slate-400': selectedShift.status !== 'OPEN'
                }"
                class="w-2 h-2 rounded-full"
              ></div>
              {{ selectedShift.status }}
            </span>
          </div>
        </div>
        
        <div class="flex gap-3">
          <button 
            v-if="authStore.userRole === 'worker' && selectedShift.status === 'OPEN'"
            @click="applyForShift(selectedShift.id)"
            class="flex-1 bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-6 py-3 rounded-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-200 font-semibold flex items-center justify-center gap-2 shadow-lg shadow-blue-500/30"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Apply Now
          </button>
          
          <button 
            @click="showShiftModal = false"
            class="flex-1 bg-slate-100 text-slate-700 px-6 py-3 rounded-xl hover:bg-slate-200 transition-all duration-200 font-semibold"
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