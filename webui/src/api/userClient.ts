import { apiFetch } from "@/api/apiFetch.ts";

export type Session = {
	id: string;
	name: string;
	isActive: boolean;
	members: SessionMember[];
};

type SessionMember = {
	id: string;
	name: string;
	username: string;
};

function useUserClient() {
	const getSessions = async () => {
		return apiFetch<Session[]>("/user/sessions");
	};

	return { getSessions };
}

export default useUserClient;
