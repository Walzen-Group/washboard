<template>
    <v-card>
        <template v-slot:text>
            <v-text-field v-model="search" label="Search" prepend-inner-icon="mdi-magnify" single-line variant="outlined"
                hide-details></v-text-field>
            <div class="d-flex flex-row mt-4">
                <v-btn variant="tonal" class="mr-2" @click="updateSelected" color="primary">Update Selected</v-btn>
                <v-btn variant="tonal" class="" @click="selectOutdated" color="primary">Select Outdated</v-btn>
            </div>
        </template>
        <div v-if="loading">
            <v-skeleton-loader type="list-item"></v-skeleton-loader>
            <v-skeleton-loader type="table-thead"></v-skeleton-loader>
            <v-skeleton-loader type="table-tbody"></v-skeleton-loader>
            <v-skeleton-loader type="table-tfoot"></v-skeleton-loader>
        </div>
        <v-data-table v-else items-per-page="-1" v-model:sort-by="sortBy" :search=search
            @update:sortBy="updateSorting" :headers="headers" v-model="selectedRows" :items="items" item-value="id"
            @update:modelValue="bulkSelect" show-select show-expand>
            <template v-slot:item.updateStatus="{ value }">
                <div class="d-flex flex-row">
                    <v-tooltip v-for="item in value" :text="item.name" location="top">
                        <template v-slot:activator="{ props }">
                            <v-icon size="x-large" v-bind="props" class="mr-0" :icon="getIcon(item.status)"
                                :color="getColor(item.status)"></v-icon>
                        </template>
                    </v-tooltip>
                </div>
            </template>
            <template v-slot:expanded-row="{ columns, item }">
                <tr>
                    <td :colspan="columns.length">
                        <!-- @vue-ignore -->
                        <v-data-table items-per-page="-1" :headers="containerTableHeaders" :items="item.containers">
                            <template v-slot:item.upToDate="{ value }">
                                <v-chip variant="tonal" :color="getColor(value)">
                                    {{ value.length > 0 ? value : "unavailable" }}
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
import { ref, onMounted, Ref, onUnmounted } from 'vue'
import axios from 'axios';

let loading: Ref<boolean> = ref(true);
let search: Ref<string> = ref("");
let selectedRows: Ref<number[]> = ref([]);
let shiftKeyOn: boolean = false;
let items: Ref<Stack[]> = ref([]);
let keyDownHandler: any;
let keyUpHandler: any;
let sortBy: Ref<any[]> = ref([{ key: 'name', order: 'asc' }]);
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
]
onMounted(() => {
    keyDownHandler = function ({ key }: any) {
        if (key == "Shift") shiftKeyOn = true;
    };
    keyUpHandler = function ({ key }: any) {
        if (key == "Shift") shiftKeyOn = false;
    };
    window.addEventListener("keydown", keyDownHandler);
    window.addEventListener("keyup", keyUpHandler);

    axios.get('/portainer-get-stacks')
        .then((response) => {
            items.value = response.data.map((stack: Stack) => {
                stack.updateStatus = [];
                for (let container of stack.containers) {
                    stack.updateStatus.push({ status: container.upToDate, name: container.name });
                }
                return stack;
            });
            loading.value = false;
        })
        .catch((error) => {
            console.log(error);
        });
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
    items.value.sort((a: any, b: any) => {
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

function getColor(status: string) {
    if (status === "outdated") return 'yellow-darken-3'
    else if (status === "updated") return 'updated'
    else if (status === "skipped") return 'grey'
    else if (status === "error") return 'red'
    else return 'grey'
}

function getIcon(status: string) {
    if (status === "outdated") return 'mdi-chevron-up-circle'
    else if (status === "updated") return 'mdi-checkbox-marked-circle'
    else if (status === "skipped") return 'mdi-minus-circle'
    else if (status === "error") return 'mdi-close-circle'
    else return 'mdi-circle'
}

function bulkSelect(e: any) {
    if (e && e.length > 1) {
        let currentSelected = e[e.length - 1];
        let lastSelected = e[e.length - 2];

        // console.log(`current selected: ${currentSelected}`);
        // console.log(`last selected: ${lastSelected}`);
        if (shiftKeyOn) {
            let start = items.value.findIndex((item: Stack) => item.id == lastSelected);
            let end = items.value.findIndex((item: any) => item.id == currentSelected);
            // console.log(start);
            // console.log(end);
            if (start - end > 0) {
                let temp = start;
                start = end;
                end = temp;
            }
            for (let i = start; i <= end; i++) {
                if (!selectedRows.value.includes(items.value[i].id)) {
                    selectedRows.value.push(items.value[i].id);
                }
            }
        }
    }
}

function selectOutdated() {
    selectedRows.value = []
    for (let item of items.value) {
        if (item.updateStatus.some((container: any) => container.status === "outdated")) {
            selectedRows.value.push(item.id);
        }
    }
}

function updateContainer(stackId: number) {
    console.log(stackId);
    axios.put('/portainer-update-stack', {
        pullImage: true,
        prune: false,
        endpointId: 1,
        stackId: stackId
    })
        .then((response) => {
            console.log(response);
        })
        .catch((error) => {
            console.log(error);
        });
}

function updateSelected() {
    console.log(selectedRows.value);
    for (let selectedRow of selectedRows.value) {
        updateContainer(selectedRow);
    }
}

</script>

<style lang="scss">
</style>
