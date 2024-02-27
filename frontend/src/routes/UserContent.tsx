import { useContext } from 'react';
import UserLists from '../ui/userLists';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';
import ListPreview from '../ui/listPreview';

const UserContent = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;

  return (
    <div className='flex flex-col justify-center py-20'>
      {userToken && (
        <div className='grid grid-cols-1 md:grid-cols-2 mx-16'>
          <div>
            <UserLists></UserLists>
          </div>
          <ListPreview></ListPreview>
        </div>
      )}
    </div>
  );
};

export default UserContent;
