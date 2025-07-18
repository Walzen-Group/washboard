<template>
    <!--<div>{{ stacksInternal.map((stack) => stack.name).join(", ") }}</div>-->


    <v-expansion-panels v-model="panel">
        <v-expansion-panel>
            <v-expansion-panel-title>
                <v-card-title class="pl-0 ml-0">Bulk Operations</v-card-title>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
                <div class="ml-2 mb-3">
                    <v-card class="pt-2" variant="flat">
                        <v-card-title>Manage Selected Stacks</v-card-title>
                        <v-card-text>
                            <div class="d-flex flex-wrap ga-3">
                                <v-btn
                                       :loading="stackOperationInProgress"
                                       variant="tonal"
                                       prepend-icon="mdi-arrow-right-drop-circle-outline"
                                       @click="ModifySelectedContainerStatus(Action.Start)">Start
                                </v-btn>
                                <v-btn
                                       :loading="stackOperationInProgress"
                                       color="stop"
                                       variant="tonal"
                                       prepend-icon="mdi-stop-circle-outline"
                                       @click="ModifySelectedContainerStatus(Action.Stop)">Stop
                                </v-btn>
                                <v-btn
                                       :loading="stackOperationInProgress"
                                       variant="tonal"
                                       color="green-lighten-1"
                                       prepend-icon="mdi-restart"
                                       @click="ModifySelectedContainerStatus(Action.Restart)">Restart
                                </v-btn>
                            </div>

                        </v-card-text>
                        <v-card-title>
                            Select Options
                        </v-card-title>

                        <v-card-text>
                            <div class="d-flex ga-3">
                                <v-btn
                                       :loading="stackOperationInProgress"
                                       variant="tonal"
                                       prepend-icon="mdi-auto-mode"
                                       @click="handleSelectAutostartContainers">Autostart
                                </v-btn>
                                <v-btn
                                       :loading="stackOperationInProgress"
                                       variant="tonal"
                                       color="stop"
                                       prepend-icon="mdi-notification-clear-all"
                                       @click="clearSelection">Clear
                                </v-btn>
                            </div>
                        </v-card-text>
                    </v-card>
                </div>
            </v-expansion-panel-text>
        </v-expansion-panel>
    </v-expansion-panels>

    <div class="d-flex mt-4">
        <v-switch color="blue" @update:model-value="toggleEditMode" inset density="compact" hide-details
                  label="Change Autostart Order"></v-switch>
    </div>

    <v-divider class="mb-2 mt-4" />

    <div class="mt-4 mb-4">
        <!-- serach box-->
        <v-text-field v-if="!orderEditMode" v-model="search" :disabled="orderEditMode" label="Search" hide-details variant="filled"></v-text-field>
    </div>

    <div v-if="loading || searchActive">
        <v-skeleton-loader v-for="index in 10" :key="index" class="mb-2"
                           type="list-item-two-line">
        </v-skeleton-loader>
    </div>
    <div ref="sortableRoot">
        <v-fade-transition group :disabled="changeOrderTransition" v-if="!loading">
            <v-card v-for="element in stacksInternal" class="pb-2 pt-3 mb-2" :key="element.id">
                <v-row dense :class="[paddingClass[element.name]]">
                    <v-col cols="auto" class="ml-2">
                        <v-sheet class="fill-height d-flex flex-column">
                            <Transition name="slide-fade-up">
                                <v-btn
                                       v-show="element.expanded && orderEditMode"
                                       icon
                                       size="40"
                                       elevation="0"
                                       variant="text"
                                       density="compact"
                                       :ripple="false"
                                       :disabled="isFirstElement(element)"
                                       @click="moveElementAndUpdate(element, 'up')">
                                    <v-icon size="40"> mdi-chevron-up </v-icon>
                                </v-btn>
                            </Transition>

                            <v-sheet class="d-flex fill-height flex-column">
                                <v-checkbox-btn class="my-auto" v-model="element.checked" v-if="panel == 0"
                                                :inline="true"></v-checkbox-btn>
                                <v-icon v-else-if="orderEditMode" class="ml-2 pr-0 mr-2 cursor-move my-auto jannis"
                                        icon="mdi-drag"></v-icon>
                            </v-sheet>

                            <Transition name="slide-fade-down">
                                <v-btn
                                       v-show="element.expanded && orderEditMode"
                                       icon
                                       size="40"
                                       elevation="0"
                                       density="compact"
                                       variant="text"
                                       :ripple="false"
                                       :disabled="isLastElement(element)"
                                       class="mt-auto"
                                       @click="moveElementAndUpdate(element, 'down')">

                                    <v-icon size="40">mdi-chevron-down</v-icon>
                                </v-btn>
                            </Transition>
                        </v-sheet>
                    </v-col>
                    <v-col>
                        <v-row dense>
                            <v-col>
                                <SortableContainer
                                                   @click:expand="showOrderArrows"
                                                   :name="element.name">
                                    <template #title>
                                        <v-row align="center" dense>
                                            <v-col align-self="end" cols="10" sm="6" md="4" lg="3" xl="3">
                                                <div class="mr-2 d-flex align-end">
                                                    <div class="text-caption ">P{{ element.priority }}</div>

                                                    <v-menu v-if="checkWebUILabel(element) > 1">
                                                        <template v-slot:activator="{ props }">
                                                            <v-chip color="wurple"
                                                                    v-bind="props"
                                                                    class="ml-2 px-2"
                                                                    density="comfortable" rounded="20"
                                                                    size="small">WebUI
                                                            </v-chip>
                                                        </template>

                                                        <v-list nav density="compact">
                                                            <v-list-item :href="resolveWebUILink(container)"
                                                                         target="_blank"
                                                                         v-for="container in element.containers.filter((container) => webUILabel in container.labels)"
                                                                         :key="container.id">

                                                                <div class="d-flex flex-column align-center">
                                                                    <v-list-item-title class="pb-1">
                                                                        {{ container.name }}
                                                                    </v-list-item-title>

                                                                    <v-img
                                                                           v-if="container.labels['net.unraid.docker.icon']"
                                                                           height="25" width="25"
                                                                           :src="container.labels['net.unraid.docker.icon']"></v-img>
                                                                    <v-icon v-else size="31">mdi-docker</v-icon>
                                                                </div>
                                                            </v-list-item>
                                                        </v-list>
                                                    </v-menu>
                                                    <v-chip v-else-if="checkWebUILabel(element) == 1"
                                                            color="wurple"
                                                            v-bind="props"
                                                            @click.stop=""
                                                            target="_blank"
                                                            :href="resolveWebUILink(element.containers.filter((container) => webUILabel in container.labels)[0])"
                                                            class="ml-2 px-2"
                                                            density="comfortable" rounded="20"
                                                            size="small">WebUI
                                                    </v-chip>
                                                </div>

                                                <a :href="getPortainerUrl(element, portainerUrlTemplate)" target="_blank"
                                                   class="text-h6"
                                                   @click.stop=""
                                                   style="text-decoration: none; color: inherit;">
                                                    {{ element.name }}
                                                </a>

                                            </v-col>

                                            <v-col v-if="!element.expanded" :class="[xs ? 'pb-3' : undefined]" cols="6"
                                                   sm="4" md="6" lg="7" xl="7">
                                                <div class="d-flex ml-1">
                                                    <div v-for="container in element.containers"
                                                         v-if="element.containers.length > 0">
                                                        <v-icon class="pr-4"
                                                                :color="getContainerStatusCircleColor(container.status)"
                                                                size="14">mdi-circle</v-icon>
                                                    </div>
                                                    <div v-else>
                                                        <v-icon class="pr-4" color="grey"
                                                                size="14">mdi-help-circle</v-icon>
                                                    </div>
                                                </div>
                                            </v-col>

                                            <v-col :class="[xs ? 'pb-3' : undefined]"
                                                   :cols="element.expanded ? '2' : '6'"
                                                   :sm="element.expanded ? '6' : '2'"
                                                   :md="element.expanded ? '8' : '2'"
                                                   :lg="element.expanded ? '9' : '2'"
                                                   :xl="element.expanded ? '9' : '2'">
                                                <v-row :class="[xs ? 'pr-4' : 'pr-10']" justify="end">
                                                    <v-switch :disabled="checkWashboard(element.containers)"
                                                              :loading="loaderState[element.id]"
                                                              @click.stop="toggleAutoStart(element)" color="blue-darken-1"
                                                              density="compact" hide-details
                                                              v-model="element.autoStart" inset></v-switch>
                                                </v-row>
                                            </v-col>
                                        </v-row>
                                    </template>
                                    <template #shortcuts>
                                    </template>
                                    <template #content>
                                        <!-- start stop and restart buttons -->
                                        <div class="d-flex flex-wrap ga-3 mt-3">
                                            <v-btn
                                                   v-if="element.containers.length === 0"
                                                   :loading="loaderState[element.id]"
                                                   variant="tonal"
                                                   prepend-icon="mdi-arrow-right-drop-circle-outline"
                                                   @click="manageStack(element, Action.Start)">Start
                                            </v-btn>
                                            <v-btn
                                                   v-else
                                                   :loading="loaderState[element.id]"
                                                   color="stop"
                                                   variant="tonal"
                                                   prepend-icon="mdi-stop-circle-outline"
                                                   @click="manageStack(element, Action.Stop)">Stop
                                            </v-btn>
                                            <v-btn
                                                   :loading="loaderState[element.id]"
                                                   variant="tonal"
                                                   color="green-lighten-1"
                                                   prepend-icon="mdi-restart"
                                                   @click="manageStack(element, Action.Restart)">Restart
                                            </v-btn>
                                        </div>
                                        <StackContent :containers="element.containers">
                                        </StackContent>
                                    </template>
                                </SortableContainer>
                            </v-col>
                        </v-row>
                    </v-col>
                </v-row>
            </v-card>
        </v-fade-transition>
    </div>
