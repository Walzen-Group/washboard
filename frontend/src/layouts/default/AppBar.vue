<template>
  <v-app-bar elevation="1" density="compact" color="washboard-appbar">
    <v-app-bar-nav-icon @click.stop="switchDrawer" />
    <v-btn icon @click.stop="switchMini">
      <v-icon>mdi-{{ `chevron-${miniVariant ? "right" : "left"}` }}</v-icon>
    </v-btn>
    <v-btn icon @click.stop="switchClipped">
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
  <v-navigation-drawer width="230" floating :temporary="!clipped" mobile-breakpoint="md"
                       v-model="drawer"
                       :rail="miniVariant">
    <v-list nav density="compact">
      <v-list-item v-for="(item, i) in items" :key="i" :to="item.to" :prepend-icon="item.icon"
                   :title="item.title" router
                   exact>
      </v-list-item>
    </v-list>

    <template v-slot:append>
      <v-divider v-if="!smAndUp && !miniVariant" class="mb-1"></v-divider>
      <v-fade-transition hide-on-leave>
        <div class="d-flex align-center justify-center" v-if="!miniVariant">
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
      </v-fade-transition>


    </template>
  </v-navigation-drawer>
</template>

<script lang="ts" setup>
import { useDisplay } from 'vuetify'
const { smAndUp, mdAndUp } = useDisplay()
import { ref } from 'vue';
import { useTheme } from 'vuetify'
import { useUpdateQuelelelStore } from '@/store/updateQuelelel';
import { storeToRefs } from 'pinia';
import { onMounted } from 'vue';
import { useLocalStore } from '@/store/local';
const updateQuelelelStore = useUpdateQuelelelStore();
const localStore = useLocalStore();
const { sidebarSettings } = storeToRefs(localStore);

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
const clipped = ref(true);
const drawer = ref(false);
const miniVariant = ref(true);
const queuedStacks = ref([]);

const theme = useTheme();


onMounted(() => {
  if (sidebarSettings.value.mini) {
    miniVariant.value = sidebarSettings.value.mini;
  } else {
    miniVariant.value = false;
  }
  if (sidebarSettings.value.clipped === false) {
    clipped.value = sidebarSettings.value.clipped;
    drawer.value = false;
    return;
  }
  if (mdAndUp.value) {
    if (sidebarSettings.value.show !== undefined) {
      drawer.value = sidebarSettings.value.show;
    } else {
      drawer.value = mdAndUp.value;
    }
  } else {
    drawer.value = mdAndUp.value;
  }
});

function switchDrawer() {
  drawer.value = !drawer.value;
  localStore.updateSidebarSettings({ show: drawer.value });
}

function switchMini() {
  miniVariant.value = !miniVariant.value;
  localStore.updateSidebarSettings({ mini: miniVariant.value });
}

function switchClipped() {
  clipped.value = !clipped.value;
  localStore.updateSidebarSettings({ clipped: clipped.value });
}

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}


</script>
