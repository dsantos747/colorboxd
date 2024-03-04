import { useContext } from 'react';
import UserLists from '../ui/userLists';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import ListPreview from '../ui/listPreview';

const UserContent = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list } = useContext(ListContext) as ListContextType;

  return (
    <div className='flex flex-col justify-center py-20 w-full'>
      {userToken && (
        <div className='flex flex-col md:flex-row md:justify-between mx-8 md:mx-16 gap-6'>
          <div className='flex-grow-0'>{<UserLists />}</div>
          <div className='grow'>{list && <ListPreview />}</div>
        </div>
      )}
    </div>
  );
};

export default UserContent;
