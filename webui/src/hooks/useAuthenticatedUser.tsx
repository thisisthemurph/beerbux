import { useEffect, useState } from "react";

type AuthenticatedUser = {
	id: string;
	username: string;
};

const KEY = "authenticated-user";

export function useAuthenticatedUser() {
	const [user, setUser] = useState<AuthenticatedUser | null>(() => {
		const storedUser = localStorage.getItem("user");
		return storedUser ? JSON.parse(storedUser) : null;
	});

	useEffect(() => {
		if (user) {
			localStorage.setItem(KEY, JSON.stringify(user));
		} else {
			localStorage.removeItem(KEY);
		}
	}, [user]);

	return { user, setUser };
}
