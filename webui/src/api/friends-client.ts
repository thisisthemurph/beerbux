import { apiFetch } from "@/api/api-fetch.ts";
import type { Friend } from "@/api/types/friend.ts";
import type { Session } from "@/api/types/session.ts";
import type { User } from "./types/user";

function useFriendsClient() {
	const getFriends = async () => {
		return apiFetch<Friend[]>("/friends");
	};

	const getFriend = async (friendId: string) => {
		return apiFetch<User>(`/friend/${friendId}`);
	};

	const getJointSessions = async (memberId: string) => {
		return apiFetch<Session[]>(`/friend/${memberId}/sessions`);
	};

	return { getFriends, getFriend, getJointSessions };
}

export default useFriendsClient;
