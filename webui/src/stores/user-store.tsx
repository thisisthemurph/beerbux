import type { UserAuthDetails } from "@/api/types/user.ts";
import { create } from "zustand";
import { persist } from "zustand/middleware";

type LoggedInUserStoreState = {
	user: UserAuthDetails;
	isLoggedIn: true;
	isLoading: boolean;
};

type LoggedOutUserStoreState = {
	user: null;
	isLoggedIn: false;
	isLoading: boolean;
};

type UserStoreState = LoggedInUserStoreState | LoggedOutUserStoreState;

type UserStoreActions = {
	setIsLoading: (isLoading: boolean) => void;
	setUser: (user: UserAuthDetails) => void;
	logout: () => void;
};

type UserStore = UserStoreState & UserStoreActions;

export const useUserStore = create<UserStore>()(
	persist(
		(set) => ({
			user: null,
			isLoading: false,
			isLoggedIn: false,
			setIsLoading: (isLoading) => set({ isLoading }),
			setUser: (user: UserAuthDetails) => set({ user, isLoggedIn: true, isLoading: false }),
			logout: () => set({ user: null, isLoggedIn: false, isLoading: false }),
		}),
		{
			name: "user-store", // Key for local storage
		},
	),
);
