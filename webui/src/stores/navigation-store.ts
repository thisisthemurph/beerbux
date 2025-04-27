import { create } from "zustand";

type NavigationStore = {
	isOpen: boolean;
	open: () => void;
	close: () => void;
};

export const useNavigationStore = create<NavigationStore>((set) => ({
	isOpen: false,
	open: () => set({ isOpen: true }),
	close: () => set({ isOpen: false }),
}));
