export default {
  '**/*.{js,jsx,ts,tsx,css,md,json}': ['prettier --write'],
  '**/*.{js,jsx,ts,tsx}': ['eslint --fix'],
};
