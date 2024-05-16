// Utilities
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', () => {
  const renderInitCompleted: Ref<boolean> = ref(false);
  const webSocketStacksUpdate = ref<WebSocket | null>(null);

  function completeRenderInit() {
      renderInitCompleted.value = true;
  }

  return {
      completeRenderInit,
      renderInitCompleted,
      webSocketStacksUpdate
  };
})
