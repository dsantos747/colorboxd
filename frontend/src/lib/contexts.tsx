import { createContext, ReactNode, useState } from 'react';
import { EntryWithImage, ListSummary, UserToken } from './definitions';

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
export type ListSummaryContextType = {
  listSummary: ListSummary[] | null;
  setListSummary: (lists: ListSummary[] | null) => void;
};
export const ListSummaryContext = createContext<ListSummaryContextType | null>(null);
export const ListSummaryProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [listSummary, setListSummary] = useState<ListSummary[] | null>(null);
  return <ListSummaryContext.Provider value={{ listSummary, setListSummary }}>{children}</ListSummaryContext.Provider>;
};

/**
 * Context for list currently being processed
 */
export type ListContextType = {
  list: EntryWithImage[] | null;
  setList: (lists: EntryWithImage[] | null) => void;
};
export const ListContext = createContext<ListContextType | null>(null);
export const ListProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [list, setList] = useState<EntryWithImage[] | null>(null);
  return <ListContext.Provider value={{ list, setList }}>{children}</ListContext.Provider>;
};
