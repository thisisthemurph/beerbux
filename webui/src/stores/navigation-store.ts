import { create } from "zustand";

type NavigationStore = {
	isOpen: boolean;
	open: () => void;
	close: () => void;
	toggle: (open: boolean) => void;
};

export const useNavigationStore = create<NavigationStore>((set) => ({
	isOpen: false,
	open: () => set({ isOpen: true }),
	close: () => set({ isOpen: false }),
	toggle: (open) => set({ isOpen: open }),
}));
