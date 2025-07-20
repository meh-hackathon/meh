import { createSignal, Show } from 'solid-js';

export function ToolTip({ hint }: { hint: string }) {
    const [showHintTooltip, setShowHintTooltip] = createSignal(false);

    return (
        <div class="relative ml-1">
            <div class="w-4 h-4 rounded-full bg-gray-200 flex items-center justify-center text-xs cursor-help" onMouseEnter={() => setShowHintTooltip(true)} onMouseLeave={() => setShowHintTooltip(false)}>
                ?
            </div>

            <Show when={showHintTooltip()}>
                <div class="absolute z-10 w-48 p-2 mt-2 text-xs text-white bg-gray-800 rounded shadow-lg -left-20 top-4">{hint}</div>
            </Show>
        </div>
    );
}
