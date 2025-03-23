import type { User } from "@/api/types.ts";
import { create } from "zustand";
import { persist } from "zustand/middleware";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

type UserStore = {
	user: User | null;
	isLoading: boolean;
	isAuthenticated: boolean;
	fetchUser: () => Promise<void>;
	logout: () => void;
};

export const useUserStore = create<UserStore>()(
	persist(
		(set) => ({
			user: null,
			isLoading: false,
			isAuthenticated: false,
			fetchUser: async () => {
				set({ isLoading: true });

				try {
					const response = await fetch(`${API_BASE_URL}/user`, {
						method: "GET",
						credentials: "include",
					});

					if (!response.ok) {
						set({
							isLoading: false,
							isAuthenticated: false,
						});
						return;
					}

					const data: User = await response.json();
					set({ user: data, isLoading: false, isAuthenticated: true });
				} catch (error) {
					console.error("Error fetching current user", error);

					set({
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
