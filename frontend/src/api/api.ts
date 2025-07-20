import createClient from "openapi-fetch";
import { createRoot } from "solid-js";
import { authMiddleware } from "./auth";
import type { LoginRequest } from "./models";
import type { paths } from "./schema";

export const api = createRoot(() => {
	const client = createClient<paths>({
		baseUrl: `${import.meta.env.VITE_BACKEND_HOST}/api/`,
	});
	client.use(authMiddleware);

	const login = (body: LoginRequest) => client.POST("/login", { body });
	const me = () => client.GET("/me", {});

	return { login, me };
});
