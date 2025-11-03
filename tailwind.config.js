/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'media',
  content: ['index.html', './src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: '#0C5FA0',
        accent: '#0EBAE3',
        night: '#0B132B',
      },
      fontFamily: {
        inter: ['Inter', 'system-ui', 'sans-serif'],
      },
      maxWidth: {
        content: '1200px',
      },
      boxShadow: {
        subtle: '0 20px 45px -30px rgba(12, 95, 160, 0.35)',
      },
    },
  },
  plugins: [],
};
