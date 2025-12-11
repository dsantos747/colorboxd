import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';
import './main.css';
import RootLayout from './routes/RootLayout';
import Home from './routes/Home';
import { UserTokenProvider, ListSummaryProvider, ListProvider } from './lib/contexts';
import RouterError from './routerError';
// import UserPage from './routes/UserPage';

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    errorElement: <RouterError />,
    children: [
      { path: '/', element: <Home />, errorElement: <RouterError /> },
      {
        path: '/user',
        element: <Navigate to="/" replace={true} />,
        errorElement: <RouterError />,
      },],
  },
]);

const root = ReactDOM.createRoot(document.getElementById('root')!); // skipcq: JS-0339

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
