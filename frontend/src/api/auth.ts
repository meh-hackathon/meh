import type { Middleware } from "openapi-fetch";
import { useUser } from "../stores/user";


   export const authMiddleware: Middleware = {

        onRequest: async ({ request }) => {
            const {user} = useUser;
            const token = user()?.token;
            if (token) {
                request.headers.set('Authorization', `Bearer ${token.access_token}`);
            }
            return request;
        },
    };
