<!-- eslint-disable vue/valid-v-slot -->
<template>
  <div class="mt-2 mb-4 text-h4">Manage & Configure Stacks</div>
  <SortableStackTable :items="testDataStacks"></SortableStackTable>

</template>

<script lang="ts" setup>
import SortableStackTable from '@/components/SortableStackTable.vue';
import { Stack, Container } from '@/types/types';
import { ref, Ref, onMounted } from 'vue';

const testDataStacks: Ref<Stack[]> = ref([]);

function getRandomInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

async function generateTestData(): Promise<Stack[]> {
  await new Promise((resolve) => setTimeout(resolve, 1000));
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

onMounted(async () => {
  testDataStacks.value = await generateTestData();
});

</script>
