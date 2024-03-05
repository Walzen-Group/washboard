<template>
    <v-card class="mt-1" variant="flat">
        <template v-slot:text>
            <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify"
                          single-line
                          variant="filled"
                          density="compact" hide-details></v-text-field>
            <div class="d-flex flex-row flex-wrap mt-4 ga-2">
                <v-btn width="200" variant="tonal" @click="selectOutdated" color="primary">Select
                    Outdated
                </v-btn>
                <slot name="controls"></slot>
            </div>
        </template>
        <div>
            <v-checkbox-btn class="ml-2 mb-5" v-model="showStoppedStacks"
                            @update:model-value="showInactiveStacks(showStoppedStacks)"
                            label="Show inactive stacks"></v-checkbox-btn>
        </div>
        <v-divider class="mt-2 mb-1"></v-divider>
        <!-- <v-fade-transition mode="out-in"> -->
        <div v-if="loading || mobileChromeLoader">
            <v-skeleton-loader type="list-item"></v-skeleton-loader>
            <v-skeleton-loader type="table-thead"></v-skeleton-loader>
            <v-skeleton-loader type="table-tbody"></v-skeleton-loader>
            <v-skeleton-loader type="table-tfoot"></v-skeleton-loader>
        </div>
        <v-data-table v-else :items-per-page="itemsPerPage" v-model:sort-by="sortBy" :search=search
                      @update:items-per-page="setItemsPerPage" @update:sortBy="updateSorting"
                      :headers="headers"
                      v-model="selectedRows" :items="itemsInternal" density="comfortable"
                      item-value="id"
                      @update:modelValue="bulkSelect" show-select show-expand>

            <template v-slot:item.name="{ item }">
                <v-row align="center" no-gutters dense class="d-flex flex-nowrap">
                    <v-col cols="1" class="d-flex align-center stack-icon">
                        <v-img v-if="getFirstContainerIcon(item.containers)" height="25"
                               :src="getFirstContainerIcon(item.containers)"></v-img>
                        <v-icon v-else-if="item.containers.length == 0"
                                size="29">mdi-power-plug-off-outline</v-icon>
                        <v-icon v-else size="31">mdi-docker</v-icon>
                    </v-col>
                    <v-col cols="auto" class="d-flex align-center">
                        <span class="ml-2">{{ item.name }}</span>
                    </v-col>
                </v-row>
            </template>

            <template v-slot:item.link="{ item }">
                <v-btn elevation="0" size="x-small" icon variant="text"
                       :href="getPortainerUrl(item)"
                       target="_blank"
                       class="mr-2">
                    <v-icon>mdi-open-in-new</v-icon>
                </v-btn>
            </template>

            <template v-slot:item.updateStatus="{ item }">
                <div class="d-flex flex-row">
                    <v-tooltip v-for="elem in item.containers" :text="elem.name" location="top"
                               :key="elem.name">
                        <template v-slot:activator="{ props }">
                            <v-icon class="clickable-indicator" size="x-large" v-bind="props"
                                    @click="indicatorClicked(elem)" :icon="getIcon(elem)"
                                    :color="getColor(elem)"></v-icon>
                        </template>
                    </v-tooltip>
                </div>
            </template>

            <template v-slot:expanded-row="{ columns, item }">
                <tr>
                    <td :colspan="columns.length">
                        <slot name="inner-actions" :item="item"></slot>
                        <v-card class="mb-3" border flat>
                            <v-data-table :loading="loaderState[item.id]" items-per-page="-1"
                                          density="comfortable"
                                          :headers="containerTableHeaders" :items="item.containers">
                                <template v-slot:item.upToDate="{ item }">
                                    <v-chip variant="tonal" :color="getColor(item)">
                                        {{ item.upToDate.length > 0 ? item.upToDate : 'n/a' }}
                                    </v-chip>
                                </template>

                                <template v-slot:item.name="{ item }">
                                    <v-row align="center" no-gutters dense
                                           class="d-flex flex-nowrap">
                                        <!-- Image/Icon placeholder -->
                                        <v-col cols="1" class="d-flex align-center stack-icon">
                                            <v-img v-if="item.labels['net.unraid.docker.icon']"
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

                                <template #bottom></template>
                            </v-data-table>

                        </v-card>
                    </td>
                </tr>
            </template>

        </v-data-table>
        <!-- </v-fade-transition>-->
    </v-card>
