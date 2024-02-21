<template>
  <v-sheet color="#121212">
    <v-img
           aspect-ratio="16/9"
           cover
           src="/img/washboard_v4.webp">
    </v-img>
  </v-sheet>
  <!-- <Landing /> -->
</template>

<script lang="ts" setup>
import { onMounted } from 'vue';

onMounted(setup);

function setup() {
}
</script>

<route lang="yaml">
  meta:
    layout: index
  </route>
