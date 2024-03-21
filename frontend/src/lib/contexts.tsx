import { createContext, ReactNode, useEffect, useMemo, useState } from 'react';
import { List, ListSummary, UserToken } from './definitions';

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
  const memoValue = useMemo(() => ({ userToken, setUserToken }), [userToken, setUserToken]);
  return <UserTokenContext.Provider value={memoValue}>{children}</UserTokenContext.Provider>;
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
  const memoValue = useMemo(() => ({ listSummary, setListSummary }), [listSummary, setListSummary]);
  return <ListSummaryContext.Provider value={memoValue}>{children}</ListSummaryContext.Provider>;
};

/**
 * Context for list currently being processed
 */
export type ListContextType = {
  list: List | null;
  setList: (lists: List | null) => void;
};

export const ListContext = createContext<ListContextType | null>(null);

export const ListProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [list, setList] = useState<List | null>(null);
  const memoValue = useMemo(() => ({ list, setList }), [list, setList]);
  useEffect(() => {
    const timer = setTimeout(() => {
      setList(null);
    }, 3600000); // Clear list after an hour

    return () => clearTimeout(timer);
  }, []);

  return <ListContext.Provider value={memoValue}>{children}</ListContext.Provider>;
};
