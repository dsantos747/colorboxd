import { useState } from 'react';
import Modal from './modal';

function HelpModal() {
  const [isHelpOpen, setHelpOpen] = useState(false);

  const handleOpenHelp = () => {
    setHelpOpen(true);
  };

  return (
    <>
      <div>
        <button
          onClick={handleOpenHelp}
          className='text-sm sm:text-base py-2 px-3 sm:px-5 h-min rounded-lg bg-gray-700 bg-opacity-50 text-gray-400 enabled:hover:text-white transition-colors duration-200'>
          Algorithm Help
        </button>
      </div>
      <Modal
        isOpen={isHelpOpen}
        onClose={() => {
          setHelpOpen(false);
        }}>
        <div className='text-sm md:text-base space-y-4'>
          <h3 className='text-lg font-bold'>A note on Sorting Colors</h3>
          <p>
            Sorting images by color is not a simple task - especially for a computer! Given the variability in any collection of movie
            posters, I&apos;ve provided a range of sorting algorithms to play with. Test them out, and choose whichever works best for your
            list! A summary of each algorithm&apos;s strengths and logic is as follows:
          </p>
          <ul className='space-y-2'>
            <li>
              <strong>Hue:</strong> This sorts posters based on the tone of the dominant color. Works best if your posters all feature vivid
              colors.
            </li>
            <li>
              <strong>Luminosity:</strong> This sorts posters based on the lightness of the dominant color. Works best with monochrome or
              muted tones.
            </li>
            <li>
              <strong>Inverse Step:</strong> These algorithms sort by color, varying back and forth between light and dark within that. The
              number indicates the number of steps. The 8-step versions will work best for shorter lists.
            </li>
            <li>
              <strong>BRBW:</strong> The BRBW algorithms are short for Black-Red-Blue-White. These place all the dark tones first, then the
              vivid colors from red to blue, and finally the light tones.
            </li>
          </ul>
        </div>
      </Modal>
    </>
  );
}

export default HelpModal;
