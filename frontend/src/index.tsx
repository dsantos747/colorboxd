import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import './main.css';
import Error from './error';
import ComingSoon from './routes/ComingSoon';
import Dashboard from './routes/Dashboard';
import RootLayout from './routes/RootLayout';
import { Auth0Provider } from '@auth0/auth0-react';
import Home from './routes/Home';

const auth0Domain = process.env.AUTH0_DOMAIN;
const auth0Client = process.env.AUTH0_CLIENT_ID;

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    errorElement: <Error />,
    children: [
      { path: '/', element: <ComingSoon />, errorElement: <Error /> },
      { path: '/user', element: <Dashboard />, errorElement: <Error /> },
      { path: '/home', element: <Home />, errorElement: <Error /> },
    ],
  },
]);

// if (!auth0Domain || !auth0Client) {
//   console.error('Error loading environment variables - Application cannot be loaded.');
// } else {
const root = ReactDOM.createRoot(document.getElementById('root')!);
root.render(
  <React.StrictMode>
    <Auth0Provider
      domain={auth0Domain ?? ''}
      clientId={auth0Client ?? ''}
      authorizationParams={{
        redirect_uri: '/user',
      }}>
      <RouterProvider router={router} />
    </Auth0Provider>
  </React.StrictMode>
);
// }

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals();
