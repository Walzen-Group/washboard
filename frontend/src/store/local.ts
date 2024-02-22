import { DockerUpdateManagerSettings, SidebarSettings } from '@/types/types'
import { defineStore } from 'pinia'

const STORE_NAME = 'local';
const SIDEBAR_SETTINGS = 'sidebarSettings';

const getDockerUpdateManagerDefaultSettings = (): DockerUpdateManagerSettings => ({
    ignoredImages: {},
})

const getSidebarDefaultSettings = (): SidebarSettings => ({
    show: undefined,
    mini: undefined,
    clipped: undefined,
})

export const useLocalStore = defineStore(STORE_NAME, {
    state: () => ({
        dockerUpdateManagerSettings: getDockerUpdateManagerSettings(),
        sidebarSettings: getSidebarSettings(),
    }),
    actions: {
        updateDockerUpdateManagerSettings(partialSettings: any) {
            this.dockerUpdateManagerSettings = {
                ...this.dockerUpdateManagerSettings,
                ...partialSettings,
            }
            localStorage.setItem(STORE_NAME, JSON.stringify(this.dockerUpdateManagerSettings))
        },
        updateSidebarSettings(partialSettings: any) {
            this.sidebarSettings = {
                ...this.sidebarSettings,
                ...partialSettings,
            }
            localStorage.setItem(SIDEBAR_SETTINGS, JSON.stringify(this.sidebarSettings))
        },
    },
})

const getSidebarSettings = (): SidebarSettings => {
    const settings = localStorage.getItem(SIDEBAR_SETTINGS)
    return settings ? JSON.parse(settings) : getSidebarDefaultSettings()
}

const getDockerUpdateManagerSettings = (): DockerUpdateManagerSettings => {
    const settings = localStorage.getItem(STORE_NAME)
    return settings ? JSON.parse(settings) : getDockerUpdateManagerDefaultSettings()
}
