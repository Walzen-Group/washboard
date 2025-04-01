// Plugins
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import Fonts from "unplugin-fonts/vite";
import Layouts from "vite-plugin-vue-layouts";
import Vue from "@vitejs/plugin-vue";
import VueRouter from "unplugin-vue-router/vite";
import Vuetify, { transformAssetUrls } from "vite-plugin-vuetify";

// Utilities
import { defineConfig } from "vite";
import { fileURLToPath, URL } from "node:url";
import { execSync } from 'child_process';

// Get the shortened Git commit ID
// const gitCommitIdShort = execSync('git rev-parse --short HEAD').toString().trim();


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    VueRouter(),
    Layouts(),
    Vue({
      template: { transformAssetUrls },
    }),
    // https://github.com/vuetifyjs/vuetify-loader/tree/master/packages/vite-plugin#readme
    Vuetify({
      autoImport: true,
      styles: {
        configFile: "src/styles/settings.scss",
      },
    }),
    Components(),
    Fonts({
      google: {
        families: [
          {
            name: "Roboto",
            styles: "wght@100;300;400;500;700;900",
          },
          {
            name: "Inter",
            styles: "wght@100;300;400;500;700;900",
          },
        ],
      },
    }),
    AutoImport({
      imports: ["vue", "vue-router"],
      dts: true,
      eslintrc: {
        enabled: true,
      },
      vueTemplate: true,
    }),
  ],
  define: {
    "process.env": {
      PORTAINER_BASE_URL: "http://10.10.194.5:5221",
      PORTAINER_QUERY_TEMPlATE:
        "/#!/${endpointId}/docker/stacks/${stackName}?id=${stackId}&type=2&regular=true&external=false&orphaned=false",
      PORTAINER_DEFAULT_ENDPOINT_ID: 1,
      //commitHash: gitCommitIdShort,
      commitHash: process.env.COMMIT_ID_SHORT
    },
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
    extensions: [".js", ".json", ".jsx", ".mjs", ".ts", ".tsx", ".vue"],
  },
  server: {
    port: 3000,
  },
  build: {
    modulePreload: { polyfill: true },
  },
  css: {
    preprocessorOptions: {
      sass: {
        api: 'modern-compiler' // or "modern"
      }
    }
  }
});
