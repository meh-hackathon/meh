import createClient from "openapi-fetch";
import { createRoot } from "solid-js";
import { authMiddleware } from "./auth";
import type { CreqteQrCodeRequest, LoginRequest } from "./models";
import type { paths } from "./schema";



export const api = createRoot(() => {
	const client = createClient<paths>({
		baseUrl: `${import.meta.env.VITE_BACKEND_HOST}/api/`,
	});
	client.use(authMiddleware);

	const login = (body: LoginRequest, opts?: {headers?: Record<string,string>}) => client.POST("/login", { body, headers: opts?.headers });
	const me = (opts?: {headers?: Record<string,string>}) => client.GET("/me", {headers: opts?.headers});

	const getQrCodes = (opts?: {headers?: Record<string,string>}) => client.GET("/qrcodes", {headers: opts?.headers});
	const createQrCode = (body: CreqteQrCodeRequest, opts?: {headers?: Record<string,string>}) => client.POST("/qrcodes", { body, headers: opts?.headers });

	return { login, me, getQrCodes, createQrCode };
});
