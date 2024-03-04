import Cookies from 'js-cookie';
import { useCallback, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';

const LoginButton = () => {
  return (
    <Link to={'./user'} className='py-2 px-6 h-min rounded-sm bg-blue-800'>
      Get Started
    </Link>
  );
};

const SignOutButton = () => {
  const { setUserToken } = useContext(UserTokenContext) as UserTokenContextType;
  const navigate = useNavigate();

  const handleSignOut = useCallback(() => {
    navigate('/');
    Cookies.remove('userToken');
    setUserToken(null);
  }, [navigate, setUserToken]);

  return (
    <button
      onClick={handleSignOut}
      className='text-sm sm:text-base py-2 px-3 sm:px-5 h-min rounded-lg bg-gray-700 bg-opacity-50 text-orange-500 hover:bg-opacity-20 hover:text-gray-400 transition-colors duration-300'>
      Sign Out
    </button>
  );
};

export { LoginButton, SignOutButton };
