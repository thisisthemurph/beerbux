import { apiFetch } from "@/api/apiFetch.ts";
import type { Session } from "@/api/types.ts";

function useUserClient() {
	const getSessions = async (pageSize = 0, pageToken: string | null = null) => {
		const params = new URLSearchParams();
		params.append("page_size", pageSize.toString());
		if (pageToken) {
			params.append("page_token", pageToken);
		}

		return apiFetch<Session[]>(`/user/sessions?${params.toString()}`);
	};

	return { getSessions };
}

export default useUserClient;
