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

	type UpdateUserResponse = {
		username: string;
		name: string;
	};

	const updateUser = async (username: string, name: string) => {
		return apiFetch<UpdateUserResponse>("/user", {
			method: "PUT",
			body: JSON.stringify({ username, name }),
		});
	};

	return { getSessions, logout, getBalance, updateUser };
}

export default useUserClient;
