import { CoffeeButton } from './buttons';

const NoList = () => {
  return (
    <div className='text-center text-gray-500 space-y-10'>
      <div className='text-xl my-auto text-center'>Choose a list and let&apos;s sort!</div>
      <div className=''>
        Or, <CoffeeButton />
      </div>
    </div>
  );
};
export default NoList;