</template>

<script lang="ts" setup>
import { startStack, stopStack, handleStackStateChange, webUILabel, webUIAddressKey } from "@/api/lib";
import { Stack, StackInternal, Action, PaddingClass, Container, StackSettingsDto } from "@/types/types";
import { ref, Ref, watch } from "vue";
import { moveArrayElement, useSortable } from "@vueuse/integrations/useSortable";
import { useSnackbarStore } from "@/store/snackbar";
import { useLocalStore } from "@/store/local";
import { getPortainerUrl, getContainerStatusCircleColor, getFirstContainerIcon } from "@/api/lib";
import { useDisplay } from "vuetify";
import axios from "axios";
import { storeToRefs } from "pinia";
import debounce from 'lodash.debounce';

const { xs } = useDisplay();

const snackbarsStore = useSnackbarStore();
const localStore = useLocalStore();
const { urlConfig } = storeToRefs(localStore);

const search: Ref<string> = ref("");
const searchActive: Ref<boolean> = ref(false)
const changeOrderTransition: Ref<boolean> = ref(false);
const orderEditMode: Ref<boolean> = ref(false);
const sortableRoot = ref<HTMLElement | null>(null);
const panel: Ref<Number | undefined> = ref(undefined);
const props = defineProps<{ items: Stack[], portainerUrlTemplate: string, loading: boolean }>();
const stacksInternal: Ref<StackInternal[]> = ref([]);
const stacksUnfiltered: Ref<StackInternal[]> = ref([]);
const paddingClass: Ref<PaddingClass> = ref({});
const stackOperationInProgress: Ref<boolean> = ref(false);
let loaderState: Record<string, boolean> = reactive({});

