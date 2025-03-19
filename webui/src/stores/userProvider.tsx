import { type ReactNode, useEffect } from "react";
import { useUserStore } from "./userStore";

export const UserProvider = ({ children }: { children: ReactNode }) => {
	const { user, fetchUser } = useUserStore();

	useEffect(() => {
		if (!user) {
			fetchUser().catch(console.error);
		}
	}, [user, fetchUser]);

	return <>{children}</>;
};
