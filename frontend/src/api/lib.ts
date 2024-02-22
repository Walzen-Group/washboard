import axios from 'axios';

async function updateStack(stackId: number, endpointId: number = 1) {
    const response = axios.put('/portainer-update-stack', {
        pullImage: true,
        prune: false,
        endpointId: endpointId,
        stackId: stackId
    });
    return response;
}

async function stopStack(stackId: number, endpointId: number = 1) {
    console.log('stopStack', stackId, endpointId)
    const response = axios.post('/portainer-stop-stack', {
        endpointId: endpointId,
        stackId: stackId
    });
    return response;
}

async function startStack(stackId: number, endpointId: number = 1) {
    const response = axios.post('/portainer-start-stack', {
        endpointId: endpointId,
        stackId: stackId
    });
    return response;
}

async function getContainers(stackName: string, endpointId: number = 1) {
    const response = axios.get('/portainer-get-containers', {
        params: {
            endpointId: endpointId,
            stackName: stackName
        }
    });
    return response;

}

export { updateStack, stopStack, startStack, getContainers};
