<template>
    <h2 class="mt-2 mb-4">Manage Stack Images</h2>
    <v-alert v-if="loading" variant="tonal" type="info"
             title="Refreshing...">
        <template v-slot:prepend>
            <v-progress-circular size="26" color="deep-blue-lighten-2" indeterminate></v-progress-circular>
        </template>
    </v-alert>
    <v-alert v-else-if="!containersNeedUpdate" variant="tonal" type="success" color="blue"
             title="You're all good"></v-alert>
    <v-alert v-else variant="tonal" type="warning" title="Updates available"></v-alert>

    <div class="d-flex justify-center">
        <v-row dense class="mt-2 mb-2">
            <v-col>
                <v-hover>
                    <template v-slot:default="{ isHovering, props }">
                        <v-skeleton-loader v-if="loading"
                                           class="mx-auto border"
                                           type="image">
                        </v-skeleton-loader>
                        <v-card v-else v-bind="props"
                                :color="isHovering ? 'surface-variant' : undefined"
                                elevation="0"
                                variant="tonal" class="fill-height" min-width="220">
                            <template v-slot:append>
                                <v-icon icon="mdi-autorenew" size="x-large"
                                        color="warning"></v-icon>
                            </template>
                            <template v-slot:title>
                                Can Be Updated
                            </template>

                            <v-card-text>
                                <h2>{{ updateStatusCounts.outdated }}</h2>
                            </v-card-text>
                        </v-card>
                    </template>
                </v-hover>
            </v-col>
            <v-col>
                <v-hover>
                    <template v-slot:default="{ isHovering, props }">
                        <v-skeleton-loader v-if="loading"
                                           class="mx-auto border"
                                           type="image">
                        </v-skeleton-loader>
                        <v-card v-else v-bind="props"
                                :color="isHovering ? 'surface-variant' : undefined"
                                elevation="0" variant="tonal" class="fill-height" min-width="220">
                            <template v-slot:append>
                                <v-icon icon="mdi-hand-okay" size="x-large" color="success"></v-icon>
                            </template>
                            <template v-slot:title>
                                Gucci
                            </template>

                            <v-card-text>
                                <h2>{{ updateStatusCounts.upToDate }}</h2>
                            </v-card-text>
                        </v-card>
                    </template>
                </v-hover>
            </v-col>
        </v-row>

    </div>
    <StackTable @click:indicator="handleIndicatorClick" @update:selectedRows="updateSelectedRows"
                :items="items" :loading="loading">
        <template v-slot:controls>
            <v-btn variant="tonal" @click="confirmUpdateSelected" color="primary"
                   :disabled="!selectedRows.length"
                   :loading="loadingUpdateButton">Update
                Selected</v-btn>
        </template>
    </StackTable>
    <v-dialog transition="dialog-top-transition" :scrim="false" v-model="dialogUpdate" width="auto">
        <v-card>
            <v-toolbar
                       color="primary"
                       class="d-flex justify-end"
                       density="compact"
                       title="Update stacks">
                <v-btn density="compact" icon="mdi-close" @click="dialogUpdate = false"></v-btn>
            </v-toolbar>
            <v-card-text class="mt-2">
                Do you want to update {{ totalStacksToUpdate }} stack{{ totalStacksToUpdate > 1 ? "s" :
                    "" }}?
                <v-list lines="one">
                    <v-list-item
                                 v-for="name in selectedStackNames"
                                 :key="name"
                                 :title="name"
                                 density="compact"></v-list-item>
                </v-list>
            </v-card-text>
            <v-card-actions class="mb-2 mr-2">
                <v-spacer></v-spacer>
                <v-btn color="primary" variant="tonal" @click="updateSelected">Enqueue</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts" setup>
import { updateStack } from '@/api/lib';
import StackTable from '@/components/StackTable.vue';
import axios from 'axios';
import { useLocalStore } from '@/store/local';
import { useSnackbarStore } from '@/store/snackbar';
import { ref, Ref, onMounted, computed } from 'vue';
import { storeToRefs } from 'pinia';


const localStore = useLocalStore();
const { dockerUpdateManagerSettings: dockerUpdateManagerSettings } = storeToRefs(localStore);

const snackbarsStore = useSnackbarStore();

