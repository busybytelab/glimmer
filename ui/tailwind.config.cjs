/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class', // Changed from 'media' to 'class' for manual theme toggling
  theme: {
    extend: {
      colors: {
        primary: '#2c3e50',
        secondary: '#3498db',
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
} 