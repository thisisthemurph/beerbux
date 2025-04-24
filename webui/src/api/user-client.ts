import { apiFetch } from "@/api/api-fetch.ts";
import type { Session } from "@/api/types/session.ts";
import type { UserBalance } from "@/api/types/user.ts";

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
