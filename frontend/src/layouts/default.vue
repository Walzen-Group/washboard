<template>
  <DefaultBar />
  <v-container class="max-container-width">
    <DefaultViewContainer />
  </v-container>
  <DefaultFooter />
</template>

<script lang="ts" setup>
import DefaultBar from './default/AppBar.vue'
import DefaultViewContainer from './default/ViewContainer.vue'
import DefaultFooter from './default/Footer.vue'

</script>


<style>
.max-container-width {
  max-width: 1800px !important;
}
</style>
