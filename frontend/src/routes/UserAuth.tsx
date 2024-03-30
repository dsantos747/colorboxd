import { useContext, useEffect } from 'react';
import { GetAccessTokenAndUser } from '../actions/actions';
import { useLocation, useNavigate } from 'react-router-dom';
import UserContent from './UserContent';
import Cookies from 'js-cookie';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';

function UserAuth() {
  const homeUrl = process.env.REACT_APP_BASE_URL ?? 'https://colorboxd.com/';

  const { userToken, setUserToken } = useContext(UserTokenContext) as UserTokenContextType;

  const navigate = useNavigate();
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const authCode = queryParams.get('code');

  useEffect(() => {
    const handleTokenStatus = async () => {
      if (!userToken) {
        await handleUserToken();
      } else if (authCode) {
        handleAuthCode();
      }
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

    const handleAuthCode = () => {
      const { pathname, search } = location;
      const updatedQueryParams = new URLSearchParams(search);
      updatedQueryParams.delete('code');
      navigate(`${pathname}?${updatedQueryParams.toString()}`);
    };

    handleTokenStatus().catch((e) => {
      console.error('Error handling user authorisation:', e);
    });
  }, [authCode, userToken, setUserToken, homeUrl, location, navigate]);

  return userToken && <UserContent />;
}

export default UserAuth;
