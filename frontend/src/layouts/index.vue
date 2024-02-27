<template>
    <DefaultBar />
    <v-sheet color="#121212" class="full-height">
        <v-img
               aspect-ratio="16/9"
               cover
               src="/img/washboard_v4.webp">
            <DefaultViewContainer />
        </v-img>
    </v-sheet>
    <DefaultFooter />
</template>

<script lang="ts" setup>
import DefaultBar from './default/AppBar.vue'
import DefaultViewContainer from './default/ViewContainer.vue'
import DefaultFooter from './default/Footer.vue'
</script>

<style scoped lang="scss">
.full-height {
    height: 100%;
}
</style>
