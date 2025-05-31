import { apiFetch } from "@/api/api-fetch.ts";
import type { Friend } from "@/api/types/friend.ts";

function useFriendsClient() {
	const getFriends = async () => {
		return apiFetch<Friend[]>("/friends");
	};

	return { getFriends };
}

export default useFriendsClient;
