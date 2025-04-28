/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        primary: "#FFCE1A",
        'primary-dark': "#e0b600",
        secondary: "#0D0842",
        blackBG: "#FFFFF8",
        favorite: "#FF5841",
      },
    },
  },
  plugins: [],
}
