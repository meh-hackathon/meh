export default function Home() {

    return (
        <div class="w-full min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50 py-4">
            {/* Hero Section */}
            <div class="flex flex-col items-center justify-center min-h-screen px-4 sm:px-6 lg:px-8">

                {/* Main Content */}
                <div class="text-center max-w-4xl mx-auto">
                    {/* Logo/Icon */}
                    <div class="mb-8">
                        <div class="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-2xl shadow-lg mb-6">
                            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 16h4.01M16 20h4.01" />
                            </svg>
                        </div>
                    </div>

                    {/* Heading */}
                    <h1 class="text-5xl sm:text-6xl lg:text-7xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent mb-6">
                        Dynamic QR Codes
                    </h1>

                    <p class="text-xl sm:text-2xl text-gray-600 mb-4 max-w-3xl mx-auto leading-relaxed">
                        Create QR codes that surprise and delight with
                        <span class="text-indigo-600 font-semibold"> random messages</span>,
                        <span class="text-purple-600 font-semibold"> images</span>, and
                        <span class="text-pink-600 font-semibold"> links</span>
                    </p>

                    <p class="text-lg text-gray-500 mb-12 max-w-2xl mx-auto">
                        From heartwarming messages to playful rickrolls – make every scan an adventure
                    </p>

                    {/* CTA Buttons */}
                    <div class="flex flex-col sm:flex-row gap-4 justify-center items-center mb-16">
                        <button class="group relative px-8 py-4 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-1 transition-all duration-300">
                            <span class="relative z-10">Create Your First QR Code</span>
                            <div class={`absolute inset-0 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-700 opacity-0 group-hover:opacity-100 transition-opacity duration-300`}></div>
                        </button>

                        <button class="px-8 py-4 border-2 border-gray-300 text-gray-700 font-semibold rounded-xl hover:border-indigo-500 hover:text-indigo-600 transition-all duration-300">
                            View Examples
                        </button>
                    </div>
                </div>

                {/* Feature Cards */}
                <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto mb-16">

                    {/* Feature 1 */}
                    <div class="bg-white/70 backdrop-blur-sm rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all duration-300 border border-white/20">
                        <h3 class="text-xl font-semibold text-gray-800 mb-2">Dynamic Content</h3>
                        <p class="text-gray-600">Add multiple messages, images, and links. Each scan reveals a random surprise with customizable weights.</p>
                    </div>

                    {/* Feature 2 */}
                    <div class="bg-white/70 backdrop-blur-sm rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all duration-300 border border-white/20">
                        <div class="w-12 h-12 bg-gradient-to-r from-purple-400 to-pink-500 rounded-lg flex items-center justify-center mb-4">
                            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 011-1h1a2 2 0 100-4H7a1 1 0 01-1-1V7a1 1 0 011-1h3a1 1 0 001-1V4z" />
                            </svg>
                        </div>
                        <h3 class="text-xl font-semibold text-gray-800 mb-2">Easy Management</h3>
                        <p class="text-gray-600">Update your QR codes anytime. Remove that rickroll or add new content – all without reprinting.</p>
                    </div>

                    {/* Feature 3 */}
                    <div class="bg-white/70 backdrop-blur-sm rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all duration-300 border border-white/20">
                        <div class="w-12 h-12 bg-gradient-to-r from-orange-400 to-red-500 rounded-lg flex items-center justify-center mb-4">
                            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                            </svg>
                        </div>
                        <h3 class="text-xl font-semibold text-gray-800 mb-2">Smart Analytics</h3>
                        <p class="text-gray-600">Track scan statistics, popular messages, and engagement metrics to optimize your content.</p>
                    </div>
                </div>

                {/* Example Use Case */}
                <div class="bg-white/60 backdrop-blur-sm rounded-2xl p-8 max-w-2xl mx-auto shadow-lg border border-white/20">
                    <div class="text-center">
                        <h3 class="text-2xl font-semibold text-gray-800 mb-4">Perfect For</h3>
                        <div class="grid grid-cols-2 gap-4 text-sm text-gray-600">
                            <div class="flex items-center">
                                <div class="w-2 h-2 bg-indigo-500 rounded-full mr-2"></div>
                                T-shirt designs
                            </div>
                            <div class="flex items-center">
                                <div class="w-2 h-2 bg-purple-500 rounded-full mr-2"></div>
                                Business cards
                            </div>
                            <div class="flex items-center">
                                <div class="w-2 h-2 bg-pink-500 rounded-full mr-2"></div>
                                Gift messages
                            </div>
                            <div class="flex items-center">
                                <div class="w-2 h-2 bg-blue-500 rounded-full mr-2"></div>
                                Marketing campaigns
                            </div>
                        </div>
                    </div>
                </div>

                {/* Registration Note */}
                <div class="mt-12 text-center mb-6">
                    <p class="text-gray-500 text-sm">
                        New here?
                        <button class="text-indigo-600 hover:text-indigo-800 font-medium ml-1 underline decoration-dotted">
                            Request an invitation
                        </button>
                    </p>
                </div>
            </div>
        </div>
    );
}
