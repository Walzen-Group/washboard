<template>
  <v-app-bar elevation="1" density="compact" color="washboard-appbar">
    <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
    <v-btn icon @click.stop="miniVariant = !miniVariant">
      <v-icon>mdi-{{ `chevron-${miniVariant ? "right" : "left"}` }}</v-icon>
    </v-btn>
    <v-btn icon @click.stop="clipped = !clipped">
      <v-icon>mdi-washing-machine</v-icon>
    </v-btn>

    <v-toolbar-title class="ml-4 mr-4" v-text="title" />
    <v-tooltip v-if="smAndUp"
               location="bottom">
      <template v-slot:activator="{ props }">
        <span rounded="xl" elevation="0" class="ma-2 text-none" :ripple="false" v-bind="props">
          Queued Stacks: {{ queuedStacks.length }}
        </span>
      </template>
      <span v-if="queuedStacks.length > 0">{{ queuedStacks.join(", ") }}</span>
      <span v-else>Empty</span>
    </v-tooltip>
    <v-spacer />


    <v-btn icon @click.stop="toggleTheme">
      <v-icon>mdi-theme-light-dark</v-icon>
    </v-btn>
  </v-app-bar>
  <v-navigation-drawer width="230" floating mobile-breakpoint="sm" v-model="drawer"
                       :rail="miniVariant">
    <v-list nav density="compact">
      <v-list-item v-for="(item, i) in items" :key="i" :to="item.to" :prepend-icon="item.icon"
                   :title="item.title" router
                   exact>
      </v-list-item>
    </v-list>

    <template v-slot:append>
      <v-divider v-if="!smAndUp" class="mb-1"></v-divider>
      <div class="d-flex align-center justify-center">
        <v-tooltip v-if="!smAndUp"
                   location="bottom">
          <template v-slot:activator="{ props }">
            <span elevation="0" class="ma-2 text-none" :ripple="false" v-bind="props">
              Queued Stacks: {{ queuedStacks.length }}
            </span>
          </template>
          <span v-if="queuedStacks.length > 0">{{ queuedStacks.join(", ") }}</span>
          <span v-else>Empty</span>
        </v-tooltip>
      </div>

    </template>
  </v-navigation-drawer>
</template>

<script lang="ts" setup>
import { useDisplay } from 'vuetify'
const { smAndUp } = useDisplay()
import { ref } from 'vue';
import { useTheme } from 'vuetify'
import { useUpdateQuelelelStore } from '@/store/updateQuelelel';
import { storeToRefs } from 'pinia';
const updateQuelelelStore = useUpdateQuelelelStore();
const { queue: updateQueue} = storeToRefs(updateQuelelelStore);

const title = "Washboard"

const items: any[] = [
  {
    icon: "mdi-home-account",
    title: "Home",
    to: "/",
  },
  {
    icon: "mdi-docker",
    title: "Update Stacks",
    to: "/docker-update-manager",
  },
  {
    icon: "mdi-toolbox",
    title: "Manage Stacks",
    to: "/docker-manager",
  }
];
const clipped = ref(false);
const drawer = ref(true);
const miniVariant = ref(false);
const queuedStacks = ref([]);

const theme = useTheme();

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}


</script>
