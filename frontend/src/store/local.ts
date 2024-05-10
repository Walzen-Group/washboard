import { DockerUpdateManagerSettings, SidebarSettings, URLConfig } from '@/types/types'
import { defineStore } from 'pinia'

const STORE_NAME = 'local';
const SIDEBAR_SETTINGS = 'sidebarSettings';
const URL_STORE_NAME = 'urlConfig';

const getDockerUpdateManagerDefaultSettings = (): DockerUpdateManagerSettings => ({
    ignoredImages: {},
})

const getSidebarDefaultSettings = (): SidebarSettings => ({
    show: undefined,
    mini: undefined,
    clipped: undefined,
})

const getDefaultURLConfig = (): URLConfig => ({
    defaultHost: location.hostname,
    defaultPortainerAddress: process.env.PORTAINER_BASE_URL || `${location.hostname}:9000`
})

export const useLocalStore = defineStore(STORE_NAME, {
    state: () => ({
        dockerUpdateManagerSettings: getDockerUpdateManagerSettings(),
        sidebarSettings: getSidebarSettings(),
        urlConfig: getURLConfig(),
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
        updateURLConfig(partialSettings: any) {
            this.urlConfig = {
                ...this.urlConfig,
                ...partialSettings,
            }
            localStorage.setItem(URL_STORE_NAME, JSON.stringify(this.urlConfig))
        },
        resetURLConfig() {
            resetURLConfig()
            this.urlConfig = getDefaultURLConfig()
        }
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

const getURLConfig = (): URLConfig => {
    const settings = localStorage.getItem(URL_STORE_NAME)
    return settings ? JSON.parse(settings) : getDefaultURLConfig()
}

const resetURLConfig = () => {
    localStorage.removeItem(URL_STORE_NAME)
}
