import { useContext } from 'react';
import UserLists from '../ui/userLists';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';

const UserContent = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;

  return (
    <div className='text-center flex flex-col justify-center'>
      {/* <header className='bg-gray-900 min-h-screen flex flex-col items-center justify-center text-white text-2xl'> */}
      <img src='logo512_clear.png' className='h-[40vmin] pointer-events-none' alt='logo' />
      <p>This will be the user dashboard</p>
      <br></br>
      <p>Please check back soon!</p>
      {userToken && <UserLists></UserLists>}
      {/* </header> */}
    </div>
  );
};

export default UserContent;
