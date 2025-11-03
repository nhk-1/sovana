/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'media',
  content: ['index.html', './src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        midnight: '#0B132B',
        ocean: '#3A506B',
      },
      fontFamily: {
        inter: ['Inter', 'system-ui', 'sans-serif'],
      },
      maxWidth: {
        content: '1200px',
      },
      boxShadow: {
        subtle: '0 20px 45px -30px rgba(11, 19, 43, 0.45)',
      },
    },
  },
  plugins: [],
};
