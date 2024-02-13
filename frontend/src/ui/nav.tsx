import { useAuth0 } from '@auth0/auth0-react';
import React from 'react';
import { Link } from 'react-router-dom';
import { LoginButton, SignOutButton } from './authButtons';
import ColorboxdLogo from './colorboxdLogo';

type NavItem = {
  text: string;
  path: string;
};

const navEntries: NavItem[] = [
  { text: 'Home', path: '/' },
  { text: 'Dashboard', path: '/user' },
];

const Nav: React.FC<{}> = () => {
  const { isAuthenticated } = useAuth0();

  return (
    <nav className='fixed w-screen'>
      <div className='flex justify-between items-center mx-8 h-14'>
        <div className='text-xl'>
          <ColorboxdLogo />
        </div>
        {/* {isAuthenticated && (
          <ul className='flex mx-8 my-4 gap-3 text-white'>
            {navEntries.map((item) => {
              return (
                <li key={item.path}>
                  <Link to={item.path}>{item.text}</Link>
                </li>
              );
            })}
          </ul>
        )} */}
        {isAuthenticated && (
          <div className='text-white'>
            <SignOutButton />
          </div>
        )}
        {!isAuthenticated && (
          <div className='text-white'>
            <LoginButton />
          </div>
        )}
      </div>
    </nav>
  );
};

export default Nav;