const dialogUpdate: Ref<boolean> = ref(false);
const loadingUpdateButton: Ref<boolean> = ref(false);
const currentProgress: Ref<number> = ref(0);
const totalStacksToUpdate: Ref<number> = ref(0);
const selectedRows: Ref<number[]> = ref([]);
const selectedStackNames: Ref<string[]> = ref([]);
const items: Ref<Stack[]> = ref([]);
const loading: Ref<boolean> = ref(true);
const containersNeedUpdate = computed(() => {
    for (let stack of items.value) {
        if (stack.containers.some((container: any) => container.upToDate === "outdated")) {
            return true;
        }
    }
});
const updateStatusCounts = computed(() => {
    let outdated = 0;
    let upToDate = 0;
    for (let stack of items.value) {
        for (let container of stack.containers as Container[]) {
            if (container.upToDate === "outdated") {
                outdated += 1;
            } else {
                upToDate += 1;
            }
        }
    }
    return { outdated, upToDate };
});

onMounted(() => {
    axios.get('/portainer-get-stacks')
        .then((response) => {
            items.value = response.data;
            loading.value = false;

            for (let [ignoredImage] of Object.entries(dockerUpdateManagerSettings.value.ignoredImages)) {
                let found = true;
                for (let stack of items.value) {
                    for (let container of stack.containers) {
                        container.upToDateIgnored = false;
                        if (container.image === ignoredImage) {
                            found = true;
                            container.upToDateIgnored = true;
                            break;
                        }
                    }
                }
                if (!found) {
                    console.log(`Removing orphaned ignored image from ${ignoredImage}`)
                    delete dockerUpdateManagerSettings.value.ignoredImages[ignoredImage];
                }
                setIgnoreData();
            }

            // TODO: remove orphaned containers from dockerUpdateManagerSettings.value.ignoredImages
            // iterate through dockeruPdateManagerSettings.value.ignoredImages and remove all that are not present in the current stacks
        })
        .catch((error) => {
            loading.value = false;
            console.log(error);
        });
});

function updateSelectedRows(data: number[]) {
    selectedRows.value = data;
    selectedStackNames.value = items.value.filter((item: Stack) => data.includes(item.id)).map((item: Stack) => item.name);
}

async function updateSelected() {
    loadingUpdateButton.value = true;
    const selectedRowsValue = selectedRows.value;
    dialogUpdate.value = false;
    for (let idx in selectedRowsValue) {
        const stackId = selectedRowsValue[idx];
        const stack = items.value.find((item: Stack) => item.id === stackId);

        if (true || stack?.containers.some((container: any) => container.upToDate === "outdated" && !container.upToDateIgnored)) {
            try {
                const response = await updateStack(stackId);
                let data = response.data;
                switch (response.status) {
                    case 200:
                        currentProgress.value += 1;
                        snackbarsStore.addSnackbar(stackId, `Stack ${stack?.name} enqueued successfully`, "success");
                        break;
                    case 202:
                        snackbarsStore.addSnackbar(stackId, `Stack ${stack?.name} already queued`, "warning");
                        break;
                    default:
                        snackbarsStore.addSnackbar(stackId, `Failed to enqueue stack ${stack?.name}: ${data.error}`, "error");
                }
            } catch (error) {
                console.error(error);
            }
        } else {
            console.log(`No update necessary for stack ${stack?.name}`);
            snackbarsStore.addSnackbar(stackId, `No update necessary for stack ${stack?.name}`, "info");
        }
    }

    loadingUpdateButton.value = false;
    currentProgress.value = 0;
}

function confirmUpdateSelected() {
    totalStacksToUpdate.value = selectedRows.value.length;
    dialogUpdate.value = true;
}

function handleIndicatorClick(container: any) {
    if (container.image in dockerUpdateManagerSettings.value.ignoredImages) {
        delete dockerUpdateManagerSettings.value.ignoredImages[container.image];
    } else {
        dockerUpdateManagerSettings.value.ignoredImages[container.image] = true;
    }
    localStore.updateDockerUpdateManagerSettings({
        ignoredImages: dockerUpdateManagerSettings.value.ignoredImages
    });
    setIgnoreData();
}

function setIgnoreData() {
    for (let stack of items.value) {
        for (let container of stack.containers) {
            if (container.image in dockerUpdateManagerSettings.value.ignoredImages) {
                container.upToDateIgnored = true;
            } else {
                container.upToDateIgnored = false;
            }
        }
    }
}
</script>
