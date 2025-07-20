import { createSignal, type JSX, onCleanup, Show } from 'solid-js';
import { ToolTip } from './ToolTip'; // Assuming ToolTip is in the same directory

type ButtonProps = {
    ref?: (ref: HTMLButtonElement) => void;
    label?: string;
    icon?: JSX.Element | (() => JSX.Element); // Function that returns JSX.Element instead of direct JSX.Element
    onClick?: () => void;
    disabled?: boolean | string; // if string: displays the string as a tooltip on hover
    loading?: boolean | string; // if boolean: shows spinner, if string: shows spinner and replaces label
    hint?: string; // renders a ToolTip next to the button
    class?: string; // Additional classes for customization
    type?: 'button' | 'submit' | 'reset'; // Button type
    color?: 'primary' | 'secondary' | 'red' | 'green' | 'blue' | 'yellow'; // Button color variant
    outline?: boolean; // Whether to use outline style instead of filled
    collapse?: boolean; // Whether to collapse the button to icon only when space is limited. Defaults to true
};

// Display modes for content
type DisplayMode = 'both' | 'label-only' | 'icon-only';

export function Button(props: ButtonProps) {
    // Validation: throw error if neither label nor icon is provided
    if (!(props.label || props.icon)) {
        throw new Error('Button must have either a label or an icon');
    }

    const [showDisabledTooltip, setShowDisabledTooltip] = createSignal(false);
    const [tooltipPosition, setTooltipPosition] = createSignal({ x: 0, y: 0 });
    const [displayMode, setDisplayMode] = createSignal<DisplayMode>(props.label && props.icon ? 'both' : props.label ? 'label-only' : 'icon-only');

    // Determine if button is disabled
    const isDisabled = () => !!props.disabled || !!props.loading;

    // Check if currently loading
    const isLoading = () => !!props.loading;

    // Get loading message if loading is a string
    const loadingMessage = () => (typeof props.loading === 'string' ? props.loading : '');

    // Get the effective label (loading message if loading, otherwise original label)
    const effectiveLabel = () => {
        if (isLoading() && loadingMessage()) {
            return loadingMessage();
        }
        return props.label;
    };

    // Get button color classes based on color prop and outline
    const getColorClasses = () => {
        if (isDisabled()) {
            return props.outline ? 'border border-gray-300 text-gray-400 bg-transparent' : 'bg-gray-300 text-gray-500';
        }

        // Default is primary if not specified
        const color = props.color || 'primary';

        if (props.outline) {
            switch (color) {
                case 'primary':
                    return 'border border-amber-500 text-amber-600 bg-transparent hover:bg-amber-50 focus:ring-amber-500';
                case 'secondary':
                    return 'border border-gray-400 text-gray-600 bg-transparent hover:bg-gray-50 focus:ring-gray-400';
                case 'red':
                    return 'border border-red-500 text-red-600 bg-transparent hover:bg-red-50 focus:ring-red-500';
                case 'green':
                    return 'border border-green-500 text-green-600 bg-transparent hover:bg-green-50 focus:ring-green-500';
                case 'blue':
                    return 'border border-blue-500 text-blue-600 bg-transparent hover:bg-blue-50 focus:ring-blue-500';
                case 'yellow':
                    return 'border border-yellow-500 text-yellow-600 bg-transparent hover:bg-yellow-50 focus:ring-yellow-500';
                default:
                    return 'border border-amber-500 text-amber-600 bg-transparent hover:bg-amber-50 focus:ring-amber-500';
            }
        }
        switch (color) {
            case 'primary':
                return 'bg-amber-500 text-white hover:bg-amber-600 focus:ring-amber-500';
            case 'secondary':
                return 'bg-gray-500 text-white hover:bg-gray-600 focus:ring-gray-400';
            case 'red':
                return 'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500';
            case 'green':
                return 'bg-green-500 text-white hover:bg-green-600 focus:ring-green-500';
            case 'blue':
                return 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500';
            case 'yellow':
                return 'bg-yellow-500 text-white hover:bg-yellow-600 focus:ring-yellow-500';
            default:
                return 'bg-amber-500 text-white hover:bg-amber-600 focus:ring-amber-500';
        }
    };

    // Get tooltip message when disabled is a string
    const disabledMessage = () => (typeof props.disabled === 'string' ? props.disabled : '');

    // Content display helpers
    const showIcon = () => {
        if (!props.icon) return false;
        return displayMode() === 'both' || displayMode() === 'icon-only';
    };

    const showLabel = () => {
        // Always show label if we have a loading message
        if (isLoading() && loadingMessage()) return true;

        if (!props.label) return false;
        return displayMode() === 'both' || displayMode() === 'label-only';
    };

    // Handle mouse events for disabled tooltip
    const handleMouseEnter = (e: MouseEvent) => {
        if (typeof props.disabled === 'string') {
            setTooltipPosition({ x: e.clientX, y: e.clientY });
            setShowDisabledTooltip(true);
        }
    };

    const handleMouseLeave = () => {
        setShowDisabledTooltip(false);
    };

    // Additional mouse event handler for the button
    const handleMouseOver = (e: MouseEvent) => {
        if (typeof props.disabled === 'string') {
            setTooltipPosition({ x: e.clientX, y: e.clientY });
            setShowDisabledTooltip(true);
        }
    };

    // Handle mouse movement for tooltip positioning
    const handleMouseMove = (e: MouseEvent) => {
        if (showDisabledTooltip()) {
            setTooltipPosition({ x: e.clientX, y: e.clientY + 10 });
        }
    };

    // References for element measurements
    let buttonRef: HTMLButtonElement | undefined;
    let buttonContainerRef: HTMLDivElement | undefined;
    let iconRef: HTMLSpanElement | undefined;
    let labelRef: HTMLSpanElement | undefined;
    let resizeObserver: ResizeObserver | undefined;

    // Check for content overflow and adjust display mode
    const checkOverflow = () => {
        if (!(buttonRef && (props.icon || props.label))) return;

        // Skip overflow checking if collapse is disabled
        if (props.collapse === false) return;

        const containerWidth = buttonRef.clientWidth;

        // Only one content type, no need to check overflow
        if (!(props.icon && props.label)) return;

        // Don't adjust display mode when loading with a custom message
        if (isLoading() && loadingMessage()) return;

        // Create temporary elements to measure content without affecting layout
        const tempContainer = document.createElement('div');
        tempContainer.style.position = 'absolute';
        tempContainer.style.visibility = 'hidden';
        tempContainer.style.display = 'inline-flex';
        tempContainer.style.alignItems = 'center';
        tempContainer.style.padding = window.getComputedStyle(buttonRef).padding;
        tempContainer.style.gap = '0.5rem'; // 2 in Tailwind
        document.body.appendChild(tempContainer);

        // Clone the icon and label to measure
        if (props.icon !== undefined && iconRef) {
            const iconClone = iconRef.cloneNode(true);
            tempContainer.appendChild(iconClone);
        }

        if (props.label && labelRef) {
            const labelClone = labelRef.cloneNode(true);
            tempContainer.appendChild(labelClone);
        }

        // Check if both elements fit
        const bothWidth = tempContainer.offsetWidth;

        // Reset and check label only
        tempContainer.innerHTML = '';
        if (props.label && labelRef) {
            const labelClone = labelRef.cloneNode(true);
            tempContainer.appendChild(labelClone);
        }
        const labelWidth = tempContainer.offsetWidth;

        // Reset and check icon only
        tempContainer.innerHTML = '';
        if (props.icon !== undefined && iconRef) {
            const iconClone = iconRef.cloneNode(true);
            tempContainer.appendChild(iconClone);
        }

        // Clean up
        document.body.removeChild(tempContainer);

        // Determine what fits
        if (bothWidth <= containerWidth) {
            setDisplayMode('both');
        } else if (labelWidth <= containerWidth) {
            setDisplayMode('label-only');
        } else {
            setDisplayMode('icon-only');
        }
    };

    // Set up refs and resize observer
    const setButtonRef = (el: HTMLButtonElement) => {
        buttonRef = el;
        buttonRef.addEventListener('mousemove', handleMouseMove);

        // Only set up resize observer if collapse is enabled (default true)
        if (props.collapse !== false) {
            // Set up resize observer to detect container size changes
            resizeObserver = new ResizeObserver(() => {
                checkOverflow();
            });

            resizeObserver.observe(el);

            // Also observe the parent element to catch constraints
            const parentElement = el.parentElement;
            if (parentElement?.parentElement) {
                resizeObserver.observe(parentElement.parentElement);
            }

            // Initial overflow check with a slight delay to ensure styles are applied
            setTimeout(checkOverflow, 10);
        }
    };

    const setIconContainerRef = (el: HTMLSpanElement) => {
        iconRef = el;
        if (props.collapse !== false) {
            setTimeout(checkOverflow, 0);
        }
    };

    const setLabelRef = (el: HTMLSpanElement) => {
        labelRef = el;
        if (props.collapse !== false) {
            setTimeout(checkOverflow, 0);
        }
    };

    // Clean up event listeners and observers
    onCleanup(() => {
        if (buttonRef) {
            buttonRef.removeEventListener('mousemove', handleMouseMove);
        }
        if (resizeObserver) {
            resizeObserver.disconnect();
        }
    });

    return (
        <div class="relative flex w-full items-center" ref={buttonContainerRef}>
            <button
                ref={r => {
                    setButtonRef(r);
                    props.ref?.(r);
                }}
                type={props.type !== undefined ? props.type : props.onClick ? 'button' : 'submit'}
                class={`
          w-full inline-flex items-center justify-center gap-2 rounded-md transition-colors overflow-hidden
          ${isDisabled() ? 'cursor-not-allowed' : 'focus:outline-none focus:ring-2 focus:ring-offset-2'}
          ${getColorClasses()}
          ${props.class || 'px-4 py-2 font-medium'}
        `}
                onClick={() => {
                    if (props.onClick && !isDisabled()) props.onClick();
                }}
                disabled={isDisabled()}
                onMouseEnter={handleMouseEnter}
                onMouseLeave={handleMouseLeave}
                onMouseOver={handleMouseOver}
                title={typeof props.disabled === 'string' ? props.disabled : undefined}>
                <Show when={isLoading()}>
                    <svg class="animate-spin h-4 w-4 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                </Show>

                <Show when={showIcon() && !isLoading()}>
                    <span ref={setIconContainerRef}>{typeof props.icon === 'function' ? props.icon() : props.icon}</span>
                </Show>

                <Show when={showLabel()}>
                    <span ref={setLabelRef} class="whitespace-nowrap overflow-hidden text-ellipsis">
                        {effectiveLabel()}
                    </span>
                </Show>
            </button>

            {/* Disabled tooltip */}
            <Show when={showDisabledTooltip() && disabledMessage()}>
                <div
                    class="absolute z-10 px-2 py-1 text-xs font-medium text-white bg-gray-900 rounded-md shadow-sm"
                    style={{
                        left: `${tooltipPosition().x}px`,
                        top: `${tooltipPosition().y}px`,
                        transform: 'translate(-50%, 8px)',
                    }}>
                    {disabledMessage()}
                    <div class="absolute w-2 h-2 bg-gray-900 transform rotate-45 -translate-y-1 -translate-x-1 left-1/2" />
                </div>
            </Show>

            {/* Hint tooltip */}
            <Show when={props.hint}>
                <div class="ml-2">
                    <ToolTip hint={props.hint!} />
                </div>
            </Show>
        </div>
    );
}
