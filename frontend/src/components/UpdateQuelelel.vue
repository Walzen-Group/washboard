<template>
    <!-- class="fill-height" -->
    <v-sheet class="fill-height" rounded>
        <h3 class="pt-4 pl-3 pb-4">Update History</h3>

        <v-fade-transition mode="out-in">
            <div v-if="loading || hide" class="pa-2">
                <v-card class="mb-2" v-for="_ in 4" variant="tonal">
                    <v-skeleton-loader class="mx-auto"
                                       type="article">
                    </v-skeleton-loader>
                </v-card>
            </div>
            <div class="px-2 pb-2" v-else>
                <v-data-iterator
                                 :items="items"
                                 :items-per-page="itemsPerPage">
                    <template v-slot:default="{ items }">
                        <v-fade-transition group hide-on-leave>
                            <div
                                 v-for="item in (items as any)"
                                 :key="item.raw.stackId"
                                 cols="auto"
                                 md="4">
                                <v-card class="pb-1 pt-2 mb-2" border flat>
                                    <v-tooltip v-if="item.raw.details"
                                               activator="parent"
                                               location="start">{{ item.raw.details }}</v-tooltip>
                                    <v-list-item :subtitle="item.raw.status">
                                        <template v-slot:append>
                                            <v-icon v-if="item.raw.status == QueueStatus.done" size="35"
                                                    icon="mdi-robot-happy" class="mr-2" />
                                            <v-icon v-else-if="item.raw.status == QueueStatus.queued"
                                                    size="35"
                                                    icon="mdi-robot-confused-outline"
                                                    class="mr-2 loader" />
                                            <v-icon v-else-if="item.raw.status == QueueStatus.error"
                                                    size="35"
                                                    icon="mdi-robot-dead-outline" class="mr-2" />
                                            <v-icon v-else size="35" icon="mdi-robot" class="mr-2" />
                                        </template>
                                        <template v-slot:title>
                                            <strong class="text-h6">{{ item.raw.stackName }}</strong>
                                        </template>
                                        <template v-slot:subtitle>
                                            <v-row dense>
                                                <v-col cols="auto" xl="6">
                                                    <span :class="getColorClass(item.raw.status)">
                                                        {{ item.raw.status == QueueStatus.queued ?
                                                            'in progress' : item.raw.status }}
                                                    </span>
                                                </v-col>
                                                <v-spacer></v-spacer>
                                                <v-col cols="auto" xl="6">
                                                    <span> {{ timeAgo(item.raw.timestamp) }}</span>
                                                </v-col>
                                            </v-row>
                                        </template>
                                        <p class="mt-2"></p>
                                    </v-list-item>
                                </v-card>
                            </div>
                        </v-fade-transition>
                    </template>

                    <template v-slot:header="{ page, pageCount, prevPage, nextPage }">
                        <div class="d-flex align-center justify-center px-4 pb-4">
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
        </v-fade-transition>

    </v-sheet>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, ComputedRef } from 'vue';
import { VDataIterator } from 'vuetify/components';
import { QueueStatus, QueueItem, UpdateQueue } from '@/types/types';


function getColorClass(status: QueueStatus) {
    return "text-" + getColor(status);
}

function getColor(status: QueueStatus) {
    switch (status) {
        case QueueStatus.queued:
            return 'blue';
        case QueueStatus.done:
            return 'light-green-darken-2';
        case QueueStatus.error:
            return 'red';
        default:
            return undefined;
    }
}

function timeAgo(unixTimestampInSeconds: number) {
    const timestampInMilliseconds = unixTimestampInSeconds * 1000;
    const now = new Date().getTime();
    const secondsPast = (now - timestampInMilliseconds) / 1000;

    if (secondsPast < 60) { // Less than a minute
        const seconds = Math.floor(secondsPast);
        if (seconds === 1) {
            return `${seconds} second ago`;
        } else {
            return `${seconds} seconds ago`;
        }
    } else if (secondsPast < 3600) { // Less than an hour
        const minutes = Math.floor(secondsPast / 60);
        if (minutes === 1) {
            return `${minutes} minute ago`;
        } else {
            return `${minutes} minutes ago`;
        }
    } else if (secondsPast < 86400) { // Less than a day
        const hours = Math.floor(secondsPast / 3600);
        if (hours === 1) {
            return `${hours} hour ago`;
        } else {
            return `${hours} hours ago`;
        }
    } else {
        const days = Math.floor(secondsPast / 86400);
        if (days === 1) {
            return `${days} day ago`;
        } else {
            return `${days} days ago`;
        }
    }
}

const items: ComputedRef<QueueItem[]> = computed(() => {
    const items: QueueItem[] = [];
    for (let record of Object.values(props.queue) as Record<string, QueueItem>[]) {
        Object.values(record).forEach((item) => {
            items.push(item);
        });
    }
    items.sort((a, b) => {
        // Prioritize items with status 'queued'
        if (a.status === QueueStatus.queued && b.status !== QueueStatus.queued) {
            return -1;
        } else if (a.status !== QueueStatus.queued && b.status === QueueStatus.queued) {
            return 1;
        } else if (a.status === QueueStatus.error && b.status !== QueueStatus.error) {
            return -1;
        } else if (a.status !== QueueStatus.error && b.status === QueueStatus.error) {
            return 1;
        }

        // Then sort alphabetically by stackName
        if (a.timestamp == b.timestamp) {
            return a.stackName.localeCompare(b.stackName);
        }
        return a.timestamp > (b.timestamp) ? -1 : 1;
    });

    return items;
});

const props = defineProps<{
    loading: boolean
    queue: UpdateQueue
    itemsPerPage: number
    hide?: boolean
}>();
</script>

<style scoped>
.loader {
    animation: spin 2s linear infinite;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
}
</style>
