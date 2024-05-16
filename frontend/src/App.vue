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
import { useUpdateQuelelelStore } from '@/store/updateQuelelel';
import { storeToRefs } from 'pinia';
import { useTheme } from 'vuetify'
import { onMounted } from 'vue';
import { callRefreshTokenRoute, connectWebSocket } from '@/api/lib';
import { useLocalStore } from "@/store/local";

const snackbarStore = useSnackbarStore();
const { snackbars: snackbars } = storeToRefs(snackbarStore);


const theme = useTheme();

let mediaEvent;

onMounted(async () => {
  mediaEvent = window.matchMedia("(prefers-color-scheme: dark)");
  mediaEvent.addEventListener("change", handleSystemThemeUpdate);
  if (mediaEvent.matches) {
    theme.global.name.value = "dark";
  } else {
    theme.global.name.value = "light";
  }

  callRefreshTokenRoute().then((result) => {
     if (result) {
       connectWebSocket();
     }
  });
});



function handleSystemThemeUpdate(e: any) {
  console.log(`updating theme based on system preference ${e.matches ? "dark" : "light"}`);
  if (e.matches) {
    theme.global.name.value = "dark";
  } else {
    theme.global.name.value = "light";
  }
}

</script>

<style lang="scss"></style>