watch(props, () => {
    stacksInternal.value = props.items.map((stack) => {
        const stackInternal: StackInternal = {
            ...stack,
            expanded: false,
            checked: false
        };
        paddingClass.value[stack.name] = "row-sst";
        return stackInternal;
    })
    stacksUnfiltered.value = stacksInternal.value;
});

watch(xs, (value) => {
    for (const stack of stacksInternal.value) {
        if (stack.expanded && value) {
            if (orderEditMode.value) {
                paddingClass.value[stack.name] = "row-sst-pad";
            } else {
                paddingClass.value[stack.name] = "row-sst-pad-noorder";
            }
        } else {
            paddingClass.value[stack.name] = "row-sst";
        }
    }
});

watch(search, (newVal, oldVal) => {
    if (newVal == "" && !searchActive.value) {
        searchActive.value = true;
    }
});

watch(search, debounce(() => {
    changeOrderTransition.value = true;
    filterStacks();
    searchActive.value = false;
    setTimeout(() => {
        changeOrderTransition.value = false;
    }, 150);
}, 50))

useSortable(sortableRoot, stacksInternal, {
    handle: ".jannis",
    animation: 150,
    forceFallback: true,
    onUpdate: (e: any) => {
        setTimeout(() => {
            if (e.oldIndex < e.newIndex) {
                console.log(`direction: down, positions: ${e.newIndex - e.oldIndex}`)
                moveElement(stacksInternal.value[e.oldIndex], "down", e.newIndex - e.oldIndex);
            } else {
                console.log(`direction: up, positions: ${e.oldIndex - e.newIndex}`)
                moveElement(stacksInternal.value[e.oldIndex], "up", e.oldIndex - e.newIndex);
            }

            nextTick(async () => {
                for (let stack of stacksInternal.value) {
                    stack.priority = stacksInternal.value.indexOf(stack);
                }
                try {
                    const newStack: StackInternal = stacksInternal.value[e.newIndex];
                    const dto: StackSettingsDto = {
                        stackName: newStack.name,
                        priority: newStack.priority,
                        autoStart: newStack.autoStart,
                        stackId: newStack.id
                    };
                    const response = await axios.put(`/api/db/stacks/${newStack.name}`, dto, { params: { updatePrio: true } });
                    snackbarsStore.addSnackbar("stacks_order", "Stack order updated", "success");
                } catch (error: any) {
                    snackbarsStore.addSnackbar("stacks_order", `Failed to update stack order: ${error.message}`, "error");
                }
            });

        }, 100);
    }
});

