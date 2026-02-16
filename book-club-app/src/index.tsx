import { setupIonicReact } from '@ionic/react';
import React from 'react';
import ReactDOM from 'react-dom/client';

import App from './App';

setupIonicReact();

import '@ionic/react/css/core.css';
import '@ionic/react/css/normalize.css';
import '@ionic/react/css/structure.css';
import '@ionic/react/css/typography.css';
import '@ionic/react/css/ionic.bundle.css';
import { AuthProvider } from './context/AuthContext';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement,
);
root.render(
  <React.StrictMode>
    <AuthProvider>
      <App />
    </AuthProvider>
  </React.StrictMode>,
);
