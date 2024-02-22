<template>
  <v-sheet color="#121212">
    <v-img
           aspect-ratio="16/9"
           cover
           src="/img/washboard_v4.webp">
      <div class="text-h1 pa-4 ma-16">Hello!</div>
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

<style lang="scss" scoped>
.transparent-header-card {
  opacity: 0.6;
}
</style>

<route lang="yaml">
  meta:
    layout: index
  </route>