function toggleEditMode() {
    orderEditMode.value = !orderEditMode.value;
    search.value = "";
}

async function ModifySelectedContainerStatus(action: Action) {
    stackOperationInProgress.value = true;
    for (const stack of stacksInternal.value) {
        if (stack.checked) {
            for (const container of stack.containers) {
                if (container.image.includes("washboard")) {
                    continue;
                }
            }
            await manageStack(stack, action);
        }
    }
    stackOperationInProgress.value = false;
}

function checkWashboard(containers: Container[]) {
    return containers.some((container) => container.image.includes("washboard"));
}

function s(arr: any) {
    return arr.length > 1 ? "s" : "";
}

async function toggleAutoStart(stack: Stack) {
    loaderState[stack.id] = true;
    if (!stack) {
        return;
    }

    stack.autoStart = !stack.autoStart;
    let dto: StackSettingsDto = {
        stackName: stack.name,
        priority: stack.priority,
        autoStart: stack.autoStart,
        stackId: stack.id
    };
    try {
        const result = await axios.put(`/api/db/stacks/${stack.name}`, dto);
        snackbarsStore.addSnackbar(`${stack.name}_autoStart`, `Auto start for ${stack.name} ${stack.autoStart ? "enabled" : "disabled"}`, "success");
    } catch (error: any) {
        stack.autoStart = !stack.autoStart;
        snackbarsStore.addSnackbar(`${stack.name}_autoStart`, `Failed to update auto start for ${stack.name}: ${error.message}`, "error");
    }
    loaderState[stack.id] = false;
}

function checkWebUILabel(stack: Stack) {
    const webUIContainers: Container[] | undefined = stack.containers.filter((container) => webUILabel in container.labels);
    return webUIContainers.length;
}

function resolveWebUILink(container: Container) {
    const ports = container.ports.map((port) => port.split(":")).map((p) => [p[1], p[0]]);
    const portDict = Object.fromEntries(ports);

    let target = urlConfig.value.defaultHost;
    if (!container.labels[webUILabel].includes(webUIAddressKey)) {
        //console.log(`container: ${container.name} has no address key`);
        return container.labels[webUILabel];
    }

    let outputUrl = container.labels[webUILabel].replace(webUIAddressKey, target);

    // iterate to portdict and replace all private with public ports
    for (const [privatePort, publicPort] of Object.entries(portDict) as [string, string][]) {
        outputUrl = outputUrl.replace(privatePort, publicPort);
    }

    //console.log(`container: ${container.name} output url: ${outputUrl}`);

    if (!target.startsWith("http://") && !target.startsWith("https://")) {
        return "//" + outputUrl;
    }
    return outputUrl;
}

