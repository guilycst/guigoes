/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./web/templates/**/*.{html,js,templ,md}','./posts/**/*.{html,js,templ,md}'],
  theme: {
    screens: {
      sm: '480px',
      md: '768px',
      lg: '976px',
      xl: '1440px',
    },
    fontFamily: {
      'sans': ['Montserrat', 'sans-serif'],
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

