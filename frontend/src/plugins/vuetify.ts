/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify, type ThemeDefinition } from 'vuetify'
import { md3 } from 'vuetify/blueprints'


const washboardThemeLight: ThemeDefinition = {
  colors: {
    background: '#F3F4F6',
    primary: '#1E88E5',
    updated: '#42A5F5',
    surface: '#fbfbfb',
    wurple: '#CE93D8',
    "washboard-appbar": "#0277BD",
    "stop": "#D81B60"
  },
}

const washboardThemeDark: ThemeDefinition = {
  colors: {
    primary: '#64B5F6',
    updated: '#0277BD',
    wurple: '#BA68C8',
    "surface-variant": "#c2c2c2",
    "washboard-appbar": "#01579B",
    "stop": "#FF4081"
  },
}

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  blueprint: md3,
  theme: {
    themes: {
      light: washboardThemeLight,
      dark: washboardThemeDark
    },
  },
})
