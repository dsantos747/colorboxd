import Cookies from 'js-cookie';
import { useContext } from 'react';
import { Link, redirect } from 'react-router-dom';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';

const LoginButton = () => {
  return (
    <Link to={'./user'} className='py-2 px-6 h-min rounded-sm bg-blue-800'>
      Get Started
    </Link>
  );
};

const SignOutButton = () => {
  // Broken
  // Need to be able to redirect the user back to the home page

  const { setUserToken } = useContext(UserTokenContext) as UserTokenContextType;

  return (
    <button
      onClick={() => {
        Cookies.remove('userToken');
        setUserToken(null);
        return redirect('/');
      }}
      className='py-2 px-6 h-min rounded-sm bg-red-800'>
      Sign Out
    </button>
  );
};

export { LoginButton, SignOutButton };
