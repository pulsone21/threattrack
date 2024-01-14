/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./frontend/**/*.{templ,html}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")],
  corePlugins: {
    preflight: true,
  },
};
