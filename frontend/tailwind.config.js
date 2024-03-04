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
        'scroll': 'scroll 10s ease-in-out alternate infinite',
      },
    },
  },
  plugins: [],
};
