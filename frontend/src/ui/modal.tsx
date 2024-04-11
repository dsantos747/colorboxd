import { useEffect, useRef, useState } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';

interface Props {
  readonly isOpen: boolean;
  readonly hasCloseBtn?: boolean;
  readonly onClose?: () => void;
  readonly children: React.ReactNode;
}

export default function Modal({ isOpen, hasCloseBtn = true, onClose, children }: Props) {
  const [isModalOpen, setModalOpen] = useState<boolean>(isOpen);
  const modalRef = useRef<HTMLDialogElement | null>(null);

  const handleCloseModal = () => {
    if (onClose) {
      onClose();
    }
    setModalOpen(false);
  };

  const handleKeyDown = (event: React.KeyboardEvent<HTMLDialogElement>) => {
    if (event.key === 'Escape') {
      handleCloseModal();
    }
  };

  useEffect(() => {
    setModalOpen(isOpen);
  }, [isOpen]);

  useEffect(() => {
    const modalElement = modalRef.current;

    if (modalElement) {
      if (isModalOpen) {
        modalElement.showModal();
      } else {
        modalElement.close(); // This is causing tests to fail, see https://github.com/jsdom/jsdom/issues/3294
      }
    }
  }, [isModalOpen]);

  return (
    <dialog
      ref={modalRef}
      onKeyDown={handleKeyDown}
      className='bg-gradient-to-r from-blue-500 via-teal-500 to-lime-500 p-[2px] max-h-[80vh] relative rounded-md shadow-md shadow-black'>
      <div className='w-full h-full p-10 rounded-md relative bg-gray-100 bg-opacity-85 max-w-prose '>
        {hasCloseBtn && (
          <button className='absolute right-3 top-3' onClick={handleCloseModal}>
            <XMarkIcon className='h-6 w-6 hover:text-lime-800 transition-colors duration-200' />
          </button>
        )}
        {children}
      </div>
    </dialog>
  );
}
