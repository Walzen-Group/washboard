<template>
    <v-footer :absolute="true" app>
        <span class="mx-2">&copy; {{ new Date().getFullYear() }} Walzen Group</span>
        <v-icon color="grey lighten-1" size="20">mdi-git</v-icon>
        <a class="code-font ml-1" :href="`https://github.com/Walzen-Group/washboard/commit/${commitHash}`">
            {{ commitHash }}
        </a>
        <v-spacer></v-spacer>
        <a class="ml-3" href="https://www.flaticon.com/free-icon/washboard_7029393" title="washboard icons"
            >Icon</a
        >
    </v-footer>
</template>

<script lang="ts" setup>
import { computed } from "vue";

const commitHash = computed(() => {
    if (process.env.commitHash && process.env.commitHash.length > 0) {
        return process.env.commitHash.substring(0, 7);
    } else {
        return "jannis";
    }
});
</script>

<style lang="scss" scoped>
.code-font {
    font-family: "Fira Code", Fira Code, monospace;
}
</style>
