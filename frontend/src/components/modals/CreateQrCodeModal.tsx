import { useNavigate } from "@solidjs/router";
import { createEffect, createSignal, Show } from "solid-js";
import { api } from "../../api/api";
import { QrCodeIcon } from "../../assets/icons";
import { useNotification } from "../../stores/notification";
import { Button } from "../Button";
import { TextInput } from "../TextInput";

export function CreateQrCodeModal({ close }: { close: () => void }) {
	const [isLoading, setIsLoading] = createSignal(false);
	const { s, e } = useNotification;
	const navigate = useNavigate();

	const [name, setName] = createSignal("");
	const [generatedSlug, setGeneratedSlug] = createSignal("");
	const [isGeneratingSlug, setIsGeneratingSlug] = createSignal(false);

	// Generate slug from name with some randomness
	const generateSlug = (inputName: string): string => {
		if (!inputName.trim()) return "";

		// Clean up the name: lowercase, replace spaces/special chars with hyphens
		const slug = inputName
			.toLowerCase()
			.trim()
			.replace(/[^a-z0-9\s-]/g, "")
			.replace(/\s+/g, "-")
			.replace(/-+/g, "-")
			.replace(/^-|-$/g, "");

		// Add 11 random characters to ensure uniqueness
		const randomSuffix = Math.random().toString(36).substring(2, 13);
		return `${slug}-${randomSuffix}`;
	};

	// Auto-generate slug when name changes
	createEffect(() => {
		const currentName = name();
		if (currentName.trim()) {
			setIsGeneratingSlug(true);
			// Small delay to show the generation animation
			setTimeout(() => {
				setGeneratedSlug(generateSlug(currentName));
				setIsGeneratingSlug(false);
			}, 300);
		} else {
			setGeneratedSlug("");
			setIsGeneratingSlug(false);
		}
	});

	// Generate QR Code URL for preview
	const getQrCodePreviewUrl = (slug: string) => {
		if (!slug) return "";
		const baseUrl = window.location.origin;
		const qrUrl = `${baseUrl}/qr/${slug}`;
		return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(qrUrl)}`;
	};

	const regenerateSlug = () => {
		if (name().trim()) {
			setIsGeneratingSlug(true);
			setTimeout(() => {
				setGeneratedSlug(generateSlug(name()));
				setIsGeneratingSlug(false);
			}, 300);
		}
	};

	const createQrCode = async () => {
		if (!name().trim()) {
			e({ title: "Error", message: "Please enter a name for your QR code" });
			return;
		}

		setIsLoading(true);
		try {
			const response = await api.createQrCode({
				name: name(),
				slug: generatedSlug(),
			});

			if (response.error) {
				e({ title: "Error", message: response.error.code });
				return;
			}

			s({
				title: "QR Code Created",
				message: "Your QR code has been created successfully!",
			});
			navigate(`/qrcodes/${response.data.id}`);
			close();
		} catch (err) {
			e({
				title: "Error",
				message: "Failed to create QR code. Please try again.",
			});
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div class="w-3xl mx-auto">
			<div class="bg-white/70 backdrop-blur-sm rounded-2xl p-8 border border-white/20 shadow-lg">
				{/* Header */}
				<div class="text-center mb-8">
					<h2 class="text-2xl font-bold text-gray-800 mb-2">
						Create New QR Code
					</h2>
					<p class="text-gray-600">
						Give your QR code a name and we'll generate everything else for you
					</p>
				</div>

				<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
					{/* Left Column - Form */}
					<div class="space-y-6">
						{/* Name Input */}
						<TextInput
							label="QR Code Name"
							value={name}
							setValue={setName}
							placeholder="e.g., Funny Rickroll, Coffee Shop Promo"
							hint="This helps you identify your QR code"
							disabled={isLoading()}
							attributes={{ maxLength: 50 }}
						/>
						<div class="text-right -mt-1">
							<span class="text-xs text-gray-400">{name().length}/50</span>
						</div>

						{/* Generated Slug Display */}
						<div>
							<div class="flex items-center justify-between mb-2">
								<label class="block text-sm font-semibold text-gray-700">
									Generated Slug
								</label>
								<Show when={generatedSlug() && !isGeneratingSlug()}>
									<Button
										onClick={regenerateSlug}
										disabled={isLoading()}
										color="secondary"
										outline={true}
										icon={() => <QrCodeIcon/>}
										label="Regenerate"
										class="text-xs"
									/>
								</Show>
							</div>
							<div class="relative">
								<div class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-700 font-mono text-sm min-h-[3rem] flex items-center">
									<Show
										when={!isGeneratingSlug() && generatedSlug()}
										fallback={
											<Show
												when={isGeneratingSlug()}
												fallback={
													<span class="text-gray-400 italic">
														Enter a name to generate slug...
													</span>
												}
											>
												<div class="flex items-center gap-2 text-gray-400">
													<div class="w-4 h-4 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
													Generating slug...
												</div>
											</Show>
										}
									>
										<span class="break-all">{generatedSlug()}</span>
									</Show>
								</div>
							</div>
							<p class="text-xs text-gray-500 mt-1">
								This will be your QR code's unique identifier
							</p>
						</div>

						{/* URL Preview */}
						<Show when={generatedSlug()}>
							<div>
								<label class="block text-sm font-semibold text-gray-700 mb-2">
									QR Code URL
								</label>
								<div class="w-full px-4 py-3 bg-gradient-to-r from-indigo-50 to-purple-50 border border-indigo-200 rounded-lg">
									<span class="text-sm text-indigo-800 font-mono break-all">
										{window.location.origin}/qr/{generatedSlug()}
									</span>
								</div>
								<p class="text-xs text-gray-500 mt-1">
									This is where your QR code will redirect
								</p>
							</div>
						</Show>
					</div>

					{/* Right Column - QR Code Preview */}
					<div class="flex flex-col items-center justify-start">
						<div class="w-full">
							<label class="block text-sm font-semibold text-gray-700 mb-4 text-center">
								QR Code Preview
							</label>

							<div class="bg-white rounded-xl p-6 shadow-sm border border-gray-200 text-center">
								<Show
									when={generatedSlug() && !isGeneratingSlug()}
									fallback={
										<div class="w-48 h-48 mx-auto bg-gray-100 rounded-lg flex items-center justify-center border-2 border-dashed border-gray-300">
											<div class="text-center">
											    <QrCodeIcon class="w-12 h-12 text-gray-400 mx-auto mb-2" />
												<p class="text-sm text-gray-500">
													QR code will appear here
												</p>
											</div>
										</div>
									}
								>
									<img
										src={getQrCodePreviewUrl(generatedSlug())}
										alt="QR Code Preview"
										class="w-48 h-48 mx-auto rounded-lg shadow-sm"
										style={{ "image-rendering": "pixelated" }}
									/>
								</Show>

								<Show when={name()}>
									<div class="mt-4 pt-4 border-t border-gray-100">
										<p class="text-sm font-medium text-gray-700">{name()}</p>
										<Show when={generatedSlug()}>
											<p class="text-xs text-gray-500 font-mono mt-1">
												{generatedSlug()}
											</p>
										</Show>
									</div>
								</Show>
							</div>

							<Show when={generatedSlug()}>
								<div class="mt-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
									<div class="flex items-start gap-2">
									    <QrCodeIcon class="w-6 h-6 text-blue-500 mt-0.5 flex-shrink-0" />
										<div class="text-xs text-blue-700">
											<p class="font-medium mb-1">Ready to create!</p>
											<p>
												After creating, you can add dynamic content like
												messages, images, and links to this QR code.
											</p>
										</div>
									</div>
								</div>
							</Show>
						</div>
					</div>
				</div>

				{/* Action Buttons */}
				<div class="flex flex-col sm:flex-row gap-3 justify-end mt-8 pt-6 border-t border-gray-200">
					<Button
						onClick={close}
						disabled={isLoading()}
						color="secondary"
						outline={true}
						label="Cancel"
					/>
					<Button
						onClick={createQrCode}
						disabled={isLoading() || !name().trim() || !generatedSlug()}
						loading={isLoading() ? "Creating..." : false}
						color="primary"
						label="Create QR Code"
						icon={() => <QrCodeIcon class="size-4"/>}
					/>
				</div>
			</div>
		</div>
	);
}
