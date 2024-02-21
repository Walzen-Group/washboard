import axios from 'axios';

async function updateStack(stackId: number) {
    const response = axios.put('/portainer-update-stack', {
        pullImage: true,
        prune: false,
        endpointId: 1,
        stackId: stackId
    });
    return response;
}

export { updateStack };
