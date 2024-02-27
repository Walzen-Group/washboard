<!-- eslint-disable vue/valid-v-slot -->
<template>
  <div class="mt-2 mb-4 text-h4">Manage & Configure Stacks</div>
  <div>{{ testDataStacks.map((stack) => stack.name).join(', ') }}</div>
  <v-data-iterator :items="testDataStacks" :items-per-page="itemsPerPage">
    <template v-slot:default="{ items }">
      <v-fade-transition group hide-on-leave>
        <div ref="el">
          <div v-for="item in items" :key="item.raw.id">
            <v-card class="pb-1 pt-2 mb-2">
              <v-row align="center">
                <v-col cols="auto" class="ml-2 pr-0"><v-icon icon="mdi-drag" class="jannis"></v-icon></v-col>
                <v-col class="pl-0">
                  <v-card-title>{{ item.raw.name }}</v-card-title>
                  <v-card-text class="pl-0">
                    <v-container>
                      <v-row>
                        <v-col v-for="container in item.raw.containers" :key="container.id">
                          <v-card>
                            <v-card-title>{{ container.name }}</v-card-title>
                            <v-card-text>{{ container.ports.join(", ") }}</v-card-text>
                          </v-card>
                        </v-col>
                      </v-row>
                    </v-container>
                  </v-card-text>
                </v-col>
              </v-row>
            </v-card>
          </div>
        </div>
      </v-fade-transition>
    </template>
    <template v-slot:header="{ page, pageCount, prevPage, nextPage }">
      <div class="d-flex align-center justify-center px-4 pb-4">
        <v-btn class="mr-1" :disabled="page === 1" icon="mdi-arrow-left" density="comfortable" variant="tonal"
          @click="prevPage"></v-btn>

        <div class="mx-2 text-caption">
          Page {{ page }} of {{ pageCount }}
        </div>

        <v-btn class="ml-1" :disabled="page >= pageCount" icon="mdi-arrow-right" density="comfortable" variant="tonal"
          @click="nextPage"></v-btn>
      </div>
    </template>
  </v-data-iterator>
</template>

<script lang="ts" setup>

import { useSortable } from '@vueuse/integrations/useSortable';
import { Stack, Container } from '@/types/types';
import { ref, Ref } from 'vue';

const testDataStacks: Ref<Stack[]> = ref(generateTestData());
const el = ref<HTMLElement | null>(null)
const itemsPerPage = ref(10)
const { option } = useSortable(el, testDataStacks, {
  handle: '.jannis',
  animation: 150,
});

function setAnimation() {
  option("animation", 300);
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
