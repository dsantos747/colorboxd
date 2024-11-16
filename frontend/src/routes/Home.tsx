import ColorScroll from '../ui/colorScroll';
import ColorboxdLogo from '../ui/colorboxdLogo';

function colorHighlight(text: string) {
  return (
    <span className='before:block before:absolute before:-inset-1 before:-skew-y-3 before:-skew-x-3 before:bg-gradient-to-r before:from-sky-500 before:via-teal-500 before:to-green-500 relative inline-block mb-4 ml-1'>
      <span className='relative text-black font-semibold italic'>{text}</span>
    </span>
  );
}

const Home = () => {
  return (
    <div className='pt-20 sm:pt-0 text-center flex flex-col sm:flex-row gap-2 md:justify-between md:grow w-full md:px-20 max-w-screen-lg items-center justify-center text-white'>
      <div className='space-y-10 '>
        <div className='hidden sm:block space-y-10'>
          <h1 className='text-5xl font-bold tracking-widest'>
            <ColorboxdLogo />
          </h1>
          <p className=''>Letterboxd lists, but {colorHighlight('prettier.')}</p>
        </div>
        <div className='hidden sm:block bg-blue-50 bg-opacity-20 rounded-2xl py-6 px-6'>
          <p className='text-sm'><strong>EDIT 11/2024:</strong> Colorboxd is under maintenance, as it is in dire need of some performance and stability improvements. I will work hard to get it back up and running as soon as I can!<br></br>
            If you&apos;d like to be notified when it&apos;s ready, you can either fill out <a className='underline font-semibold decoration-2' target="_blank" href="https://forms.gle/hqiqCyknKMLMNasU9">this form</a>, or periodically check this <a className='underline font-semibold decoration-2' target="_blank" href="https://www.reddit.com/r/Letterboxd/comments/1cd2mqx/i_created_colorboxd_a_website_that_lets_you_sort/">reddit post</a>.</p>
        </div>
      </div>

      <h1 className='sm:hidden text-xl mb-2'>Letterboxd lists, but {colorHighlight('prettier.')}</h1>
      <div>
        <ColorScroll />
      </div>
      <div className='sm:hidden bg-blue-50 bg-opacity-20 rounded-2xl py-4 px-4 scale-75 '>
        <p className='text-sm'><strong>EDIT 11/2024:</strong> Colorboxd is under maintenance, as it is in dire need of some performance and stability improvements. I will work hard to get it back up and running as soon as I can!<br></br>
          If you&apos;d like to be notified when it&apos;s ready, you can either fill out <a className='underline font-semibold decoration-2' target="_blank" href="https://forms.gle/hqiqCyknKMLMNasU9">this form</a>, or periodically check this <a className='underline font-semibold decoration-2' target="_blank" href="https://www.reddit.com/r/Letterboxd/comments/1cd2mqx/i_created_colorboxd_a_website_that_lets_you_sort/">reddit post</a>.</p>
      </div>
    </div>
  );
};

export default Home;
