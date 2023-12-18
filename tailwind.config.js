/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./web/templates/**/*.{html,js,templ}'],
  theme: {
    fontFamily: {
      'sans': ['"Montserrat"'],
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

