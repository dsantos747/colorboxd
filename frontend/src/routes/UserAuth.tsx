import { useContext, useEffect } from 'react';
import { GetAccessTokenAndUser } from '../actions/actions';
import { useLocation, useNavigate } from 'react-router-dom';
import UserContent from './UserContent';
import Cookies from 'js-cookie';
import { UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { UserToken } from '../lib/definitions';

function UserAuth() {
  const authorisationUrl = process.env.REACT_APP_LBOXD_AUTH_URL ?? 'https://colorboxd.com/';

  const { userToken, setUserToken } = useContext(UserTokenContext) as UserTokenContextType;

  const navigate = useNavigate();
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const authCode = queryParams.get('code');

  useEffect(() => {
    const handleTokenStatus = async () => {
      if (!userToken) {
        const cookieUserToken = Cookies.get('userToken');
        if (cookieUserToken) {
          const cookieToken: UserToken = JSON.parse(cookieUserToken);
          setUserToken(cookieToken);
        } else if (authCode) {
          try {
            const fetchUserToken = await GetAccessTokenAndUser(authCode);
            setUserToken(fetchUserToken);
          } catch (error) {
            console.error('Error getting access token:', error);
          }
        } else {
          window.location.href = authorisationUrl;
        }
      } else if (authCode) {
        // Remove authcode from params
        const { pathname, search } = location;
        const updatedQueryParams = new URLSearchParams(search);
        updatedQueryParams.delete('code');
        navigate(`${pathname}?${updatedQueryParams.toString()}`);
      }
    };
    handleTokenStatus();
  }, [authCode, userToken, setUserToken, authorisationUrl, location, navigate]);

  return userToken && <UserContent />;
}

export default UserAuth;
