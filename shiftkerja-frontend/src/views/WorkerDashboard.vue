<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const applications = ref([]);
const loading = ref(true);

const fetchMyApplications = async () => {
  try {
    const res = await fetch('http://localhost:8080/my-applications', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    
    if (res.ok) {
      applications.value = await res.json() || [];
    }
  } catch (error) {
    console.error('Error fetching applications:', error);
  } finally {
    loading.value = false;
  }
};

const logout = () => {
  authStore.logout();
};

const goToMap = () => {
  router.push('/');
};

const getStatusColor = (status) => {
  switch (status) {
    case 'PENDING':
      return 'bg-yellow-100 text-yellow-800';
    case 'ACCEPTED':
      return 'bg-green-100 text-green-800';
    case 'REJECTED':
      return 'bg-red-100 text-red-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

onMounted(() => {
  fetchMyApplications();
});
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-blue-900 text-white p-4 shadow-lg">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-bold">My Applications</h1>
        <div class="flex gap-4">
          <button 
            @click="goToMap" 
            class="px-4 py-2 bg-blue-700 hover:bg-blue-600 rounded-lg transition"
          >
            🗺️ Find More Shifts
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
      <!-- Loading State -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Loading your applications...</p>
      </div>

      <!-- Applications List -->
      <div v-else-if="applications.length > 0" class="space-y-4">
        <div 
          v-for="app in applications" 
          :key="app.id"
          class="bg-white rounded-lg shadow-md p-6"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <h3 class="text-xl font-bold text-gray-800 mb-2">
                {{ app.shift_title }}
              </h3>
              
              <div class="space-y-1">
                <p class="text-green-600 font-semibold">
                  💰 Rp {{ app.shift_pay_rate?.toLocaleString() }}
                </p>
                <p class="text-sm text-gray-600">
                  📅 Applied: {{ new Date(app.created_at).toLocaleDateString('en-US', { 
                    year: 'numeric', 
                    month: 'long', 
                    day: 'numeric' 
                  }) }}
                </p>
              </div>
            </div>
            
            <span 
              :class="getStatusColor(app.status)"
              class="px-4 py-2 rounded-full text-sm font-semibold"
            >
              {{ app.status }}
            </span>
          </div>

          <!-- Status-specific messages -->
          <div v-if="app.status === 'ACCEPTED'" class="mt-4 p-3 bg-green-50 rounded-lg">
            <p class="text-green-800 text-sm">
              🎉 Congratulations! Your application has been accepted. The business will contact you soon.
            </p>
          </div>
          
          <div v-if="app.status === 'REJECTED'" class="mt-4 p-3 bg-red-50 rounded-lg">
            <p class="text-red-800 text-sm">
              Unfortunately, your application was not selected. Keep applying to other shifts!
            </p>
          </div>
          
          <div v-if="app.status === 'PENDING'" class="mt-4 p-3 bg-yellow-50 rounded-lg">
            <p class="text-yellow-800 text-sm">
              ⏳ Your application is under review. The business will respond soon.
            </p>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-12 bg-white rounded-lg shadow">
        <div class="mb-6">
          <svg class="mx-auto h-24 w-24 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-gray-700 mb-2">No Applications Yet</h3>
        <p class="text-gray-600 mb-6">Start exploring shifts on the map and apply to the ones you like!</p>
        <button 
          @click="goToMap"
          class="px-6 py-3 bg-blue-900 text-white rounded-lg hover:bg-blue-800 transition font-semibold"
        >
          🗺️ Browse Available Shifts
        </button>
      </div>
    </div>
  </div>
</template>
