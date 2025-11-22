/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: '#0a0a0a',
        'brand-light': '#666',
        'brand-gray': '#f8f8f8',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      animation: {
        'delay-100': 'delay 100ms',
        'delay-200': 'delay 200ms',
      },
      keyframes: {
        delay: {
          '0%, 80%, 100%': {
            opacity: '0.3',
            transform: 'scale(0.8)',
          },
          '40%': {
            opacity: '1',
            transform: 'scale(1)',
          },
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
