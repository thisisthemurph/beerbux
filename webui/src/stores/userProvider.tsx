import { type ReactNode, useEffect } from "react";
import { useUserStore } from "./userStore";

export const UserProvider = ({ children }: { children: ReactNode }) => {
	const user = useUserStore((state) => state.user);
	const fetchUser = useUserStore((state) => state.fetchUser);

	useEffect(() => {
		if (!user) {
			fetchUser().catch(console.error);
		}
	}, [user, fetchUser]);

	return <>{children}</>;
};
