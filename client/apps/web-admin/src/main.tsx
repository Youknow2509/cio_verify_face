// src/main.tsx

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import '@/styles/index.scss';
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