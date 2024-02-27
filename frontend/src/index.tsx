import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import './main.css';
import Error from './error';
import RootLayout from './routes/RootLayout';
import Home from './routes/Home';
import UserAuth from './routes/UserAuth';
import { UserTokenProvider, ListSummaryProvider, ListProvider } from './lib/contexts';

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    errorElement: <Error />,
    children: [
      { path: '/', element: <Home />, errorElement: <Error /> },
      { path: '/user', element: <UserAuth />, errorElement: <Error /> },
    ],
  },
]);

const root = ReactDOM.createRoot(document.getElementById('root')!);

root.render(
  <React.StrictMode>
    <UserTokenProvider>
      <ListSummaryProvider>
        <ListProvider>
          <RouterProvider router={router} />
        </ListProvider>
      </ListSummaryProvider>
    </UserTokenProvider>
  </React.StrictMode>
);
