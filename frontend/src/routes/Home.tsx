import ColorboxdLogo from '../ui/colorboxdLogo';

const ComingSoon = () => {
  return (
    <div className='text-center flex flex-col md:flex-row md:justify-between md:grow w-full md:px-20 max-w-screen-lg items-center justify-center text-white'>
      <div className='space-y-10'>
        <h1 className='text-5xl font-bold tracking-widest hidden sm:block'>
          <ColorboxdLogo />
        </h1>
        <p className=''>Letterboxd lists, but even prettier.</p>
      </div>
      <div>
        <img src='logo512_clear.png' className='h-[40vmin] pointer-events-none mx-auto' alt='logo' />
        <p>animation of coloured rectangles being sorted here</p>
        <p>alternatively, animated scroll of a big sorted list</p>
      </div>
    </div>
  );
};

export default ComingSoon;
