const plugin = require("tailwindcss/plugin")
const path = require("path")

module.exports = {
  content: [
    "./js/**/*.js",
    "../lib/test_app_web.ex",
    "../lib/test_app_web/**/*.*ex",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
