/// <reference types="vitest/config" />

import fs from "node:fs";
import path from "node:path";
import babel from "@rolldown/plugin-babel";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [
    devtools(),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    react(),
    babel({ presets: [reactCompilerPreset({ target: "19" })] }),
    {
      name: "inject-splash",
      transformIndexHtml(html) {
        const splashHtml = fs.readFileSync(
          path.resolve(__dirname, "src/splash/splash.html"),
          "utf-8",
        );
        const splashCss = fs.readFileSync(
          path.resolve(__dirname, "src/splash/splash.css"),
          "utf-8",
        );
        const splashJs = fs.readFileSync(
          path.resolve(__dirname, "src/splash/splash.js"),
          "utf-8",
        );
        return html
          .replace("<!-- VITE_SPLASH_PLACEHOLDER -->", splashHtml)
          .replace(
            "<!-- VITE_SPLASH_CSS_PLACEHOLDER -->",
            `<style>\n${splashCss}\n</style>`,
          )
          .replace(
            "<!-- VITE_SPLASH_JS_PLACEHOLDER -->",
            `<script>\n${splashJs}\n</script>`,
          );
      },
    },
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
    // Under Bazel, test/source files are exposed to the process via a
    // runfiles symlink tree. Vite's dev-server-style module graph (used by
    // vitest) follows symlinks to their real path by default, which resolves
    // outside that tree and isn't visible to the sandbox running the test.
    // Scoped to `vitest` (which sets process.env.VITEST) since enabling this
    // for `vite build` breaks Rolldown's resolution of pnpm's nested,
    // symlinked transitive dependencies.
    preserveSymlinks: !!process.env.VITEST,
  },
  test: {
    environment: "jsdom",
  },
});
