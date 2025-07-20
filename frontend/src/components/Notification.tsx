import { For, Show } from "solid-js";
import { type Notification, useNotification } from "../stores/notification";

export default function NotificationContainer() {
	const { notifications, removeNotification } = useNotification;

	const getNotificationStyles = (type: Notification["type"]) => {
		const baseStyles =
			"relative flex items-start gap-3 p-4 rounded-lg shadow-lg border backdrop-blur-sm transition-all duration-300 hover:shadow-xl";

		switch (type) {
			case "success":
				return `${baseStyles} bg-green-50/90 border-green-200 text-green-800`;
			case "error":
				return `${baseStyles} bg-red-50/90 border-red-200 text-red-800`;
			case "warning":
				return `${baseStyles} bg-yellow-50/90 border-yellow-200 text-yellow-800`;
			case "info":
				return `${baseStyles} bg-blue-50/90 border-blue-200 text-blue-800`;
			case "loading":
				return `${baseStyles} bg-gray-50/90 border-gray-200 text-gray-800`;
			case "aurevoir":
				return `${baseStyles} bg-purple-50/90 border-purple-200 text-purple-800`;
			default:
				return `${baseStyles} bg-gray-50/90 border-gray-200 text-gray-800`;
		}
	};

	const getIcon = (type: Notification["type"]) => {
		switch (type) {
			case "success":
				return (
					<svg
						class="w-5 h-5 text-green-500 mt-0.5 flex-shrink-0"
						fill="currentColor"
						viewBox="0 0 20 20"
					>
						<path
							fill-rule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
							clip-rule="evenodd"
						/>
					</svg>
				);
			case "error":
				return (
					<svg
						class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0"
						fill="currentColor"
						viewBox="0 0 20 20"
					>
						<path
							fill-rule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
							clip-rule="evenodd"
						/>
					</svg>
				);
			case "warning":
				return (
					<svg
						class="w-5 h-5 text-yellow-500 mt-0.5 flex-shrink-0"
						fill="currentColor"
						viewBox="0 0 20 20"
					>
						<path
							fill-rule="evenodd"
							d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
							clip-rule="evenodd"
						/>
					</svg>
				);
			case "info":
				return (
					<svg
						class="w-5 h-5 text-blue-500 mt-0.5 flex-shrink-0"
						fill="currentColor"
						viewBox="0 0 20 20"
					>
						<path
							fill-rule="evenodd"
							d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
							clip-rule="evenodd"
						/>
					</svg>
				);
			case "loading":
				return (
					<svg
						class="w-5 h-5 text-gray-500 mt-0.5 flex-shrink-0 animate-spin"
						fill="none"
						viewBox="0 0 24 24"
					>
						<circle
							class="opacity-25"
							cx="12"
							cy="12"
							r="10"
							stroke="currentColor"
							stroke-width="4"
						></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						></path>
					</svg>
				);
			case "aurevoir":
				return (
					<svg
						class="w-5 h-5 text-purple-500 mt-0.5 flex-shrink-0"
						fill="currentColor"
						viewBox="0 0 20 20"
					>
						<path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				);
			default:
				return null;
		}
	};

	const handleClose = (id: number) => {
		removeNotification(id);
	};

	return (
		<div class="fixed top-4 right-4 z-50 flex flex-col gap-2 max-w-sm w-full">
			<For each={notifications()}>
				{(notification) => (
					<div
						class={`${getNotificationStyles(notification.type)} animate-in slide-in-from-right-full duration-300`}
						style={{
							animation: "slideInRight 0.3s ease-out",
						}}
					>
						{getIcon(notification.type)}

						<div class="flex-1 min-w-0">
							<div class="font-semibold text-sm leading-5">
								{notification.title}
							</div>
							<Show when={notification.message}>
								<div class="text-sm opacity-90 mt-1 leading-relaxed">
									{notification.message}
								</div>
							</Show>
						</div>

						<button
							onClick={() => handleClose(notification.id)}
							class="flex-shrink-0 ml-2 p-1 rounded-md hover:bg-black/5 transition-colors duration-200"
							aria-label="Close notification"
						>
							<svg
								class="w-4 h-4 opacity-60 hover:opacity-100"
								fill="currentColor"
								viewBox="0 0 20 20"
							>
								<path
									fill-rule="evenodd"
									d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
									clip-rule="evenodd"
								/>
							</svg>
						</button>
					</div>
				)}
			</For>
	</div>
	);
}
