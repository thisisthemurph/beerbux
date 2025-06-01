import { apiFetch } from "@/api/api-fetch.ts";
import type { Friend } from "@/api/types/friend.ts";
import type { Session } from "@/api/types/session.ts";

function useFriendsClient() {
	const getFriends = async () => {
		return apiFetch<Friend[]>("/friends");
	};

	const getJointSessions = async (memberId: string) => {
		return apiFetch<Session[]>(`/friend/${memberId}/sessions`);
	};

	return { getFriends, getJointSessions };
}

export default useFriendsClient;
