import Cookies from 'js-cookie';
import { EntryWithImage, List, ListSummary, UserToken } from '../lib/definitions';

const BACKEND_URL1 = process.env.REACT_APP_BACKEND_URL1;
const BACKEND_URL2 = process.env.REACT_APP_BACKEND_URL2;

async function GetAccessTokenAndUser(authCode: string, refresh = false): Promise<UserToken> {
  const cookieUserToken = Cookies.get('userToken');
  if (cookieUserToken !== undefined) {
    const cookieData: UserToken = JSON.parse(cookieUserToken);
    return cookieData;
  }

  /**
   *
   *
   * IMPORTANT: Caching isn't really working, because a new authCode is being generated each time. This
   * invalidates any method of caching tbh. Might be necessary to:
   * - Upon first attempt, create authCode cookie. Need expiry time, check docs?
   *             =Apparently, most use 30-60 second expiration times
   *                  - In this case, is there any use in creating a cookie?
   * - Next time you try to sign in, need to check if an authcode cookie exists
   *     - If so, directly run this function using the existing authcode
   *     - If not, redirect to sign in page (that then redirects with authCode in url query)
   *
   *
   */
  const cacheMode: RequestCache = refresh ? 'reload' : 'default';

  // Fetch access token from backend
  const response = await fetch(`${BACKEND_URL1}AuthUser?authCode=${encodeURIComponent(authCode)}`, {
    method: 'GET',
    cache: cacheMode,
    credentials: 'include',
  });
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }

  const data: UserToken = await response.json();
  Cookies.set('userToken', JSON.stringify(data), { expires: data.TokenExpiresIn / 86400 }); // expiresIn converted from seconds to days
  return data;
}

async function GetLists(accessToken: string, userId: string, refresh = false): Promise<ListSummary[]> {
  const cacheMode: RequestCache = refresh ? 'reload' : 'default';

  const response = await fetch(`${BACKEND_URL2}GetLists?accessToken=${encodeURIComponent(accessToken)}&userId=${userId}`, {
    method: 'GET',
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

async function SortList(accessToken: string, listSummary: ListSummary, refresh = false): Promise<List> {
  const cacheMode: RequestCache = refresh ? 'reload' : 'default';

  const response = await fetch(`${BACKEND_URL1}SortList?accessToken=${encodeURIComponent(accessToken)}&listId=${listSummary.id}`, {
    method: 'GET',
    cache: cacheMode,
    credentials: 'include',
  });
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }

  const data = await response.json();
  const entryListWithImages: EntryWithImage[] = data.items as EntryWithImage[];
  const list: List = {
    ...listSummary,
    entries: entryListWithImages,
  };

  return list;
}

async function WriteSortedList(accessToken: string, list: List, offset: number, refresh = false): Promise<string> {
  const cacheMode: RequestCache = refresh ? 'reload' : 'default';

  const requestBody = { accessToken, list, offset };

  const response = await fetch(`${BACKEND_URL2}WriteList`, {
    method: 'POST',
    cache: cacheMode,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(requestBody),
  });
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }

  const message = await response.json();
  return message;
}

export { GetAccessTokenAndUser, GetLists, SortList, WriteSortedList };
