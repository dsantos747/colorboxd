import Cookies from 'js-cookie';
import { EntryWithImage, List, ListSummary, UserToken } from '../lib/definitions';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;
const BACKEND_URL2 = process.env.REACT_APP_BACKEND_URL2;

async function GetAccessTokenAndUser(authCode: string, refresh: boolean = false): Promise<UserToken> {
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
  let cacheMode: RequestCache = refresh ? 'reload' : 'default';

  // Fetch access token from backend
  const response = await fetch(`${BACKEND_URL}AuthUser?authCode=${encodeURIComponent(authCode)}`, {
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

async function GetLists(accessToken: string, userId: string, refresh: boolean = false): Promise<ListSummary[]> {
  let cacheMode: RequestCache = refresh ? 'reload' : 'default';

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

async function SortList(accessToken: string, listSummary: ListSummary, refresh: boolean = false): Promise<List> {
  let cacheMode: RequestCache = refresh ? 'reload' : 'default';

  const response = await fetch(`${BACKEND_URL}SortList?accessToken=${encodeURIComponent(accessToken)}&listId=${listSummary.id}`, {
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
  const entryListWithImages: EntryWithImage[] = data['items'] as any as EntryWithImage[];
  const list: List = {
    ...listSummary,
    entries: entryListWithImages,
  };

  return list;
}

async function WriteSortedList(accessToken: string, list: List, offset: number): Promise<string> {
  // This will send the processed sortMap to the backend. The sortMap tells the backend how we would like to re-sort the list. That should then be written to the user's letterboxd list.
  console.log('time to write the sorted list');

  const requestBody = {
    accessToken: accessToken,
    list: list,
    offset: offset,
  };

  console.log(requestBody);

  const response = await fetch(`${BACKEND_URL2}SortList`, {
    method: 'POST',
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
