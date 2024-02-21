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

const washboardThemeLight: ThemeDefinition = {
  colors: {
    primary: '#039BE5',
    updated: '#42A5F5',
    surface: '#fafafa',
    "washboard-appbar": "#0288D1",
  },
}

const washboardThemeDark: ThemeDefinition = {
  colors: {
    primary: '#64B5F6',
    updated: '#0277BD',
    "surface-variant": "#c2c2c2",
    "washboard-appbar": "#01579B",
  },
}

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    themes: {
      light: washboardThemeLight,
      dark: washboardThemeDark
    },
  },
})
