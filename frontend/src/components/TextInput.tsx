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
    isPassword?: boolean; // enables password toggle functionality
    onFocus?: (e: FocusEvent) => void;
    onBlur?: (e: FocusEvent) => void;
    onChange?: (e: InputEvent) => void;
    attributes?: JSX.InputHTMLAttributes<HTMLInputElement>;
}) {
    const [showTooltip, setShowTooltip] = createSignal(false);
    const [showPassword, setShowPassword] = createSignal(false);

    // Handle input change
    const handleChange = (e: InputEvent) => {
        const target = e.target as HTMLInputElement;
        props.setValue(target.value);
        if (props.onChange) props.onChange(e);
    };

    // Determine if input should be disabled
    const isDisabled = () => props.loading || !!props.disabled;

    // Determine input type
    const inputType = () => {
        if (props.isPassword) {
            return showPassword() ? 'text' : 'password';
        }
        return 'text';
    };

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
                <div class="relative w-full">
                    <input
                        {...props.attributes}
                        ref={props.ref}
                        type={inputType()}
                        value={props.value()}
                        onInput={handleChange}
                        onFocus={props.onFocus}
                        onBlur={props.onBlur}
                        disabled={isDisabled()}
                        placeholder={props.placeholder || ''}
                        class={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 transition ${
                            props.isPassword ? 'pr-10' : ''
                        }
                ${props.error ? 'border-red-500 focus:border-red-500 focus:ring-red-400/90' : 'border-gray-300 focus:border-primary-500 focus:ring-primary-500/60'}
                ${props.success ? 'border-green-500 focus:border-green-500 focus:ring-green-400/90' : ''}
                ${isDisabled() ? 'bg-gray-100 cursor-not-allowed' : 'bg-white'}
              `}
                        onMouseEnter={() => typeof props.disabled === 'string' && setShowTooltip(true)}
                        onMouseLeave={() => setShowTooltip(false)}
                    />

                    {/* Password Toggle Button */}
                    <Show when={props.isPassword}>
                        <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword())}
                            disabled={isDisabled()}
                            class="absolute right-2 top-1/2 transform -translate-y-1/2 p-1 text-gray-500 hover:text-gray-700 disabled:cursor-not-allowed disabled:opacity-50"
                        >
                            <Show when={showPassword()} fallback={
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                </svg>
                            }>
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
                                </svg>
                            </Show>
                        </button>
                    </Show>
                </div>

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
