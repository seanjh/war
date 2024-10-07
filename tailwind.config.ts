import type { Config } from 'tailwindcss'

export default {
    darkMode: 'media',
    content: ["./templates/**/*.{html,js}"],
    theme: {
        extend: {
            fontFamily: { sans: ['Inter', 'sans-serif'] },
            gridTemplateColumns: { game: "1fr 3fr 1fr" },
        },
    },
    plugins: [],
} satisfies Config
