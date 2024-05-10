<!-- eslint-disable vue/valid-v-slot -->
<template>
    <div class="mt-2 mb-4 text-h4">Manage & Configure Stacks</div>
    <SortableStackTable :loading="loading" :items="dataStacks" :portainer-url-template="portainerStackUrl"></SortableStackTable>
</template>

<script lang="ts" setup>
import { Stack, Container } from "@/types/types";
import { ref, Ref, onMounted } from "vue";
import axios from "axios";

const dataStacks: Ref<Stack[]> = ref([]);
const loading: Ref<boolean> = ref(true);
const defaultEndpointId = process.env.PORTAINER_DEFAULT_ENDPOINT_ID || "1";

// computed properties
const portainerStackUrl = computed(() => {
    return (
        process.env.PORTAINER_BASE_URL?.replace("${endpointId}", defaultEndpointId) || process.env.BASE_URL || ""
    );
});

onMounted(async () => {
    //testDataStacks.value = await generateTestData();
    const request = axios.get("/portainer/stacks", { params: { skeletonOnly: true } });
    const response = await request;
    // Change sorting to priority
    dataStacks.value = response.data.sort((a: Stack, b: Stack) => a.name.localeCompare(b.name));
    dataStacks.value = response.data.sort((a: Stack, b: Stack) => a.priority - b.priority);
    loading.value = false;
});
</script>
