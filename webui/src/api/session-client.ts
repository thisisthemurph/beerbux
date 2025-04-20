import { apiFetch } from "@/api/api-fetch.ts";
import type { Session } from "@/api/types.ts";

type SessionCreatedResponse = {
	id: string;
	name: string;
};

function useSessionClient() {
	const getSession = async (sessionId: string) => {
		return apiFetch<Session>(`/session/${sessionId}`);
	};

	const createSession = async (name: string) => {
		return apiFetch<SessionCreatedResponse>("/session", {
			method: "POST",
			body: JSON.stringify({ name }),
		});
	};

	const addMemberToSession = async (sessionId: string, username: string) => {
		return apiFetch<void>(`/session/${sessionId}/member`, {
			method: "POST",
			body: JSON.stringify({ username }),
		});
	};

	const updateSessionMemberAdmin = async (
		sessionId: string,
		memberId: string,
		newAdminState: boolean,
	) => {
		return apiFetch<void>(`/session/${sessionId}/member/${memberId}/admin`, {
			method: "POST",
			body: JSON.stringify({ newAdminState }),
		});
	};

	return {
		getSession,
		createSession,
		addMemberToSession,
		updateSessionMemberAdmin,
	};
}

export default useSessionClient;
