import { createContext, ReactNode, useState } from 'react';
import { ListSummary, UserToken } from './definitions';

/**
 * Context for currently authed User and their access token
 */
export type UserTokenContextType = {
  userToken: UserToken | null;
  setUserToken: (userToken: UserToken | null) => void;
};
export const UserTokenContext = createContext<UserTokenContextType | null>(null);
export const UserTokenProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [userToken, setUserToken] = useState<UserToken | null>(null);

  return <UserTokenContext.Provider value={{ userToken, setUserToken }}>{children}</UserTokenContext.Provider>;
};

/**
 * Context for currently authed user's lists
 */
const dummyList: ListSummary = { description: 'test description', filmCount: 1, id: 'test id', name: 'test name', version: 1 };
export type ListsContextType = {
  lists: ListSummary[];
  setLists: (lists: ListSummary[]) => void;
};
export const ListsContext = createContext<ListsContextType | null>(null);
export const ListsProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [lists, setLists] = useState<ListSummary[]>([dummyList]);
  return <ListsContext.Provider value={{ lists, setLists }}>{children}</ListsContext.Provider>;
};
