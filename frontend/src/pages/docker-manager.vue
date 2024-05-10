<!-- eslint-disable vue/valid-v-slot -->
<template>
    <div class="mt-2 mb-4 text-h4">Manage & Configure Stacks</div>
    <SortableStackTable :loading="loading" :items="dataStacks" :portainer-url-template="portainerStackUrl">
    </SortableStackTable>
</template>

<script lang="ts" setup>
import { Stack, Container } from "@/types/types";
import { ref, Ref, onMounted } from "vue";
import axios from "axios";
import { useSnackbarStore } from "@/store/snackbar";
import { useLocalStore } from "@/store/local";
import { storeToRefs } from "pinia";
const snackbarsStore = useSnackbarStore();
const localStore = useLocalStore();
const { urlConfig } = storeToRefs(localStore);

const dataStacks: Ref<Stack[]> = ref([]);
const loading: Ref<boolean> = ref(true);
const defaultEndpointId = process.env.PORTAINER_DEFAULT_ENDPOINT_ID || "1";

// computed properties
const portainerStackUrl = computed(() => {
    const prebuildUrl = urlConfig.value.defaultPortainerAddress + process.env.PORTAINER_QUERY_TEMPlATE;
    return prebuildUrl.replace("${endpointId}", defaultEndpointId);
});

onMounted(async () => {
    //testDataStacks.value = await generateTestData();
    try {
        await axios.post("/api/db/sync", { endpointIds: [defaultEndpointId] });
        const request = axios.get("/api/portainer/stacks", { params: { skeletonOnly: true } });
        const response = await request;
        // Change sorting to priority
        dataStacks.value = response.data.sort((a: Stack, b: Stack) => a.name.localeCompare(b.name));
        dataStacks.value = response.data.sort((a: Stack, b: Stack) => a.priority - b.priority);
    } catch (error) {
        console.error("Error fetching stacks", error);
        snackbarsStore.addSnackbar("load_containers", "Could not fetch stacks: " + error, "error");
    }

    loading.value = false;
});
</script>
