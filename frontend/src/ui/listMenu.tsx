import { Dispatch, SetStateAction, useCallback, useContext, useEffect, useState } from 'react';
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

type Props = {
  readonly setError: Dispatch<SetStateAction<string | null>>;
  readonly loading: boolean;
  readonly setLoading: Dispatch<SetStateAction<boolean>>;
};

function ListMenu({ setError, loading, setLoading }: Props) {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { listSummary, setListSummary } = useContext(ListSummaryContext) as ListSummaryContextType;
  const { list, setList } = useContext(ListContext) as ListContextType;
  const [chosenListIndex, setChosenListIndex] = useState<number>();
  const [listLengthMessage, setListLengthMessage] = useState<listTypes | null>();
  const [menuOpen, setMenuOpen] = useState<boolean>(true);

  if (!userToken) {
    throw new Error('Cannot render user lists - no authenticated user.');
  }

  // Get user lists on mount
  useEffect(() => {
    GetLists(userToken.Token, userToken.UserId)
      .then((ls) => {
        setListSummary(ls);
        setChosenListIndex(0);
      })
      .catch((error: string) => {
        setError(error);
      });
  }, [userToken, setListSummary, setError]);

  // "list too short" message timer
  useEffect(() => {
    const fadeTimer = setTimeout(() => {
      setListLengthMessage(null);
    }, 3000);
    return () => clearTimeout(fadeTimer);
  }, [listLengthMessage]);

  const handleFormSubmit = useCallback(
    (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      const isMobile = window.matchMedia('(max-width: 768px)');

      if (chosenListIndex !== undefined && userToken && listSummary) {
        if (listSummary[chosenListIndex].filmCount < minListLength) {
          setListLengthMessage(listTooShortMessages[Math.floor(Math.random() * listTooShortMessages.length)]);
          return;
        }
        setLoading(true);
        SortList(userToken.Token, listSummary[chosenListIndex])
          .then((lwi) => {
            setList(lwi);
            if (isMobile.matches) {
              setMenuOpen(false);
            }
          })
          .catch((error) => {
            setError(error);
          })
          .finally(() => {
            setLoading(false);
          });
      }
    },
    [chosenListIndex, userToken, listSummary, setList, setError, setLoading]
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
        setError(error);
      });
  }, [setListSummary, userToken, setList, setError]);

  const handleMenuState = useCallback(() => {
    setMenuOpen(!menuOpen);
  }, [menuOpen]);

  useEffect(() => {
    if (list === null) {
      setMenuOpen(true);
    }
  }, [list]);

  return (
    <div className='w-max bg-white bg-opacity-5 rounded-2xl py-6 px-8 outline-indigo-400 mx-auto md:mx-0 min-w-48 md:min-w-56'>
      <div className='flex justify-between'>
        <button type='button' onClick={handleMenuState}>
          <h2 className='text-lg'>{menuOpen ? 'Your lists:' : 'Show Lists'}</h2>
        </button>
        {menuOpen && (
          <button type='button' onClick={handleRefreshLists} title='Refresh Lists'>
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
            disabled={loading}
            className='enabled:hover:text-teal-400 enabled:hover:translate-x-0.5 transition-all disabled:text-gray-500 font-semibold border-b-[1px] border-indigo-500'>
            {loading ? 'Please wait...' : "Let's Sort!"}
          </button>
          <p className='text-xs text-red-400 absolute'>{listLengthMessage}</p>
        </form>
      )}
    </div>
  );
}

export default ListMenu;
