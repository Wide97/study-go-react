import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Durante lo sviluppo il browser richiede /api/health a Vite (porta 5173).
// Il proxy inoltra quella richiesta al backend Go (porta 8080) e rimuove
// il prefisso /api. Il backend riceve quindi la sua rotta reale: /health.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
