import { useContext, useState } from 'react';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { WriteSortedList } from '../actions/actions';

type Props = {};

const sorts = [
  { id: 'hue', name: 'Hue-Based Sort' },
  { id: 'step', name: 'Alternating Step Sort' },
  { id: 'hilbert', name: 'Hilbert' },
  { id: 'cie2000', name: 'CIELAB2000' },
] as const;

type SortTypes = (typeof sorts)[number];

export default function ListPreview({}: Props) {
  const { userToken } = useContext(UserTokenContext) as UserTokenContextType;
  const { list, setList } = useContext(ListContext) as ListContextType;
  const [sortMode, setSortMode] = useState<SortTypes>({ id: 'hue', name: 'Hue-Based Sort' });
  const [startIndex, setStartIndex] = useState<number>(0);

  function handleSaveList() {
    if (userToken && list) {
      WriteSortedList(userToken.Token, list, startIndex)
        .then((message) => {
          console.log(message);
          setStartIndex(0);
          // setList(null);
        })
        .catch((error) => {
          console.error('Error writing sorted list to letterboxd account:', error);
        });
    }
  }

  return (
    <div>
      <div className='flex justify-between flex-wrap'>
        <div>
          <h3 className='text-lg'>List Sorting results</h3>
          <p className='text-sm'>Hint: Click an item to make it the first in the list.</p>
        </div>
        <form className='flex gap-2 h-8'>
          <select name='sortMethod' id='sortMethod' className='h-8 bg-gray-900 w-max'>
            {sorts.map((mode) => {
              return (
                <option key={mode.id} value={mode.id}>
                  {mode.name}
                </option>
              );
            })}
          </select>
          <button
            onClick={() => {
              handleSaveList();
            }}
            className='bg-blue-700 px-3 rounded-sm'
            type='button'>
            Save List
          </button>
        </form>
      </div>
      {/* Need to determine a better method of defining the height of the frame */}
      <div className='grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 overflow-y-auto scrollbar-hide h-[700px]'>
        {list?.entries.map((l, i) => {
          const ind = (i + startIndex) % list.entries.length; // Use this to determine starting image
          return (
            <div key={l.entryId} className='m-2'>
              <button
                type='button'
                onClick={() => {
                  setStartIndex(ind);
                }}>
                <img src={list.entries[ind].posterUrl} alt={list.entries[ind].name}></img>
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}
