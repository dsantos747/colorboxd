import { useContext, useEffect, useState } from 'react';
import {
  ListContext,
  ListContextType,
  ListSummaryContext,
  ListSummaryContextType,
  UserTokenContext,
  UserTokenContextType,
} from '../lib/contexts';
import { GetLists, SortList } from '../actions/actions';

function UserLists() {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { listSummary, setListSummary } = useContext(ListSummaryContext) as ListSummaryContextType;
  const { setList } = useContext(ListContext) as ListContextType;
  const [chosenListIndex, setChosenListIndex] = useState<number>();

  if (!userToken) {
    throw new Error('Cannot render user lists - no authenticated user.');
  }

  useEffect(() => {
    GetLists(userToken.Token, userToken.UserId)
      .then((ls) => {
        setListSummary(ls);
        setChosenListIndex(0);
      })
      .catch((error) => {
        console.error('Error getting user lists:', error);
      });
  }, [userToken]);

  function handleFormSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (chosenListIndex !== undefined && userToken && listSummary) {
      SortList(userToken.Token, listSummary[chosenListIndex].id)
        .then((lwi) => {
          setList(lwi);
        })
        .catch((error) => {
          console.error('Error getting sorted list:', error);
        });
    }
  }

  // A nice touch would be to have a limit which only allows the links to work if list.filmCount is greater than e.g. 5.
  // If it's less than 5, give a message like "don't waste my time"
  //
  return (
    <form
      onSubmit={(e) => {
        handleFormSubmit(e);
      }}>
      <h2 className='text-lg'>Your lists:</h2>

      {listSummary?.map((list, ind) => {
        return (
          <div key={list.id} className='has-[:checked]:bg-indigo-500 has-[:checked]:bg-opacity-20 w-full py-2 px-2 rounded-lg '>
            <input
              type='radio'
              id={list.id}
              value={list.id}
              name='userList'
              checked={chosenListIndex === ind}
              onChange={() => setChosenListIndex(ind)}
              className='hidden peer'
            />
            <label htmlFor={list.id} className='cursor-pointer'>
              {list.name}{' '}
              <span className='text-xs text-gray-400'>
                ({list.filmCount}
                {list.filmCount > 1 ? ' entries' : ' entry'})
              </span>
            </label>
          </div>
        );
      })}
      <button type='submit'>Submit</button>
    </form>
  );
}

export default UserLists;
