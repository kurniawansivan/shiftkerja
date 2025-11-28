<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

const authStore = useAuthStore();
const router = useRouter();

const shifts = ref([]);
const applications = ref({});
const loading = ref(true);
const error = ref('');

const newShift = ref({
  title: '',
  description: '',
  pay_rate: '',
  lat: '',
  lng: ''
});

const editingShift = ref(null);
const showCreateForm = ref(false);
const showEditModal = ref(false);
const showLocationPicker = ref(false);
const locationSearch = ref('');
const locationResults = ref([]);
const pickingLocationFor = ref('create'); // 'create' or 'edit'
const mapInstance = ref(null);
const tempMarker = ref(null);
const showMapPicker = ref(false);

const fetchMyShifts = async () => {
  try {
    const res = await fetch('http://localhost:8080/shifts/my-shifts', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    
    if (res.ok) {
      const data = await res.json();
      shifts.value = data || [];
      
      for (const shift of shifts.value) {
        await fetchShiftApplications(shift.id);
      }
    } else {
      shifts.value = [];
      if (res.status === 401) {
        authStore.logout();
      }
    }
  } catch (err) {
    console.error('Error fetching shifts:', err);
    shifts.value = [];
    error.value = 'Failed to load shifts';
  } finally {
    loading.value = false;
  }
};

const fetchShiftApplications = async (shiftId) => {
  try {
    const res = await fetch(`http://localhost:8080/shifts/applications?shift_id=${shiftId}`, {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    
    if (res.ok) {
      const apps = await res.json();
      applications.value[shiftId] = apps || [];
    } else {
      applications.value[shiftId] = [];
    }
  } catch (error) {
    console.error('Error fetching applications:', error);
    applications.value[shiftId] = [];
  }
};

const createShift = async () => {
  if (!newShift.value.title || !newShift.value.pay_rate || !newShift.value.lat || !newShift.value.lng) {
    alert('Please fill in all required fields');
    return;
  }

  try {
    const res = await fetch('http://localhost:8080/shifts/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({
        title: newShift.value.title,
        description: newShift.value.description,
        pay_rate: parseFloat(newShift.value.pay_rate),
        lat: parseFloat(newShift.value.lat),
        lng: parseFloat(newShift.value.lng)
      })
    });
    
    if (res.ok) {
      newShift.value = {
        title: '',
        description: '',
        pay_rate: '',
        lat: '',
        lng: ''
      };
      showCreateForm.value = false;
      await fetchMyShifts();
    } else {
      const errorText = await res.text();
      alert('Failed to create shift: ' + errorText);
    }
  } catch (err) {
    console.error('Error creating shift:', err);
    alert('Error creating shift');
  }
};

const updateApplicationStatus = async (applicationId, status) => {
  try {
    const res = await fetch('http://localhost:8080/shifts/applications/update', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({
        application_id: applicationId,
        status: status
      })
    });
    
    if (res.ok) {
      // Refresh the shifts and applications
      await fetchMyShifts();
    } else {
      alert('Failed to update application');
    }
  } catch (error) {
    console.error('Error updating application:', error);
  }
};

const logout = () => {
  authStore.logout();
};

const goToMap = () => {
  router.push('/');
};

