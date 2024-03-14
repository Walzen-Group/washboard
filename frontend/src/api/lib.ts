import axios, { AxiosError } from "axios";

async function updateStack(stackId: number, endpointId: number = 1) {
    const response = axios.put(`/portainer/stacks/${stackId}/update`, {
        pullImage: true,
        prune: true,
        endpointId: endpointId,
    });
    return response;
}

async function stopStack(stackId: number, endpointId: number = 1) {
    console.log("stopStack", stackId, endpointId);
    const response = axios.post(`/portainer/stacks/${stackId}/stop`, {
        endpointId: endpointId,
    });
    return response;
}

async function startStack(stackId: number, endpointId: number = 1) {
    const response = axios.post(`/portainer/stacks/${stackId}/start`, {
        endpointId: endpointId,
    });
    return response;
}

async function getContainers(stackName: string, endpointId: number = 1) {
    const response = axios.get("/portainer/containers", {
        params: {
            endpointId: endpointId,
            stackName: stackName,
        },
    });
    return response;
}

async function isAuthorized() {
    try {
        await axios.get("/");
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
        await axios.post("/auth/refresh_token");
        return true;
    } catch (e) {
        const error = e as AxiosError;
        if (error.response?.status === 401) {
            console.log("refresh token invalid, redirecting to login");
        }
    }
    return false;
}
export { updateStack, stopStack, startStack, getContainers, isAuthorized, callRefreshTokenRoute };
