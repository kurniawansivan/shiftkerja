<script setup>
import MapData from './components/MapData.vue'
</script>

<template>
  <div class="h-screen w-screen flex flex-col">
    <header class="bg-blue-900 text-white p-4 shadow-lg z-10">
      <h1 class="text-xl font-bold">ShiftKerja 🇮🇩</h1>
    </header>

    <main class="flex-grow relative">
      <MapData />
    </main>
  </div>
</template>