import { useNavigate } from '@solidjs/router';
import { Button } from './components/Button';

export  function NotFound() {
    const navigate = useNavigate();

    return (
        <div class="w-full min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4">
            <div class="max-w-lg w-full text-center">
                {/* 404 Animation */}
                <div class="relative mb-8">
                    <div class="text-8xl md:text-9xl font-black text-gray-200 select-none">
                        404
                    </div>
                </div>

                {/* Content */}
                <div class="mb-8">
                    <h1 class="text-3xl md:text-4xl font-bold text-gray-800 mb-4">
                        Oops! Page not found
                    </h1>
                    <p class="text-lg text-gray-600 mb-6 leading-relaxed">
                        The page you're looking for seems to have wandered off into the digital void.
                        Don't worry though, even the best explorers sometimes take a wrong turn!
                    </p>
                </div>

                {/* Action Buttons */}
                <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
                    <Button
                        label="Go Home"
                        color="primary"
                        icon={() => (
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                            </svg>
                        )}
                        onClick={() => navigate('/', { replace: true })}
                        class="min-w-32"
                    />

                    <Button
                        label="Go Back"
                        color="secondary"
                        outline={true}
                        icon={() => (
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                            </svg>
                        )}
                        onClick={() => window.history.back()}
                        class="min-w-32"
                    />
                </div>

                {/* Fun Facts */}
                <div class="mt-12 p-6 bg-white/70 backdrop-blur-sm rounded-xl shadow-sm border border-white/50">
                    <h3 class="text-sm font-semibold text-gray-700 mb-3 uppercase tracking-wide">
                        Did you know?
                    </h3>
                    <p class="text-sm text-gray-600 italic">
                        The 404 error code was named after room 404 at CERN, where the original web servers were located.
                        When researchers couldn't find a file, they'd literally go to room 404 to check the servers!
                    </p>
                </div>

                {/* Floating Elements */}
                <div class="absolute top-20 left-10 w-4 h-4 bg-indigo-300 rounded-full animate-pulse opacity-60"></div>
                <div class="absolute top-40 right-16 w-6 h-6 bg-blue-300 rounded-full animate-pulse opacity-40" style="animation-delay: 1s"></div>
                <div class="absolute bottom-32 left-20 w-3 h-3 bg-purple-300 rounded-full animate-pulse opacity-50" style="animation-delay: 2s"></div>
                <div class="absolute bottom-20 right-10 w-5 h-5 bg-pink-300 rounded-full animate-pulse opacity-60" style="animation-delay: 0.5s"></div>
            </div>
        </div>
    );
}
