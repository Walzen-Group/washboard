<template>
    <v-footer :absolute="true" app>
        <span class="mx-2">&copy; {{ new Date().getFullYear() }} Walzen Group</span>
        <v-icon color="grey lighten-1" size="20">mdi-git</v-icon>
        <a class="code-font ml-1"
           :href="`https://github.com/Walzen-Group/washboard/commit/${commitHash}`">
            {{ commitHash }}
        </a>
    </v-footer>
</template>

<script lang="ts" setup>
import { reactive, computed } from 'vue';


const commitHash = computed(() => {
    if (process.env.commitHash) {
        return process.env.commitHash.substring(0, 7);
    } else {
        return "jannis";
    }
});
</script>
