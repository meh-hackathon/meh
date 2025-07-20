import { createRoot } from "solid-js";
import type { Token, User } from "../api/models";
import { createLocalSignal } from "../signals/createLocalStorage";

export type UserWithToken = User & {token: Token};

export const useUser = createRoot(() => {
	const [user, setUser] = createLocalSignal<UserWithToken | null>("user", null);

	const updateUser = (user: User) =>
		setUser((old) => {
			if (!old) return old;
			return { ...old, ...user };
		});

	const login = (user: UserWithToken) => setUser(user);
	const logout = () => setUser(null);

	return { user, login, logout, updateUser };
});
