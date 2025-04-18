import { apiFetch } from "@/api/api-fetch.ts";
import type { Session, UserBalance } from "@/api/types.ts";

function useUserClient() {
	const logout = async () => {
		return apiFetch("/auth/logout", { method: "POST" });
	};

	const getSessions = async (pageSize = 0, pageToken: string | null = null) => {
		const params = new URLSearchParams();
		params.append("page_size", pageSize.toString());
		if (pageToken) {
			params.append("page_token", pageToken);
		}

		return apiFetch<Session[]>(`/user/sessions?${params.toString()}`);
	};

	const getBalance = async (userId: string) => {
		return apiFetch<UserBalance>(`/user/${userId}/balance`);
	};

	return { getSessions, logout, getBalance };
}

export default useUserClient;
