<template>
    <div>{{ stacksInternal.map((stack) => stack.name).join(', ') }}</div>
    <div ref="sortableRoot">
        <v-fade-transition group>
            <v-card v-for="element in stacksInternal" class="pb-1 pt-2 mb-2" :key="element.id">
                <v-row dense>
                    <v-col cols="auto" class="ml-2">
                        <v-sheet
                                 class="fill-height d-flex flex-column">
                            <v-btn icon size="40" elevation="0" variant="text" density="compact"
                                   :ripple="false"
                                   :disabled="isFirstElement(element)"
                                   @click="moveElement(element, 'up')">
                                <v-icon size="40"> mdi-chevron-up </v-icon>
                            </v-btn>
                            <v-sheet class="d-flex fill-height flex-column">
                                <v-icon class="ml-2 pr-0 cursor-move my-auto jannis"
                                        icon="mdi-drag"></v-icon>
                            </v-sheet>
                            <v-btn icon size="40" elevation="0" density="compact" variant="text"
                                   :ripple="false"
                                   :disabled="isLastElement(element)"
                                   class="mt-auto"
                                   @click="moveElement(element, 'down')">
                                <v-icon size="40">mdi-chevron-down</v-icon>
                            </v-btn>
                        </v-sheet>
                    </v-col>
                    <v-col>
                        <StackContainer>
                            <template #title>
                                <div class="d-flex">
                                    <div class="text-h6">
                                        Manuel
                                    </div>
                                    <div>
                                        Banuel
                                    </div>
                                </div>
                            </template>
                            <template #shortcuts>Bannis</template>
                        </StackContainer>
                    </v-col>
                </v-row>
            </v-card>
        </v-fade-transition>
    </div>
</template>

<script lang="ts" setup>
import { Stack, Container } from '@/types/types';
import { ref, Ref, watch } from 'vue';
import { useSortable } from '@vueuse/integrations/useSortable';

const sortableRoot = ref<HTMLElement | null>(null);
const props = defineProps<{ items: Stack[] }>();
const stacksInternal: Ref<Stack[]> = ref(props.items);


watch(props, () => {
    stacksInternal.value = props.items;
});

useSortable(sortableRoot, stacksInternal, {
    handle: '.jannis',
    animation: 250,
    forceFallback: true,
});

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
</style>
