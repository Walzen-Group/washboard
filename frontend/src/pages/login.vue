<template>
    <v-container class="fill-height">
        <div class="mx-auto login-card">
            <v-card>
                <div>
                    <v-card-title class="py-5">It's time to wash</v-card-title>
                    <v-card-text>
                        <v-row align="center" justify="center">
                            <v-col cols="auto">
                                <WashingMachine animate :scale="0.8"></WashingMachine>
                            </v-col>
                            <v-col class="login-fields">
                                <v-alert
                                    v-if="loginError"
                                    type="error"
                                    color="deep-purple-lighten-2"
                                    class="mb-6"
                                    dense
                                    >{{ loginErrorMessage }}</v-alert
                                >
                                <v-text-field
                                    v-model="username"
                                    label="Name"
                                    class="mb-3"
                                    :rules="usernameRules"
                                    variant="filled"
                                >
                                </v-text-field>
                                <v-text-field
                                    @keydown.enter="submitLogin"
                                    v-model="password"
                                    label="Password"
                                    type="password"
                                    variant="filled"
                                >
                                </v-text-field>
                            </v-col>
                        </v-row>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn :loading="loading" class="mx-1 my-1" variant="tonal" @click="submitLogin"
                            >Start Laundry</v-btn
                        >
                    </v-card-actions>
                </div>
            </v-card>
        </div>
    </v-container>
</template>

<script lang="ts" setup>
import axios, { AxiosError } from "axios";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();

const usernameRules = [
    (value: string) => {
        if (/^[a-zA-Z0-9_.-]*$/.test(value)) return true;

        return "Last name can not contain digits.";
    },
];

const username: Ref<string> = ref("");
const password: Ref<string> = ref("");
const loginError: Ref<boolean> = ref(false);
const loading: Ref<boolean> = ref(false);
const loginErrorMessage: Ref<string> = ref("");

import { storeToRefs } from "pinia";
import { useAppStore } from "@/store/app";
import { connectWebSocket } from "@/api/lib";

const appStore = useAppStore();
const { renderInitCompleted } = storeToRefs(appStore);

async function submitLogin() {
    loginError.value = false;
    loading.value = true;
    try {
        const response = await axios.post("/api/auth/login", {
            username: username.value,
            password: password.value,
        });
        renderInitCompleted.value = false;
        if (route.query.redirect) {
            router.push({ path: route.query.redirect as string });
        } else {
            router.push({ path: "/" });
        }
        connectWebSocket();
    } catch (e) {
        if ((e as AxiosError).response?.status === 401) {
            loginError.value = true;
            loginErrorMessage.value = "Invalid username or password";
        } else {
            loginError.value = true;
            loginErrorMessage.value = "An error occurred while logging in";
        }
    }
    loading.value = false;
}
</script>

<style scoped>
.login-card {
    max-width: 450px;
    /* 'initial width' */
    width: 100%;
    min-width: 200px;
    /* 'minimum width' */
}

.login-fields {
    min-width: 200px;
    width: 100%;
}
</style>

<route lang="yaml">
meta:
    layout: login
</route>
