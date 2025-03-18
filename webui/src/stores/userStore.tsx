import { create } from "zustand";
import { persist } from "zustand/middleware";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export type User = {
	id: string;
	name: string;
	username: string;
	netBalance: number;
};

type UserStore = {
	user: User | null;
	isLoading: boolean;
	isAuthenticated: boolean;
	error: string | null;
	fetchUser: () => Promise<void>;
	logout: () => void;
};

export const useUserStore = create<UserStore>()(
	persist(
		(set) => ({
			user: null,
			isLoading: false,
			isAuthenticated: false,
			error: null,
			fetchUser: async () => {
				set({ isLoading: true, error: null });
				try {
					const response = await fetch(`${API_BASE_URL}/user`, {
						method: "GET",
						credentials: "include",
					});

					if (!response.ok) {
						set({
							error: "Failed to load the user",
							isLoading: false,
							isAuthenticated: false,
						});
						return;
					}

					const data: User = await response.json();
					set({ user: data, isLoading: false, isAuthenticated: true });
				} catch (error) {
					set({
						error: (error as Error).message,
						isLoading: false,
						isAuthenticated: false,
					});
				}
			},
			logout: () => set({ user: null, isAuthenticated: false }),
		}),
		{
			name: "user-store", // Key for local storage
		},
	),
);
