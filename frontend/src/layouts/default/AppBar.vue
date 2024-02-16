<template>
  <v-navigation-drawer v-model="drawer" :rail="miniVariant">
    <v-list>
      <v-list-item v-for="(item, i) in items" :key="i" :to="item.to" :prepend-icon="item.icon" :title="item.title" router
        exact>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
  <v-app-bar :absolute="clipped">
    <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
    <v-btn icon @click.stop="miniVariant = !miniVariant">
      <v-icon>mdi-{{ `chevron-${miniVariant ? "right" : "left"}` }}</v-icon>
    </v-btn>
    <v-btn icon @click.stop="clipped = !clipped">
      <v-icon>mdi-washing-machine</v-icon>
    </v-btn>

    <v-toolbar-title class="ml-4" v-text="title" />
    <v-spacer />
    <v-btn icon @click.stop="toggleTheme">
      <v-icon>mdi-theme-light-dark</v-icon>
    </v-btn>
  </v-app-bar>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useTheme } from 'vuetify'
const title = "Washboard"
let mediaEvent;
const items: any[] = [
  {
    icon: "mdi-home-account",
    title: "Home",
    to: "/",
  },
  {
    icon: "mdi-docker",
    title: "Stack Updater",
    to: "/docker-update-manager",
  }
];
const clipped = ref(false);
const drawer = ref(false);
const miniVariant = ref(false);

const theme = useTheme();

onMounted(() => {
  mediaEvent = window.matchMedia("(prefers-color-scheme: dark)");
  mediaEvent.addEventListener("change", handleSystemThemeUpdate);
  if (mediaEvent.matches) {
    theme.global.name.value = "dark";
  } else {
    theme.global.name.value = "light";
  }
  console.log(mediaEvent);
})

function handleSystemThemeUpdate(e: any) {
      console.log(`updating theme based on system preference ${e.matches ? "dark" : "light"}`);
      if (e.matches) {
        theme.global.name.value = "dark";
      } else {
        theme.global.name.value = "light";
      }
}

function toggleTheme () {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

</script>
