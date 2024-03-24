import { useContext, useState } from 'react';
import UserLists from '../ui/userLists';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import ListPreview from '../ui/listPreview';
import Error from '../errorDiv';

const UserContent = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list } = useContext(ListContext) as ListContextType;
  const [error, setError] = useState<string | null>(null);

  return (
    <div className='flex flex-col justify-center pt-20 pb-4 w-full'>
      {userToken && (
        <div className='flex flex-col md:flex-row md:justify-between items-center mx-8 md:mx-16 gap-6 2xl:mx-32'>
          {error && <Error message={error.toString()} />}
          {!error && (
            <>
              <div className='flex-grow-0'>{<UserLists setError={setError} />}</div>
              <div className='grow'>
                {list && <ListPreview setError={setError} />}
                {!list && <div className='text-xl text-gray-500 my-auto text-center'>Choose a list and let&apos;s sort!</div>}
              </div>
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default UserContent;
