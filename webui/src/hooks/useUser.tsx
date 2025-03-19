import { useUserStore } from "@/stores/userStore";

export const useUser = () => {
	const { user, isLoading, isAuthenticated, error, logout } = useUserStore();
	return { user, isLoading, isAuthenticated, error, logout };
};
