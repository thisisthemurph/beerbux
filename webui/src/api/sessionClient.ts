import { apiFetch } from "@/api/apiFetch.ts";
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

	return { getSession, createSession, addMemberToSession };
}

export default useSessionClient;
