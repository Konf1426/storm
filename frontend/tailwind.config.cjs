/** @type {import("tailwindcss").Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js}", "./src/components/**/*.{vue,js}"],
  theme: {
    extend: {
      fontFamily: {
        display: ["Space Grotesk", "ui-sans-serif", "system-ui"],
        mono: ["IBM Plex Mono", "ui-monospace", "SFMono-Regular"],
      },
      colors: {
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        border: "hsl(var(--border))",
        ring: "hsl(var(--ring))",
      },
      boxShadow: {
        glow: "0 0 0 1px hsl(var(--ring) / 0.2), 0 12px 40px hsl(var(--ring) / 0.15)",
      },
    },
  },
  plugins: [],
}
