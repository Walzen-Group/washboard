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
    updated: '#42A5F5'
  },
}

const washboardThemeDark: ThemeDefinition = {
  colors: {
    primary: '#64B5F6',
    updated: '#64B5F6'
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
