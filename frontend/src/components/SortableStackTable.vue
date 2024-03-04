<template>
    <div>{{ stacksInternal.map((stack) => stack.name).join(', ') }}</div>
    <div ref="sortableRoot">
        <v-card v-for="element in stacksInternal" class="pb-1 pt-2 mb-2">
            <v-row dense>
                <v-col cols="auto" class="ml-2">
                    <v-sheet
                             class="fill-height d-flex flex-column">
                        <v-btn icon size="40" elevation="0" variant="text" density="compact"
                               :ripple="false"
                               @click="moveElement(element, stacksInternal.findIndex((i) => i.name == element.name) - 1)">
                            <v-icon size="40"> mdi-menu-up-outline </v-icon>
                        </v-btn>
                        <v-sheet class="d-flex fill-height flex-column">
                            <v-icon class="ml-2 pr-0 cursor-move my-auto jannis"
                                    icon="mdi-drag"></v-icon>
                        </v-sheet>
                        <v-btn icon size="40" elevation="0" density="compact" variant="text"
                               :ripple="false"
                               class="mt-auto"
                               @click="moveElement(element, stacksInternal.findIndex((i) => i.name == element.name) + 1)">
                            <v-icon size="40">mdi-menu-down-outline</v-icon>
                        </v-btn>
                    </v-sheet>
                </v-col>
                <v-col class="pl-0">
                    <v-card-title>{{ element.name }}</v-card-title>
                    <v-card-text class="pl-0">
                        <v-container>
                            <v-row>
                                <v-col v-for="container in element.containers"
                                       :key="container.id">
                                    <v-card>
                                        <v-card-title>{{ container.name
                                            }}</v-card-title>
                                        <v-card-text>{{ container.ports.join(", ")
                                            }}</v-card-text>
                                    </v-card>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                </v-col>
            </v-row>
        </v-card>
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
    scroll: true,
    forceFallback: true,
    bubbleScroll: true,
});

function moveElement(element: Stack, toIndex: number) {
    const fromIndex = stacksInternal.value.findIndex((i) => i.id == element.id);
    stacksInternal.value.splice(fromIndex, 1);
    stacksInternal.value.splice(toIndex, 0, element);
}

function getRandomInt(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function generateTestData(): Stack[] {
    const testData: Stack[] = [];

    for (let i = 1; i <= 15; i++) {
        const containers: Container[] = [];
        const numberOfContainers = getRandomInt(1, 4); // Random number of containers between 1 and 4

        for (let j = 1; j <= numberOfContainers; j++) {
            const container: Container = {
                id: `container-${i}-${j}`,
                name: `Container ${j}`,
                image: `image${j}:latest`,
                upToDate: `2024-02-${j < 10 ? '0' + j : j}`,
                upToDateIgnored: j % 2 === 0,
                status: j % 2 === 0 ? 'Running' : 'Stopped',
                ports: [8080 + j, 9090 + j],
                labels: {
                    'owner': `owner${i}`,
                    'project': `project${i}`
                }
            };
            containers.push(container);
        }

        const stack: Stack = {
            id: i,
            name: `Stack ${i}`,
            containers: containers,
            updateStatus: [{
                'lastChecked': `2024-02-${i < 10 ? '0' + i : i}`,
                'isUpToDate': i % 2 === 0
            }]
        };

        testData.push(stack);
    }

    return testData;
}

</script>

<style lang="scss">
.draggable {
    cursor: move;
}
</style>
