/// <reference types="vitest/config" />

import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'
import react, { reactCompilerPreset } from '@vitejs/plugin-react'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
import babel from '@rolldown/plugin-babel'
import path from 'node:path';
import fs from 'node:fs';

export default defineConfig({

  plugins: [
    devtools(),
    tanstackRouter({
      target: 'react',
      autoCodeSplitting: true,
    }),
    react(),
    babel({ presets: [reactCompilerPreset({ target: '19' })] }),
    {
      name: 'inject-splash',
      transformIndexHtml(html) {
        const splashHtml = fs.readFileSync(path.resolve(__dirname, 'src/splash/splash.html'), 'utf-8');
        const splashCss = fs.readFileSync(path.resolve(__dirname, 'src/splash/splash.css'), 'utf-8');
        const splashJs = fs.readFileSync(path.resolve(__dirname, 'src/splash/splash.js'), 'utf-8');
        return html
          .replace('<!-- VITE_SPLASH_PLACEHOLDER -->', splashHtml)
          .replace('<!-- VITE_SPLASH_CSS_PLACEHOLDER -->', `<style>\n${splashCss}\n</style>`)
          .replace('<!-- VITE_SPLASH_JS_PLACEHOLDER -->', `<script>\n${splashJs}\n</script>`);
      }
    }
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  test: {
    environment: 'jsdom',
  },
})
