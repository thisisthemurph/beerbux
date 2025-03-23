import { create } from "zustand";

export type NavigationView = "primary" | "user";

type NavigationStore = {
	view: NavigationView;
	isOpen: boolean;
	setView: (view: NavigationView) => void;
	open: (view?: NavigationView) => void;
	close: () => void;
};

export const useNavigationStore = create<NavigationStore>((set) => ({
	view: "primary",
	isOpen: false,
	setView: (view) => set({ view }),
	open: (view) => set({ view: view ?? "primary", isOpen: true }),
	close: () => set({ isOpen: false, view: "primary" }),
}));
