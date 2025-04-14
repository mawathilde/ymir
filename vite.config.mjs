import react from "@vitejs/plugin-react";
import sass from "sass";

/** @type {import('vite').UserConfig} */
export default {
  plugins: [react()],
  server: {
    host: true,
  },
  css: {
    preprocessorOptions: {
      scss: {
        implementation: sass,
      },
    },
  },
};
