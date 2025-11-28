<script setup>
import { RouterView } from 'vue-router';
</script>

<template>
  <div class="h-screen w-screen flex flex-col">
    <!-- Header -->
    <header class="bg-blue-900 text-white p-4 shadow-lg z-10">
      <h1 class="text-xl font-bold">ShiftKerja 🇮🇩</h1>
    </header>

    <!-- Main Content (Dynamic Router View) -->
    <main class="flex-grow relative bg-gray-100">
      <!-- 
        RouterView is a placeholder. 
        If path is '/', it puts MapData here.
        If path is '/login', it puts LoginView here.
      -->
      <RouterView />
    </main>
  </div>
</template>