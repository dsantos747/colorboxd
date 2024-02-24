import { useContext, useEffect } from 'react';
import { ListsContext, ListsContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { GetLists } from '../actions/actions';

function UserLists() {
  const { lists, setLists } = useContext(ListsContext) as ListsContextType;
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;

  if (!userToken) {
    throw new Error('Cannot render user lists - no authenticated user.');
  }

  useEffect(() => {
    GetLists(userToken.Token, userToken.UserId)
      .then((ls) => {
        setLists(ls);
      })
      .catch((error) => {
        console.error('Error getting user lists:', error);
      });
  }, [userToken, setLists]);

  return (
    <ul>
      {lists.map((list) => {
        return <li key={list.id}>{list.name}</li>;
      })}
    </ul>
  );
}

export default UserLists;
