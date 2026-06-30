/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        accent: '#ecad0a',
        primary: '#209dd7',
        secondary: '#753991',
        navy: '#032147',
        graytext: '#888888',
      },
    },
  },
  plugins: [],
};
