import React from 'react';
import { Link } from 'react-router-dom';

type NavItem = {
  text: string;
  path: string;
};

const navEntries: NavItem[] = [
  { text: 'Home', path: '/' },
  { text: 'Dashboard', path: '/user' },
];

const Nav: React.FC = () => {
  return (
    <nav className='fixed'>
      <ul className='flex mx-8 my-4 gap-3 text-white'>
        {navEntries.map((item) => {
          return (
            <li key={item.path}>
              <Link to={item.path}>{item.text}</Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
};

export default Nav;
