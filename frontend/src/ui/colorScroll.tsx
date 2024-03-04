export default function ColorScroll() {
  return (
    <div className='h-80 w-80 overflow-hidden relative'>
      <img src='colorboxd_scroll.png' alt='A list sorted with Colorboxd' className='object-cover motion-safe:animate-scroll' />
      <div className='absolute inset-0 bg-gradient-to-t from-gray-900 to-transparent to-5%'></div>
    </div>
  );
}
