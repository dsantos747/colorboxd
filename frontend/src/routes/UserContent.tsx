import { useContext } from 'react';
import UserLists from '../ui/userLists';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';

const UserContent = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;

  return (
    <div className='text-center flex flex-col justify-center'>
      <p>Your lists:</p>
      {userToken && <UserLists></UserLists>}
    </div>
  );
};

export default UserContent;
