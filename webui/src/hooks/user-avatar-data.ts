import { userAvatarStore } from "@/stores/user-avatar-store.ts";

export function useUserAvatarDataBySession(sessionId: string) {
	const avatarData = userAvatarStore((state) => state.avatarDataBySession[sessionId]);
	return avatarData ?? {};
}
