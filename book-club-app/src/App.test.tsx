import { render } from '@testing-library/react';

import App from './App';
import { AuthProvider } from './context/AuthContext';

test('renders App without crashing', () => {
  expect(() =>
    render(
      <AuthProvider>
        <App />
      </AuthProvider>,
    ),
  ).not.toThrow();
});
