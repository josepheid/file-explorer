import '@testing-library/jest-dom/vitest';

import { cleanup } from '@testing-library/react';
import { afterAll, afterEach, beforeAll, vi } from 'vitest';

import { server } from './src/test/server';

beforeAll(() => server.listen());

afterEach(() => {
  cleanup();

  vi.resetAllMocks();

  server.resetHandlers();
});

afterAll(() => server.close());
