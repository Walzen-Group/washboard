<template>
    <h2 class="mt-2 mb-4">Manage Stack Images</h2>
    <v-alert v-if="!containersNeedUpdate" variant="tonal" type="success" color="blue"
             title="You're all good"></v-alert>
    <v-alert v-else variant="tonal" type="warning" title="Updates available"></v-alert>

    <div class="d-flex justify-center">
        <v-row dense class="mt-2 mb-2">
            <v-col>
                <v-hover>
                    <template v-slot:default="{ isHovering, props }">
                        <v-card v-bind="props" :color="isHovering ? 'surface-variant' : undefined"
                                elevation="0"
                                variant="tonal" class="fill-height" min-width="220">
                            <template v-slot:append>
                                <v-icon icon="mdi-autorenew" size="x-large" color="warning"></v-icon>
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
                        <v-card v-bind="props" :color="isHovering ? 'surface-variant' : undefined"
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
    <StackTable @click:indicator="handleIndicatorClick" @update:selectedRows="updateSelectedRows" :items="items" :loading="loading">
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
                Do you want to update {{ totalStacksToUpdate }} stack{{ totalStacksToUpdate > 1 ? "s" : ""}}?
                <v-list lines="one">
                    <v-list-item
                        v-for="name in selectedStackNames"
                        :key="name"
                        :title="name"
                        density="compact"
                    ></v-list-item>
                </v-list>
            </v-card-text>
            <v-card-actions class="mb-2 mr-2">
                <v-spacer></v-spacer>
                <v-btn color="primary" variant="tonal" @click="updateSelected">Enqueue</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
    <div>
        <v-snackbar v-for="(snackbar, index) in snackbars" :key="index" :model-value="snackbar.show"
                    :color="snackbar.color"
                    :timeout=-1 @update:model-value="value => closeSnackbar(snackbar.id, value)"
                    :style="{ 'margin-bottom': calcSnackbarMargin(index) }" location="bottom right"
                    close-on-content-click
                    multiline>
            {{ snackbar.message }}
        </v-snackbar>
    </div>
</template>

<script lang="ts" setup>
import { updateStack } from '@/api/lib';
import StackTable from '@/components/StackTable.vue';
import axios from 'axios';
import { useLocalStore } from '@/store/local';
import { ref, Ref, onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia';


const localStore = useLocalStore();
const { dockerUpdateManagerSettings: dockerUpdateManagerSettings } = storeToRefs(localStore);

const dialogUpdate: Ref<boolean> = ref(false);
const loadingUpdateButton: Ref<boolean> = ref(false);
const currentProgress: Ref<number> = ref(0);
const totalStacksToUpdate: Ref<number> = ref(0);
const selectedRows: Ref<number[]> = ref([]);
const selectedStackNames: Ref<string[]> = ref([]);
const items: Ref<Stack[]> = ref([]);
const loading: Ref<boolean> = ref(true);
const snackbars: Ref<Snackbar[]> = ref([]);
const timeout = 5000;
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

            // TODO: remove orphaned containers from dockerUpdateManagerSettings.value.ignoredImages
            // iterate through dockeruPdateManagerSettings.value.ignoredImages and remove all that are not present in the current stacks
        })
        .catch((error) => {
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
    // let wsAddr = `${axios.defaults.baseURL}/ws/stacks-update`.replace('http://', 'ws://').replace('https://', 'wss://');
    // let socket = new WebSocket(wsAddr);
    // socket.onmessage = function (event) {
    //     let data = JSON.parse(event.data);
    //     console.log(data);
    // };
    for (let idx in selectedRowsValue) {
        const stackId = selectedRowsValue[idx];
        const stack = items.value.find((item: Stack) => item.id === stackId);
        if (stack?.containers.some((container: any) => container.upToDate === "outdated")) {
            try {
                const response = await updateStack(stackId);
                console.log(response);
                currentProgress.value += 1;
                addSnackbar(stackId, `Stack ${stack?.name} enqueued successfully`);
            } catch (error) {
                console.error(error);
            }
        } else {
            console.log(`No update necessary for stack ${stack?.name}`);
            addSnackbar(stackId, `No update necessary for Stack ${stack?.name}`, "warning");
        }
    }

    loadingUpdateButton.value = false;
    currentProgress.value = 0;
}

function confirmUpdateSelected() {
    totalStacksToUpdate.value = selectedRows.value.length;
    dialogUpdate.value = true;
}

function addSnackbar(id: number, message: string, color: string = "success", show: boolean = true) {
    snackbars.value.push({ id, message, color, show });
    setTimeout(() => {
        closeSnackbar(id, false);
    }, timeout);
}

function closeSnackbar(id: number, value: boolean) {
    if (!value) {
        setTimeout(() => {
            const index = snackbars.value.findIndex((snackbar: Snackbar) => snackbar.id === id);
            if (index !== -1) {
                snackbars.value.splice(index, 1);
            }
        }, 300);
    }
}

function calcSnackbarMargin(index: number) {
    return `${index * 60 + 10}px`;
}

function handleIndicatorClick(container: any) {
    if (container.image in dockerUpdateManagerSettings.value.ignoredImages) {
        delete dockerUpdateManagerSettings.value.ignoredImages[container.image];
        //let stack = items.value.find((stack: Stack) => stack.id === item.);
    } else {
        dockerUpdateManagerSettings.value.ignoredImages[container.image] = true;
    }
    localStore.updateDockerUpdateManagerSettings({
       ignoredImages: dockerUpdateManagerSettings.value.ignoredImages
    });
    console.log(dockerUpdateManagerSettings.value.ignoredImages);
}
</script>
