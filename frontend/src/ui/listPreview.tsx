import { useContext, useState } from 'react';
import { ListContext, ListContextType, UserTokenContext, UserTokenContextType } from '../lib/contexts';
import { WriteSortedList } from '../actions/actions';

type Props = {};

const sorts = [
  { id: 'hue', name: 'Hue-Based Sort' },
  { id: 'step', name: 'Alternating Step Sort' },
  { id: 'hilbert', name: 'Hilbert Sort' },
  { id: 'cie2000', name: 'CIELAB2000 Sort' },
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
          setList(null);
        })
        .catch((error) => {
          console.error('Error writing sorted list to letterboxd account:', error);
        });
    }
  }

  function handleCancel() {
    if (userToken && list) {
      setList(null);
    }
  }

  return (
    <div>
      <div className=''>
        {/* <h3 className='text-lg'>Preview Results</h3> */}

        <div className='flex justify-between flex-wrap align-middle items-center'>
          {/* <form className='flex gap-2 h-8'> */}
          <p className=''>Hint: Click an item to make it the start of the list.</p>
          <select name='sortMethod' id='sortMethod' className='h-8 bg-gray-900 w-max ml-auto'>
            {sorts.map((mode) => {
              return (
                <option key={mode.id} value={mode.id} className=''>
                  {mode.name}
                </option>
              );
            })}
          </select>
        </div>

        {/* </form> */}
      </div>
      {/* Need to determine a better method of defining the height of the frame */}
      <div className='grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 overflow-y-auto scrollbar-hide h-[700px] my-4'>
        {list?.entries.map((l, i) => {
          const ind = (i + startIndex) % list.entries.length; // Use this to determine starting image
          return (
            <div key={l.entryId} className='m-1'>
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
      <div className='bg-gradient-to-r from-blue-600 via-teal-500 to-lime-500 h-0.5 w-full'></div>

      <div className='w-max mx-auto mt-4 space-x-2'>
        <button
          onClick={() => {
            handleCancel();
          }}
          className='px-3 py-2 rounded-sm mx-auto bg-gray-700 bg-opacity-50 text-orange-500 hover:bg-opacity-20 hover:text-gray-400 transition-colors duration-300'
          type='button'>
          Cancel
        </button>
        <button
          onClick={() => {
            handleSaveList();
          }}
          className='px-3 py-2 rounded-sm mx-auto bg-gray-700 bg-opacity-50 text-teal-400 hover:bg-teal-800 hover:text-white transition-colors duration-300'
          type='button'>
          Save List
        </button>
      </div>
    </div>
  );
}
