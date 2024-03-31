import Cookies from 'js-cookie';
import { EntryWithImage, List, ListSummary, SortModeType, UserToken } from '../lib/definitions';

const BACKEND_URL1 = process.env.REACT_APP_BACKEND_URL1;
const BACKEND_URL2 = process.env.REACT_APP_BACKEND_URL2;

async function GetAccessTokenAndUser(authCode: string, refresh = false): Promise<UserToken> {
  const cookieUserToken = Cookies.get('userToken');
  if (cookieUserToken !== undefined) {
    const cookieData: UserToken = JSON.parse(cookieUserToken);
    return cookieData;
  }

  const cacheMode: RequestCache = refresh ? 'reload' : 'default';

  // Fetch access token from backend
  const response = await fetch(`${BACKEND_URL1}AuthUser?authCode=${encodeURIComponent(authCode)}`, {
    method: 'GET',
    cache: cacheMode,
    credentials: 'include',
  });
  if (!response.ok) {
    let errorText;
    try {
      errorText = await response.text();
    } catch (error) {
      errorText = `Error code: ${response.status}; Message: ${response.statusText}`;
    }
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
    let errorText;
    try {
      errorText = await response.text();
    } catch (error) {
      errorText = `Error code: ${response.status}; Message: ${response.statusText}`;
    }
    throw new Error(errorText);
  }

  const data: ListSummary[] = await response.json();
  return data;
}

let listCache: Record<string, List> = {};

function ClearListCache() {
  listCache = {};
}

async function SortList(accessToken: string, listSummary: ListSummary, refresh = false): Promise<List> {
  const cacheMode: RequestCache = refresh || !listCache[listSummary.id] ? 'reload' : 'default';

  // const cacheMode: RequestCache = 'reload';
  // listSummary.id = 'tqjLE'; // Has >1500 entries!!!

  const response = await fetch(`${BACKEND_URL1}SortList?accessToken=${encodeURIComponent(accessToken)}&listId=${listSummary.id}`, {
    method: 'GET',
    cache: cacheMode,
    credentials: 'include',
  });
  if (!response.ok) {
    let errorText;
    try {
      errorText = await response.text();
    } catch (error) {
      errorText = `Error code: ${response.status}; Message: ${response.statusText}`;
    }
    throw new Error(errorText);
  }

  const data = await response.json();
  const entryListWithImages: EntryWithImage[] = data.items as EntryWithImage[];
  const list: List = {
    ...listSummary,
    entries: entryListWithImages,
  };

  listCache[listSummary.id] = list;
  return list;
}

async function WriteSortedList(
  accessToken: string,
  list: List,
  offset: number,
  sortMethod: SortModeType['sortMode']['id'],
  reverse: boolean,
  refresh = false
): Promise<string[]> {
  const cacheMode: RequestCache = refresh ? 'reload' : 'default';

  const requestBody = { accessToken, list, offset, sortMethod, reverse };

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
    let errorText;
    try {
      errorText = await response.text();
    } catch (error) {
      errorText = `Error code: ${response.status}; Message: ${response.statusText}`;
    }
    throw new Error(errorText);
  }

  const message: string[] = await response.json();
  return message;
}

export { GetAccessTokenAndUser, GetLists, SortList, WriteSortedList, ClearListCache };
