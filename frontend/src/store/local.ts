import { DockerUpdateManagerSettings, SidebarSettings, URLConfig } from '@/types/types'
import { defineStore } from 'pinia'

const STORE_NAME = 'local';
const SIDEBAR_SETTINGS = 'sidebarSettings';
const URL_STORE_NAME = 'urlConfig';
const PORT_SETTING = 'defaultStartPort';

export const useLocalStore = defineStore(STORE_NAME, () => {
    const dockerUpdateManagerSettings: Ref<DockerUpdateManagerSettings> = ref(getDockerUpdateManagerSettings())
    const sidebarSettings: Ref<SidebarSettings> = ref(getSidebarSettings())
    const urlConfig: Ref<URLConfig> = ref(getURLConfig())
    const defaultStartPort: Ref<Number> = ref(getDefaultStartPort())

    function updateDockerUpdateManagerSettings(partialSettings: any) {
        dockerUpdateManagerSettings.value = {
            ...dockerUpdateManagerSettings.value,
            ...partialSettings,
        }
        localStorage.setItem(STORE_NAME, JSON.stringify(dockerUpdateManagerSettings.value))
    }

    function updateSidebarSettings(partialSettings: any) {
        sidebarSettings.value = {
            ...sidebarSettings.value,
            ...partialSettings,
        }
        localStorage.setItem(SIDEBAR_SETTINGS, JSON.stringify(sidebarSettings.value))
    }

    function updateURLConfig(partialSettings: any) {
        urlConfig.value = {
            ...urlConfig.value,
            ...partialSettings,
        }
        localStorage.setItem(URL_STORE_NAME, JSON.stringify(urlConfig.value))
    }

    function resetURLConfig() {
        localStorage.removeItem(URL_STORE_NAME)
        urlConfig.value = getDefaultURLConfig()
    }

    function updateDefaultStartPort(newPort: number) {
        defaultStartPort.value = newPort
    }

    return {
        urlConfig,
        dockerUpdateManagerSettings,
        sidebarSettings,
        defaultStartPort,

        updateDockerUpdateManagerSettings,
        updateSidebarSettings,
        updateURLConfig,
        resetURLConfig,
        updateDefaultStartPort
    }
})

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

const getDefaultStartPort = (): number => {
    const port = localStorage.getItem(PORT_SETTING)
    return port ? parseInt(port) : 10000
}

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
