import { createSignal } from "solid-js";
import type { AppError } from "../api/models";

export function ErrorComponent(props: {
    error: AppError;
    onRetry?: () => void;
    showDetails?: boolean;
}) {
    const [showDetails, setShowDetails] = createSignal(props.showDetails || false);

    const getErrorTitle = (statusCode: number): string => {
        switch (statusCode) {
            case 400:
                return "Invalid Request";
            case 401:
                return "Unauthorized Access";
            case 403:
                return "Access Forbidden";
            case 404:
                return "Not Found";
            case 500:
                return "Server Error";
            case 503:
                return "Service Unavailable";
            default:
                return "Something Went Wrong";
        }
    };

    const getErrorDescription = (statusCode: number): string => {
        if (props.error.api_message) {
            return props.error.api_message;
        }

        switch (statusCode) {
            case 400:
                return "The request contains invalid data. Please check your input and try again.";
            case 401:
                return "You need to log in to access this resource.";
            case 403:
                return "You don't have permission to perform this action.";
            case 404:
                return "The requested resource could not be found.";
            case 500:
                return "An internal server error occurred. Please try again later.";
            case 503:
                return "The service is temporarily unavailable. Please try again later.";
            default:
                return "An unexpected error occurred. Please try again.";
        }
    };

    const getErrorIcon = (statusCode: number) => {
        if (statusCode >= 500) {
            return (
                <svg class="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.99-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
            );
        } else if (statusCode === 404) {
            return (
                <svg class="w-8 h-8 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
            );
        } else {
            return (
                <svg class="w-8 h-8 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            );
        }
    };

    return (
        <div class="flex flex-col items-center justify-center p-8 max-w-md mx-auto">
            {/* Error Icon */}
            <div class="mb-4">
                {getErrorIcon(props.error.status_code)}
            </div>

            {/* Error Title */}
            <h3 class="text-xl font-semibold text-gray-800 mb-2 text-center">
                {getErrorTitle(props.error.status_code)}
            </h3>

            {/* Error Description */}
            <p class="text-gray-600 text-center mb-6 leading-relaxed">
                {getErrorDescription(props.error.status_code)}
            </p>

            {/* Action Buttons */}
            <div class="flex flex-col sm:flex-row gap-3 w-full">
                {props.onRetry && (
                    <button
                        onClick={props.onRetry}
                        class="flex items-center justify-center px-4 py-2 bg-indigo-500 text-white font-medium rounded-lg hover:bg-indigo-600 transition-colors duration-200"
                    >
                        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                        Try Again
                    </button>
                )}

                <button
                    onClick={() => setShowDetails(!showDetails())}
                    class="flex items-center justify-center px-4 py-2 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 transition-colors duration-200"
                >
                    <svg
                        class={`w-4 h-4 mr-2 transform transition-transform duration-200 ${showDetails() ? 'rotate-180' : ''}`}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                    </svg>
                    {showDetails() ? 'Hide' : 'Show'} Details
                </button>
            </div>

            {/* Error Details */}
            {showDetails() && (
                <div class="mt-6 w-full bg-gray-50 rounded-lg p-4 border">
                    <div class="text-sm space-y-2">
                        <div class="flex justify-between">
                            <span class="font-medium text-gray-700">Error Code:</span>
                            <span class="text-gray-600 font-mono">{props.error.code}</span>
                        </div>
                        <div class="flex justify-between">
                            <span class="font-medium text-gray-700">Status Code:</span>
                            <span class="text-gray-600 font-mono">{props.error.status_code}</span>
                        </div>
                        {props.error.values && Object.keys(props.error.values).length > 0 && (
                            <div class="mt-3">
                                <span class="font-medium text-gray-700 block mb-2">Additional Info:</span>
                                <div class="bg-white rounded border p-2 space-y-1">
                                    {Object.entries(props.error.values).map(([key, value]) => (
                                        <div class="flex justify-between text-xs">
                                            <span class="text-gray-600">{key}:</span>
                                            <span class="text-gray-800 font-mono max-w-32 truncate">
                                                {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                                            </span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}