// Location search
const searchLocations = async () => {
  if (!locationSearch.value || locationSearch.value.length < 3) {
    locationResults.value = [];
    return;
  }

  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(locationSearch.value)}&limit=5&countrycodes=id`
    );
    locationResults.value = await response.json();
  } catch (error) {
    console.error('Search error:', error);
  }
};

const selectLocation = (result) => {
  const lat = parseFloat(result.lat);
  const lng = parseFloat(result.lon);
  
  if (pickingLocationFor.value === 'create') {
    newShift.value.lat = lat.toFixed(6);
    newShift.value.lng = lng.toFixed(6);
  } else {
    editingShift.value.lat = lat;
    editingShift.value.lng = lng;
  }
  
  locationSearch.value = result.display_name;
  locationResults.value = [];
  showLocationPicker.value = false;
};

const openLocationPicker = (forType) => {
  pickingLocationFor.value = forType;
  locationSearch.value = '';
  locationResults.value = [];
  showLocationPicker.value = true;
};

// Map-based location picker
const openMapPicker = (forType) => {
  pickingLocationFor.value = forType;
  showMapPicker.value = true;
  
  // Initialize map after modal opens
  setTimeout(() => {
    if (!mapInstance.value) {
      mapInstance.value = L.map('locationMap').setView([-8.6478, 115.1385], 13);
      
      L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '© OpenStreetMap'
      }).addTo(mapInstance.value);
      
      // Click to place marker
      mapInstance.value.on('click', handleMapClick);
    }
  }, 100);
};

const handleMapClick = (e) => {
  const lat = e.latlng.lat;
  const lng = e.latlng.lng;
  
  // Remove previous marker
  if (tempMarker.value) {
    mapInstance.value.removeLayer(tempMarker.value);
  }
  
  // Create pin marker
  const pinIcon = L.divIcon({
    className: 'temp-pin-marker',
    html: '<div style=\"background: #EF4444; width: 30px; height: 30px; border-radius: 50% 50% 50% 0; transform: rotate(-45deg); border: 4px solid white; box-shadow: 0 4px 8px rgba(0,0,0,0.3);\"><div style=\"width: 10px; height: 10px; background: white; border-radius: 50%; position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%) rotate(45deg);\"></div></div>',
    iconSize: [30, 30],
    iconAnchor: [15, 30]
  });
  
  tempMarker.value = L.marker([lat, lng], { icon: pinIcon })
    .addTo(mapInstance.value)
    .bindPopup(`
      <div style=\"text-align: center; font-family: sans-serif;\">
        <p style=\"font-weight: 600; margin: 0 0 8px 0;\">📍 Selected Location</p>
        <p style=\"font-size: 11px; color: #666; margin: 0 0 12px 0; font-family: monospace;\">${lat.toFixed(6)}, ${lng.toFixed(6)}</p>
        <button onclick=\"document.getElementById('confirmLocation').click()\" style=\"background: #3B82F6; color: white; padding: 8px 16px; border-radius: 8px; border: none; cursor: pointer; font-weight: 600;\">Use This Location</button>
      </div>
    `)
    .openPopup();
};

const confirmMapSelection = () => {
  if (!tempMarker.value) {
    alert('Please click on the map to select a location');
    return;
  }
  
  const latlng = tempMarker.value.getLatLng();
  
  if (pickingLocationFor.value === 'create') {
    newShift.value.lat = latlng.lat.toFixed(6);
    newShift.value.lng = latlng.lng.toFixed(6);
  } else {
    editingShift.value.lat = latlng.lat;
    editingShift.value.lng = latlng.lng;
  }
  
  closeMapPicker();
};

const closeMapPicker = () => {
  showMapPicker.value = false;
  if (tempMarker.value && mapInstance.value) {
    mapInstance.value.removeLayer(tempMarker.value);
    tempMarker.value = null;
  }
};

// Edit shift
const openEditModal = (shift) => {
  editingShift.value = { ...shift };
  showEditModal.value = true;
};

const updateShift = async () => {
  try {
    const res = await fetch('http://localhost:8080/shifts/update', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({
        id: editingShift.value.id,
        title: editingShift.value.title,
        description: editingShift.value.description,
        pay_rate: parseFloat(editingShift.value.pay_rate),
        lat: parseFloat(editingShift.value.lat),
        lng: parseFloat(editingShift.value.lng),
        status: editingShift.value.status
      })
    });
    
    if (res.ok) {
      showEditModal.value = false;
      editingShift.value = null;
      await fetchMyShifts();
    } else {
      const errorText = await res.text();
      alert('Failed to update shift: ' + errorText);
    }
  } catch (err) {
    console.error('Error updating shift:', err);
    alert('Error updating shift');
  }
};

// Delete shift
const deleteShift = async (shiftId) => {
  if (!confirm('Are you sure you want to delete this shift? This will also delete all applications.')) {
    return;
  }

  try {
    const res = await fetch(`http://localhost:8080/shifts/delete?shift_id=${shiftId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    
    if (res.ok) {
      await fetchMyShifts();
    } else {
      const errorText = await res.text();
      alert('Failed to delete shift: ' + errorText);
    }
  } catch (err) {
    console.error('Error deleting shift:', err);
    alert('Error deleting shift');
  }
};

