import Cookies from 'js-cookie';
import { AccessTokenResponse, EntryWithImage, ListSummary, UserToken } from '../lib/definitions';

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
   * - Next time you try to sign in, need to check if an authcode cookie exists
   *     - If so, directly run this function using the existing authcode
   *     - If not, redirect to sign in page (that then redirects with authCode in url query)
   *
   *
   */
  let cacheMode: RequestCache = refresh ? 'reload' : 'default';

  // Fetch access token from backend
  const response = await fetch(`${BACKEND_URL}AuthUser?authCode=${encodeURIComponent(authCode)}`, {
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

async function SortList(accessToken: string, listId: string, refresh: boolean = false): Promise<EntryWithImage[]> {
  console.log('received request');
  // This list will fetch the list of images in the list, run them through the sorting process, then return the list of images as well as an array which shows which order they should be in - maybe a map?

  let cacheMode: RequestCache = refresh ? 'reload' : 'default';

  const response = await fetch(`${BACKEND_URL}SortList?accessToken=${encodeURIComponent(accessToken)}&listId=${listId}`, {
    cache: 'reload', // Update this after testing
    credentials: 'include',
  });
  if (!response.ok) {
    const errorText = `Error code: ${response.status}; message: ${response.statusText}`;
    console.error(errorText);
    throw new Error(errorText);
  }

  const data = await response.json();

  const entryListWithImages: EntryWithImage[] = data['items'] as any as EntryWithImage[];

  return entryListWithImages;
}

function WriteSortedList(accessToken: string, listId: string, sortMap: object) {
  // This will send the processed sortMap to the backend. The sortMap tells the backend how we would like to re-sort the list. That should then be written to the user's letterboxd list.
}

export { GetAccessTokenAndUser, GetLists, SortList, WriteSortedList };
