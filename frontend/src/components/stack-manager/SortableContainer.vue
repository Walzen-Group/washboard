<template>
    <!-- Content title -->
    <v-row class="fill-height row-sc" dense>
        <v-col>
            <div class="clickyboi"  @click="show = !show" >
                <div>
                    <slot name="title"></slot>
                </div>
                <div class="d-flex flex-wrap ga-2">
                    <slot name="shortcuts"></slot>
                </div>
            </div>
            <v-expand-transition>
                <!-- Content -->
                <div v-if="show">
                    <slot name="content"></slot>
                </div>
            </v-expand-transition>
        </v-col>
        <v-col cols="auto" class="pr-3 ml-1 pb-0 fill-height">
            <v-hover v-slot="{ isHovering, props }">
                <v-card
                        rounded="md"
                        @click.stop="show = !show"
                        variant="flat"
                        :color="isHovering ? 'rgba(0, 0, 0, 0.04)' : undefined"
                        v-bind="props"
                        class="d-flex align-center fill-height">
                    <v-icon :size="35">
                        {{ !show ? "mdi-chevron-down" : "mdi-chevron-up" }}
                    </v-icon>
                </v-card>
            </v-hover>
        </v-col>
    </v-row>
</template>

<script lang="ts" setup>

const show: Ref<boolean> = ref(false);
const emit = defineEmits([
    "click:expand"
]);
const props = defineProps<{
    name: string;
}>();
watch(show, () => {
    emit("click:expand", { name: props.name, expandState: show.value });
});
</script>

<style scoped lang="scss">
.row-sc {
    flex-wrap: nowrap;
}

.clickyboi {
    cursor: pointer;
}

</style>
