import { type Accessor, createSignal, type JSX, Show } from 'solid-js';
import { ToolTip } from './ToolTip';

export function TextInput(props: {
    ref?: (ref: HTMLInputElement) => void;
    label?: string;
    value: Accessor<string>;
    setValue: (val: string) => void;
    error?: string;
    hint?: string; // renders a ToolTip next to the input
    disabled?: boolean | string; // if string: displays the string as a tooltip on hover
    loading?: boolean; // renders an indicator and disables the input
    placeholder?: string;
    success?: boolean;
    onFocus?: (e: FocusEvent) => void;
    onBlur?: (e: FocusEvent) => void;
    onChange?: (e: InputEvent) => void;
    attributes?: JSX.InputHTMLAttributes<HTMLInputElement>;
}) {
    const [showTooltip, setShowTooltip] = createSignal(false);

    // Handle input change
    const handleChange = (e: InputEvent) => {
        const target = e.target as HTMLInputElement;
        props.setValue(target.value);
        if (props.onChange) props.onChange(e);
    };

    // Determine if input should be disabled
    const isDisabled = () => props.loading || !!props.disabled;

    return (
        <div class="w-full">
            {/* Label and Hint */}
            <Show when={props.label || props.hint}>
                <div class="flex items-center mb-1">
                    <Show when={props.label}>
                        <label class="text-sm font-medium text-gray-700">{props.label}</label>
                    </Show>

                    <Show when={props.label && props.hint}>
                        <ToolTip hint={props.hint!} />
                    </Show>
                </div>
            </Show>

            {/* Input with Loading Indicator */}
            <div class="flex items-center gap-2">
                <input
                    {...props.attributes}
                    ref={props.ref}
                    type={'text'}
                    value={props.value()}
                    onInput={handleChange}
                    onFocus={props.onFocus}
                    onBlur={props.onBlur}
                    disabled={isDisabled()}
                    placeholder={props.placeholder || ''}
                    class={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 transition
            ${props.error ? 'border-red-500 focus:border-red-500 focus:ring-red-400/90' : 'border-gray-300 focus:border-primary-500 focus:ring-primary-500/60'}
            ${props.success ? 'border-green-500 focus:border-green-500 focus:ring-green-400/90' : ''}
            ${isDisabled() ? 'bg-gray-100 cursor-not-allowed' : 'bg-white'}
          `}
                    onMouseEnter={() => typeof props.disabled === 'string' && setShowTooltip(true)}
                    onMouseLeave={() => setShowTooltip(false)}
                />

                {/* Loading Spinner */}
                <Show when={props.loading}>
                    <div class="flex items-center">
                        <div class="w-4 h-4 border-2 border-gray-300 border-t-primary-500 rounded-full animate-spin" />
                    </div>
                </Show>

                <Show when={!props.label && props.hint}>
                    <ToolTip hint={props.hint!} />
                </Show>
            </div>

            {/* Error Message */}
            <Show when={props.error}>
                <p class="mt-1 text-xs text-red-500">{props.error}</p>
            </Show>

            {/* Disabled Tooltip */}
            <Show when={typeof props.disabled === 'string' && showTooltip()}>
                <div class="absolute z-10 p-2 mt-1 text-xs text-white bg-gray-800 rounded shadow-lg">{props.disabled as string}</div>
            </Show>
        </div>
    );
}
