<template>
  <v-app>
    <v-main>
      <router-view />
      <div>
        <v-snackbar v-for="(snackbar, index) in snackbars" :key="index" :model-value="snackbar.show"
                    :color="snackbar.color"
                    :timeout=-1
                    @update:model-value="value => snackbarStore.removeSnackbar(snackbar.id)"
                    :style="{ 'margin-bottom': snackbarStore.calcSnackbarMargin(index) }"
                    location="bottom right"
                    close-on-content-click
                    multiline>
          {{ snackbar.message }}
        </v-snackbar>
      </div>
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { useSnackbarStore } from '@/store/snackbar';
import { useUpdateQuelelelStore } from './store/updateQuelelel';
import { storeToRefs } from 'pinia';
import { useTheme } from 'vuetify'
import { onMounted } from 'vue';
import axios from 'axios';
import { UpdateQueue, QueueStatus, QueueItem } from './types/types';

const snackbarStore = useSnackbarStore();
const { snackbars: snackbars } = storeToRefs(snackbarStore);


const updateQuelelelStore = useUpdateQuelelelStore();
const { queue: stackQueue } = storeToRefs(updateQuelelelStore);
const snackbarsStore = useSnackbarStore();

const theme = useTheme();

let mediaEvent;

onMounted(() => {
  mediaEvent = window.matchMedia("(prefers-color-scheme: dark)");
  mediaEvent.addEventListener("change", handleSystemThemeUpdate);
  if (mediaEvent.matches) {
    theme.global.name.value = "dark";
  } else {
    theme.global.name.value = "light";
  }

  let wsAddr = `${axios.defaults.baseURL}/ws/stacks-update`.replace('http://', 'ws://').replace('https://', 'wss://');
  let socket = new WebSocket(wsAddr);
  socket.onmessage = function (event) {
    let data: UpdateQueue = JSON.parse(event.data);


    for (let [newStatus, newItems] of Object.entries(data) as [QueueStatus, Record<string, QueueItem>][]) {
      for (let stackName in newItems) {
        const queueItem = newItems[stackName];

        let previousBucket: string | undefined = undefined;
        for (let [oldStatus, oldItems] of Object.entries(stackQueue.value)) {
          if (queueItem.stackName in oldItems) {
            previousBucket = oldStatus;
            break;
          }
        }

        switch (newStatus) {
          case "queued":
            break;
          case "done":
            if (previousBucket && previousBucket != newStatus) {
              snackbarsStore.addSnackbar(queueItem.stackId, `Stack ${queueItem.stackName} updated successfully`, "success");
            }
            break;
          case "error":
            if (previousBucket && previousBucket != newStatus) {
              snackbarsStore.addSnackbar(queueItem.stackId, `Stack ${queueItem.stackName} update failed`, "error");
            }
            break;
        }
      }
    }
    updateQuelelelStore.update(data);
  };
})

function handleSystemThemeUpdate(e: any) {
  console.log(`updating theme based on system preference ${e.matches ? "dark" : "light"}`);
  if (e.matches) {
    theme.global.name.value = "dark";
  } else {
    theme.global.name.value = "light";
  }
}

</script>

<style lang="scss">
</style>
