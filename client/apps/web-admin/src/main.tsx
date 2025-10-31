// src/main.tsx

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import '@/styles/theme.scss';
import '@/styles/index.scss';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from '@/App';

const container = document.getElementById('root');
if (!container) {
  throw new Error('Root element not found');
}

const root = createRoot(container);

root.render(
  <StrictMode>
    <App />
  </StrictMode>
);