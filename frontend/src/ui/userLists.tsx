import { useCallback, useContext, useEffect, useState } from 'react';
import {
  ListContext,
  ListContextType,
  ListSummaryContext,
  ListSummaryContextType,
  UserTokenContext,
  UserTokenContextType,
} from '../lib/contexts';
import { ClearListCache, GetLists, SortList } from '../actions/actions';
import { ArrowPathIcon } from '@heroicons/react/16/solid';

const minListLength = 20;

const listTooShortMessages = [
  "Don't waste my time...",
  'You call that a list?',
  'You can sort that yourself!',
  `Give me at least ${minListLength} films.`,
  `Give me at least ${minListLength} films.`,
] as const;

type listTypes = (typeof listTooShortMessages)[number];

function UserLists() {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { listSummary, setListSummary } = useContext(ListSummaryContext) as ListSummaryContextType;
  const { list, setList } = useContext(ListContext) as ListContextType;
  const [chosenListIndex, setChosenListIndex] = useState<number>();
  const [listLengthMessage, setListLengthMessage] = useState<listTypes | null>();
  const [menuOpen, setMenuOpen] = useState<boolean>(true);
  const [formSubmitting, setFormSubmitting] = useState<boolean>(false);

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
  }, [userToken, setListSummary]);

  useEffect(() => {
    const fadeTimer = setTimeout(() => {
      setListLengthMessage(null);
    }, 3000);

    return () => clearTimeout(fadeTimer);
  }, [listLengthMessage]);

  const handleFormSubmit = useCallback(
    (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();

      if (chosenListIndex !== undefined && userToken && listSummary) {
        if (listSummary[chosenListIndex].filmCount < minListLength) {
          setListLengthMessage(listTooShortMessages[Math.floor(Math.random() * listTooShortMessages.length)]);
          return;
        }
        setFormSubmitting(true);
        SortList(userToken.Token, listSummary[chosenListIndex])
          .then((lwi) => {
            setList(lwi);
            setMenuOpen(false);
          })
          .catch((error) => {
            console.error('Error getting sorted list:', error);
          })
          .finally(() => {
            setFormSubmitting(false);
          });
      }
    },
    [chosenListIndex, userToken, listSummary, setList]
  );

  const handleRefreshLists = useCallback(() => {
    ClearListCache();
    setList(null);

    GetLists(userToken.Token, userToken.UserId, true)
      .then((ls) => {
        setListSummary(ls);
        setChosenListIndex(0);
      })
      .catch((error) => {
        console.error('Error getting user lists:', error);
      });
  }, [setListSummary, userToken, setList]);

  const handleMenuState = useCallback(() => {
    setMenuOpen(!menuOpen);
  }, [menuOpen]);

  useEffect(() => {
    if (list === null) {
      setMenuOpen(true);
    }
  }, [list]);

  return (
    <div className='w-max bg-white bg-opacity-5 rounded-2xl py-6 px-8 outline-indigo-400 mx-auto md:mx-0 min-w-40 md:min-w-56'>
      <div className='flex justify-between'>
        <button type='button' onClick={handleMenuState}>
          <h2 className='text-lg'>{menuOpen ? 'Your lists:' : 'Show Lists'}</h2>
        </button>
        {menuOpen && (
          <button type='button' onClick={handleRefreshLists}>
            <ArrowPathIcon className='h-5 w-5 text-gray-500 mr-1 hover:text-teal-400 transition-colors duration-150' />
          </button>
        )}
      </div>

      <div className='bg-gradient-to-r from-blue-600 via-teal-500 to-lime-500 h-0.5 w-full' />
      {menuOpen && (
        <form
          onSubmit={(e) => {
            handleFormSubmit(e);
          }}>
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
          <button
            type='submit'
            disabled={formSubmitting}
            className='enabled:hover:text-teal-400 enabled:hover:translate-x-0.5 transition-all disabled:text-gray-500 font-semibold border-b-[1px] border-indigo-500'>
            {formSubmitting ? 'Please wait...' : "Let's Sort!"}
          </button>
          <p className='text-xs text-red-400 absolute'>{listLengthMessage}</p>
        </form>
      )}
    </div>
  );
}

export default UserLists;
