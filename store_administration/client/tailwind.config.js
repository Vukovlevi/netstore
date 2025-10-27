module.exports = {
  darkMode: "class",
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}", // include all source files
  ],
  theme: {
    extend: {
      colors: {
        primary: "#1173d4",
        "background-light": "#f6f7f8",
        "background-dark": "#101922",
        "foreground-light": "#111418",
        "foreground-dark": "#f6f7f8",
        "input-light": "#ffffff",
        "input-dark": "#1a2734",
        "border-light": "#dbe0e6",
        "border-dark": "#34495e",
        "placeholder-light": "#617589",
        "placeholder-dark": "#9ab0c4",
      },
      fontFamily: {
        display: ["Work Sans", "Noto Sans", "sans-serif"],
      },
      borderRadius: {
        DEFAULT: "0.25rem",
        lg: "0.5rem",
        xl: "0.75rem",
        full: "9999px",
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/container-queries"),
  ],
};
