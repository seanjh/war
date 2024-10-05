import type { Config } from 'tailwindcss'

export default {
    content: ["./templates/**/*.{html,js}"],
    theme: {
        extend: {
            gridTemplateColumns: {
                game: "1fr 3fr 1fr"
            },
        },
    },
    plugins: [],
} satisfies Config
