import axios, { AxiosError } from "axios";
import { Stack, StackInternal, Action, Container } from "@/types/types";
import { Store } from "pinia";

const webUILabel = "org.walzen.washb.webui";
const webUIAddressKey = "${ADDRESS}"

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
            console.log("refresh token invalid, redirecting to login");
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

export {
    updateStack,
    stopStack,
    startStack,
    getContainers,
    isAuthorized,
    callRefreshTokenRoute,
    webUILabel,
    webUIAddressKey,
    handleStackStateChange,
    getPortainerUrl,
    getContainerStatusCircleColor,
};
