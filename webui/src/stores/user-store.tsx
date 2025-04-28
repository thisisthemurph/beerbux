import type { User } from "@/api/types/user.ts";
import { create } from "zustand";
import { persist } from "zustand/middleware";

type UserStoreState = {
	user: User | null;
	isLoading: boolean;
	isAuthenticated: boolean;
};

type UserStoreActions = {
	setIsLoading: (isLoading: boolean) => void;
	setUser: (user: User) => void;
	logout: () => void;
};

type UserStore = UserStoreState & UserStoreActions;

export const useUserStore = create<UserStore>()(
	persist(
		(set) => ({
			user: null,
			isLoading: false,
			isAuthenticated: false,
			setIsLoading: (isLoading) => set({ isLoading }),
			setUser: (user: User) => set({ user, isAuthenticated: true, isLoading: false }),
			logout: () => set({ user: null, isAuthenticated: false }),
		}),
		{
			name: "user-store", // Key for local storage
		},
	),
);
