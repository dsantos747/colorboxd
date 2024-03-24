type Props = {
  message?: string;
};

const titles = ['Jinkies!', 'Oops!', 'This is embarassing...', 'Can you forgive me?'];

const ErrorDiv = ({ message }: Props) => {
  return (
    <div id='error' className='w-full h-full bg-gray-900 flex flex-col items-center justify-center text-gray-500 text-center'>
      <h1 className='text-3xl text-white'>{titles[Math.floor(Math.random() * titles.length)]}</h1>
      <p className='text-lg mb-4'>Sorry, an unexpected error has occurred.</p>
      <p>{message}</p>
    </div>

    /**
     * Could be good to have a "start again" button here, which signs the user out then automatically signs them back in
     */
  );
};

export default ErrorDiv;
