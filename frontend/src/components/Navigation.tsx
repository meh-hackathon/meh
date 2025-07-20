import { createSignal, onCleanup, Show } from "solid-js";
import { useUser } from "../stores/user";
import { Modal } from "./Modal";
import LoginModal from "./modals/LoginModal";

export default function  Navbar () {
    const {user, logout} = useUser

	const [isUserDropdownOpen, setIsUserDropdownOpen] = createSignal(false);
	const [isMobileMenuOpen, setIsMobileMenuOpen] = createSignal(false);

	const [showLoginModal, setShowLoginModal] = createSignal(false);

	// Close dropdowns when clicking outside
	const handleClickOutside = (e: MouseEvent) => {
		const target = e.target as Element;
		if (!target.closest(".dropdown-container")) {
			setIsUserDropdownOpen(false);
		}
	};

	document.addEventListener("click", handleClickOutside);
	onCleanup(() => document.removeEventListener("click", handleClickOutside));

	const navItems = [
		{ name: "Home", href: "#" },
		{ name: "Products", href: "#", hasDropdown: true },
		{ name: "About", href: "#" },
		{ name: "Contact", href: "#" },
	];

	const productDropdownItems = [
		{ name: "Web Apps", href: "#" },
		{ name: "Mobile Apps", href: "#" },
		{ name: "Desktop Apps", href: "#" },
		{ name: "APIs", href: "#" },
	];

	return (
		<nav class="w-full bg-white/95 backdrop-blur-md border-b border-gray-200/20 sticky top-0 z-50 shadow-sm">
		    <Modal isOpen={showLoginModal()} onClose={() => setShowLoginModal(false)}>
                <LoginModal />
            </Modal>
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex justify-between items-center h-16">
					{/* Logo */}
					<div class="flex-shrink-0 flex items-center">
						<div class="group cursor-pointer">
							<div class="flex items-center space-x-2 transition-all duration-300 group-hover:scale-105">
								<div class="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg group-hover:shadow-xl transition-all duration-300">
									<span class="text-white font-bold text-xl">A</span>
								</div>
								<span class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent">
									AppName
								</span>
							</div>
						</div>
					</div>

					{/* Desktop Navigation */}
					<div class="hidden md:block">
						<div class="ml-10 flex items-baseline space-x-8">
							{navItems.map((item) => (
								<div class="relative group dropdown-container">
									<a
										href={item.href}
										class="text-gray-700 hover:text-blue-600 px-3 py-2 text-sm font-medium transition-all duration-200 hover:scale-105 flex items-center"
									>
										{item.name}
										<Show when={item.hasDropdown}>
											<svg
												class="ml-1 h-4 w-4 transition-transform duration-200 group-hover:rotate-180"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M19 9l-7 7-7-7"
												></path>
											</svg>
										</Show>
									</a>

									{/* Products Dropdown */}
									<Show when={item.hasDropdown}>
										<div class="absolute left-0 mt-2 w-48 rounded-xl bg-white shadow-lg ring-1 ring-black ring-opacity-5 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 transform translate-y-2 group-hover:translate-y-0">
											<div class="py-1" role="menu">
												{productDropdownItems.map((dropdownItem) => (
													<a
														href={dropdownItem.href}
														class="block px-4 py-3 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
														role="menuitem"
													>
														{dropdownItem.name}
													</a>
												))}
											</div>
										</div>
									</Show>
								</div>
							))}
						</div>
					</div>

					{/* User Section */}
					<div class="flex items-center space-x-4">
						<Show
							when={user()}
							fallback={
								<button
									onClick={() => setShowLoginModal(true)}
									class="bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white px-6 py-2 rounded-full text-sm font-medium transition-all duration-300 transform hover:scale-105 hover:shadow-lg"
								>
									Sign In
								</button>
							}
						>
							{user => (
							<div class="relative dropdown-container">
								<button
									onClick={() => setIsUserDropdownOpen(!isUserDropdownOpen())}
									class="flex items-center space-x-2 text-gray-700 hover:text-blue-600 transition-colors duration-200 focus:outline-none"
								>
									<div class="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white text-sm font-medium">
										{user().username?.charAt(0).toUpperCase() || "U"}
									</div>
									<span class="hidden sm:block text-sm font-medium">
										{user().username}
									</span>
									<svg
										class={`h-4 w-4 transition-transform duration-200 ${isUserDropdownOpen() ? "rotate-180" : ""}`}
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M19 9l-7 7-7-7"
										></path>
									</svg>
								</button>

								{/* User Dropdown */}
								<div
									class={`absolute right-0 mt-2 w-48 rounded-xl bg-white shadow-lg ring-1 ring-black ring-opacity-5 transition-all duration-300 transform ${
										isUserDropdownOpen()
											? "opacity-100 visible translate-y-0"
											: "opacity-0 invisible translate-y-2"
									}`}
								>
									<div class="py-1" role="menu">
										<a
											href="#"
											class="block px-4 py-3 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors duration-200 rounded-t-xl"
										>
											Profile
										</a>
										<a
											href="#"
											class="block px-4 py-3 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors duration-200"
										>
											Settings
										</a>
										<hr class="my-1 border-gray-200" />
										<button
											onClick={() => logout()}
											class="block w-full text-left px-4 py-3 text-sm text-red-600 hover:bg-red-50 transition-colors duration-200 rounded-b-xl"
										>
											Sign Out
										</button>
									</div>
								</div>
							</div>
							)}
						</Show>

						{/* Mobile menu button */}
						<button
							onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen())}
							class="md:hidden inline-flex items-center justify-center p-2 rounded-md text-gray-700 hover:text-blue-600 hover:bg-gray-100 transition-colors duration-200"
						>
							<svg
								class="h-6 w-6"
								stroke="currentColor"
								fill="none"
								viewBox="0 0 24 24"
							>
								<path
									class={`${isMobileMenuOpen() ? "hidden" : "inline-flex"}`}
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M4 6h16M4 12h16M4 18h16"
								/>
								<path
									class={`${isMobileMenuOpen() ? "inline-flex" : "hidden"}`}
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</button>
					</div>
				</div>

				{/* Mobile Menu */}
				<div
					class={`md:hidden transition-all duration-300 ease-in-out ${
						isMobileMenuOpen()
							? "max-h-96 opacity-100"
							: "max-h-0 opacity-0 overflow-hidden"
					}`}
				>
					<div class="px-2 pt-2 pb-3 space-y-1 bg-white/90 backdrop-blur-sm rounded-b-xl">
						{navItems.map((item) => (
							<a
								href={item.href}
								class="text-gray-700 hover:text-blue-600 hover:bg-blue-50 block px-3 py-2 text-base font-medium rounded-lg transition-colors duration-200"
							>
								{item.name}
							</a>
						))}
					</div>
				</div>
			</div>
		</nav>
	);
};
