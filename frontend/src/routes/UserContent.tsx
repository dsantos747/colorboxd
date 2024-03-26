import { useContext, useState } from 'react';
import ListMenu from '../ui/listMenu';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import ListPreview from '../ui/listPreview';
import Error from '../errorDiv';
import NoList from '../ui/noList';

const UserContent = () => {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list } = useContext(ListContext) as ListContextType;
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  return (
    <div className='flex flex-col justify-center pt-20 pb-4 w-full'>
      {userToken && (
        <div className='flex flex-col md:flex-row md:justify-between items-center mx-8 md:mx-16 gap-6 2xl:mx-32'>
          {error && <Error message={error.toString()} />}
          {!error && (
            <>
              <div className='flex-grow-0'>{<ListMenu setError={setError} loading={loading} setLoading={setLoading} />}</div>
              <div className='grow'>{list ? <ListPreview setError={setError} /> : <NoList loading={loading} />}</div>
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default UserContent;
