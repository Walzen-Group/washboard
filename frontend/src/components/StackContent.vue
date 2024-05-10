<template>
    <v-data-table
                  :headers="headers"
                  :items-per-page=-1
                  :items="containers"
                  no-data-text="Stack inactive"
                  item-key="name">
        <template #bottom></template>
        <template v-slot:item.name="{ item }">
            <v-row align="center" no-gutters dense class="d-flex flex-nowrap">
                <!-- Image/Icon placeholder -->
                <v-col cols="1" class="d-flex align-center stack-icon">
                    <v-img
                           v-if="item.labels['net.unraid.docker.icon']"
                           height="25"
                           :src="item.labels['net.unraid.docker.icon']"></v-img>
                    <v-icon v-else size="31">mdi-docker</v-icon>
                </v-col>
                <!-- Text placeholder next to the image -->
                <v-col cols="auto" class="d-flex align-center">
                    <span class="ml-2">{{ item.name }}</span>
                </v-col>
            </v-row>
        </template>

        <template v-slot:item.status="{ item }">
            <v-chip
                    :color="getContainerStatusCircleColor(item.status)"
                    rounded="20">
                {{ item.status }}
            </v-chip>
        </template>

        <template v-slot:item.ports="{ item }">
            {{ item.ports.map(p => p.split(":")[0]).join(", ") }}
        </template>

        <template v-slot:item.networks="{ item }">
            {{ item.networks.join(", ") }}
        </template>
    </v-data-table>
</template>

<script lang="ts" setup>
import { Container } from '@/types/types';
import { title } from 'process';
import { getContainerStatusCircleColor } from '@/api/lib';

const headers: any[] = [
    { title: 'Container', key: 'name' },
    { title: 'Status', key: 'status' },
    { title: 'Ports', key: 'ports' },
    { title: 'Networks', align: "end", key: 'networks' }
];

const props = defineProps<{ containers: Container[] }>();


</script>

<style scoped lang="scss"></style>
