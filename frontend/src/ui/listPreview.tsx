import { useContext, useState } from 'react';
import { ListContext, ListContextType } from '../lib/contexts';

type Props = {};

export default function ListPreview({}: Props) {
  const { list } = useContext(ListContext) as ListContextType;
  const [startIndex, setStartIndex] = useState<number>(0);

  return (
    <div>
      <div className='grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 overflow-y-auto scrollbar-hide'>
        {list?.map((l, i) => {
          const ind = (i + startIndex) % list.length; // Use this to determine starting image
          return (
            <div key={l.entryId} className='m-2'>
              <img
                src={list[ind].posterUrl}
                onClick={() => {
                  setStartIndex(ind);
                }}
                className='cursor-pointer'></img>
            </div>
          );
        })}
      </div>
    </div>
  );
}
