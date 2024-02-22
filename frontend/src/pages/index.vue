<template>
  <v-sheet color="#121212">
    <v-img
           aspect-ratio="16/9"
           cover
           src="/img/washboard_v4.webp">
      <div class="stroke text-h1 pa-4 ma-16">Hello!</div>
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

@mixin stroke($color: #000, $size: 1px) {
  text-shadow: -#{$size} -#{$size} 0 $color,
  0 -#{$size} 0 $color,
  #{$size} -#{$size} 0 $color,
  #{$size} 0 0 $color,
  #{$size} #{$size} 0 $color,
  0 #{$size} 0 $color,
  -#{$size} #{$size} 0 $color,
  -#{$size} 0 0 $color;
}

.stroke {
  @include stroke(black, 2px);
}
</style>

<route lang="yaml">
  meta:
    layout: index
  </route>
