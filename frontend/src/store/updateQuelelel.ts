import { QueueItem, QueueStatus, UpdateQueue } from "@/types/types";
import { defineStore } from "pinia";
import { ref, Ref } from 'vue';

const STORE_NAME = "updateQuelelel";

export const useUpdateQuelelelStore = defineStore(STORE_NAME, () => {
    const queue: Ref<UpdateQueue> = ref({ [QueueStatus.Done]: {}, [QueueStatus.Queued]: {}, [QueueStatus.Error]: {} });
    const queueCount: Ref<number> = ref(0);
    const queueItems: Ref<QueueItem[]> = ref([]);

    function update(queueItems: UpdateQueue) {
        queue.value = queueItems;
        if (QueueStatus.Queued in queueItems) {
            queueCount.value = Object.keys(queueItems[QueueStatus.Queued]).length;
        } else {
            queueCount.value = 0;
        }
    }

    function enqueue(
        stack: QueueItem
    ) {
        if (!queue.value[stack.status]) {
            queue.value[stack.status] = {};
        }
        queue.value[stack.status][stack.stackName] = stack;
        if (stack.status === QueueStatus.Queued) {
            queueCount.value++;
        }
    }

    function dequeue(stackName: number) {
        for (const status in queue.value) {
            if (status in QueueStatus) {
                const queueStatusKey = status as QueueStatus;
                if (queue.value[queueStatusKey][stackName]) {
                    if (queueStatusKey === QueueStatus.Queued) {
                        queueCount.value--;
                    }
                    delete queue.value[queueStatusKey][stackName];
                }
            }
        }
    }

    return {
        queue,
        queueItems,
        queueCount,
        update,
        enqueue,
        dequeue
    };
})
