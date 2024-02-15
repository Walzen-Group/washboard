<template>
    <template>
    <v-data-table
        v-model="selected"
        :items="items"
        item-value="name"
        show-select
    ></v-data-table>
    </template>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const items = ref([]);
const selected: any = ref([]);


</script>
