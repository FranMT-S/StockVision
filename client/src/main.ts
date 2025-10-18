import { createApp } from 'vue'
import './assets/style.css'
import './assets/main.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import vuetify from './plugins/vuetify'
import router from './router'

const pinia = createPinia()

createApp(App).use(vuetify).use(pinia).use(router).mount('#app')
