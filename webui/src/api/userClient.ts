import { apiFetch } from "@/api/apiFetch.ts";
import type { Session } from "@/api/types.ts";

function useUserClient() {
	const getSessions = async () => {
		return apiFetch<Session[]>("/user/sessions");
	};

	return { getSessions };
}

export default useUserClient;
