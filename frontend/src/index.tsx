import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { API_BASE_PATH } from './values';
import axios from 'axios';


axios.defaults.baseURL = API_BASE_PATH;

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
    <App />
);
