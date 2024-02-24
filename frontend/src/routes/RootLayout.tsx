import { Outlet } from 'react-router-dom';
import Nav from '../ui/nav';

const RootLayout = () => {
  return (
    <>
      <Nav />
      <div className='min-h-screen bg-gray-900 flex flex-col items-center justify-between text-white'>
        <div className='flex grow w-full justify-center'>
          <Outlet />
        </div>
        <div className='w-full'>
          <div className=' text-right mx-8 mb-4 mt-2 text-slate-400 text-xs'>
            Created by Daniel Santos |{' '}
            <a href='https://github.com/dsantos747' className='decoration-none underline'>
              Github
            </a>{' '}
            |{' '}
            <a href='https://danielsantosdev.vercel.app/' className='decoration-none underline'>
              Website
            </a>
          </div>
        </div>
      </div>
    </>
  );
};

export default RootLayout;
