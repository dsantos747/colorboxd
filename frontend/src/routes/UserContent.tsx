import { useContext } from 'react';
import UserLists from '../ui/userLists';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import ListPreview from '../ui/listPreview';

const UserContent = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list } = useContext(ListContext) as ListContextType;

  return (
    <div className='flex flex-col justify-center pt-20 pb-4 w-full'>
      {userToken && (
        <div className='flex flex-col md:flex-row md:justify-between items-center mx-8 md:mx-16 gap-6'>
          <div className='flex-grow-0'>{<UserLists />}</div>
          <div className='grow'>
            {list ? <ListPreview /> : <div className='text-xl text-gray-500 my-auto text-center'>Choose a list and let&apos;s sort!</div>}
          </div>
        </div>
      )}
    </div>
  );
};

export default UserContent;
