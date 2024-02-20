<template>
    <v-sheet class="fill-height">
        <div v-if="loading" class="pa-2">
            <v-card class="mb-2" v-for="_ in 5" variant="tonal">
                <v-skeleton-loader class="mx-auto"
                                   max-width="360"
                                   type="article">
                </v-skeleton-loader>
            </v-card>
        </div>
        <div class="pa-2" v-else>
            <v-data-iterator
                             :items="items"
                             :items-per-page="3">
                <template v-slot:default="{ items }">
                    <div
                         v-for="item in (items as any)"
                         :key="item.raw.stackId"
                         cols="auto"
                         md="4">
                        <v-card class="pb-3" border flat>
                            <v-list-item class="mb-2" :subtitle="item.raw.status">
                                <template v-slot:title>
                                    <strong class="text-h6 mb-2">{{ item.raw.stackName }}</strong>
                                </template>
                            </v-list-item>

                            <div class="d-flex justify-space-between px-4">
                                <div
                                     class="d-flex align-center text-caption text-medium-emphasis me-1">
                                    <v-icon icon="mdi-clock" start></v-icon>
                                    <div class="text-truncate">{{ item.raw.details }}</div>
                                </div>

                                <v-btn
                                       border
                                       flat
                                       size="small"
                                       class="text-none"
                                       text="Read">
                                </v-btn>
                            </div>
                        </v-card>
                    </div>
                </template>

                <template v-slot:footer="{ page, pageCount, prevPage, nextPage }">
                    <div class="d-flex align-center justify-center pa-4">
                        <v-btn
                               :disabled="page === 1"
                               icon="mdi-arrow-left"
                               density="comfortable"
                               variant="tonal"
                               rounded
                               @click="prevPage"></v-btn>

                        <div class="mx-2 text-caption">
                            Page {{ page }} of {{ pageCount }}
                        </div>

                        <v-btn
                               :disabled="page >= pageCount"
                               icon="mdi-arrow-right"
                               density="comfortable"
                               variant="tonal"
                               rounded
                               @click="nextPage"></v-btn>
                    </div>
                </template>
            </v-data-iterator>
        </div>
    </v-sheet>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, ComputedRef } from 'vue';
import { VDataIterator } from 'vuetify/components';

const items: ComputedRef<QueueItem[]> = computed(() => {
    const items: QueueItem[] = [];
    for (let record of Object.values(props.queue) as Record<string, QueueItem>[]) {
        Object.values(record).forEach((item) => {
            items.push(item);
        });
    }
    return items;
});

const props = defineProps<{
    loading: boolean
    queue: UpdateQueue
}>();
</script>

<style scoped>
</style>
