<template>
    <div class="mt-2 mb-4 text-h4">Update Stack Images</div>
    <div class="mb-2">
        <v-alert v-if="loading" variant="tonal" type="info"
                 title="Refreshing...">
            <template v-slot:prepend>
                <v-progress-circular size="26" color="deep-blue-lighten-2"
                                     indeterminate></v-progress-circular>
            </template>
        </v-alert>
        <v-alert v-else-if="connectionFailed" variant="tonal" type="error"
                 title="No data"></v-alert>
        <v-alert v-else-if="!containersNeedUpdate" variant="tonal" type="success" color="blue"
                 title="You're all good"></v-alert>

        <v-alert v-else variant="tonal" type="warning" title="Updates available"></v-alert>
    </div>
    <v-row dense>
        <v-col cols="12" lg="9">
            <div class="d-flex justify-center">
                <v-row dense class="mb-0">
                    <v-col>
                        <v-hover>
                            <template v-slot:default="{ isHovering, props }">
                                <v-skeleton-loader v-if="loading"
                                                   class="mx-auto border"
                                                   type="image">
                                </v-skeleton-loader>
                                <v-card v-else v-bind="props"
                                        :color="isHovering ? undefined : 'surface-variant'"
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
                                        <h2>{{ tweeenedOutdated.number.toFixed(0) }}</h2>
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
                                        :color="isHovering ? undefined : 'surface-variant'"
                                        elevation="0" variant="tonal" class="fill-height"
                                        min-width="220">
                                    <template v-slot:append>
                                        <v-icon icon="mdi-hand-okay" size="x-large"
                                                color="success"></v-icon>
                                    </template>
                                    <template v-slot:title>
                                        Gucci
                                    </template>

                                    <v-card-text>
                                        <h2>{{ tweeenedUpToDate.number.toFixed(0) }}</h2>
                                    </v-card-text>
                                </v-card>
                            </template>
                        </v-hover>
                    </v-col>
                </v-row>

            </div>
            <StackTable @click:indicator="handleIndicatorClick"
                        @update:selectedRows="updateSelectedRows"
                        @update:items-per-page="calculateItemsPerPage"
                        @update:stack-modified="leeroad"
                        :item-url="portainerStackUrl"
                        :items="items" :loading="loading">
                <template v-slot:controls>
                    <v-btn width="200" variant="tonal" @click="confirmUpdateSelected" color="primary"
                           :disabled="!selectedRows.length"
                           :loading="loadingUpdateButton">
                        Update Selected
                    </v-btn>
                </template>
                <template v-slot:inner-actions="{ item }">
                    <div class="d-flex flex-row flex-wrap mt-4">
                            <v-btn v-if="item.containers.length === 0"
                                   :loading="loaderState[item.id]" class="mr-2 mb-2" variant="tonal"
                                   prepend-icon="mdi-arrow-right-drop-circle-outline"
                                   @click="startOrStopStack(item, 'start')">Start Stack</v-btn>
                            <v-btn v-else :loading="loaderState[item.id]" class="mr-2 mb-2"
                                   variant="tonal" prepend-icon="mdi-stop-circle-outline"
                                   @click="startOrStopStack(item, 'stop')">Stop Stack</v-btn>
                        </div>
                </template>
            </StackTable>
        </v-col>
        <v-col cols="12" lg="3">
            <UpdateQuelelel :loading="loading" :queue="queue"
                            :itemsPerPage="updateWidgetItemsPerPage" :hide="hideWidget">
            </UpdateQuelelel>
        </v-col>
    </v-row>


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
                <v-virtual-scroll
                                  class="mt-2"
                                  :max-height="200"
                                  :items="selectedStackNames">
                    <template v-slot:default="{ item }">
                        {{ item }}
                    </template>
                </v-virtual-scroll>

            </v-card-text>
            <v-card-actions class="mb-2 mr-2">
                <v-spacer></v-spacer>
                <v-btn color="primary" variant="tonal" @click="updateSelected">Enqueue</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts" setup>