async function manageStack(stack: Stack, action: Action) {
    loaderState[stack.id] = true;
    if (![Action.Start, Action.Stop, Action.Restart].includes(action)) {
        throw new Error(
            `Action should be "${Action.Start}", "${Action.Stop}", or "${Action.Restart}", got "${action}"`
        );
    }

    try {
        if (action === Action.Restart) {
            await stopStack(stack.id);
            const startResponse = await startStack(stack.id);
            await handleStackStateChange(stack, Action.Restart, startResponse, snackbarsStore, stacksInternal, false);
        } else {
            const response = await (action === Action.Start ? startStack(stack.id) : stopStack(stack.id));
            await handleStackStateChange(stack, action, response, snackbarsStore, stacksInternal, false);
        }
    } catch (error: any) {
        snackbarsStore.addSnackbar(
            `${stack.id}_startstop`,
            `Failed to ${action} ${stack?.name}: ${error.message}`,
            "error"
        );
    } finally {
        loaderState[stack.id] = false;
    }
}

function showOrderArrows(expand: any) {
    const elem: StackInternal | undefined = stacksInternal.value.find((stack) => stack.name === expand.name);
    if (!elem) {
        return;
    }

    if (elem.name === expand.name) {
        elem.expanded = !elem.expanded;
        if (xs.value && elem.expanded) {
            if (orderEditMode.value) {
                paddingClass.value[elem.name] = "row-sst-pad";
            } else {
                paddingClass.value[elem.name] = "row-sst-pad-noorder";
            }
        } else {
            paddingClass.value[elem.name] = "row-sst";
        }
    }
}

function isFirstElement(element: Stack) {
    return stacksInternal.value[0].id === element.id;
}

function isLastElement(element: Stack) {
    return stacksInternal.value[stacksInternal.value.length - 1].id === element.id;
}

function moveElement(element: Stack, action: string, amount: number = 1) {
    let toIndex;
    const currIdx = stacksInternal.value.findIndex((i) => i.id == element.id);
    if (action === "up") {
        if (currIdx === 0) return;
        toIndex = currIdx - amount;
    } else if (action === "down") {
        if (currIdx === stacksInternal.value.length - 1) return;
        toIndex = currIdx + amount;
    } else {
        toIndex = currIdx;
    }
    const [removed] = stacksInternal.value.splice(currIdx, 1);
    removed.priority = toIndex;
    stacksInternal.value.splice(toIndex, 0, removed);
    return removed;
}

function moveElementAndUpdate(element: Stack, action: string) {
    const stack = moveElement(element, action, 1);
    if (stack) {
        updatePriority(stack);
    }
}

async function updatePriority(stack: StackInternal) {
    try {
        const newStack: StackInternal = stack;
        const dto: StackSettingsDto = {
            stackName: newStack.name,
            priority: newStack.priority,
            autoStart: newStack.autoStart,
            stackId: newStack.id
        };
        const response = await axios.put(`/api/db/stacks/${newStack.name}`, dto, { params: { updatePrio: true } });
    } catch (error: any) {
        snackbarsStore.addSnackbar("stacks_order", `Failed to update stack order: ${error.message}`, "error");
    }
}

function handleSelectAutostartContainers() {
    stacksInternal.value.forEach((stack) => {
        if (stack.autoStart) {
            stack.checked = true;
        }
    });
}

function clearSelection() {
    stacksInternal.value.forEach((stack) => {
        stack.checked = false;
    });
}

function filterStacks() {
    stacksInternal.value = stacksUnfiltered.value.filter((stack) => {
        return stack.name.toLowerCase().includes(search.value.toLowerCase());
    });
}

</script>

<style lang="scss" scoped>
.draggable {
    cursor: move;
}

/*
  Enter and leave animations can use different
  durations and timing functions.
*/
.slide-fade-up-enter-active {
    transition: all 0.3s ease-out;
}

.slide-fade-up-leave-active {
    transition: all 0.1s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-up-enter-from,
.slide-fade-up-leave-to {
    transform: translateY(15px);
    opacity: 0;
}

.slide-fade-down-enter-active {
    transition: all 0.2ds ease-out;
}

.slide-fade-down-leave-active {
    transition: all 0.1s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-down-enter-from,
.slide-fade-down-leave-to {
    transform: translateY(-15px);
    opacity: 0;
}

.row-sst-pad {
    flex-wrap: nowrap;
    padding-right: 110px;
}

.row-sst-pad-noorder {
    flex-wrap: nowrap;
    padding-right: 72px;
}

.row-sst {
    flex-wrap: nowrap;
}

.stack-icon {
    min-width: 20px;
    max-width: 20px;
}
</style>
