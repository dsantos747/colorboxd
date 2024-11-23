import { useContext } from 'react';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { LoginButton, SignOutButton } from './buttons';
import ColorboxdLogo from './colorboxdLogo';
import { Link, useLocation } from 'react-router-dom';

const Nav = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const location = useLocation();
  const isOnUserPage = location.pathname.toLowerCase().startsWith('/user');

  return (
    <nav className='fixed w-screen bg-gradient-to-b from-gray-900 from-70% to-transparent z-10'>
      <div className='flex justify-between items-center mx-8 h-14 my-4'>
        <Link to={'/'} className='text-xl'>
          <h1>
            <ColorboxdLogo />
          </h1>
        </Link>
        {isOnUserPage && userToken ? <SignOutButton /> : <LoginButton />}
      </div>
    </nav>
  );
};

export default Nav;
