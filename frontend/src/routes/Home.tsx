import ColorScroll from '../ui/colorScroll';
import ColorboxdLogo from '../ui/colorboxdLogo';

const ComingSoon = () => {
  return (
    <div className='text-center flex flex-col md:flex-row md:justify-between md:grow w-full md:px-20 max-w-screen-lg items-center justify-center text-white'>
      <div className='space-y-10 hidden sm:block'>
        <h1 className='text-5xl font-bold tracking-widest'>
          <ColorboxdLogo />
        </h1>
        <p className=''>Letterboxd lists, but even prettier.</p>
      </div>
      <h1 className='sm:hidden text-xl'>
        Letterboxd lists, but{'  '}
        <span className='before:block before:absolute before:-inset-1 before:-skew-y-3 before:-skew-x-3 before:bg-gradient-to-r before:from-sky-500 before:via-teal-500 before:to-green-500 relative inline-block mb-4 ml-1'>
          <span className='relative text-black font-semibold italic'>prettier.</span>
        </span>
      </h1>
      <div>
        <ColorScroll></ColorScroll>
      </div>
    </div>
  );
};

export default ComingSoon;