import StackTable from '@/components/StackTable.vue';
import UpdateQuelelel from '@/components/UpdateQuelelel.vue';
import axios from 'axios';
import gsap from 'gsap';
import { startStack, stopStack, updateStack, getContainers } from '@/api/lib';
import { useLocalStore } from '@/store/local';
import { useSnackbarStore } from '@/store/snackbar';
import { useUpdateQuelelelStore } from '@/store/updateQuelelel';
import { storeToRefs } from 'pinia';
import { ref, Ref, onMounted, computed, watch, reactive } from 'vue';
import { Stack, Container, UpdateQueue, QueueStatus, ImageStatus, ContainerStatus } from '@/types/types';

const defaultEndpointId = process.env.PORTAINER_DEFAULT_ENDPOINT_ID || "1";

// stores
const localStore = useLocalStore();
const { dockerUpdateManagerSettings: dockerUpdateManagerSettings } = storeToRefs(localStore);

const updateQuelelelStore = useUpdateQuelelelStore();
const { queue, queueCount } = storeToRefs(updateQuelelelStore);

const snackbarsStore = useSnackbarStore();


let loaderState: Record<string, boolean> = reactive({});

// update card values
const tweeenedOutdated = reactive({ number: 0 });
const tweeenedUpToDate = reactive({ number: 0 });

// widget controls
const hideWidget: Ref<boolean> = ref(false);
const updateWidgetItemsPerPage: Ref<number> = ref(4);

const connectionFailed: Ref<boolean> = ref(false);
const dialogUpdate: Ref<boolean> = ref(false);
const loadingUpdateButton: Ref<boolean> = ref(false);
const currentProgress: Ref<number> = ref(0);
const totalStacksToUpdate: Ref<number> = ref(0);
const selectedRows: Ref<number[]> = ref([]);
const selectedStackNames: Ref<string[]> = ref([]);
const items: Ref<Stack[]> = ref([]);
const loading: Ref<boolean> = ref(true);


// computed properties
const portainerStackUrl = computed(() => {
    return process.env.PORTAINER_BASE_URL?.replace("${endpointId}", defaultEndpointId) || process.env.BASE_URL || "";
});

const containersNeedUpdate = computed(() => {
    for (let stack of items.value) {
        if (stack.containers.some((container: any) => container.upToDate === ImageStatus.Outdated)) {
            return true;
        }
    }
    return false;
});


// hooks
onMounted(() => {
    leeroad();
});

// watches
watch(queueCount, (newVal, oldVal) => {
    // TODO: check if this needs a delay
    if (newVal === 0 && oldVal !== 0) {
        leeroad();
    }
});

// functions

async function startOrStopStack(stack: Stack, action: string) {
    loaderState[stack.id] = true;
    if (!['start', 'stop'].includes(action)) {
        throw new Error(`Action should be "start" or "stop", got "${action}"`);
    }
    try {
        const response = await (action === 'start' ? startStack(stack.id) : stopStack(stack.id));
        await handleResponse(stack, action, response);

    } catch (error: any) {
        snackbarsStore.addSnackbar(`${stack.id}_startstop`, `Failed to ${action} ${stack?.name}: ${error.message}`, "error");
    } finally {
        loaderState[stack.id] = false;
    }
}

async function handleResponse(stack: Stack, action: string, response: any) {
    if (response.status === 200) {
        if (action === 'start') {
            await updateContainersOnStart(stack);
        } else {
            clearStackContainers(stack);
        }
        snackbarsStore.addSnackbar(`${stack.id}_startstop`, `Successfully ${action}ed ${stack?.name}`, "success");
    } else {
        throw new Error(`Received unexpected response status: ${response.status}`);
    }
}

