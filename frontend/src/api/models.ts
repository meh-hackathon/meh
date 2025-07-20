import type { components, paths } from "./schema";

export type User = components["schemas"]["User"];
export type Token = components["schemas"]["Token"];

export type LoginRequest = paths['/login']['post']['requestBody']['content']['application/json'];
