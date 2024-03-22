import { useCallback, useContext, useEffect, useState } from 'react';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { WriteSortedList } from '../actions/actions';
import { Button } from './buttons';
import { ArrowUpIcon, ArrowDownIcon } from '@heroicons/react/16/solid';
import { SortModeType, sorts } from '../lib/definitions';

function calcIndex(i: number, startIndex: number, len: number, reverse: boolean) {
  const ind = reverse ? (len - i) % len : i;
  return (ind + startIndex) % len;
}

export default function ListPreview() {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list, setList } = useContext(ListContext) as ListContextType;
  const [currSort, setCurrSort] = useState<SortModeType>({
    sortMode: { id: 'hue', name: 'Hue' },
    visible: true,
    reverse: false,
  });
  const [startIndex, setStartIndex] = useState<number>(0);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [imgLoading, setImgLoading] = useState<boolean>(false);

  const handleSaveList = useCallback(() => {
    if (userToken && list) {
      setSubmitting(true);
      // setTimeout(() => {
      //   setSubmitting(false);
      // }, 2000);
      WriteSortedList(userToken.Token, list, startIndex, currSort.sortMode.id, currSort.reverse)
        .then((message) => {
          if (message[0].startsWith('List updated successfully')) {
            setStartIndex(0);
            setList(null);
          } else {
            console.error(message);
          }
        })
        .catch((error) => {
          console.error('Error writing sorted list to letterboxd account:', error);
        })
        .finally(() => {
          setSubmitting(false);
        });
    }
  }, [userToken, list, setList, startIndex, currSort]);

  const handleCancel = useCallback(() => {
    if (userToken && list) {
      setList(null);
    }
  }, [userToken, list, setList]);

  const handleShowOriginal = useCallback(() => {
    setCurrSort({ sortMode: currSort.sortMode, visible: !currSort.visible, reverse: false });
    if (currSort.visible) {
      list?.entries.sort((a, b) => {
        return Number(a.entryId) - Number(b.entryId);
      });
    } else {
      list?.entries.sort((a, b) => {
        return a.sorts[currSort.sortMode.id] - b.sorts[currSort.sortMode.id];
      });
    }
  }, [currSort, list]);

  const handleReverseOrder = useCallback(() => {
    setCurrSort({ sortMode: currSort.sortMode, visible: currSort.visible, reverse: !currSort.reverse });
  }, [currSort]);

  const handleSortChange = useCallback(
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      const selectedMode = sorts.find((s) => s.id === e.target.value);
      if (selectedMode) {
        setStartIndex(0);
        setCurrSort({ sortMode: selectedMode, visible: true, reverse: false });
        list?.entries.sort((a, b) => {
          return a.sorts[selectedMode.id] - b.sorts[selectedMode.id];
        });
      }
    },
    [list]
  );

  useEffect(() => {
    async function loadImages() {
      if (list?.entries) {
        await Promise.all(
          list?.entries.map(
            (entry) =>
              new Promise((resolve, reject) => {
                const img = new Image();
                img.src = entry.posterUrl;
                img.onload = resolve;
                img.onerror = reject;
              })
          )
        );
      }
      setImgLoading(true);
    }

    setStartIndex(0);
    setCurrSort({ sortMode: sorts[0], visible: true, reverse: false });
    setImgLoading(false);
    loadImages().catch((e) => {
      console.error('Error loading images:', e);
    });
  }, [list, setImgLoading]);

  return (
    <div className='mx-auto max-w-6xl'>
      <div className='flex flex-wrap justify-between'>
        <p className='my-auto'>Click a poster to make it the start of the list.</p>
        <form className='flex justify-end flex-wrap align-middle items-center select-none gap-2 ml-auto'>
          <div className='mx-auto'>
            <input type='checkbox' id='showOriginal' className='hidden peer' checked={!currSort.visible} onChange={handleShowOriginal} />
            <label
              htmlFor='showOriginal'
              className='peer-checked:decoration-teal-500 underline cursor-pointer transition-colors duration-200'>
              View {currSort.visible ? 'Original' : 'Sorted'}
            </label>
          </div>
          <select
            name='sortMethod'
            id='sortMethod'
            className='h-8 bg-gray-900 w-max'
            value={currSort.sortMode.id}
            onChange={handleSortChange}>
            {sorts.map((mode) => {
              return (
                <option key={mode.id} value={mode.id} className=''>
                  {mode.name}
                </option>
              );
            })}
          </select>
          <div className='mx-auto'>
            <input type='checkbox' id='reverseOrder' className='hidden peer' checked={currSort.reverse} onChange={handleReverseOrder} />
            <label htmlFor='reverseOrder' className='cursor-pointer'>
              {currSort.reverse ? <ArrowUpIcon className='h-5 w-5' /> : <ArrowDownIcon className='h-5 w-5' />}
            </label>
          </div>
        </form>
      </div>

      <div className='grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 overflow-y-auto scrollbar-hide md:h-[60vh] my-4 mx-auto'>
        {imgLoading &&
          list?.entries.map((l, i) => {
            const ind = calcIndex(i, startIndex, list.entries.length, currSort.reverse);
            return (
              <div key={l.entryId} className='m-1 text-center'>
                <button
                  type='button'
                  onClick={() => {
                    setStartIndex(ind);
                  }}>
                  <img src={list.entries[ind].posterUrl} alt={list.entries[ind].name} />
                </button>
              </div>
            );
          })}
        {!imgLoading &&
          Array.from({ length: list?.entries.length ?? 10 }).map((_, i) => {
            return (
              <div key={i} className='m-1 text-center max-h-56 max-w-36 aspect-[2/3] '>
                <div className=' w-full h-full  bg-gray-500 bg-opacity-20 animate-pulse'> </div>
              </div>
            );
          })}
      </div>

      <div className='bg-gradient-to-r from-blue-600 via-teal-500 to-lime-500 h-0.5 w-full' />
      <div className='w-max mx-auto mt-4 space-x-2'>
        <Button theme='sad' handleClick={handleCancel}>
          Cancel
        </Button>
        <Button theme='happy' handleClick={handleSaveList} disabled={submitting}>
          Save List
        </Button>
      </div>
    </div>
  );
}
