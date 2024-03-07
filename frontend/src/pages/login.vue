<template>
    <v-container class="fill-height">
        <v-card class="mx-auto" width="400">
            <v-card-title class="py-5">It's time to wash</v-card-title>
            <v-card-text>
                <v-alert v-if="loginError" type="error" color="deep-purple-lighten-2" class="mb-6" dense>{{ loginErrorMessage
                    }}</v-alert>
                <v-text-field v-model="username" label="Name"
                              class="mb-3"
                              :rules="usernameRules"
                              variant="filled">
                </v-text-field>
                <v-text-field @keydown.enter="submitLogin" v-model="password" label="Password"
                              type="password"
                              variant="filled">
                </v-text-field>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn :loading="loading" class="mx-1 my-1" variant="tonal"
                       @click="submitLogin">Start Laundry</v-btn>
            </v-card-actions>
        </v-card>
    </v-container>
</template>

<script lang="ts" setup>
import axios, { AxiosError } from 'axios';
import { useRouter, useRoute } from 'vue-router';


const router = useRouter();
const route = useRoute();

const usernameRules = [
    (value: string) => {
        if (/^[a-zA-Z0-9_.-]*$/.test(value)) return true;

        return 'Last name can not contain digits.'
    },
];

const username: Ref<string> = ref('');
const password: Ref<string> = ref('');
const loginError: Ref<boolean> = ref(false);
const loading: Ref<boolean> = ref(false);
const loginErrorMessage: Ref<string> = ref('');

import { storeToRefs } from 'pinia';
import { useAppStore } from '@/store/app';

const appStore = useAppStore();
const { renderInitCompleted } = storeToRefs(appStore);

async function submitLogin() {
    loginError.value = false;
    loading.value = true;
    try {
        const response = await axios.post("/auth/login", {
            username: username.value,
            password: password.value
        });
        renderInitCompleted.value = false;
        if (route.query.redirect) {
            router.push({ path: route.query.redirect as string });
        } else {
            router.push({ path: '/' });
        }

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

<style scoped></style>

<route lang="yaml">meta: layout: login </route>
