<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const shifts = ref([]);
const applications = ref({});
const loading = ref(true);

const newShift = ref({
  title: '',
  description: '',
  pay_rate: '',
  lat: -8.6478,
  lng: 115.1385
});

const showCreateForm = ref(false);

const fetchMyShifts = async () => {
  try {
    const res = await fetch('http://localhost:8080/shifts/my-shifts', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    
    if (res.ok) {
      shifts.value = await res.json();
      
      // Fetch applications for each shift
      for (const shift of shifts.value) {
        await fetchShiftApplications(shift.id);
      }
    }
  } catch (error) {
    console.error('Error fetching shifts:', error);
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
    }
  } catch (error) {
    console.error('Error fetching applications:', error);
  }
};

const createShift = async () => {
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
        lat: newShift.value.lat,
        lng: newShift.value.lng
      })
    });
    
    if (res.ok) {
      // Reset form
      newShift.value = {
        title: '',
        description: '',
        pay_rate: '',
        lat: -8.6478,
        lng: 115.1385
      };
      showCreateForm.value = false;
      
      // Refresh list
      await fetchMyShifts();
    } else {
      alert('Failed to create shift');
    }
  } catch (error) {
    console.error('Error creating shift:', error);
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

onMounted(() => {
  fetchMyShifts();
});
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-blue-900 text-white p-4 shadow-lg">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-bold">Business Dashboard</h1>
        <div class="flex gap-4">
          <button 
            @click="goToMap" 
            class="px-4 py-2 bg-blue-700 hover:bg-blue-600 rounded-lg transition"
          >
            🗺️ View Map
          </button>
          <button 
            @click="logout" 
            class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg transition"
          >
            Logout
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="container mx-auto p-6">
      <!-- Create New Shift Button -->
      <div class="mb-6">
        <button 
          @click="showCreateForm = !showCreateForm"
          class="px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition font-semibold"
        >
          {{ showCreateForm ? '❌ Cancel' : '➕ Post New Shift' }}
        </button>
      </div>

      <!-- Create Shift Form -->
      <div v-if="showCreateForm" class="bg-white p-6 rounded-lg shadow-md mb-6">
        <h2 class="text-xl font-bold mb-4 text-gray-800">Create New Shift</h2>
        <form @submit.prevent="createShift" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Job Title</label>
            <input 
              v-model="newShift.title" 
              type="text" 
              required
              placeholder="e.g., Barista at Canggu Coffee"
              class="w-full border border-gray-300 p-2 rounded-lg"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea 
              v-model="newShift.description" 
              rows="3"
              placeholder="Describe the job requirements..."
              class="w-full border border-gray-300 p-2 rounded-lg"
            ></textarea>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Pay Rate (Rp)</label>
              <input 
                v-model="newShift.pay_rate" 
                type="number" 
                required
                placeholder="75000"
                class="w-full border border-gray-300 p-2 rounded-lg"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Latitude</label>
              <input 
                v-model="newShift.lat" 
                type="number" 
                step="0.000001"
                required
                class="w-full border border-gray-300 p-2 rounded-lg"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Longitude</label>
              <input 
                v-model="newShift.lng" 
                type="number" 
                step="0.000001"
                required
                class="w-full border border-gray-300 p-2 rounded-lg"
              />
            </div>
          </div>

          <button 
            type="submit"
            class="w-full bg-blue-900 text-white p-3 rounded-lg hover:bg-blue-800 transition font-semibold"
          >
            Create Shift
          </button>
        </form>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Loading your shifts...</p>
      </div>

      <!-- Shifts List -->
      <div v-else-if="shifts.length > 0" class="space-y-6">
        <div 
          v-for="shift in shifts" 
          :key="shift.id"
          class="bg-white rounded-lg shadow-md p-6"
        >
          <!-- Shift Info -->
          <div class="border-b pb-4 mb-4">
            <div class="flex justify-between items-start">
              <div>
                <h3 class="text-xl font-bold text-gray-800">{{ shift.title }}</h3>
                <p class="text-gray-600 mt-1">{{ shift.description }}</p>
                <p class="text-green-600 font-semibold mt-2">Rp {{ shift.pay_rate.toLocaleString() }}</p>
              </div>
              <span 
                :class="{
                  'bg-green-100 text-green-800': shift.status === 'OPEN',
                  'bg-blue-100 text-blue-800': shift.status === 'FILLED',
                  'bg-gray-100 text-gray-800': shift.status === 'CANCELLED'
                }"
                class="px-3 py-1 rounded-full text-sm font-semibold"
              >
                {{ shift.status }}
              </span>
            </div>
          </div>

          <!-- Applications -->
          <div>
            <h4 class="font-semibold text-gray-700 mb-3">
              Applications ({{ applications[shift.id]?.length || 0 }})
            </h4>
            
            <div v-if="applications[shift.id]?.length > 0" class="space-y-3">
              <div 
                v-for="app in applications[shift.id]" 
                :key="app.id"
                class="flex justify-between items-center bg-gray-50 p-4 rounded-lg"
              >
                <div>
                  <p class="font-semibold text-gray-800">{{ app.worker_name }}</p>
                  <p class="text-sm text-gray-600">{{ app.worker_email }}</p>
                  <p class="text-xs text-gray-500 mt-1">
                    Applied: {{ new Date(app.created_at).toLocaleDateString() }}
                  </p>
                </div>
                
                <div v-if="app.status === 'PENDING'" class="flex gap-2">
                  <button 
                    @click="updateApplicationStatus(app.id, 'ACCEPTED')"
                    class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition text-sm"
                  >
                    ✓ Accept
                  </button>
                  <button 
                    @click="updateApplicationStatus(app.id, 'REJECTED')"
                    class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition text-sm"
                  >
                    ✗ Reject
                  </button>
                </div>
                
                <span 
                  v-else
                  :class="{
                    'bg-green-100 text-green-800': app.status === 'ACCEPTED',
                    'bg-red-100 text-red-800': app.status === 'REJECTED'
                  }"
                  class="px-3 py-1 rounded-full text-sm font-semibold"
                >
                  {{ app.status }}
                </span>
              </div>
            </div>
            
            <p v-else class="text-gray-500 text-sm">No applications yet</p>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-12 bg-white rounded-lg shadow">
        <p class="text-gray-600 mb-4">You haven't posted any shifts yet</p>
        <button 
          @click="showCreateForm = true"
          class="px-6 py-3 bg-blue-900 text-white rounded-lg hover:bg-blue-800 transition"
        >
          Post Your First Shift
        </button>
      </div>
    </div>
  </div>
</template>
