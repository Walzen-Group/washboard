// Utilities
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', () => {
  const renderInitCompleted: Ref<boolean> = ref(false);

  function completeRenderInit() {
      renderInitCompleted.value = true;
  }

  return {
      completeRenderInit,
      renderInitCompleted
  };
})
