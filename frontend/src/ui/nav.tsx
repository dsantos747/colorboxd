import { useContext } from 'react';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { LoginButton, SignOutButton } from './authButtons';
import ColorboxdLogo from './colorboxdLogo';
import { Link } from 'react-router-dom';

const Nav = () => {
  // Theres a glitch where the sign out button stays visible when navigating to home,
  // but then when you refresh it's back to log in. To make things easier, it's probably
  // easier to just have a "Get started" button (navigate to users page) and a "sign out"
  // button (clear auth token) - with the sign out button only visible on the user page

  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;

  return (
    <nav className='fixed w-screen bg-gradient-to-b from-gray-900 from-70% to-transparent z-10'>
      <div className='flex justify-between items-center mx-8 h-14 my-4'>
        <Link to={'/'} className='text-xl'>
          <h1>
            <ColorboxdLogo />
          </h1>
        </Link>
        {userToken && (
          <div className='text-white'>
            <SignOutButton />
          </div>
        )}
        {!userToken && (
          <div className='text-white'>
            <LoginButton />
          </div>
        )}
      </div>
    </nav>
  );
};

export default Nav;
