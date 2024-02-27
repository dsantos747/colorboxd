import Cookies from 'js-cookie';
import { useContext } from 'react';
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

  return (
    <button
      onClick={() => {
        navigate('/');
        Cookies.remove('userToken');
        setUserToken(null);
      }}
      className='py-2 px-6 h-min rounded-sm bg-red-800'>
      Sign Out
    </button>
  );
};

export { LoginButton, SignOutButton };
