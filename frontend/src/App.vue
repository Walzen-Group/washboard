<template>
  <v-app>
    <v-main>
      <router-view />
      <div>
          <v-snackbar v-for="(snackbar, index) in snackbars" :key="index" :model-value="snackbar.show"
                      :color="snackbar.color"
                      :timeout=-1 @update:model-value="value => snackbarStore.removeSnackbar(snackbar.id)"
                      :style="{ 'margin-bottom': snackbarStore.calcSnackbarMargin(index) }" location="bottom right"
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
import { storeToRefs } from 'pinia';

const snackbarStore = useSnackbarStore();
const { snackbars: snackbars } = storeToRefs(snackbarStore);
for (let snackbar of snackbars.value) {
    console.log(snackbar);
}
</script>