async function updateContainersOnStart(stack: Stack) {
    const containersResponse = await getContainers(stack.name, parseInt(defaultEndpointId, 10));
    let containers: Container[] = containersResponse.data;
    containers = await Promise.all(containers.map(async container => updateContainerStatus(container)));
    items.value = items.value.map(item => (item.id === stack.id ? { ...item, containers } : item));
}

async function updateContainerStatus(container: Container) {
    const containerImageStatusResponse = await axios.get(`/portainer-get-image-status`, {
        params: { endpointId: defaultEndpointId, containerId: container.id }
    });
    const containerImageStatus = containerImageStatusResponse.data;
    return { ...container, upToDate: containerImageStatus.status };
}

function clearStackContainers(stack: Stack) {
    items.value = items.value.map(item => (item.id === stack.id ? { ...item, containers: [] } : item));
}


function updateStatusCounts() {
    let outdated = 0;
    let upToDate = 0;
    for (let stack of items.value) {
        //console.log(`stack name: ${stack.name}`);
        for (let container of stack.containers as Container[]) {
            if (container.upToDate === "outdated") {
                outdated += 1;
            } else {
                upToDate += 1;
            }
        }
    }
    setTimeout(() => {
        gsap.to(tweeenedOutdated, { duration: 0.5, number: Number(outdated) || 0 })
        gsap.to(tweeenedUpToDate, { duration: 0.5, number: Number(upToDate) || 0 })
    }, 200);
    return { outdated, upToDate };
}



function leeroad() {
    axios.get('/portainer-get-stacks')
        .then((response) => {
            console.log("leeroaded");
            items.value = response.data;
            loading.value = false;
            connectionFailed.value = false;

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
            updateStatusCounts();


            // TODO: remove orphaned containers from dockerUpdateManagerSettings.value.ignoredImages
            // iterate through dockeruPdateManagerSettings.value.ignoredImages and remove all that are not present in the current stacks
        })
        .catch((error) => {
            loading.value = false;
            connectionFailed.value = true;
            console.log(error);
        });
}




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

        // eslint-disable-next-line no-constant-condition
        if (true || stack?.containers.some((container: any) => container.upToDate === "outdated" && !container.upToDateIgnored)) {
            try {
                const response = await updateStack(stackId);
                let data = response.data;
                switch (response.status) {
                    case 200:
                        currentProgress.value += 1;
                        // snackbarsStore.addSnackbar(stackId, `Stack ${stack?.name} enqueued successfully`, "info");
                        break;
                    case 202:
                        snackbarsStore.addSnackbar(`${stackId}_queued`, `Stack ${stack?.name} already queued`, "warning");
                        break;
                    default:
                        snackbarsStore.addSnackbar(`${stackId}_error`, `Failed to enqueue stack ${stack?.name}: ${data.error}`, "error");
                }
            } catch (error) {
                console.error(error);
            }
        } else {
            console.log(`No update necessary for stack ${stack?.name}`);
            snackbarsStore.addSnackbar(`${stackId}_noup`, `No update necessary for stack ${stack?.name}`, "info");
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

function calculateItemsPerPage(itemsPerPage: number) {
    let val = Math.round(itemsPerPage * 0.55);
    updateWidgetItemsPerPage.value = val + 2;

    if (!loading.value) {
        hideWidget.value = true;
        setTimeout(() => {
            hideWidget.value = false;
        }, 200);
    }
    return
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function generateTestValues(num: number) {
    const testItems: UpdateQueue = {
        "done": {},
        "error": {},
        "queued": {}
    };

    for (let i = 0; i < num; i++) {
        const status = Math.random() > 0.5 ? QueueStatus.Done : Math.random() > 0.5 ? QueueStatus.Error : QueueStatus.Queued;
        const stackName = `test-${i}`;
        const stackId = i;
        const timestamp = Math.floor(Math.random() * 1000000000);
        testItems[status][stackName] = {
            endpointId: 1,
            stackId: stackId,
            stackName: stackName,
            status: status,
            details: "Ich bin der Hämmerer",
            timestamp: timestamp
        };
    }
    return testItems;
}
</script>
