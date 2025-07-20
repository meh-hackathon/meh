import { createEffect, type JSX, onCleanup, Show } from 'solid-js';
import { Portal } from 'solid-js/web';

interface ModalProps {
    children: JSX.Element;
    isOpen: boolean;
    onClose: () => void;
    closable?: boolean;
    ignoreBackdropClick?: boolean;
    title?: string;
}

export function Modal(props: ModalProps) {
    const closable = () => props.closable !== false;

    // Handle ESC key to close modal
    const handleKeyDown = (event: KeyboardEvent) => {
        if (event.key === 'Escape' && closable() && props.isOpen) {
            props.onClose();
        }
    };

    // Prevent scrolling on body when modal is open
    createEffect(() => {
        if (props.isOpen) {
            document.body.style.overflow = 'hidden';
            document.addEventListener('keydown', handleKeyDown);
        } else {
            document.body.style.overflow = '';
            document.removeEventListener('keydown', handleKeyDown);
        }
    });

    // Clean up event listeners when component unmounts
    onCleanup(() => {
        document.body.style.overflow = '';
        document.removeEventListener('keydown', handleKeyDown);
    });

    // Close when clicking outside content area
    const handleBackdropClick = (event: MouseEvent) => {
        if (!props.ignoreBackdropClick && closable() && event.target === event.currentTarget) {
            props.onClose();
        }
    };

    return (
        <Show when={props.isOpen}>
            <Portal>
                <div class="fixed inset-0 z-30 flex items-center justify-center p-4 bg-gray-700 bg-opacity-50 animate-fadeIn" onClick={handleBackdropClick}>
                    <div class="bg-white rounded-lg shadow-xl w-fit max-h-[90vh] flex flex-col animate-scaleIn" onClick={e => e.stopPropagation()}>
                        {/* Modal header */}
                        <div class="flex items-center justify-between p-4 border-b">
                            {props.title ? <h3 class="text-lg font-medium">{props.title}</h3> : <div />}
                            <Show when={closable()}>
                                <button type="button" class="text-gray-400 hover:text-gray-500 focus:outline-none" onClick={props.onClose} aria-label="Close">
                                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </button>
                            </Show>
                        </div>

                        {/* Modal content */}
                        <div class="p-4 grow overflow-auto">{props.children}</div>
                    </div>
                </div>
            </Portal>
        </Show>
    );
}
