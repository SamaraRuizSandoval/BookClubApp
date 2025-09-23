module.exports = {
    root: true,
    parser: '@typescript-eslint/parser',
    parserOptions: { ecmaVersion: 'latest', sourceType: 'module' },
    settings: { react: { version: 'detect' } },
    env: { browser: true, es2021: true, node: true },
    plugins: ['@typescript-eslint', 'react', 'react-hooks', 'import', 'jsx-a11y'],
    extends: [
        'eslint:recommended',
        'plugin:@typescript-eslint/recommended',
        'plugin:react/recommended',
        'plugin:react-hooks/recommended',
        'plugin:import/recommended',
        'plugin:jsx-a11y/recommended',
        'prettier'
    ],
    overrides: [
        {
            files: ['**/*.{test,spec}.{ts,tsx}'],
            env: { 'vitest-globals/env': true } // optional if you install eslint-plugin-vitest
        }
    ],
    rules: {
        'react/react-in-jsx-scope': 'off' // React 17+ JSX transform
    }
};
