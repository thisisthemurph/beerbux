import useAuthClient from "@/api/auth-client.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { type ReactNode, useCallback, useEffect } from "react";
import { useUserStore } from "./user-store.tsx";

// UseProvider uses an optimistic check to ensure there is always a user object set.
// This API call will be made on every render if there is no user set but if the call
// fails, nothing happens.
export const UserProvider = ({ children }: { children: ReactNode }) => {
	const { getCurrentUser } = useAuthClient();
	const user = useUserStore((state) => state.user);
	const setUser = useUserStore((state) => state.setUser);
	const setIsLoading = useUserStore((state) => state.setIsLoading);

	const tryRefreshUser = useCallback(async () => {
		if (user) return;
		setIsLoading(true);

		const { data, err } = await tryCatch(getCurrentUser());
		if (err) {
			console.error("Error refreshing user", err);
			setIsLoading(false);
			return;
		}

		if (data) setUser(data);
	}, [user, getCurrentUser, setUser, setIsLoading]);

	useEffect(() => {
		void tryRefreshUser();
	}, [tryRefreshUser]);

	return <>{children}</>;
};