</template>

<script lang="ts" setup>
import { Container, Stack, ImageStatus, ContainerStatus } from '@/types/types';
import { ref, onMounted, Ref, onUnmounted, watch, reactive } from 'vue'
import { useDisplay } from 'vuetify';

let keyDownHandler: any;
let keyUpHandler: any;
let shiftKeyOn: boolean = false;

const { platform } = useDisplay();

const emit = defineEmits(["update:selectedRows", "click:indicator", "update:itemsPerPage", "update:stackModified"]);
const props = defineProps<{
    items: Stack[],
    loading: boolean,
    itemUrl: string
}>();

const defaultImage = "/img/container.png";
const isMobileChromeInitialized: Ref<boolean> = ref(false);
const mobileChromeLoader: Ref<boolean> = ref(false);
const initCompleted: Ref<boolean> = ref(false);
const itemsPerPage: Ref<number> = ref(-1);
const showStoppedStacks: Ref<boolean> = ref(false);
const search: Ref<string> = ref("");
const selectedRows: Ref<number[]> = ref([]);
const itemsInternal: Ref<Stack[]> = ref([]);
const sortBy: Ref<any[]> = ref([{ key: 'name', order: 'asc' }]);
const headers = [
    { title: "Stack Name", key: "name", value: "name" },
    { title: "Update Status", key: "updateStatus", value: "updateStatus" },
    { title: "ID", key: "id", value: "id" },
    { title: "Link", value: "link" }
];
const containerTableHeaders = [
    { title: "Image Status", value: "upToDate" },
    { title: "Name", key: "name", value: "name" },
    { title: "Status", value: "status" },
    { title: "Image", key: "image", value: "image" }
];
let loaderState: Record<string, boolean> = reactive({});

watch(() => props.items, async (newVal) => {
    console.log("platform is chrome mobile: ", platform.value.chrome && platform.value.android)
    // chrome mobile fix for laggy UI on page load
    if (platform.value.chrome && platform.value.android && !isMobileChromeInitialized.value) {
        mobileChromeLoader.value = true;
        await new Promise((resolve) => setTimeout(resolve, 500));
        isMobileChromeInitialized.value = true;
        mobileChromeLoader.value = false;
    }
    itemsInternal.value = newVal;
    loaderState = createLoaderState(itemsInternal.value);
    updateSorting(sortBy.value);
    showInactiveStacks(showStoppedStacks.value);
    if (!initCompleted.value) {
        emit("update:itemsPerPage", itemsInternal.value.length);
    }
    initCompleted.value = true;
});

watch(selectedRows, (newVal) => {
    emit("update:selectedRows", newVal);
});

watch(showStoppedStacks, () => {
    emitItemsPerPage(itemsPerPage.value);
});

onMounted(() => {
    keyDownHandler = function ({ key }: any) {
        if (key == "Shift") shiftKeyOn = true;
    };
    keyUpHandler = function ({ key }: any) {
        if (key == "Shift") shiftKeyOn = false;
    };
    window.addEventListener("keydown", keyDownHandler);
    window.addEventListener("keyup", keyUpHandler);
})

onUnmounted(() => {
    window.removeEventListener("keydown", keyDownHandler);
    window.removeEventListener("keyup", keyUpHandler);
}
);


function getFirstContainerIcon(containers: Container[]): string | undefined {
    // get random number between 0 and 1
    if (Math.random() < 0.995) {
        const ico = containers.find(container => container.labels['net.unraid.docker.icon'])?.labels['net.unraid.docker.icon'];
        return ico;
    } else {
        return `/img/craughing.png`;
    }
}

function createLoaderState(items: Stack[]) {
    const data: Record<string, boolean> = {};
    for (let item of items) {
        data[item.id] = false;
    }
    return reactive(data);
}

function showInactiveStacks(show: boolean) {
    if (show) {
        itemsInternal.value = props.items;
    } else {
        itemsInternal.value = props.items.filter((item: Stack) => item.containers.length > 0);
    }
}


