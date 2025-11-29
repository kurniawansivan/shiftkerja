<script setup>
import { ref, onMounted, watch, onUnmounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useSocketStore } from '@/stores/socket';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const socketStore = useSocketStore();
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

// Real-time notification state
const statusUpdateNotification = ref(null);
const showNotification = ref(false);

// Handle real-time WebSocket messages
watch(() => socketStore.messages, (newMessages) => {
  if (newMessages.length === 0) return;

  const lastMsgJson = newMessages[newMessages.length - 1];
  
  try {
    const data = JSON.parse(lastMsgJson);
    console.log("🔔 Worker Dashboard received:", data);

    if (data.type === 'application_status_updated') {
      // Check if this is one of our applications
      const app = applications.value.find(a => a.id === data.application_id);
      if (app) {
        // Show notification
        statusUpdateNotification.value = {
          application_id: data.application_id,
          new_status: data.new_status,
          shift_title: app.shift_title
        };
        showNotification.value = true;

        // Auto-hide after 5 seconds
        setTimeout(() => {
          showNotification.value = false;
        }, 5000);

        // Refresh applications
        fetchMyApplications();
        
        // Play notification sound
        playNotificationSound();
      }
    }
  } catch (e) {
    console.error("Error parsing WebSocket message", e);
  }
}, { deep: true });

// Optional: Play notification sound
const playNotificationSound = () => {
  try {
    const audio = new Audio('data:audio/wav;base64,UklGRnoGAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQoGAACBhYqFbF1fdJivrJBhNjVgodDbq2EcBj+a2/LDciUFLIHO8tiJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmwhBSuBzvLZiTYIG2m98OScTgwOUKvl8bllHgU2jdXyzn0vBSV3yPDej0QKFF+07+yrWBQLR6Hh8sFuIwUqgdDy2Ik2CBtpvfDmnE4MDlCs5fG6Zh4FNo3V8s9+MAUleMnw35FGChVftO/trFoVDEah4fHCcCQFKoLQ8tqKNwgcasLw6Z1PDAxOq+XxumgeBS2O1fLPfzIFJXnK8OCTRwoVX7Xv7a1cFQxGoeLxxG8lBSqC0PLaijcIHGrD8OueUAwMT6vl8btpHwUujtbyz4AzBSV5yvDglEgKFV+17+2uXRUMRqHi8cVwJgUrgtDy24o4CBxqw/DsnlENDFCr5fG8aR8FLo7W8tCANAUlecvw4JVJChZgtO/urlkVDEag4vHFcSYFK4LQ8tuKOAgcasLw7J5SDQxQq+XxvGkfBS+O1vLRgDQFJXnL8OCVSQoWYLTv766aFg1GoOLxx3EnBSuC0PLaiTgIHGrC8OydUg0MUKrl8bxpHwUvjtXy0YA0BSV5y/DglkkKFmC07++vmhYNRqDi8cdxJwUrgtDy2ok4CBxqwvDsnFINDFCq5fG8aB8FL47V8tGAMwUlecvw4JVICBZgtO/vrpoWDUWg4vHHcSYFK4LQ8tqJOAgcasLw65xSDAxQquXxvGgfBS+O1fLRgDMFJXnL8OCUSAoWYLTv766ZFg1FoOLxxnElBSuB0PLaiTgIG2rB8OqcUgwMT6rl8btpHwUvjtby0IAzBSV5y/DglEgKFmC07++umRYNRaDi8cZxJQUrgdDy2ok4CBtqwfDqnFIMDE+q5fG7aR8FL47V8tCAMwUlecvw4JNHChZftO/vrpkWDUWg4fHGcSUFK4HQ8tmJOAgbasHw6pxSDAx');
    audio.play().catch(() => {});
  } catch (e) {
    // Ignore errors
  }
};

onMounted(() => {
  // Connect to WebSocket
  socketStore.connect();
  
  // Fetch initial data
  fetchMyApplications();
});

