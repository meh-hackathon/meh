import {
	createMemo,
	createSignal,
	For,
	Match,
	onMount,
	Show,
	Switch,
} from "solid-js";
import { api } from "./api/api";
import type { AppError, QrCode } from "./api/models";
import { SearchIcon } from "./assets/icons";
import { ErrorComponent } from "./components/ErrorComponent";
import { LoadingComponent } from "./components/Loading";
import { Modal } from "./components/Modal";
import { CreateQrCodeModal } from "./components/modals/CreateQrCodeModal";

// Generate QR Code URL (using a free QR code service)
const generateQrCodeUrl = (slug: string, size: number = 150) => {
	const baseUrl = window.location.origin; // Your app's base URL
	const qrUrl = `${baseUrl}/qr/${slug}`;
	return `https://api.qrserver.com/v1/create-qr-code/?size=${size}x${size}&data=${encodeURIComponent(qrUrl)}`;
};

// Format date helper
const formatDate = (dateString: string) => {
	const date = new Date(dateString);
	const now = new Date();
	const diffTime = Math.abs(now.getTime() - date.getTime());
	const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

	if (diffDays === 1) return "Today";
	if (diffDays === 2) return "Yesterday";
	if (diffDays <= 7) return `${diffDays - 1} days ago`;

	return date.toLocaleDateString("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric",
	});
};

