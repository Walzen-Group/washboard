import axios from 'axios';

async function updateStack(stackId: number, endpointId: number = 1) {
    const response = axios.put(`/portainer/stacks/${stackId}/update`, {
        pullImage: true,
        prune: true,
        endpointId: endpointId,
    });
    return response;
}

async function stopStack(stackId: number, endpointId: number = 1) {
    console.log('stopStack', stackId, endpointId)
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
    const response = axios.get('/portainer/containers', {
        params: {
            endpointId: endpointId,
            stackName: stackName
        }
    });
    return response;

}

export { updateStack, stopStack, startStack, getContainers};
