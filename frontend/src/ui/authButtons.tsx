import { useAuth0 } from '@auth0/auth0-react';

const LoginButton: React.FC<{}> = () => {
  const { loginWithRedirect } = useAuth0();

  return (
    <button onClick={() => loginWithRedirect()} className='py-2 px-6 h-min rounded-sm bg-blue-800'>
      Log In
    </button>
  );
};

const SignOutButton: React.FC<{}> = () => {
  const { logout } = useAuth0();

  return (
    <button onClick={() => logout()} className='py-2 px-6 h-min rounded-sm bg-red-800'>
      Sign Out
    </button>
  );
};

export { LoginButton, SignOutButton };
