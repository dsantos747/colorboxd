import { isRouteErrorResponse, useRouteError } from 'react-router-dom';
import Nav from './ui/nav';

const Error: React.FC<{}> = () => {
  const error = useRouteError();
  console.error(error);

  return (
    <>
      <Nav />
      <div id='error-page' className='min-h-screen bg-gray-900 flex flex-col items-center justify-center text-white'>
        <h1 className='text-3xl'>Oops!</h1>
        <p className='text-lg'>Sorry, an unexpected error has occurred.</p>
        <p>Message: {isRouteErrorResponse(error) ? error.statusText : 'Unknown error'}</p>
      </div>
    </>
  );
};

export default Error;
