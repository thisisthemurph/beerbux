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

	const leaveSession = async (sessionId: string) => {
		return apiFetch<void>(`/session/${sessionId}/leave`, {
			method: "DELETE",
		});
	};

	const removeMemberFromSession = async (
		sessionId: string,
		memberId: string,
	) => {
		return apiFetch<void>(`/session/${sessionId}/member/${memberId}`, {
			method: "DELETE",
		});
	};

	const updateSessionActiveState = async (
		sessionId: string,
		newActiveState: boolean,
	) => {
		const command = newActiveState ? "activate" : "deactivate";
		return apiFetch<void>(`/session/${sessionId}/state/${command}`, {
			method: "PUT",
		});
	};

	return {
		getSession,
		createSession,
		addMemberToSession,
		updateSessionMemberAdmin,
		leaveSession,
		removeMemberFromSession,
		updateSessionActiveState,
	};
}

export default useSessionClient;
