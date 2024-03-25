import { isRouteErrorResponse, useRouteError } from 'react-router-dom';
import Nav from './ui/nav';
import ErrorDiv from './errorDiv';

const RouterError = () => {
  const error = useRouteError();
  const errorText = isRouteErrorResponse(error) ? error.statusText : 'Unknown error';

  return (
    <>
      <Nav />
      <div className='min-h-screen flex flex-col justify-center items-center'>
        <ErrorDiv message={`Message: ${errorText}`} />
      </div>
    </>
  );
};

export default RouterError;
