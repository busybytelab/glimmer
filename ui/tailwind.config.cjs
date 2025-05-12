/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'media', // Use 'media' for OS preference detection
  theme: {
    extend: {
      colors: {
        primary: '#2c3e50',
        secondary: '#3498db',
      }
    },
  },
  plugins: [],
} 