import { useCallback, useState } from 'react';

export default function ColorScroll() {
  const [imgLoadState, setImgLoadState] = useState(false);

  const handleImageLoad = useCallback(() => setImgLoadState(true), []);

  return (
    <div className='aspect-square h-64 sm:h-80 lg:h-96 overflow-hidden relative'>
      <img
        src='colorboxd_scroll_min.png'
        alt='A list sorted with Colorboxd'
        className={`object-cover motion-safe:animate-scroll ${!imgLoadState && 'hidden'}`}
        onLoad={handleImageLoad}
      />
      <div className='absolute inset-0 bg-gradient-to-t from-gray-900 to-transparent to-5%' />
    </div>
  );
}
