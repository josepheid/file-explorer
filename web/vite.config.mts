import react from '@vitejs/plugin-react-swc';
import { defineConfig } from 'vitest/config';

const config = defineConfig({
  clearScreen: false,
  server: {
    fs: {
      allow: ['.'],
    },
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'https://localhost:8080',
        secure: false,
        changeOrigin: true,
      },
    },
  },
  test: {
    include: ['src/**/*.test.ts(x)'],
    environment: 'jsdom',
    setupFiles: ['vitest.setup.ts'],
  },
  plugins: [
    react({
      plugins: [
        [
          '@swc/plugin-styled-components',
          {
            displayName: true,
          },
        ],
      ],
    }),
  ],
});

export { config as default };
