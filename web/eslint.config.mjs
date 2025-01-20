import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import reactPlugin from 'eslint-plugin-react';
import vitest from '@vitest/eslint-plugin';
import reactHooksPlugin from 'eslint-plugin-react-hooks';
import testingLibraryPlugin from 'eslint-plugin-testing-library';
import globals from 'globals';

export default tseslint.config(
  {
    ignores: [
      'eslint.config.mjs',
      'vite.config.mts',
      'vitest.setup.ts',
      'dist/**/*',
      '.prettierrc.js',
    ],
  },
  eslint.configs.recommended,
  tseslint.configs.recommendedTypeChecked,
  reactPlugin.configs.flat.recommended,
  {
    languageOptions: {
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
        project: './web/tsconfig.json',
      },
    },
  },
  {
    settings: {
      react: {
        version: 'detect',
      },
    },
    languageOptions: {
      ecmaVersion: 6,
      sourceType: 'module',
      parser: tseslint.parser,
    },
    plugins: {
      // There is no flat config available.
      'react-hooks': reactHooksPlugin,
    },
    rules: {
      ...reactHooksPlugin.configs.recommended.rules,
      'import/no-unresolved': 'off',
      'no-case-declarations': 'off',
      'prefer-const': 'off',
      'no-var': 'off',
      'prefer-rest-params': 'off',
      'no-console': 'warn',
      'no-trailing-spaces': 'error',
      'react/jsx-no-undef': 'error',
      'react/jsx-pascal-case': 'error',
      'react/no-danger': 'error',
      'react/jsx-no-duplicate-props': 'error',
      'react/jsx-sort-prop-types': 'off',
      'react/jsx-sort-props': 'off',
      'react/jsx-uses-vars': 'warn',
      'react/no-did-mount-set-state': 'warn',
      'react/no-did-update-set-state': 'warn',
      'react/no-unknown-property': 'warn',
      'react/prop-types': 'off',
      'react/jsx-wrap-multilines': 'warn',
      'react/jsx-no-useless-fragment': ['error', { allowExpressions: true }],
      'react/display-name': 'off',
      'react/no-children-prop': 'warn',
      'react/no-unescaped-entities': 'warn',
      'react/jsx-key': 'warn',
      'react/jsx-no-target-blank': 'warn',
      // Turned off because we use automatic runtime.
      'react/jsx-uses-react': 'off',
      'react/react-in-jsx-scope': 'off',
    },
  },
  {
    files: ['**/*.test.{ts,tsx}'],
    plugins: {
      vitest,
      'testing-library': testingLibraryPlugin,
    },
    rules: {
      ...vitest.configs.recommended.rules,
      ...testingLibraryPlugin.configs['flat/react'].rules,

      'vitest/expect-expect': 'off',

      '@typescript-eslint/no-non-null-assertion': 'off',
    },
    settings: {
      vitest: {
        typecheck: true,
      },
    },
    languageOptions: {
      globals: {
        ...vitest.environments.env.globals,
      },
    },
  }
);
