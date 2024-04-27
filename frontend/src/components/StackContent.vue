<template>
    <v-container>
        <v-row>
            <v-col>
                    <v-card-subtitle>Name</v-card-subtitle>
            </v-col>
            <v-col>
                    <v-card-subtitle>Status</v-card-subtitle>
            </v-col>
            <v-col>
                    <v-card-subtitle>Ports</v-card-subtitle>
            </v-col>
            <v-col>
                    <v-card-subtitle>Networks</v-card-subtitle>
            </v-col>
        </v-row>
    </v-container>
    <v-card v-for="container in containers" :key="container.id" class="pb-2 pt-2 mb-2">
        <v-row>
            <v-col>
                <v-card-text>{{ container.name }}</v-card-text>
            </v-col>
            <v-col>
                <v-card-text>{{ container.status }}</v-card-text>
            </v-col>
            <v-col>
                <v-card-text>{{ container.ports.join(", ") }}</v-card-text>
            </v-col>
            <v-col>
                <v-card-text>{{ container.networks.join(", ") }}</v-card-text>
            </v-col>
        </v-row>

    </v-card>
    Da mal was rein
</template>

<script lang="ts" setup>
import { Container } from '@/types/types';

const props = defineProps<{ containers: Container[] }>();


</script>
