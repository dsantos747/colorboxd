import Cookies from 'js-cookie';
import { ReactNode, useCallback, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';

const happyButtonStyle =
  'text-sm sm:text-base py-2 px-3 sm:px-5 h-min rounded-lg bg-gray-700 bg-opacity-50 text-teal-400 enabled:hover:bg-teal-800 enabled:hover:text-white enabled:active:bg-teal-700 enabled:active:text-white disabled:bg-opacity-20 disabled:text-gray-400 transition-colors duration-200';

const sadButtonStyle =
  'text-sm sm:text-base py-2 px-3 sm:px-5 h-min rounded-lg bg-gray-700 bg-opacity-50 text-orange-500 enabled:hover:bg-opacity-20 enabled:hover:text-gray-400 enabled:active:bg-opacity-10 enabled:active:text-gray-500 disabled:bg-opacity-20 disabled:text-gray-400 transition-all duration-200';

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
  readonly disabled?: boolean;
  readonly children: ReactNode;
};

const HappyButton = ({ handleClick, type = 'button', disabled = false, children }: ButtonProps) => {
  return (
    <button onClick={handleClick} className={`${happyButtonStyle} `} type={type} disabled={disabled}>
      {children}
    </button>
  );
};

const SadButton = ({ handleClick, type = 'button', disabled = false, children }: ButtonProps) => {
  return (
    <button onClick={handleClick} className={`${sadButtonStyle} `} type={type} disabled={disabled}>
      {children}
    </button>
  );
};

export { LoginButton, SignOutButton, HappyButton, SadButton };
