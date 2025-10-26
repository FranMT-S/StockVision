import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";
import { useTheme } from "vuetify";

export const useGlobalStore = defineStore('global', () => {
  const Theme = useTheme();
  Theme.change(localStorage.getItem('theme') || 'light');
  const isDark = ref(localStorage.getItem('theme') === 'dark');

  const updateTheme = () =>{
    try {
      if (isDark.value) {
        document.body.classList.add('tw-dark')
      } else {
        document.body.classList.remove('tw-dark')
      }

      localStorage.setItem('theme', isDark.value ? 'dark' : 'light');
    } catch (error) {
      console.error('error saving theme:', error);
    }
  }

  updateTheme();

  watch(() => isDark.value, () => {
    updateTheme();
  });

  const toggleTheme = () => {
    isDark.value = !isDark.value;
    Theme.global.name.value = isDark.value ? 'dark' : 'light'
    updateTheme();
  };


  return {
    toggleTheme,
    isDark,
  };
});