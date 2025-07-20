import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { api } from "../../api/api";
import { useNotification } from "../../stores/notification";
import { Button } from "../Button";
import { TextInput } from "../TextInput";

export function LoginModal() {
	const [email, setEmail] = createSignal("");
	const [password, setPassword] = createSignal("");
	const [emailError, setEmailError] = createSignal("");
	const [passwordError, setPasswordError] = createSignal("");
	const [isSubmitting, setIsSubmitting] = createSignal(false);

	const notification = useNotification;
	const navigate = useNavigate();

	// Email validation
	const validateEmail = (value: string): string => {
		if (!value.trim()) return "Email is required";
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		if (!emailRegex.test(value)) return "Please enter a valid email address";
		return "";
	};

	// Password validation
	const validatePassword = (value: string): string => {
		if (!value.trim()) return "Password is required";
		if (value.length < 6) return "Password must be at least 6 characters";
		return "";
	};

	// Handle form submission
	const handleSubmit = async (e: Event) => {
		e.preventDefault();

		// Clear previous errors
		setEmailError("");
		setPasswordError("");

		// Validate fields
		const emailErr = validateEmail(email());
		const passwordErr = validatePassword(password());

		if (emailErr) setEmailError(emailErr);
		if (passwordErr) setPasswordError(passwordErr);

		// Stop if validation fails
		if (emailErr || passwordErr) return;

		// Submit login request
		setIsSubmitting(true);

		// Show loading notification
		const loadingId = notification.l({
			title: "Signing in...",
			message: "Please wait while we authenticate you",
		});

		try {
			const response = await api.login({
				grant_type: "password",
				username: email(),
				password: password(),
			});

			// Remove loading notification
			notification.removeNotification(loadingId);

			if (response.error) {
				// Handle API error with error notification
				const error = response.error;
				notification.e({
					title: "Login Failed",
					message: error.api_message || `Authentication failed: ${error.code}`,
				});
			} else if (response.data) {
				// Handle successful login
				const token = response.data;

				// Store token (you might want to use a proper auth store)
				localStorage.setItem("access_token", token.access_token);
				localStorage.setItem("refresh_token", token.refresh_token);

				// Show success notification
				notification.s({
					title: "Welcome back!",
					message: "You have been successfully signed in",
				});

				// Navigate to dashboard or home page after successful login
				navigate("/dashboard", { replace: true });
			}
		} catch (error) {
			// Remove loading notification
			notification.removeNotification(loadingId);

			console.error("Login error:", error);
			notification.e({
				title: "Connection Error",
				message:
					"Unable to connect to the server. Please check your internet connection and try again.",
			});
		} finally {
			setIsSubmitting(false);
		}
	};

	// Real-time email validation
	const handleEmailChange = (value: string) => {
		setEmail(value);
		if (emailError()) {
			setEmailError(validateEmail(value));
		}
	};

	// Real-time password validation
	const handlePasswordChange = (value: string) => {
		setPassword(value);
		if (passwordError()) {
			setPasswordError(validatePassword(value));
		}
	};

	return (
		<div class="w-3xl mx-auto">
				{/* Login Form */}
				<form onSubmit={handleSubmit} class="space-y-4">
					{/* Email Input */}
					<TextInput
						label="Email"
						value={email}
						setValue={handleEmailChange}
						error={emailError()}
						placeholder="Enter your email"
						disabled={isSubmitting()}
						attributes={{
							type: "email",
							autocomplete: "email",
							required: true,
						}}
					/>

					{/* Password Input */}
					<TextInput
						label="Password"
						value={password}
						setValue={handlePasswordChange}
						error={passwordError()}
						placeholder="Enter your password"
						isPassword={true}
						disabled={isSubmitting()}
						attributes={{
							autocomplete: "current-password",
							required: true,
						}}
					/>

					{/* Submit Button */}
					<Button
						type="submit"
						label="Sign In"
						color="primary"
						disabled={isSubmitting() || !email() || !password()}
						loading={isSubmitting() ? "Signing in..." : false}
						class="w-full"
					/>
				</form>
			</div>
	);
}
