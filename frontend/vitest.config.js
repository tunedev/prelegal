import { svelte } from '@sveltejs/vite-plugin-svelte';
import { svelteTesting } from '@testing-library/svelte/vite';
import { defineConfig } from 'vitest/config';
import { fileURLToPath } from 'url';

export default defineConfig({
  plugins: [svelte(), svelteTesting()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.js'],
    include: ['src/**/*.test.{js,ts}'],
  },
  resolve: {
    alias: {
      $lib: fileURLToPath(new URL('./src/lib', import.meta.url)),
      '$app/navigation': fileURLToPath(new URL('./src/test/mocks/app-navigation.js', import.meta.url)),
      '$app/state': fileURLToPath(new URL('./src/test/mocks/app-state.svelte.js', import.meta.url)),
    },
  },
});
