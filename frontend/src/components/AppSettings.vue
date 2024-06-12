<template>
    <div class="mt-2 mb-4 text-h4">Settings</div>

    <v-card>
        <v-card-title class="mt-4">Default Endpoint Host</v-card-title>
        <v-card-subtitle style="text-wrap: wrap;">Change the URL washboard uses to connect to a stack's web ui.</v-card-subtitle>
        <v-card-text class="pb-0">
            <v-text-field
                          v-model="defaultHost"
                          @update:model-value="handleEndpointFieldChange"
                          variant="solo-filled"
                          dense>
            </v-text-field>
        </v-card-text>

        <v-card-title class="pt-0">Portainer Link URL</v-card-title>
        <v-card-subtitle style="text-wrap: wrap;">Change the URL washboard uses when linking to Portainer.</v-card-subtitle>
        <v-card-text>
            <v-text-field
                          v-model="portainerAddress"
                          @update:model-value="handlePortainerAddressChange"
                          variant="solo-filled"
                          dense>
            </v-text-field>
        </v-card-text>
        <v-card-actions>
            <v-btn color="primary" @click="resetUrlConfig">Reset</v-btn>
        </v-card-actions>
    </v-card>
</template>

<script lang="ts" setup>
import { useLocalStore } from '@/store/local';
import { storeToRefs } from 'pinia';

const localStore = useLocalStore();
const { urlConfig } = storeToRefs(localStore);

const defaultHost = ref(urlConfig.value.defaultHost);
const portainerAddress = ref(urlConfig.value.defaultPortainerAddress);

function handleEndpointFieldChange() {
    localStore.updateURLConfig({
        defaultHost: defaultHost.value,
        defaultPortainerAddress: portainerAddress.value,
    });
}

function handlePortainerAddressChange() {
    localStore.updateURLConfig({
        defaultHost: defaultHost.value,
    });
}

function resetUrlConfig() {
    localStore.resetURLConfig()
    defaultHost.value = urlConfig.value.defaultHost;
    portainerAddress.value = urlConfig.value.defaultPortainerAddress;
}


</script>

<style lang="scss" scoped></style>
