<template>
    <div>{{ stacksInternal.map((stack) => stack.name).join(", ") }}</div>

    <div class="mb-3">
        <v-card class="pt-2" variant="flat">
            <v-card-title>Manage Stack State</v-card-title>
            <v-card-text>
                <div class="d-flex flex-wrap ga-3">
                    <v-btn
                           :loading="false"
                           variant="tonal"
                           prepend-icon="mdi-arrow-right-drop-circle-outline"
                           @click="">Stacks
                    </v-btn>
                    <v-btn
                           :loading="false"
                           color="stop"
                           variant="tonal"
                           prepend-icon="mdi-stop-circle-outline"
                           @click="">Stacks
                    </v-btn>

                </div>

            </v-card-text>

            <v-card-title>
                Manage Container State
            </v-card-title>

            <v-card-text>
                <div class="d-flex flex-wrap ga-3">
                    <v-btn
                           :loading="false"
                           variant="tonal"
                           prepend-icon="mdi-arrow-right-drop-circle-outline"
                           @click="">Containers
                    </v-btn>
                    <v-btn
                           :loading="false"
                           color="stop"
                           variant="tonal"
                           prepend-icon="mdi-stop-circle-outline"
                           @click="">Containers
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

    <v-divider class="mb-3" />

    <div ref="sortableRoot">
        <v-fade-transition group>
            <v-card v-for="element in stacksInternal" class="pb-2 pt-2 mb-2" :key="element.id">
                <v-row dense>
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
                                <v-icon class="ml-2 mr-2 pr-0 cursor-move my-auto jannis" icon="mdi-drag"></v-icon>
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
                        <SortableContainer
                                           @click:expand="showOrderArrows"
                                           :name="element.name">
                            <template #title>
                                <div class="d-flex align-center">
                                    <v-checkbox-btn :inline="true"></v-checkbox-btn>
                                    <div class="d-flex flex-wrap ga-1">
                                        <div class="text-h6">{{ element.name }} {{ element.id }}</div>
                                        <div class="">Priority: {{ element.priority }}</div>
                                    </div>
                                </div>
                            </template>
                            <template #shortcuts>

                            </template>
                            <template #content>
                                <StackContent :containers="element.containers">
                                </StackContent>
                            </template>
                        </SortableContainer>
                    </v-col>
                </v-row>
            </v-card>
        </v-fade-transition>
    </div>
</template>

<script lang="ts" setup>
import { Stack, StackInternal } from "@/types/types";
import { ref, Ref, watch } from "vue";
import { useSortable } from "@vueuse/integrations/useSortable";

const sortableRoot = ref<HTMLElement | null>(null);
const props = defineProps<{ items: Stack[] }>();
const stacksInternal: Ref<StackInternal[]> = ref([]);
let loaderState: Record<string, boolean> = reactive({});

watch(props, () => {
    stacksInternal.value = props.items.map((stack) => {
        const stackInternal: StackInternal = {
            ...stack,
            expanded: false,
        };
        return stackInternal;
    })
});

useSortable(sortableRoot, stacksInternal, {
    handle: ".jannis",
    animation: 250,
    forceFallback: true,
});

function showOrderArrows(expand: any) {
    stacksInternal.value = stacksInternal.value.map((stack) => {
        if (stack.name === expand.name) {
            stack.expanded = !stack.expanded;
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
</style>
