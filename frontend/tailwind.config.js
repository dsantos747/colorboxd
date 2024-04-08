/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      keyframes: {
        scroll: {
          '0%': { transform: 'translateY(0%)' },
          '100%': { transform: 'translateY(-85%)' },
        },
      },
      animation: {
        'scroll': 'scroll 18s ease-in-out alternate infinite',
      },
      colors: {
        'cb1': '#2563eb',
        'cb2': '#3170d0',
        'cb3': '#3d7db6',
        'cb4': '#498a9b',
        'cb5': '#549880',
        'cb6': '#60a566',
        'cb7': '#6cb24b',
        'cb8': '#78bf31',
        'cb9': '#84cc16',
      },
    },
  },
  plugins: [],
};
