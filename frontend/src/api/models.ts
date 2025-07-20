import type { components, paths } from "./schema";

export type User = components["schemas"]["User"];
export type Token = components["schemas"]["Token"];
export type AppError = components["schemas"]["AppError"];

export type LoginRequest = paths['/login']['post']['requestBody']['content']['application/json'];
