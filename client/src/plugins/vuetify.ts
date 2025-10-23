import '@mdi/font/css/materialdesignicons.css' // Ensure you are using css-loader
import { createVuetify } from 'vuetify'
import 'vuetify/styles'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'


export default createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi', // This is already the default value - only for display purposes
  },
  theme: {
    defaultTheme: 'dark',
    themes: {
      light: {
        dark: false,
        colors: {
          primary: '#1976D2', 
          secondary: '#424242',
          accent: '#82B1FF',
          background: '#FFFFFF',
          surface: '#FFFFFF',
          active: '#1976D2',
          inactive: '#424242',
          error: '#FF5252',
          info: '#2196F3',
          success: '#4CAF50',
          warning: '#FFC107',
        },
      },
      dark: {
        dark: true,
        colors: {
          primary: '#42b883', 
          secondary: '#00FF9C',
          accent: '#42b883',
          background: '#0A0A0A', 
          surface: '#121212',
          active: '#2dcf78',
          inactive: '#294D3F',
          error: '#FF5252',
          info: '#39FF14',
          success: '#39FF14',
          warning: '#F5FF14',
        },
      },
    },
  },
})


