<template>
    <router-view v-if="renderInitCompleted" />
</template>

<script lang="ts" setup>
import { useAppStore } from "@/store/app";
import { storeToRefs } from "pinia";

const appStore = useAppStore();
const { renderInitCompleted } = storeToRefs(appStore);

onMounted(() => {
    if (!renderInitCompleted.value) {
        setTimeout(() => {
            appStore.completeRenderInit();
        }, 200);
    }
});
</script>
