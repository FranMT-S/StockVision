/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  prefix: 'tw-',
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    screens: {
      sm: '480px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
      '2xl': '1536px',
    },
    extend: {},
  },  
  plugins: [],
}
