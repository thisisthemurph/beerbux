import { create } from "zustand";

type BackNavigationStore = {
	backLink: string | undefined;
	setBackLink: (link: string) => void;
	clearBackLink: () => void;
};

export const useBackNavigationStore = create<BackNavigationStore>((set) => ({
	backLink: undefined,
	setBackLink: (backLink) => set({ backLink }),
	clearBackLink: () => set({ backLink: undefined }),
}));