function updateSorting(sortByRequest: any) {
    if (sortByRequest.length === 0) {
        sortByRequest = [{ key: 'name', order: 'asc' }];
    }
    if (sortByRequest[0].key == "updateStatus") {
        // check if any container has an image with status outdated, if so it should be at the top of th elist
        itemsInternal.value.sort((a: any, b: any) => {
            if (sortByRequest[0].order === "asc") {
                if (a.containers.some((container: Container) => container.upToDate === ContainerStatus.Outdated)) return -1;
                if (b.containers.some((container: Container) => container.upToDate !== ContainerStatus.Outdated)) return 1;
                return 0;
            } else {
                if (a.containers.some((container: Container) => container.upToDate === ContainerStatus.Outdated)) return 1;
                if (b.containers.some((container: Container) => container.upToDate !== ContainerStatus.Outdated)) return -1;
                return 0;
            }
        });
        return;
    }
    itemsInternal.value.sort((a: any, b: any) => {
        if (sortByRequest[0].order === "asc") {
            if (a[sortByRequest[0].key] < b[sortByRequest[0].key]) return -1;
            if (a[sortByRequest[0].key] > b[sortByRequest[0].key]) return 1;
            return 0;
        } else {
            if (a[sortByRequest[0].key] > b[sortByRequest[0].key]) return -1;
            if (a[sortByRequest[0].key] < b[sortByRequest[0].key]) return 1;
            return 0;
        }
    });
}

function getColor(elem: Container) {
    if (elem.upToDateIgnored) {
        return 'light-green-lighten-1';
    }
    if (elem.upToDate === ContainerStatus.Outdated) return 'yellow-darken-3'
    else if (elem.upToDate === ContainerStatus.Updated) return 'updated'
    else if (elem.upToDate === ContainerStatus.Skipped) return 'grey'
    else if (elem.upToDate === ContainerStatus.Error) return 'red'
    else return 'grey'
}

function getIcon(elem: Container) {
    if (elem.upToDateIgnored) {
        return 'mdi-pause-circle-outline';
    }
    if (elem.upToDate === ContainerStatus.Outdated) return 'mdi-chevron-up-circle-outline'
    else if (elem.upToDate === ContainerStatus.Updated) return 'mdi-check-circle-outline'
    else if (elem.upToDate === ContainerStatus.Skipped) return 'mdi-minus-circle-outline'
    else if (elem.upToDate === ContainerStatus.Error) return 'mdi-close-circle-outline'
    else return 'mdi-circle-outline'
}

function bulkSelect(e: any) {
    if (e && e.length > 1) {
        let currentSelected = e[e.length - 1];
        let lastSelected = e[e.length - 2];

        // console.log(`current selected: ${currentSelected}`);
        // console.log(`last selected: ${lastSelected}`);
        if (shiftKeyOn) {
            let start = itemsInternal.value.findIndex((item: Stack) => item.id == lastSelected);
            let end = itemsInternal.value.findIndex((item: any) => item.id == currentSelected);
            // console.log(start);
            // console.log(end);
            if (start - end > 0) {
                let temp = start;
                start = end;
                end = temp;
            }
            for (let i = start; i <= end; i++) {
                if (!showInactiveStacks && itemsInternal.value[i].containers.length === 0) continue;
                if (!selectedRows.value.includes(itemsInternal.value[i].id)) {
                    selectedRows.value.push(itemsInternal.value[i].id);
                }
            }
        }
    }
}

function indicatorClicked(elem: Container) {
    emit("click:indicator", elem);
}

function setItemsPerPage(e: number) {
    itemsPerPage.value = e;
    emitItemsPerPage(e);
}

function emitItemsPerPage(e: number) {
    if (itemsPerPage.value === -1 || itemsPerPage.value > itemsInternal.value.length) {
        emit("update:itemsPerPage", itemsInternal.value.length);
    } else {
        emit("update:itemsPerPage", e);
    }
}

function selectOutdated() {
    selectedRows.value = []
    for (let stack of itemsInternal.value as Stack[]) {
        if (stack.containers.some((container: Container) => container.upToDate === "outdated")) {
            selectedRows.value.push(stack.id);
        }
    }
}

function getPortainerUrl(item: Stack) {
    return props.itemUrl.replace('${stackId}', item.id.toString()).replace('${stackName}', item.name);
}
</script>

<style lang="scss" scoped>
.clickable-indicator {
    cursor: pointer;
}

.stack-icon {
    min-width: 28px;
    max-width: 28px;
}
</style>
