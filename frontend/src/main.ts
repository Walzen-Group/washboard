/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'

// Axios
import axios from 'axios'

const app = createApp(App)


axios.defaults.withCredentials = true
if (process.env.NODE_ENV !== 'production') {
  axios.defaults.baseURL = `http://${location.hostname}:8080`;
  console.log("Development mode, base url: " + axios.defaults.baseURL);
}
else {
  axios.defaults.baseURL = `http://${location.hostname}:8080`;
  //axios.defaults.baseURL = location.protocol + '//' + location.host;
  console.log("Production mode, base urL: " + axios.defaults.baseURL);
}

registerPlugins(app)

app.mount('#app')
