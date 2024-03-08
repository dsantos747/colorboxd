import { Outlet } from 'react-router-dom';
import Nav from '../ui/nav';
import Footer from '../ui/footer';

const RootLayout = () => {
  return (
    <>
      <Nav />
      <div className='min-h-screen bg-gray-900 flex flex-col items-center justify-between text-white'>
        <div className='flex grow w-full justify-center'>
          <Outlet />
        </div>
        <Footer />
      </div>
    </>
  );
};

export default RootLayout;
