import Cookies from 'js-cookie';
import { ReactNode, useCallback, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';

const happyButtonStyle =
  'text-sm sm:text-base py-2 px-3 sm:px-5 h-min rounded-lg bg-gray-700 bg-opacity-50 text-teal-400 hover:bg-teal-800 hover:text-white active:bg-teal-700 active:text-white transition-colors duration-200';

const sadButtonStyle =
  'text-sm sm:text-base py-2 px-3 sm:px-5 h-min rounded-lg bg-gray-700 bg-opacity-50 text-orange-500 hover:bg-opacity-20 hover:text-gray-400 active:bg-opacity-10 active:text-gray-500 transition-all duration-200';

const LoginButton = () => {
  return (
    <Link to={'./user'} className={`${happyButtonStyle}`}>
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
    <button onClick={handleSignOut} className={`${sadButtonStyle}`}>
      Sign Out
    </button>
  );
};

type ButtonProps = {
  readonly handleClick: any;
  readonly type?: 'button' | 'submit' | 'reset';
  readonly children: ReactNode;
};

const HappyButton = ({ handleClick, type = 'button', children }: ButtonProps) => {
  return (
    <button onClick={handleClick} className={`${happyButtonStyle} `} type={type}>
      {children}
    </button>
  );
};

const SadButton = ({ handleClick, type = 'button', children }: ButtonProps) => {
  return (
    <button onClick={handleClick} className={`${sadButtonStyle} `} type={type}>
      {children}
    </button>
  );
};

export { LoginButton, SignOutButton, HappyButton, SadButton };
