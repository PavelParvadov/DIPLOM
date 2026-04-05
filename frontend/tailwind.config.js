export default {
    content: ["./index.html", "./src/**/*.{ts,tsx}"],
    theme: {
        extend: {
            colors: {
                sand: "#f5efe4",
                ink: "#1d252f",
                ember: "#d86a3b",
                moss: "#385b57",
                cloud: "#f8fafc",
            },
            boxShadow: {
                soft: "0 20px 55px rgba(31, 41, 55, 0.12)",
            },
            fontFamily: {
                sans: ["Manrope", "Segoe UI", "sans-serif"],
                display: ["Fraunces", "Georgia", "serif"],
            },
        },
    },
    plugins: [],
};
