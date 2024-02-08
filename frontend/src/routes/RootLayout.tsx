import { Outlet, Link } from 'react-router-dom';
import Nav from '../ui/nav';

const RootLayout = () => {
  return (
    <>
      <Nav />
      <div className='min-h-screen bg-gray-900 flex flex-col items-center justify-center'>
        <Outlet />
      </div>
    </>
  );
};

export default RootLayout;
