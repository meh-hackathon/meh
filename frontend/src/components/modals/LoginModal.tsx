import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { api } from "../../api/api";
import { useNotification } from "../../stores/notification";
import { useUser } from "../../stores/user";
import { Button } from "../Button";
import { TextInput } from "../TextInput";

export function LoginModal({ close}: {close: () => void}) {
	const { login } = useUser;

	const [username, setUsername] = createSignal("");
	const [usernameError, setUsernameError] = createSignal("");
	const [password, setPassword] = createSignal("");
	const [passwordError, setPasswordError] = createSignal("");
	const [isSubmitting, setIsSubmitting] = createSignal(false);

	const notification = useNotification;
	const navigate = useNavigate();

	// Password validation
	const validatePassword = (value: string): string => {
		if (!value.trim()) return "Password is required";
		if (value.length < 6) return "Password must be at least 6 characters";
		return "";
	};

	const validateUsername = (value: string): string => {
		if (!value.trim()) return "username is required";
		return "";
	};

	// Handle form submission
	const handleSubmit = async (e: Event) => {
		e.preventDefault();

		// Clear previous errors
		setUsernameError("");
		setPasswordError("");

		// Validate fields
		const usernameErr = validateUsername(username());
		const passwordErr = validatePassword(password());

		if (usernameErr) setUsernameError(usernameErr);
		if (passwordErr) setPasswordError(passwordErr);

		// Stop if validation fails
		if (usernameErr || passwordErr) return;

		// Submit login request
		setIsSubmitting(true);

		// Show loading notification
		const loadingId = notification.l({
			title: "Signing in...",
			message: "Please wait while we authenticate you",
		});

		try {
			const tokenResponse = await api.login({
				grant_type: "password",
				username: username(),
				password: password(),
			});

			// Remove loading notification
			notification.removeNotification(loadingId);

			if (tokenResponse.error) {
				// Handle API error with error notification
				const error = tokenResponse.error;
				notification.e({
					title: "Login Failed",
					message: error.api_message || `Authentication failed: ${error.code}`,
				});
				return;
			}

				// Handle successful login
				const token = tokenResponse.data;

				// Store token (you might want to use a proper auth store)
				localStorage.setItem("access_token", token.access_token);
				localStorage.setItem("refresh_token", token.refresh_token);

				// Show success notification
				notification.s({
					title: "Welcome back!",
					message: "You have been successfully signed in",
				});

				const userResponse = await api.me({headers: { Authorization: `Bearer ${token.access_token}` }});
				if (userResponse.error) {
                    notification.e({
                        title: "Error",
                        message: "Failed to fetch user data after login.",
                    });
                    return;
                }
				const userData = userResponse.data;
				login({...userData, token: tokenResponse.data});

				// Navigate to dashboard or home page after successful login
				navigate("/dashboard", { replace: true });
				close()
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

	const handleUsernameChange = (value: string) => {
		setUsername(value);
		if (usernameError()) {
			setUsernameError(validateUsername(value));
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
				{/* Username Input */}
				<TextInput
					label="Username"
					value={username}
					setValue={handleUsernameChange}
					error={usernameError()}
					placeholder="Enter your username"
					disabled={isSubmitting()}
					attributes={{
						type: "text",
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
					disabled={isSubmitting() || !username() || !password()}
					loading={isSubmitting() ? "Signing in..." : false}
					class="w-full"
				/>
			</form>
		</div>
	);
}
