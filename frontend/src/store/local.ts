import { DockerUpdateManagerSettings } from '@/types/types'
import { defineStore } from 'pinia'

const STORE_NAME = 'local'

const getDockerUpdateManagerDefaultSettings = (): DockerUpdateManagerSettings => ({
    ignoredImages: {},
})

export const useLocalStore = defineStore(STORE_NAME, {
    state: () => ({
        dockerUpdateManagerSettings: getDockerUpdateManagerSettings(),
    }),
    actions: {
        updateDockerUpdateManagerSettings(partialSettings: any) {
            this.dockerUpdateManagerSettings = {
                ...this.dockerUpdateManagerSettings,
                ...partialSettings,
            }

            localStorage.setItem(STORE_NAME, JSON.stringify(this.dockerUpdateManagerSettings))
        },
    },
})

const getDockerUpdateManagerSettings = (): DockerUpdateManagerSettings => {
    const settings = localStorage.getItem(STORE_NAME)
    return settings ? JSON.parse(settings) : getDockerUpdateManagerDefaultSettings()
}
