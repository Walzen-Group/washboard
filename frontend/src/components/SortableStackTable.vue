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
                                       :loading="false"
                                       variant="tonal"
                                       prepend-icon="mdi-arrow-right-drop-circle-outline"
                                       @click="">Start
                                </v-btn>
                                <v-btn
                                       :loading="false"
                                       color="stop"
                                       variant="tonal"
                                       prepend-icon="mdi-stop-circle-outline"
                                       @click="">Stop
                                </v-btn>
                                <v-btn
                                       :loading="false"
                                       variant="tonal"
                                       color="green-lighten-1"
                                       prepend-icon="mdi-restart"
                                       @click="">Restart
                                </v-btn>
                            </div>

                        </v-card-text>
                        <v-card-title>
                            Select Options
                        </v-card-title>

                        <v-card-text>
                            <div>
                                <v-btn
                                       :loading="false"
                                       variant="tonal"
                                       prepend-icon="mdi-auto-mode"
                                       @click="">Autostart
                                </v-btn>
                            </div>
                        </v-card-text>
                    </v-card>
                </div>
            </v-expansion-panel-text>
        </v-expansion-panel>
    </v-expansion-panels>


    <v-divider class="mb-3 mt-4" />

    <div ref="sortableRoot">
        <v-fade-transition group>
            <v-card v-for="element in stacksInternal" class="pb-2 pt-3 mb-2" :key="element.id">
                <v-row dense :class="[paddingClass[element.name]]">
                    <v-col cols="auto" class="ml-2">
                        <v-sheet class="fill-height d-flex flex-column">
                            <Transition name="slide-fade-up">
                                <v-btn
                                       v-show="element.expanded"
                                       icon
                                       size="40"
                                       elevation="0"
                                       variant="text"
                                       density="compact"
                                       :ripple="false"
                                       :disabled="isFirstElement(element)"
                                       @click="moveElement(element, 'up')">
                                    <v-icon size="40"> mdi-chevron-up </v-icon>
                                </v-btn>
                            </Transition>

                            <v-sheet class="d-flex fill-height flex-column">
                                <v-checkbox-btn class="my-auto" v-model="element.checked" v-if="panel == 0"
                                                :inline="true"></v-checkbox-btn>
                                <v-icon v-else class="ml-2 pr-0 mr-2 cursor-move my-auto jannis" icon="mdi-drag"></v-icon>
                            </v-sheet>

                            <Transition name="slide-fade-down">
                                <v-btn
                                       v-show="element.expanded"
                                       icon
                                       size="40"
                                       elevation="0"
                                       density="compact"
                                       variant="text"
                                       :ripple="false"
                                       :disabled="isLastElement(element)"
                                       class="mt-auto"
                                       @click="moveElement(element, 'down')">

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
                                            <v-col align-self="end" cols="12" sm="6" md="4" lg="3" xl="3">

                                                <div class="text-caption mr-2">P{{ element.priority }}
                                                    </div>

                                                    <a :href="getPortainerUrl(element, portainerUrlTemplate)" target="_blank"
                                                       class="text-h6"
                                                       style="text-decoration: none; color: inherit;">
                                                        {{ element.name }}
                                                    </a>

                                            </v-col>

                                            <v-col :class="[xs ? 'pb-3' : undefined]" cols="6" sm="4" md="6" lg="7" xl="7">
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

                                            <v-col :class="[xs ? 'pb-3' : undefined]" cols="6" sm="2" md="2" lg="2" xl="2">
                                                <v-row :class="[xs ? 'pr-4' : 'pr-10']" justify="end">
                                                    <v-switch color="blue-darken-1" density="compact" hide-details
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
import { startStack, stopStack, handleResponse } from "@/api/lib";
import { Stack, StackInternal, Action, PaddingClass } from "@/types/types";
import { ref, Ref, watch } from "vue";
import { useSortable } from "@vueuse/integrations/useSortable";
import { useSnackbarStore } from "@/store/snackbar";
import { getPortainerUrl, getContainerStatusCircleColor } from "@/api/lib";
import { useDisplay } from "vuetify";

const { xs } = useDisplay();

const snackbarsStore = useSnackbarStore();
const sortableRoot = ref<HTMLElement | null>(null);
const panel: Ref<Number | undefined> = ref(undefined);
const props = defineProps<{ items: Stack[], portainerUrlTemplate: string }>();
const stacksInternal: Ref<StackInternal[]> = ref([]);
const paddingClass: Ref<PaddingClass> = ref({});
let loaderState: Record<string, boolean> = reactive({});

watch(props, () => {
    stacksInternal.value = props.items.map((stack) => {
        const stackInternal: StackInternal = {
            ...stack,
            expanded: false,
            checked: false
        };
        paddingClass.value[stack.name] = undefined;
        return stackInternal;
    })
});

watch(xs, (value) => {
    for (const stack of stacksInternal.value) {
        if (stack.expanded && value) {
            paddingClass.value[stack.name] = "row-sst-pad";
        } else {
            paddingClass.value[stack.name] = undefined;
        }
    }
});

useSortable(sortableRoot, stacksInternal, {
    handle: ".jannis",
    animation: 250,
    forceFallback: true,
});

function s(arr: any) {
    return arr.length > 1 ? "s" : "";
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
            await handleResponse(stack, Action.Restart, startResponse, snackbarsStore, stacksInternal);
        } else {
            const response = await (action === Action.Start ? startStack(stack.id) : stopStack(stack.id));
            await handleResponse(stack, action, response, snackbarsStore, stacksInternal);
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
    stacksInternal.value = stacksInternal.value.map((stack) => {
        if (stack.name === expand.name) {
            stack.expanded = !stack.expanded;
            if (xs.value && stack.expanded) {
                paddingClass.value[stack.name] = "row-sst-pad";
            } else {
                paddingClass.value[stack.name] = "row-sst";
                setTimeout(() => {
                    paddingClass.value[stack.name] = undefined;
                }, 500);
            }
        }
        return stack;
    });
}

function isFirstElement(element: Stack) {
    return stacksInternal.value[0].id === element.id;
}

function isLastElement(element: Stack) {
    return stacksInternal.value[stacksInternal.value.length - 1].id === element.id;
}

function moveElement(element: Stack, action: string) {
    let toIndex;
    const currIdx = stacksInternal.value.findIndex((i) => i.id == element.id);
    if (action === "up") {
        if (currIdx === 0) return;
        toIndex = currIdx - 1;
    } else if (action === "down") {
        if (currIdx === stacksInternal.value.length - 1) return;
        toIndex = currIdx + 1;
    } else {
        toIndex = currIdx;
    }
    const [removed] = stacksInternal.value.splice(currIdx, 1);
    stacksInternal.value.splice(toIndex, 0, removed);
}
</script>

<style lang="scss">
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

.row-sst {
    flex-wrap: nowrap;
}

.stack-icon {
    min-width: 28px;
    max-width: 28px;
}
</style>
