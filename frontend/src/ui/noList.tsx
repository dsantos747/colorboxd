import { useEffect, Dispatch, SetStateAction } from 'react';
import { CoffeeButton } from './buttons';

type Props = {
  readonly loading: boolean;
  readonly loadingIndex: number;
  readonly setLoadingIndex: Dispatch<SetStateAction<number>>;
};

const loadingMessages: string[] = [
  'Hang tight! Working my magic...',
  'Ooh, this is exciting...',
  "You've got some great films in here!",
  'Wow, some very bold picks indeed...',
  "I'd say your taste is definitely above the average I've dealt with.",
  'Seriously, you should see what the last guy had in their list...',
  "Enough chatter from me. I'm almost done!",
];

const NoList = ({ loading, loadingIndex, setLoadingIndex }: Props) => {
  useEffect(() => {
    const loadingInterval = setInterval(() => {
      setLoadingIndex((prevIndex) => (prevIndex + 1) % loadingMessages.length);
    }, 10000);

    return () => clearInterval(loadingInterval);
  }, [setLoadingIndex]);

  return (
    <div className='text-center text-gray-500'>
      {loading && (
        <div className='space-y-2 md:space-y-4'>
          <div className='text-sm my-auto'>Loading - Please Wait</div>
          <div className='text-xl my-auto'>{loadingMessages[loadingIndex]}</div>
        </div>
      )}
      {!loading && (
        <div className='space-y-10'>
          <div className='text-xl my-auto'>Choose a list and let&apos;s sort!</div>
          <div className=''>
            Or, <CoffeeButton />
          </div>
        </div>
      )}
    </div>
  );
};
export default NoList;