onUnmounted(() => {
  // Clean up notification
  showNotification.value = false;
});
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
    <!-- Real-time Status Update Notification -->
    <transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="transform translate-y-full opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="transform translate-y-0 opacity-100"
      leave-to-class="transform translate-y-full opacity-0"
    >
      <div 
        v-if="showNotification && statusUpdateNotification"
        class="fixed bottom-6 right-6 z-[9999] max-w-md"
      >
        <div 
          :class="statusUpdateNotification.new_status === 'ACCEPTED' 
            ? 'bg-gradient-to-r from-green-500 to-emerald-600' 
            : 'bg-gradient-to-r from-red-500 to-pink-600'"
          class="rounded-2xl shadow-2xl p-6 text-white animate-bounce-slow"
        >
          <div class="flex items-start gap-4">
            <div class="flex-shrink-0">
              <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center animate-pulse">
                <svg v-if="statusUpdateNotification.new_status === 'ACCEPTED'" class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
            <div class="flex-1">
              <h3 class="text-lg font-bold mb-1">
                {{ statusUpdateNotification.new_status === 'ACCEPTED' ? '🎉 Congratulations!' : '😔 Application Update' }}
              </h3>
              <p class="text-sm mb-2" :class="statusUpdateNotification.new_status === 'ACCEPTED' ? 'text-green-50' : 'text-red-50'">
                Your application for 
                <span class="font-semibold block mt-1">{{ statusUpdateNotification.shift_title }}</span>
                has been <span class="font-bold">{{ statusUpdateNotification.new_status }}</span>!
              </p>
              <button 
                @click="showNotification = false"
                class="text-xs bg-white/20 hover:bg-white/30 px-3 py-1 rounded-lg transition-colors"
              >
                Dismiss
              </button>
            </div>
            <button 
              @click="showNotification = false"
              class="flex-shrink-0 text-white/80 hover:text-white transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Modern Header -->
    <header class="bg-white border-b border-slate-200 sticky top-0 z-50 backdrop-blur-lg bg-white/90">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex justify-between items-center">
          <div>
            <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
              My Applications
            </h1>
            <p class="text-sm text-slate-500 mt-1">Track your shift applications</p>
          </div>
          <div class="flex gap-2">
            <button 
              @click="goToMap" 
              class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition-all duration-200 flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
              </svg>
              Find Shifts
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
      <!-- Loading -->
      <div v-if="loading" class="text-center py-20">
        <div class="inline-block w-12 h-12 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
        <p class="text-slate-600 mt-4">Loading applications...</p>
      </div>

      <!-- Applications List -->
      <div v-else-if="applications.length > 0" class="space-y-6">
        <div 
          v-for="app in applications" 
          :key="app.id"
          class="bg-white rounded-2xl shadow-lg border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-200"
        >
          <div class="p-6">
            <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4 mb-4">
              <div class="flex-1">
                <h3 class="text-xl font-bold text-slate-800 mb-3">
                  {{ app.shift_title }}
                </h3>
                
                <div class="space-y-2">
                  <div class="flex items-center gap-2 text-green-600 font-semibold">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Rp {{ app.shift_pay_rate?.toLocaleString() }}
                  </div>
                  <div class="flex items-center gap-2 text-sm text-slate-600">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    Applied {{ new Date(app.created_at).toLocaleDateString('en-US', { 
                      year: 'numeric', 
                      month: 'long', 
                      day: 'numeric' 
                    }) }}
                  </div>
                </div>
              </div>
              
              <span 
                :class="getStatusColor(app.status)"
                class="px-4 py-2 rounded-full text-sm font-semibold whitespace-nowrap"
              >
                {{ app.status }}
              </span>
            </div>

            <!-- Status-specific messages -->
            <div v-if="app.status === 'ACCEPTED'" class="mt-4 p-4 bg-green-50 border border-green-200 rounded-xl">
              <div class="flex items-start gap-3">
                <svg class="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p class="text-green-800 text-sm">
                  Congratulations! Your application has been accepted. The business will contact you soon.
                </p>
              </div>
            </div>
            
            <div v-if="app.status === 'REJECTED'" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-xl">
              <div class="flex items-start gap-3">
                <svg class="w-5 h-5 text-red-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p class="text-red-800 text-sm">
                  Unfortunately, your application was not selected. Keep applying to other shifts!
                </p>
              </div>
            </div>
            
            <div v-if="app.status === 'PENDING'" class="mt-4 p-4 bg-amber-50 border border-amber-200 rounded-xl">
              <div class="flex items-start gap-3">
                <svg class="w-5 h-5 text-amber-600 mt-0.5 flex-shrink-0 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p class="text-amber-800 text-sm">
                  Your application is under review. The business will respond soon.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-20">
        <svg class="w-24 h-24 mx-auto text-slate-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-xl font-semibold text-slate-700 mb-2">No Applications Yet</h3>
        <p class="text-slate-500 mb-6">Start exploring shifts on the map and apply to the ones you like!</p>
        <button 
          @click="goToMap"
          class="px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-200 shadow-lg shadow-blue-500/30 flex items-center gap-2 mx-auto"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          Browse Available Shifts
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes bounce-slow {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.animate-bounce-slow {
  animation: bounce-slow 2s ease-in-out 2;
}
</style>
