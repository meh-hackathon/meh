import { createSignal, onMount } from "solid-js";



// Loading Component
export function LoadingComponent(props: {
    message?: string;
    size?: 'sm' | 'md' | 'lg';
}) {
    const [dots, setDots] = createSignal("");

    onMount(() => {
        const interval = setInterval(() => {
            setDots(prev => prev.length >= 3 ? "" : prev + ".");
        }, 500);

        return () => clearInterval(interval);
    });

    const sizeClasses = {
        sm: "w-6 h-6",
        md: "w-8 h-8",
        lg: "w-12 h-12"
    };

    const textSizes = {
        sm: "text-sm",
        md: "text-base",
        lg: "text-lg"
    };

    const size = props.size || 'md';

    return (
        <div class="flex flex-col items-center justify-center p-8">
            {/* Animated Spinner */}
            <div class="relative mb-4">
                <div class={`${sizeClasses[size]} border-4 border-gray-200 rounded-full animate-spin border-t-indigo-500`}></div>

                {/* Pulsing Background Circle */}
                <div class={`absolute inset-0 ${sizeClasses[size]} bg-indigo-500/10 rounded-full animate-ping`}></div>
            </div>

            {/* Loading Text */}
            <div class={`${textSizes[size]} text-gray-600 font-medium flex items-center min-h-[1.5rem]`}>
                <span>{props.message || "Loading"}</span>
                <span class="w-6 text-indigo-500">{dots()}</span>
            </div>
        </div>
    );
}
