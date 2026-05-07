import { ImageRefreshState } from "@/types/types";
import { defineStore } from "pinia";
import { ref, Ref } from "vue";

const STORE_NAME = "imageRefresh";

export const useImageRefreshStore = defineStore(STORE_NAME, () => {
    const running: Ref<boolean> = ref(false);
    const startedAt: Ref<number> = ref(0);
    const finishedAt: Ref<number> = ref(0);
    const error: Ref<string> = ref("");

    function update(payload: ImageRefreshState) {
        running.value = payload.running;
        startedAt.value = payload.startedAt;
        finishedAt.value = payload.finishedAt;
        error.value = payload.error;
    }

    return {
        running,
        startedAt,
        finishedAt,
        error,
        update,
    };
});
