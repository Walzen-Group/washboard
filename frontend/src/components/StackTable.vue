<template>
    <v-card variant="flat">
        <template v-slot:text>
            <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" single-line
                          variant="filled" density="compact"
                          hide-details></v-text-field>
            <div class="d-flex flex-row mt-4">
                <v-btn variant="tonal" class="mr-2" @click="selectOutdated" color="primary">Select
                    Outdated</v-btn>
                <slot name="controls"></slot>
            </div>
        </template>
        <div v-if="loading">
            <v-skeleton-loader type="list-item"></v-skeleton-loader>
            <v-skeleton-loader type="table-thead"></v-skeleton-loader>
            <v-skeleton-loader type="table-tbody"></v-skeleton-loader>
            <v-skeleton-loader type="table-tfoot"></v-skeleton-loader>
        </div>
        <v-data-table v-else items-per-page="-1" v-model:sort-by="sortBy" :search=search
                      @update:sortBy="updateSorting"
                      :headers="headers" v-model="selectedRows" :items="itemsInternal" item-value="id"
                      @update:modelValue="bulkSelect"
                      show-select show-expand>
            <template v-slot:item.updateStatus="{ item }">
                <div class="d-flex flex-row">
                    <v-tooltip v-for="elem in item.containers" :text="elem.name" location="top">
                        <template v-slot:activator="{ props }">
                            <v-icon class="clickable-indicator" size="x-large" v-bind="props"
                                    @click="indicatorClicked(elem)"
                                    :icon="getIcon(elem)"
                                    :color="getColor(elem)"></v-icon>
                        </template>
                    </v-tooltip>
                </div>
            </template>
            <template v-slot:expanded-row="{ columns, item }">
                <tr>
                    <td :colspan="columns.length">
                        <!-- @vue-ignore -->
                        <v-data-table items-per-page="-1" :headers="containerTableHeaders"
                                      :items="item.containers">
                            <template v-slot:item.upToDate="{ item }">
                                <v-chip variant="tonal" :color="getColor(item)">
                                    {{ item.upToDate.length > 0 ? item.upToDate : "unavailable" }}
                                </v-chip>
                            </template>
                            <template #bottom></template>
                        </v-data-table>
                    </td>
                </tr>
            </template>

        </v-data-table>
    </v-card>
</template>

<script lang="ts" setup>
import { ref, onMounted, Ref, onUnmounted, watch } from 'vue'

let keyDownHandler: any;
let keyUpHandler: any;
let shiftKeyOn: boolean = false;

const emit = defineEmits(["update:selectedRows", "click:indicator"]);
const props = defineProps<{
    items: Stack[],
    loading: boolean
}>();

const search: Ref<string> = ref("");
const selectedRows: Ref<number[]> = ref([]);
const itemsInternal: Ref<Stack[]> = ref([]);
const sortBy: Ref<any[]> = ref([{ key: 'name', order: 'asc' }]);
const headers = [
    { title: "Stack Name", key: "name", value: "name" },
    { title: "Update Status", value: "updateStatus" },
    { title: "ID", key: "id", value: "id" }
];
const containerTableHeaders = [
    { title: "Image Status", value: "upToDate" },
    { title: "Name", key: "name", value: "name" },
    { title: "Status", value: "status" },
    { title: "Image", key: "image", value: "image" },
];

watch(() => props.items, (newVal, _) => {
    itemsInternal.value = newVal;
    updateSorting(sortBy.value);
});

watch(selectedRows, (newVal, _) => {
    emit("update:selectedRows", newVal);
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

function updateSorting(sortByRequest: any) {
    if (sortByRequest.length === 0) {
        sortByRequest = [{ key: 'name', order: 'asc' }];
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
    if (elem.upToDate === "outdated") return 'yellow-darken-3'
    else if (elem.upToDate === "updated") return 'updated'
    else if (elem.upToDate === "skipped") return 'grey'
    else if (elem.upToDate === "error") return 'red'
    else return 'grey'
}

function getIcon(elem: Container) {
    if (elem.upToDateIgnored) {
        return 'mdi-pause-circle-outline';
    }
    if (elem.upToDate === "outdated") return 'mdi-chevron-up-circle-outline'
    else if (elem.upToDate === "updated") return 'mdi-check-circle-outline'
    else if (elem.upToDate === "skipped") return 'mdi-minus-circle-outline'
    else if (elem.upToDate === "error") return 'mdi-close-circle-outline'
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

function selectOutdated() {
    selectedRows.value = []
    for (let stack of itemsInternal.value as Stack[]) {
        if (stack.containers.some((container: Container) => container.upToDate === "outdated")) {
            selectedRows.value.push(stack.id);
        }
    }
}
</script>

<style lang="scss" scoped>
.clickable-indicator {
    cursor: pointer;
}
</style>
