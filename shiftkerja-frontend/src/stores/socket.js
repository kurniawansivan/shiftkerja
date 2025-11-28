import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSocketStore = defineStore('socket', () => {
  const isConnected = ref(false);
  const messages = ref([]); // To store incoming data
  let socket = null;

  // 1. Connect Function
  const connect = () => {
    if (socket && socket.readyState === WebSocket.OPEN) return;

    console.log("🔌 Connecting to WebSocket...");
    socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => {
      console.log("✅ WebSocket Connected!");
      isConnected.value = true;
    };

    socket.onmessage = (event) => {
      // We will handle data here later
      console.log("📩 New Message:", event.data);
      messages.value.push(event.data);
    };

    socket.onclose = () => {
      console.log("❌ WebSocket Disconnected");
      isConnected.value = false;
      // Optional: Auto-reconnect logic could go here
    };

    socket.onerror = (error) => {
      console.error("WebSocket Error:", error);
    };
  };

  // 2. Send Function (e.g., sending GPS updates)
  const sendMessage = (msg) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(msg));
    }
  };

  return { isConnected, messages, connect, sendMessage };
});