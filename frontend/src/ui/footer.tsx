export default function Footer() {
  return (
    <div className='w-full px-8 mb-4 mt-2 text-slate-400 text-xs flex justify-between flex-wrap gap-2'>
      <div className='text-right'>
        Created by Daniel Santos |{' '}
        <a href='https://github.com/dsantos747' className='decoration-none underline'>
          Github
        </a>{' '}
        |{' '}
        <a href='https://danielsantosdev.vercel.app/' className='decoration-none underline'>
          Website
        </a>
      </div>
      <div className='md:block'>
        Please consider{' '}
        <a href='https://ko-fi.com/danielsantosdev' className='underline'>
          donating
        </a>{' '}
        to help pay my server costs.
      </div>
    </div>
  );
}
