import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',

    setupFiles: ['./src/setupTests.ts'],
    globals: true,

    coverage: {
      provider: 'istanbul', // or 'v8'
      reporter: ['text', 'lcov'],
      reportsDirectory: './coverage',
    },
  },
});
