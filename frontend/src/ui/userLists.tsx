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

const listTooShortMessages = [
  "Don't waste my time...",
  'You call that a list?',
  'I only sort lists with at least 10 films.',
  'You can sort that yourself!',
] as const;

type listTypes = (typeof listTooShortMessages)[number];

function UserLists() {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { listSummary, setListSummary } = useContext(ListSummaryContext) as ListSummaryContextType;
  const { setList } = useContext(ListContext) as ListContextType;
  const [chosenListIndex, setChosenListIndex] = useState<number>();
  const [listLengthMessage, setListLengthMessage] = useState<listTypes | null>();

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

  useEffect(() => {
    const fadeTimer = setTimeout(() => {
      setListLengthMessage(null);
    }, 3000);

    return () => clearTimeout(fadeTimer);
  }, [listLengthMessage]);

  function handleFormSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (chosenListIndex !== undefined && userToken && listSummary) {
      if (listSummary[chosenListIndex].filmCount < 10) {
        setListLengthMessage(listTooShortMessages[Math.floor(Math.random() * listTooShortMessages.length)]);
        return;
      }
      SortList(userToken.Token, listSummary[chosenListIndex])
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
    <div className='w-max'>
      <form
        onSubmit={(e) => {
          handleFormSubmit(e);
        }}>
        <h2 className='text-lg'>Your lists:</h2>
        <div className='bg-gradient-to-r from-blue-600 via-teal-500 to-lime-500 h-0.5 w-full'></div>
        <div className='w-max -mx-2 my-4'>
          {listSummary?.map((list, ind) => {
            return (
              <div
                key={list.id}
                className='has-[:checked]:bg-indigo-500 has-[:checked]:bg-opacity-20 transition-colors duration-100 w-full py-2 px-2 rounded-lg'>
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
                    {list.filmCount > 1 ? ' films' : ' film'})
                  </span>
                </label>
              </div>
            );
          })}
        </div>
        <button type='submit' className='hover:text-teal-400 hover:translate-x-0.5 transition-all'>
          Let&apos;s Sort!
        </button>
        <p className='text-xs text-red-400 absolute'>{listLengthMessage}</p>
      </form>
    </div>
  );
}

export default UserLists;
