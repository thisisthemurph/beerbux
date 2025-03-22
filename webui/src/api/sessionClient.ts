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

	return { getSession, createSession };
}

export default useSessionClient;
