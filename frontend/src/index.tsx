import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import './main.css';
import Error from './error';
import ComingSoon from './routes/ComingSoon';
import RootLayout from './routes/RootLayout';
import Home from './routes/Home';
import UserAuth from './routes/UserAuth';
import { ListsProvider, UserTokenProvider } from './lib/contexts';

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    errorElement: <Error />,
    children: [
      { path: '/', element: <ComingSoon />, errorElement: <Error /> },
      { path: '/user', element: <UserAuth />, errorElement: <Error /> },
      { path: '/home', element: <Home />, errorElement: <Error /> },
    ],
  },
]);

const root = ReactDOM.createRoot(document.getElementById('root')!);

root.render(
  <React.StrictMode>
    <UserTokenProvider>
      <ListsProvider>
        <RouterProvider router={router} />
      </ListsProvider>
    </UserTokenProvider>
  </React.StrictMode>
);
