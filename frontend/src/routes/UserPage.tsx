import { useContext, useEffect, useState } from 'react';
import { GetAccessTokenAndUser } from '../actions/actions';
import { useLocation, useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import ListMenu from '../ui/listMenu';
import ListPreview from '../ui/listPreview';
import Error from '../errorDiv';
import NoList from '../ui/noList';

function UserPage() {
  const { userToken, setUserToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list } = useContext(ListContext) as ListContextType;
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [loadingIndex, setLoadingIndex] = useState<number>(0);

  const navigate = useNavigate();
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const authCode = queryParams.get('code');
  const homeUrl = process.env.REACT_APP_BASE_URL ?? 'https://colorboxd.com/';

  useEffect(() => {
    const handleAuthCode = () => {
      const { pathname, search } = location;
      const updatedQueryParams = new URLSearchParams(search);
      updatedQueryParams.delete('code');
      navigate(`${pathname}?${updatedQueryParams.toString()}`);
    };

    const handleUserToken = async () => {
      const cookieUserToken = Cookies.get('userToken');
      if (cookieUserToken) {
        setUserToken(JSON.parse(cookieUserToken));
      } else if (authCode) {
        try {
          const fetchUserToken = await GetAccessTokenAndUser(authCode);
          setUserToken(fetchUserToken);
          handleAuthCode();
        } catch (error) {
          console.error('Error getting access token:', error);
        }
      } else {
        window.location.href = homeUrl;
      }
    };

    const handleTokenStatus = async () => {
      if (!userToken) {
        await handleUserToken();
      } else if (authCode) {
        handleAuthCode();
      }
    };

    handleTokenStatus().catch((e) => {
      console.error('Error handling user authorisation:', e);
    });
  }, [authCode, userToken, setUserToken, homeUrl, location, navigate]);

  useEffect(() => {
    setLoadingIndex(0);
  }, [list, loading]);

  return (
    userToken && (
      <div className='flex flex-col justify-center pt-20 pb-4 w-full'>
        {userToken && (
          <div className='flex flex-col md:flex-row md:justify-between items-center mx-8 md:mx-16 gap-6 2xl:mx-32'>
            {error && <Error message={error.toString()} />}
            {!error && (
              <>
                <div className='flex-grow-0'>{<ListMenu setError={setError} loading={loading} setLoading={setLoading} />}</div>
                <div className='grow'>
                  {list ? (
                    <ListPreview setError={setError} />
                  ) : (
                    <NoList loading={loading} loadingIndex={loadingIndex} setLoadingIndex={setLoadingIndex} />
                  )}
                </div>
              </>
            )}
          </div>
        )}
      </div>
    )
  );
}

export default UserPage;