export function QrCodes() {
	const [isLoading, setIsLoading] = createSignal(false);
	const [qrCodes, setQrCodes] = createSignal<QrCode[]>([]);
	const [error, setError] = createSignal<AppError | null>(null);
	const [selectedQrCodes, setSelectedQrCodes] = createSignal<Set<string>>(
		new Set(),
	);

	const [showCreateModal, setShowCreateModal] = createSignal(false);

	const [searchTerm, setSearchTerm] = createSignal("");

	// Sorted and filtered QR codes
	const filteredAndSortedQrCodes = createMemo(() => {
		let codes = qrCodes();

		// Filter by search term
		if (searchTerm()) {
			codes = codes.filter((qr) =>
				qr.slug.toLowerCase().includes(searchTerm().toLowerCase()),
			);
		}

		// Sort by updated_at (newest first)
		return codes.sort(
			(a, b) =>
				new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime(),
		);
	});

	// Selection handlers
	const toggleQrCode = (slug: string) => {
		const newSelected = new Set(selectedQrCodes());
		if (newSelected.has(slug)) {
			newSelected.delete(slug);
		} else {
			newSelected.add(slug);
		}
		setSelectedQrCodes(newSelected);
	};

	const selectAll = () => {
		setSelectedQrCodes(
			new Set(filteredAndSortedQrCodes().map((qr) => qr.slug)),
		);
	};

	const deselectAll = () => {
		setSelectedQrCodes(new Set<string>());
	};

	const isSelected = (slug: string) => selectedQrCodes().has(slug);
	const selectedCount = () => selectedQrCodes().size;

	const fetchData = async () => {
		setIsLoading(true);
		try {
			const response = await api.getQrCodes();
			if (response.error) {
				setError(response.error);
				return;
			}
			setQrCodes(response.data || []);
		} catch (err) {
			setError({
				api_message: "Failed to fetch QR codes",
				code: "FETCH_ERROR",
				status_code: 500,
			});
		} finally {
			setIsLoading(false);
		}
	};
	onMount(fetchData);

	const handleRetry = () => {
		setError(null);
		fetchData();
	};

	return (
		<div class="w-full min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50 py-8 px-4">
			<Modal
				isOpen={showCreateModal()}
				onClose={() => setShowCreateModal(false)}
				title="Create QR Code"
			>
				<CreateQrCodeModal close={() => setShowCreateModal(false)} />
			</Modal>

			<Switch>
				<Match when={isLoading()}>
					<div class="flex items-center justify-center min-h-[60vh]">
						<LoadingComponent message="Loading your QR codes" />
					</div>
				</Match>

				<Match when={error()}>
					{(error) => (
						<div class="flex items-center justify-center min-h-[60vh]">
							<ErrorComponent error={error()} onRetry={handleRetry} />
						</div>
					)}
				</Match>

				<Match when={qrCodes().length > 0}>
					<div class="max-w-6xl mx-auto">
						{/* Header */}
						<div class="mb-8">
							<h1 class="text-3xl font-bold text-gray-800 mb-2">
								Your QR Codes
							</h1>
							<p class="text-gray-600">
								Manage and track your dynamic QR codes
							</p>
						</div>

						{/* Controls Bar */}
						<div class="bg-white/70 backdrop-blur-sm rounded-xl p-4 mb-6 border border-white/20 shadow-sm">
							<div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
								{/* Search */}
								<div class="relative flex-1 max-w-md">
									<SearchIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
									<input
										type="text"
										placeholder="Search QR codes..."
										value={searchTerm()}
										onInput={(e) => setSearchTerm(e.currentTarget.value)}
										class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-colors"
									/>
								</div>

								{/* Selection Controls */}
								<div class="flex items-center gap-3">
									<Show when={selectedCount() > 0}>
										<span class="text-sm text-gray-600 bg-indigo-100 px-3 py-1 rounded-full">
											{selectedCount()} selected
										</span>
										<button
											onClick={deselectAll}
											class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
										>
											Clear
										</button>
									</Show>
									<button
										onClick={
											selectedCount() === filteredAndSortedQrCodes().length
												? deselectAll
												: selectAll
										}
										class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
									>
										{selectedCount() === filteredAndSortedQrCodes().length
											? "Deselect All"
											: "Select All"}
									</button>
								</div>
							</div>
						</div>

						{/* QR Codes Grid */}
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
							<For each={filteredAndSortedQrCodes()}>
								{(qrCode) => (
									<div
										class={`relative bg-white/70 backdrop-blur-sm rounded-xl p-6 border transition-all duration-200 cursor-pointer hover:shadow-lg ${
											isSelected(qrCode.slug)
												? "border-indigo-500 ring-2 ring-indigo-200 bg-indigo-50/50"
												: "border-white/20 hover:border-indigo-300"
										}`}
										onClick={() => toggleQrCode(qrCode.slug)}
									>
										{/* Selection Checkbox */}
										<div class="absolute top-4 right-4">
											<div
												class={`w-5 h-5 rounded border-2 flex items-center justify-center transition-colors ${
													isSelected(qrCode.slug)
														? "bg-indigo-500 border-indigo-500"
														: "border-gray-300 hover:border-indigo-400"
												}`}
											>
												<Show when={isSelected(qrCode.slug)}>
													<svg
														class="w-3 h-3 text-white"
														fill="none"
														stroke="currentColor"
														viewBox="0 0 24 24"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="3"
															d="M5 13l4 4L19 7"
														/>
													</svg>
												</Show>
											</div>
										</div>

										{/* QR Code Image */}
										<div class="flex justify-center mb-4">
											<div class="bg-white p-2 rounded-lg shadow-sm">
												<img
													src={generateQrCodeUrl(qrCode.slug, 120)}
													alt={`QR Code for ${qrCode.slug}`}
													class="w-24 h-24 object-contain"
													loading="lazy"
												/>
											</div>
										</div>

										{/* QR Code Info */}
										<div class="text-center">
											<h3
												class="font-semibold text-gray-800 mb-1 truncate"
												title={qrCode.name}
											>
												{qrCode.name}
											</h3>
											<p class="text-sm text-gray-500 mb-2">
												Updated {formatDate(qrCode.updated_at)}
											</p>
											<p class="text-xs text-gray-400">
												Created {formatDate(qrCode.created_at)}
											</p>
										</div>

										{/* Quick Actions */}
										<div class="mt-4 flex justify-center gap-2">
											<button
												onClick={(e) => {
													e.stopPropagation();
													// Handle edit
												}}
												class="p-2 text-gray-400 hover:text-indigo-500 hover:bg-indigo-50 rounded-lg transition-colors"
												title="Edit QR Code"
											>
												<svg
													class="w-4 h-4"
													fill="none"
													stroke="currentColor"
													viewBox="0 0 24 24"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
													/>
												</svg>
											</button>
											<button
												onClick={(e) => {
													e.stopPropagation();
													// Handle share/copy
												}}
												class="p-2 text-gray-400 hover:text-green-500 hover:bg-green-50 rounded-lg transition-colors"
												title="Copy Link"
											>
												<svg
													class="w-4 h-4"
													fill="none"
													stroke="currentColor"
													viewBox="0 0 24 24"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
													/>
												</svg>
											</button>
											<button
												onClick={(e) => {
													e.stopPropagation();
													// Handle delete
												}}
												class="p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
												title="Delete QR Code"
											>
												<svg
													class="w-4 h-4"
													fill="none"
													stroke="currentColor"
													viewBox="0 24 24"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
													/>
												</svg>
											</button>
										</div>
									</div>
								)}
							</For>
						</div>

						{/* Empty State for Search */}
						<Show
							when={
								filteredAndSortedQrCodes().length === 0 && qrCodes().length > 0
							}
						>
							<div class="text-center py-12">
								<svg
									class="w-16 h-16 text-gray-300 mx-auto mb-4"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
									/>
								</svg>
								<h3 class="text-lg font-medium text-gray-700 mb-2">
									No QR codes found
								</h3>
								<p class="text-gray-500">Try adjusting your search terms</p>
							</div>
						</Show>

						{/* Actions Bar for Selected Items */}
						<Show when={selectedCount() > 0}>
							<div class="fixed bottom-6 left-1/2 transform -translate-x-1/2 bg-white shadow-lg rounded-xl p-4 border">
								<div class="flex items-center gap-4">
									<span class="text-sm text-gray-600">
										{selectedCount()} item{selectedCount() > 1 ? "s" : ""}{" "}
										selected
									</span>
									<div class="flex gap-2">
										<button class="px-4 py-2 bg-indigo-500 text-white font-medium rounded-lg hover:bg-indigo-600 transition-colors">
											Bulk Edit
										</button>
										<button class="px-4 py-2 bg-red-500 text-white font-medium rounded-lg hover:bg-red-600 transition-colors">
											Delete Selected
										</button>
									</div>
								</div>
							</div>
						</Show>
					</div>
				</Match>

				{/* Empty State - No QR Codes */}
				<Match when={!isLoading() && !error() && qrCodes().length === 0}>
					<div class="flex flex-col items-center justify-center min-h-[60vh] text-center px-4">
						<div class="w-24 h-24 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center mb-6">
							<svg
								class="w-12 h-12 text-white"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 16h4.01M16 20h4.01"
								/>
							</svg>
						</div>
						<h2 class="text-2xl font-bold text-gray-800 mb-4">
							No QR Codes Yet
						</h2>
						<p class="text-gray-600 mb-8 max-w-md">
							Create your first dynamic QR code to get started. Add messages,
							images, and links that can surprise your audience!
						</p>
						<button
							onClick={() => setShowCreateModal(true)}
							class="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-1 transition-all duration-300"
						>
							Create Your First QR Code
						</button>
					</div>
				</Match>
			</Switch>
		</div>
	);
}
