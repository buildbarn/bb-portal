/// <reference types="vitest/config" />

import fs from "node:fs";
import path from "node:path";
import babel from "@rolldown/plugin-babel";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const DEV_ENV = JSON.stringify({
  featureFlags: {
    home: { fileUpload: {}, besUploadInstructions: {} },
    bes: {
      pageBuilds: {},
      pageInvocations: {},
      pageTargets: {},
      pageTests: {},
      pageTrends: {},
    },
    browser: {},
    scheduler: {},
  },
  grpcBackendUrl: "grpc://localhost:8082",
  companyName: "UrbanCompass",
  footerContent: [],
  additionalBuildColumns: [],
  additionalBuildInvocationColumns: [],
  prometheusUrl: "http://localhost:9090",
});

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
    {
      // In dev mode the Go backend doesn't serve index.html, so window.__env__
      // is never injected. Replace the placeholder with a local dev config.
      name: "inject-dev-env",
      apply: "serve",
      transformIndexHtml(html) {
        return html.replace(
          "<!-- BB_PORTAL_CONFIGURATION_PLACEHOLDER -->",
          `<script>window.__env__ = ${DEV_ENV}</script>`,
        );
      },
    },
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/api/v1/prometheus": {
        target: "http://localhost:9090",
        rewrite: (path) => path.replace(/^\/api\/v1\/prometheus/, ""),
      },
      "/api": "http://localhost:4180",
      "/graphql": "http://localhost:4180",
    },
  },
  test: {
    environment: "jsdom",
  },
});
