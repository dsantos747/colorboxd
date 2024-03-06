import { useCallback, useContext, useState } from 'react';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { WriteSortedList } from '../actions/actions';
import { HappyButton, SadButton } from './buttons';

const sorts = [
  { id: 'hue', name: 'Hue-Based Sort' },
  { id: 'step', name: 'Alternating Step Sort' },
  { id: 'hilbert', name: 'Hilbert Sort' },
  { id: 'cie2000', name: 'CIELAB2000 Sort' },
] as const;

type SortTypes = (typeof sorts)[number];

type SortModeType = {
  sortMode: SortTypes;
  visible: boolean;
};

export default function ListPreview() {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list, setList } = useContext(ListContext) as ListContextType;
  const [currSort, setCurrSort] = useState<SortModeType>({ sortMode: { id: 'hue', name: 'Hue-Based Sort' }, visible: true });
  const [startIndex, setStartIndex] = useState<number>(0);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const handleSaveList = useCallback(() => {
    if (userToken && list) {
      setSubmitting(true);
      setTimeout(() => {
        setSubmitting(false);
      }, 2000);
      WriteSortedList(userToken.Token, list, startIndex)
        .then((message) => {
          // console.log(message);
          setStartIndex(0);
          setList(null);
        })
        .catch((error) => {
          console.error('Error writing sorted list to letterboxd account:', error);
        });
      // setSubmitting(false);
    }
  }, [userToken, list, setList, startIndex]);

  const handleCancel = useCallback(() => {
    if (userToken && list) {
      setList(null);
    }
  }, [userToken, list, setList]);

  const handleShowOriginal = useCallback(() => {
    setCurrSort({ sortMode: currSort.sortMode, visible: !currSort.visible });
  }, [currSort]);

  const handleSortChange = useCallback((e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedMode = sorts.find((s) => s.id === e.target.value);
    if (selectedMode) {
      setCurrSort({ sortMode: selectedMode, visible: true });
    }
  }, []);

  return (
    <div className='mx-auto max-w-6xl'>
      <div className='flex flex-wrap justify-between'>
        <p className='my-auto'>Hint: Click an item to make it the start of the list.</p>
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
        </form>
      </div>

      {/* Need to determine a better method of defining the height of the frame */}
      <div className='grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 overflow-y-auto scrollbar-hide md:h-[60vh] my-4 mx-auto'>
        {list?.entries.map((l, i) => {
          const ind = (i + startIndex) % list.entries.length; // Use this to determine starting image
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
      </div>

      <div className='bg-gradient-to-r from-blue-600 via-teal-500 to-lime-500 h-0.5 w-full' />
      <div className='w-max mx-auto mt-4 space-x-2'>
        <SadButton handleClick={handleCancel}>Cancel</SadButton>
        <HappyButton handleClick={handleSaveList} disabled={submitting}>
          Save List
        </HappyButton>
      </div>
    </div>
  );
}