onMounted(() => {
  fetchMyShifts();
});
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
    <!-- Modern Header -->
    <header class="bg-white border-b border-slate-200 sticky top-0 z-50 backdrop-blur-lg bg-white/90">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex justify-between items-center">
          <div>
            <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
              Business Dashboard
            </h1>
            <p class="text-sm text-slate-500 mt-1">Manage your shifts and applications</p>
          </div>
          <div class="flex gap-2">
            <button 
              @click="goToMap" 
              class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition-all duration-200 flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
              </svg>
              Map
            </button>
            <button 
              @click="logout" 
              class="px-4 py-2 text-sm font-medium text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-all duration-200"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Create Button -->
      <div class="mb-6">
        <button 
          @click="showCreateForm = !showCreateForm"
          class="w-full sm:w-auto px-6 py-3 text-sm font-semibold rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
          :class="showCreateForm 
            ? 'bg-slate-200 text-slate-700 hover:bg-slate-300' 
            : 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white hover:from-blue-700 hover:to-indigo-700 shadow-lg shadow-blue-500/30'"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="showCreateForm ? 'M6 18L18 6M6 6l12 12' : 'M12 4v16m8-8H4'" />
          </svg>
          {{ showCreateForm ? 'Cancel' : 'Post New Shift' }}
        </button>
      </div>

      <!-- Create Form -->
      <transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div v-if="showCreateForm" class="bg-white rounded-2xl shadow-xl p-6 mb-8 border border-slate-200">
          <h2 class="text-xl font-bold text-slate-800 mb-6">Create New Shift</h2>
          <form @submit.prevent="createShift" class="space-y-5">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-2">Job Title *</label>
              <input 
                v-model="newShift.title" 
                type="text" 
                required
                placeholder="e.g., Barista at Canggu Coffee"
                class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-700 mb-2">Description</label>
              <textarea 
                v-model="newShift.description" 
                rows="3"
                placeholder="Describe the job requirements..."
                class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all resize-none"
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-700 mb-2">Pay Rate (Rp) *</label>
              <input 
                v-model="newShift.pay_rate" 
                type="number" 
                required
                placeholder="75000"
                class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-700 mb-2">Location *</label>
              <div class="grid grid-cols-2 gap-3">
                <button
                  type="button"
                  @click="openLocationPicker('create')"
                  class="px-4 py-3 border-2 border-slate-300 rounded-xl hover:border-blue-500 hover:bg-blue-50 transition-all flex flex-col items-center gap-1 text-slate-600"
                >
                  <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  <span class="text-sm font-medium">Search Address</span>
                </button>
                <button
                  type="button"
                  @click="openMapPicker('create')"
                  class="px-4 py-3 border-2 border-red-300 rounded-xl hover:border-red-500 hover:bg-red-50 transition-all flex flex-col items-center gap-1 text-red-600"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  <span class="text-sm font-medium">Pin on Map</span>
                </button>
              </div>
              <div v-if="newShift.lat && newShift.lng" class="mt-2 px-3 py-2 bg-green-50 border border-green-200 rounded-lg">
                <p class="text-sm text-green-800">
                  <span class="font-semibold">Selected:</span>
                  <span class="font-mono ml-1">{{ newShift.lat }}, {{ newShift.lng }}</span>
                </p>
              </div>
            </div>

            <button 
              type="submit"
              class="w-full py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-200 shadow-lg shadow-blue-500/30"
            >
              Create Shift
            </button>
          </form>
        </div>
      </transition>

      <!-- Loading -->
      <div v-if="loading" class="text-center py-20">
        <div class="inline-block w-12 h-12 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
        <p class="text-slate-600 mt-4">Loading shifts...</p>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4 text-red-700">
        {{ error }}
      </div>

      <!-- Shifts List -->
      <div v-else-if="shifts.length > 0" class="space-y-6">
        <div 
          v-for="shift in shifts" 
          :key="shift.id"
          class="bg-white rounded-2xl shadow-lg border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-200"
        >
          <div class="p-6">
            <!-- Shift Header -->
            <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4 mb-4">
              <div class="flex-1">
                <h3 class="text-xl font-bold text-slate-800 mb-2">{{ shift.title }}</h3>
                <p v-if="shift.description" class="text-slate-600 text-sm mb-3">{{ shift.description }}</p>
                <div class="flex flex-wrap items-center gap-4 text-sm">
                  <span class="flex items-center gap-1 text-green-600 font-semibold">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Rp {{ (shift.pay_rate || 0).toLocaleString() }}
                  </span>
                  <span class="flex items-center gap-1 text-slate-500">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    {{ shift.lat?.toFixed(4) }}, {{ shift.lng?.toFixed(4) }}
                  </span>
                </div>
              </div>
              <span 
                :class="{
                  'bg-green-100 text-green-700': shift.status === 'OPEN',
                  'bg-blue-100 text-blue-700': shift.status === 'FILLED',
                  'bg-slate-100 text-slate-700': shift.status === 'CANCELLED'
                }"
                class="px-4 py-2 rounded-full text-sm font-semibold whitespace-nowrap"
              >
                {{ shift.status }}
              </span>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-2 mt-4 pt-4 border-t border-slate-200">
              <button
                @click="openEditModal(shift)"
                class="flex-1 px-4 py-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-all flex items-center justify-center gap-2 text-sm font-medium"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Edit
              </button>
              <button
                @click="deleteShift(shift.id)"
                class="flex-1 px-4 py-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition-all flex items-center justify-center gap-2 text-sm font-medium"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                Delete
              </button>
            </div>

            <!-- Applications Section -->
            <div class="border-t border-slate-200 pt-4 mt-4">
              <h4 class="font-semibold text-slate-700 mb-4 flex items-center gap-2">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
                Applications <span class="text-slate-500">({{ (applications[shift.id] || []).length }})</span>
              </h4>
              
              <div v-if="(applications[shift.id] || []).length > 0" class="space-y-3">
                <div 
                  v-for="app in applications[shift.id]" 
                  :key="app.id"
                  class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 bg-slate-50 p-4 rounded-xl border border-slate-200"
                >
                  <div class="flex-1">
                    <p class="font-semibold text-slate-800">{{ app.worker_name }}</p>
                    <p class="text-sm text-slate-600">{{ app.worker_email }}</p>
                    <p class="text-xs text-slate-500 mt-1">
                      Applied {{ new Date(app.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }) }}
                    </p>
                  </div>
                  
                  <div v-if="app.status === 'PENDING'" class="flex gap-2">
                    <button 
                      @click="updateApplicationStatus(app.id, 'ACCEPTED')"
                      class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-all flex items-center gap-1"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                      </svg>
                      Accept
                    </button>
                    <button 
                      @click="updateApplicationStatus(app.id, 'REJECTED')"
                      class="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 transition-all flex items-center gap-1"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      Reject
                    </button>
                  </div>
                  
                  <span 
                    v-else
                    :class="{
                      'bg-green-100 text-green-700': app.status === 'ACCEPTED',
                      'bg-red-100 text-red-700': app.status === 'REJECTED'
                    }"
                    class="px-4 py-2 rounded-full text-sm font-semibold whitespace-nowrap"
                  >
                    {{ app.status }}
                  </span>
                </div>
              </div>
              
              <p v-else class="text-slate-500 text-sm italic">No applications yet</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-20">
        <svg class="w-24 h-24 mx-auto text-slate-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-xl font-semibold text-slate-700 mb-2">No Shifts Yet</h3>
        <p class="text-slate-500 mb-6">Start by posting your first shift to attract workers</p>
        <button 
          @click="showCreateForm = true"
          class="px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-200 shadow-lg shadow-blue-500/30"
        >
          Post Your First Shift
        </button>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal && editingShift" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[10000] p-4">
      <div class="bg-white rounded-2xl shadow-2xl p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-start mb-6">
          <h2 class="text-2xl font-bold text-slate-800">Edit Shift</h2>
          <button @click="showEditModal = false" class="text-slate-400 hover:text-slate-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <form @submit.prevent="updateShift" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Job Title *</label>
            <input v-model="editingShift.title" type="text" required class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all" />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Description</label>
            <textarea v-model="editingShift.description" rows="3" class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all resize-none"></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Pay Rate (Rp) *</label>
            <input v-model="editingShift.pay_rate" type="number" required class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all" />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Location</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                type="button"
                @click="openLocationPicker('edit')"
                class="px-4 py-3 border-2 border-slate-300 rounded-xl hover:border-blue-500 hover:bg-blue-50 transition-all flex flex-col items-center gap-1 text-slate-600"
              >
                <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <span class="text-sm font-medium">Search Address</span>
              </button>
              <button
                type="button"
                @click="openMapPicker('edit')"
                class="px-4 py-3 border-2 border-red-300 rounded-xl hover:border-red-500 hover:bg-red-50 transition-all flex flex-col items-center gap-1 text-red-600"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <span class="text-sm font-medium">Pin on Map</span>
              </button>
            </div>
            <div class="mt-2 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg">
              <p class="text-sm text-slate-700">
                <span class="font-semibold">Current:</span>
                <span class="font-mono ml-1">{{ editingShift.lat?.toFixed(6) }}, {{ editingShift.lng?.toFixed(6) }}</span>
              </p>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Status</label>
            <select v-model="editingShift.status" class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all">
              <option value="OPEN">OPEN</option>
              <option value="FILLED">FILLED</option>
              <option value="CANCELLED">CANCELLED</option>
            </select>
          </div>

          <div class="flex gap-3">
            <button type="submit" class="flex-1 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-200 shadow-lg shadow-blue-500/30">
              Update Shift
            </button>
            <button type="button" @click="showEditModal = false" class="flex-1 py-3 bg-slate-100 text-slate-700 font-semibold rounded-xl hover:bg-slate-200 transition-all duration-200">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Location Picker Modal (Search) -->
    <div v-if="showLocationPicker" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[10001] p-4">
      <div class="bg-white rounded-2xl shadow-2xl p-6 max-w-2xl w-full mx-4">
        <div class="flex justify-between items-start mb-6">
          <h2 class="text-2xl font-bold text-slate-800">Search Location</h2>
          <button @click="showLocationPicker = false" class="text-slate-400 hover:text-slate-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="relative mb-4">
          <input
            v-model="locationSearch"
            @input="searchLocations"
            type="text"
            placeholder="Search for Canggu, Denpasar, Seminyak..."
            class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          />
        </div>

        <div v-if="locationResults.length > 0" class="space-y-2 max-h-96 overflow-y-auto">
          <button
            v-for="(result, index) in locationResults"
            :key="index"
            @click="selectLocation(result)"
            type="button"
            class="w-full text-left px-4 py-3 border border-slate-200 rounded-xl hover:bg-blue-50 hover:border-blue-500 transition-all"
          >
            <p class="font-medium text-slate-800">{{ result.name || result.display_name.split(',')[0] }}</p>
            <p class="text-sm text-slate-500 mt-1">{{ result.display_name }}</p>
          </button>
        </div>
        
        <p v-else-if="locationSearch && locationResults.length === 0" class="text-center text-slate-500 py-8">
          Type at least 3 characters to search...
        </p>
        
        <div class="mt-4 pt-4 border-t border-slate-200">
          <button
            @click="showLocationPicker = false; openMapPicker(pickingLocationFor)"
            type="button"
            class="w-full bg-gradient-to-r from-red-600 to-pink-600 text-white px-4 py-3 rounded-xl hover:from-red-700 hover:to-pink-700 font-semibold flex items-center justify-center gap-2 transition-all"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
            </svg>
            Can't find it? Pin location on map manually
          </button>
        </div>
      </div>
    </div>

    <!-- Map Picker Modal -->
    <div v-if="showMapPicker" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[10001] p-4">
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-4xl mx-4 overflow-hidden">
        <div class="bg-gradient-to-r from-blue-600 to-indigo-600 p-6 text-white">
          <div class="flex justify-between items-start">
            <div>
              <h2 class="text-2xl font-bold">📍 Pin Your Location</h2>
              <p class="text-blue-100 mt-1">Click anywhere on the map to select the exact location</p>
            </div>
            <button @click="closeMapPicker" class="text-white/80 hover:text-white transition-colors">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
        
        <div id="locationMap" class="h-96 w-full"></div>
        
        <div class="p-6 bg-slate-50 border-t border-slate-200">
          <div class="flex gap-3">
            <button
              id="confirmLocation"
              @click="confirmMapSelection"
              class="flex-1 bg-gradient-to-r from-green-600 to-emerald-600 text-white px-6 py-3 rounded-xl hover:from-green-700 hover:to-emerald-700 font-semibold flex items-center justify-center gap-2 transition-all shadow-lg"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              Confirm Location
            </button>
            <button
              @click="closeMapPicker"
              class="px-6 py-3 bg-white text-slate-700 border-2 border-slate-300 rounded-xl hover:bg-slate-50 font-semibold transition-all"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
