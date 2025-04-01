<template>
    <v-container>
        <div class="pa-4 d-flex justify-center">
            <v-card style="background-color: #94949458; backdrop-filter: blur(10px);" variant="flat">
                <v-card-text class="">
                    <div class="d-flex align-center justify-center">
                        <span class="text-h1 text-white">Hello</span>
                        <span class="text-h2 ml-3">❤️</span>
                    </div>
                    <div class="d-flex mt-6">
                        <v-row>
                            <v-col>
                                <div class="d-flex justify-center">
                                    <WashingMachine class="mr-3" animate :scale="0.6" />
                                </div>
                            </v-col>
                            <v-col class="mt-1">
                                <v-card variant="tonal" color="white">
                                    <v-card-title class="">
                                        <div>
                                            Image Updates
                                        </div>
                                        <v-progress-circular class="mt-2" indeterminate color="white" v-if="refreshing"
                                                             size="23" />
                                        <div v-else class="mt-2">
                                            {{ tweeenedOutdated.number.toFixed(0) }}
                                        </div>
                                    </v-card-title>
                                </v-card>
                                <div class="d-flex">
                                    <v-spacer />
                                    <v-btn :disabled="!portButtonEnabled" @click="getNextFreePort" prepend-icon="mdi-ferry" color="white" variant="tonal"
                                           class="mt-4">Gib Hafen</v-btn>
                                </div>
                            </v-col>
                        </v-row>

                    </div>
                </v-card-text>
            </v-card>
        </div>
        <div class="d-flex justify-center">
            <v-card v-if="viewHafen" style="background-color: #94949458; backdrop-filter: blur(10px);" variant="flat"
                    width="360">
                <v-card-title class="text-white">
                    <div>Nächster Hafen ist</div>
                    <div>{{ tweeenedHafen.number.toFixed(0) }}</div>
                </v-card-title>
            </v-card>
        </div>

    </v-container>
    <!-- <Landing /> -->
</template>

<script lang="ts" setup>
import { onMounted } from "vue";

import axios from "axios";
import { Container, Stack } from "@/types/types";
import gsap from "gsap";
import { storeToRefs } from "pinia";
import { useLocalStore } from "@/store/local";
import { awaitTimeout } from "@/api/lib";

const localStore = useLocalStore();
const { defaultStartPort } = storeToRefs(localStore);
const refreshing: Ref<boolean> = ref(false);
const items: Ref<Stack[]> = ref([]);
const connectionFailed: Ref<boolean> = ref(false);
const tweeenedOutdated = reactive({ number: 0 });
const tweeenedHafen = reactive({ number: 0 });
const viewHafen: Ref<boolean> = ref(false);
const portButtonEnabled: Ref<boolean> = ref(false);

async function init() {
    const response = await axios.get("/api/portainer/stacks", { params: { skeletonOnly: true } });
    items.value = response.data;
    portButtonEnabled.value = true;
}

async function leeroad() {
    refreshing.value = true;
    try {
        const request = axios.get("/api/portainer/stacks");
        const timeout = awaitTimeout(5000);
        const first = await Promise.any([request, timeout]);
        if (first === "loading") {
            await init();
        }
        const response = await request;
        console.log("leeroaded");
        items.value = response.data;
        connectionFailed.value = false;
        updateStatusCounts();
    } catch (error) {
        connectionFailed.value = true;
        console.log(error);
    }
    refreshing.value = false;
    portButtonEnabled.value = true;
}


function updateStatusCounts() {
    let outdated = 0;
    for (let stack of items.value) {
        //console.log(`stack name: ${stack.name}`);
        for (let container of stack.containers as Container[]) {
            if (container.upToDate === "outdated") {
                outdated += 1;
            }
        }
    }
    setTimeout(() => {
        gsap.to(tweeenedOutdated, { duration: 0.5, number: Number(outdated) || 0 });
    }, 200);
    return { outdated };
}

function getNextFreePort() {
    let port = defaultStartPort.value;
    let publicPorts: number[] = [];
    for (let stack of items.value) {
        for (let container of stack.containers as Container[]) {
            publicPorts.push(...container.ports.map((port: string) => Number(port.split(":")[0])));
        }
    }
    publicPorts = publicPorts.filter((port: number) => port >= Number(defaultStartPort.value));
    publicPorts.sort((a: number, b: number) => a - b);
    console.log(publicPorts);
    for (let i = 0; i < publicPorts.length; i++) {
        if (i + 1 < publicPorts[i]) {
            if (publicPorts[i + 1] != publicPorts[i] + 1) {
                port = publicPorts[i] + 1;
                console.log(`next free port: ${port}`)
                gsap.to(tweeenedHafen, { duration: 1, number: Number(port) || 0 });
                viewHafen.value = true;
                return;
            }
        }
    }
}

onMounted(setup);

async function setup() {
    await leeroad();
    updateStatusCounts();
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
