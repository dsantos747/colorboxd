import { CoffeeButton } from './buttons';

type Props = {
  readonly loading: boolean;
};

const NoList = ({ loading }: Props) => {
  return (
    <div className='text-center text-gray-500 space-y-10'>
      {loading && <div className='text-xl my-auto text-center'>Hang tight! Working my magic...</div>}
      {!loading && (
        <>
          <div className='text-xl my-auto text-center'>Choose a list and let&apos;s sort!</div>
          <div className=''>
            Or, <CoffeeButton />
          </div>
        </>
      )}
    </div>
  );
};
export default NoList;
