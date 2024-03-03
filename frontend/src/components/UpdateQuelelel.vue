<template>
    <!-- class="fill-height" -->
    <v-sheet class="fill-height" rounded="lg">
        <h3 class="pt-4 pl-3 pb-4">Update History</h3>

        <v-fade-transition mode="out-in">
            <div v-if="loading || hide" class="pa-2">
                <v-card class="mb-2" v-for="(item, index) in 4" :key="index" variant="tonal">
                    <v-skeleton-loader class="mx-auto" type="article">
                    </v-skeleton-loader>
                </v-card>
            </div>
            <div class="px-2" v-else>
                <p class="mx-1 text-subtitle-2" v-if="items.length == 0">No stacks have been queued
                    yet...</p>
                <v-data-iterator v-else :items="items" :items-per-page="itemsPerPage">
                    <template v-slot:default="{ items }">
                        <v-fade-transition group hide-on-leave>
                            <div v-for="item in (items as any)" :key="item.raw.stackId" cols="auto"
                                 md="4">
                                <v-card class="pb-1 pt-2 mb-2" border flat>
                                    <v-tooltip v-if="item.raw.details" activator="parent"
                                               location="start">{{ item.raw.details }}</v-tooltip>
                                    <v-list-item :subtitle="item.raw.status">
                                        <template v-slot:append>
                                            <v-icon v-if="item.raw.status == QueueStatus.Done"
                                                    size="35" icon="mdi-robot-happy" color="primary"
                                                    class="mr-2" />
                                            <v-icon v-else-if="item.raw.status == QueueStatus.Queued"
                                                    size="35" icon="mdi-robot-confused-outline"
                                                    :class="`mr-2 loader`" />
                                            <v-icon v-else-if="item.raw.status == QueueStatus.Error"
                                                    size="35" icon="mdi-robot-dead-outline"
                                                    class="mr-2" />
                                            <v-icon v-else size="35" icon="mdi-robot"
                                                    class="mr-2" />
                                        </template>

                                        <template v-slot:title>
                                            <strong class="text-h6">{{ item.raw.stackName
                                                }}</strong>
                                        </template>

                                        <template v-slot:subtitle>
                                            <v-row dense>
                                                <v-col cols="auto" xl="6">
                                                    <span :class="getColorClass(item.raw.status)">
                                                        {{ item.raw.status == QueueStatus.Queued ?
                'in progress' : item.raw.status }}
                                                    </span>
                                                </v-col>
                                                <v-spacer></v-spacer>
                                                <v-col cols="auto" xl="6">
                                                    <v-tooltip activator="parent" location="start">
                                                        {{ new Date(item.raw.timestamp *
                1000).toLocaleString() }}
                                                    </v-tooltip>
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
                            <v-btn class="mr-1" :disabled="page === 1" icon="mdi-arrow-left"
                                   density="comfortable" variant="tonal" @click="prevPage"></v-btn>

                            <div class="mx-2 text-caption">
                                Page {{ page }} of {{ pageCount }}
                            </div>

                            <v-btn class="ml-1" :disabled="page >= pageCount" icon="mdi-arrow-right"
                                   density="comfortable" variant="tonal" @click="nextPage"></v-btn>
                        </div>
                    </template>
                </v-data-iterator>
            </div>
        </v-fade-transition>

    </v-sheet>
</template>

<script setup lang="ts">
import { computed, ComputedRef, Ref, ref } from 'vue';
import { VDataIterator } from 'vuetify/components';
import { QueueStatus, QueueItem, UpdateQueue } from '@/types/types';
import { onMounted } from 'vue';


const useLoaderStop: Ref<string | undefined> = ref(undefined);

function getColorClass(status: QueueStatus) {
    return "text-" + getColor(status);
}

function getColor(status: QueueStatus) {
    switch (status) {
        case QueueStatus.Queued:
            return 'blue';
        case QueueStatus.Done:
            return 'light-green-darken-2';
        case QueueStatus.Error:
            return 'red';
        default:
            return undefined;
    }
}

onMounted(() => {
    setInterval(() => {
        useLoaderStop.value = "loader-stop";
    }, 3000);
});

function timeAgo(unixTimestampInSeconds: number) {
    const timestampInMilliseconds = unixTimestampInSeconds * 1000;
    const now = new Date().getTime();
    const secondsPast = (now - timestampInMilliseconds) / 1000;

    if (secondsPast < 1) {
        return "now";
    } else if (secondsPast < 60) { // Less than a minute
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
        if (a.status === QueueStatus.Queued && b.status !== QueueStatus.Queued) {
            return -1;
        } else if (a.status !== QueueStatus.Queued && b.status === QueueStatus.Queued) {
            return 1;
        } else if (a.status === QueueStatus.Error && b.status !== QueueStatus.Error) {
            return -1;
        } else if (a.status !== QueueStatus.Error && b.status === QueueStatus.Error) {
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
    animation: spin 3s linear infinite;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
        /* left center right top bottom */
        transform-origin: center 56%;
    }

    100% {
        transform: rotate(360deg);
        transform-origin: center 56%;
    }
}
</style>
