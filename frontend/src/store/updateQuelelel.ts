import { defineStore } from "pinia";
import { ref, Ref } from 'vue';

const STORE_NAME = "updateQuelelel";

export const useUpdateQuelelelStore = defineStore(STORE_NAME, () => {
    const queue: Ref<UpdateQueue> = ref({"done": {}, "queued": {}, "error": {}});

    function update(queueItems: UpdateQueue) {
        queue.value = queueItems;
    }

    function enqueue(
        stack: QueueItem
    ) {
        if (!queue.value[stack.status]) {
            queue.value[stack.status] = {};
        }
        queue.value[stack.status][stack.stackName] = stack;
    }

    function dequeue(stackName: number) {
        for (const status in queue.value) {
            if (status in QueueStatus) {
                const queueStatusKey = status as QueueStatus;
                if (queue.value[queueStatusKey][stackName]) {
                    delete queue.value[queueStatusKey][stackName];
                }
            }
        }
    }

    return {
        queue,
        update,
        enqueue,
        dequeue
    };
})
