import Cookies from 'js-cookie';
import { AccessTokenResponse, ListSummary, UserToken } from '../lib/definitions';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;
const BACKEND_URL2 = process.env.REACT_APP_BACKEND_URL2;

// DEPRECATED
async function GetAccessTokenAndLists(authCode: string, setUserLists: (lists: ListSummary[]) => void) {
  // Check if accessToken cookie already exists
  if (Cookies.get('accessToken') !== undefined) {
    return Cookies.get('accessToken');
  }

  // Fetch access token from backend
  const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;
  const response = await fetch(`${BACKEND_URL}AuthUserGetLists?authCode=${encodeURIComponent(authCode)}`);
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }
  const data = await response.json();

  // Set accessToken cookie, and return value
  Cookies.set('accessToken', data.token.access_token, { expires: data.token.expires_in / 86400 }); // expiresIn converted from seconds to days
  // Cookies.set('userLists', JSON.stringify(data.lists));
  setUserLists(data.lists);

  return data;
}

// DEPRECATED
async function GetAccessToken(authCode: string): Promise<string> {
  const cookieToken = Cookies.get('accessToken');
  if (cookieToken !== undefined) {
    return cookieToken;
  }

  // Fetch access token from backend
  const response = await fetch(`${BACKEND_URL}AuthUser?authCode=${encodeURIComponent(authCode)}`);
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }
  const data: AccessTokenResponse = await response.json();
  Cookies.set('accessToken', data.access_token, { expires: data.expires_in / 86400 }); // expiresIn converted from seconds to days
  return data.access_token;
}

async function GetAccessTokenAndUser(authCode: string): Promise<UserToken> {
  const cookieUserToken = Cookies.get('userToken');
  if (cookieUserToken !== undefined) {
    const cookieData: UserToken = JSON.parse(cookieUserToken);
    return cookieData;
  }

  // Fetch access token from backend
  const response = await fetch(`${BACKEND_URL}AuthUser?authCode=${encodeURIComponent(authCode)}`);
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }

  const data: UserToken = await response.json();
  Cookies.set('userToken', JSON.stringify(data), { expires: data.TokenExpiresIn / 86400 }); // expiresIn converted from seconds to days
  return data;
}

async function GetLists(accessToken: string, userId: string, refresh: boolean = false): Promise<ListSummary[]> {
  let cacheMode: RequestCache = 'default';
  if (refresh) {
    // Ideally, use no-store - but not sure how to handle that server side
    cacheMode = 'reload';
  }

  const response = await fetch(`${BACKEND_URL2}GetLists?accessToken=${encodeURIComponent(accessToken)}&userId=${userId}`, {
    cache: cacheMode,
    credentials: 'include',
  });
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }

  const data: ListSummary[] = await response.json();
  return data;
}

// DEPRECATED
function HasAccessToken(): boolean {
  const accessToken = Cookies.get('userToken');
  return !!accessToken;
}

export { HasAccessToken, GetAccessTokenAndLists, GetAccessToken, GetAccessTokenAndUser, GetLists };
