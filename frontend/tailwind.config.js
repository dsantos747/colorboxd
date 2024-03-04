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
    },
  },
  plugins: [],
};
