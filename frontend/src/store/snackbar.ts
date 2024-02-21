import { Snackbar } from "@/types/types";
import { defineStore } from "pinia";
import { ref, Ref } from 'vue';

const STORE_NAME = "snackbar";
const TIMEOUT = 5000;

export const useSnackbarStore = defineStore(STORE_NAME, () => {
    const snackbars: Ref<Snackbar[]> = ref([]);

    function addSnackbar(
        id: string,
        message: string,
        color: string = "success",
        show: boolean = true
    ) {
        const snackbar = { id, message, color, show };
        snackbars.value.push(snackbar);
        setTimeout(() => {
            removeSnackbar(id);
        }, TIMEOUT);
    }

    function removeSnackbar(id: string) {
        setTimeout(function () {
            const index = snackbars.value.findIndex(
                (snackbar: Snackbar) => snackbar.id === id
            );
            if (index !== -1) {
                snackbars.value.splice(index, 1);
            }
        }, 300);
    }

    function calcSnackbarMargin(index: number) {
        return `${index * 60 + 10}px`;
    }

    return {
        snackbars,
        addSnackbar,
        removeSnackbar,
        calcSnackbarMargin,
    };
})
