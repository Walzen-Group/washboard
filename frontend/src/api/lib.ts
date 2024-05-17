import axios, { AxiosError } from "axios";
import { Stack, StackInternal, Action, Container, UpdateQueue, QueueItem, QueueStatus } from "@/types/types";
import { Store, storeToRefs } from "pinia";
import { useLocalStore } from "@/store/local";
import { useSnackbarStore } from "@/store/snackbar";
import { useUpdateQuelelelStore } from "@/store/updateQuelelel";
import { useAppStore } from "@/store/app";

const webUILabel = "org.walzen.washb.webui";
const webUIAddressKey = "${ADDRESS}";

const defaultEndpointId = process.env.PORTAINER_DEFAULT_ENDPOINT_ID || "1";

async function updateStack(stackId: number, endpointId: number = 1) {
    const response = axios.put(`/api/portainer/stacks/${stackId}/update`, {
        pullImage: true,
        prune: true,
        endpointId: endpointId,
    });
    return response;
}

async function stopStack(stackId: number, endpointId: number = 1) {
    console.log("stopStack", stackId, endpointId);
    const response = axios.post(`/api/portainer/stacks/${stackId}/stop`, {
        endpointId: endpointId,
    });
    return response;
}

async function startStack(stackId: number, endpointId: number = 1) {
    const response = axios.post(`/api/portainer/stacks/${stackId}/start`, {
        endpointId: endpointId,
    });
    return response;
}

async function getContainers(stackName: string, endpointId: number = 1) {
    const response = axios.get("/api/portainer/containers", {
        params: {
            endpointId: endpointId,
            stackName: stackName,
        },
    });
    return response;
}

async function isAuthorized() {
    try {
        await axios.get("/api");
        return true;
    } catch (e) {
        if ((e as AxiosError).response?.status == 401) {
            console.log("Unauthorized");
            return false;
        }
    }
    return false;
}

async function callRefreshTokenRoute() {
    try {
        await axios.post("/api/auth/refresh_token");
        return true;
    } catch (e) {
        const error = e as AxiosError;
        if (error.response?.status === 401) {
            console.log("refresh token invalid");
        }
    }
    return false;
}

let loaderState = {} as { [key: string]: boolean };

// Handle response of start/stop/restart stack
async function handleStackStateChange(
    stack: Stack,
    action: string,
    response: any,
    snackbarsStore: any,
    stacksInternal: Ref<Stack[]>,
    pullImageStatus: boolean = true
) {
    if (response.status === 200) {
        if (action === Action.Start || action === Action.Restart) {
            await updateContainers(stack, stacksInternal, pullImageStatus);
        } else {
            clearStackContainers(stack, stacksInternal);
        }
        snackbarsStore.addSnackbar(`${stack.id}_startstop`, `Successfully ${action}ed ${stack?.name}`, "success");
    } else {
        throw new Error(`Received unexpected response status: ${response.status}`);
    }
}

function clearStackContainers(stack: Stack, stacksInternal: Ref<Stack[]>) {
    stacksInternal.value = stacksInternal.value.map((item) =>
        item.id === stack.id ? { ...item, containers: [] } : item
    );
}

async function updateContainers(stack: Stack, stacksInternal: Ref<Stack[]>, pullImageStatus: boolean = true) {
    const containersResponse = await getContainers(stack.name, parseInt(defaultEndpointId));
    let containers: Container[] = containersResponse.data;
    if (pullImageStatus) {
        containers = await Promise.all(containers.map(async (container) => updateContainerStatus(container)));
    }
    stacksInternal.value = stacksInternal.value.map((item) =>
        item.id === stack.id ? { ...item, containers } : item
    );
}

async function updateContainerStatus(container: Container) {
    const containerImageStatusResponse = await axios.get(`/api/portainer/image-status`, {
        params: { endpointId: defaultEndpointId, containerId: container.id },
    });
    const containerImageStatus = containerImageStatusResponse.data;
    return { ...container, upToDate: containerImageStatus.status };
}

function getPortainerUrl(item: Stack, itemUrl: string) {
    return itemUrl.replace("${stackId}", item.id.toString()).replace("${stackName}", item.name);
}

function getContainerStatusCircleColor(status: string) {
    // suwitcho caseo
    switch (status) {
        case "running":
            return "snakegreen";
        case "exited":
            return "stop";
        default:
            return "grey";
    }
}

function getFirstContainerIcon(containers: Container[]): string | undefined {
    // get random number between 0 and 1
    if (Math.random() < 0.995) {
        const ico = containers.find((container) => container.labels["net.unraid.docker.icon"])?.labels[
            "net.unraid.docker.icon"
        ];
        return ico;
    } else {
        return `/img/craughing.png`;
    }
}

function connectWebSocket() {
    const updateQuelelelStore = useUpdateQuelelelStore();
    const { queue: stackQueue } = storeToRefs(updateQuelelelStore);
    const snackbarsStore = useSnackbarStore();
    const appStore = useAppStore();
    const { webSocketStacksUpdate } = storeToRefs(appStore);

    console.log("connecting to websocket...");
    let wsAddr = `${axios.defaults.baseURL}/api/ws/stacks-update`
        .replace("http://", "ws://")
        .replace("https://", "wss://");
    let socket = new WebSocket(wsAddr);
    webSocketStacksUpdate.value = socket;
    socket.onmessage = function (event) {
        let data: UpdateQueue = JSON.parse(event.data);

        for (let [newStatus, newItems] of Object.entries(data) as [QueueStatus, Record<string, QueueItem>][]) {
            for (let stackName in newItems) {
                const queueItem = newItems[stackName];

                let previousBucket: string | undefined = undefined;
                for (let [oldStatus, oldItems] of Object.entries(stackQueue.value)) {
                    if (queueItem.stackName in oldItems) {
                        previousBucket = oldStatus;
                        break;
                    }
                }

                switch (newStatus) {
                    case QueueStatus.Queued:
                        break;
                    case QueueStatus.Done:
                        if (previousBucket && previousBucket != newStatus) {
                            snackbarsStore.addSnackbar(
                                `${queueItem.stackId}_update`,
                                `Stack ${queueItem.stackName} updated successfully`,
                                "success"
                            );
                        }
                        break;
                    case QueueStatus.Error:
                        if (previousBucket && previousBucket != newStatus) {
                            snackbarsStore.addSnackbar(
                                `${queueItem.stackId}_update`,
                                `Stack ${queueItem.stackName} update failed`,
                                "error"
                            );
                        }
                        break;
                }
            }
        }
        updateQuelelelStore.update(data);
    };

    socket.onclose = function (event) {
        console.log("Socket is closed. Reconnect will be attempted in 1 second.", event.reason);
        setTimeout(function () {
            connectWebSocket();
        }, 1000);
    };

    socket.onerror = function () {
        console.error("Socket encountered error, closing socket");
        socket.close();
    };
}

// functions
async function awaitTimeout(delay: number) {
    return new Promise((resolve) => setTimeout(resolve, delay, "loading"));
}

export {
    updateStack,
    stopStack,
    startStack,
    getContainers,
    isAuthorized,
    callRefreshTokenRoute,
    webUILabel,
    webUIAddressKey,
    getFirstContainerIcon,
    handleStackStateChange,
    getPortainerUrl,
    getContainerStatusCircleColor,
    connectWebSocket,
    awaitTimeout,
};
