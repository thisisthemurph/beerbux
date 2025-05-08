import { create } from "zustand";
import { persist } from "zustand/middleware";

type AvatarDataBySessionMap = {
	[sessionId: string]: AvatarDataMap;
};

type AvatarDataMap = {
	[username: string]: UserAvatarData;
};

export type UserAvatarData = {
	username: string | undefined;
	initial: string | undefined;
	color: string | undefined;
};

type UserAvatarStoreState = {
	avatarDataBySession: AvatarDataBySessionMap;
};

type UserAvatarStoreActions = {
	// Builds the avatar data, persists it in the store, and returns it.
	setAvatarData: (sessionId: string, usernames: string[]) => AvatarDataMap;
	requiresUpdate: (sessionId: string, usernames: string[]) => boolean;
};

type UserAvatarStore = UserAvatarStoreState & UserAvatarStoreActions;

export const userAvatarStore = create<UserAvatarStore>()(
	persist(
		(set) => ({
			avatarDataBySession: {},
			setAvatarData: (sessionId: string, usernames: string[]) => {
				const avatarData = createAvatarData(usernames);
				set((state) => ({
					avatarDataBySession: {
						...state.avatarDataBySession,
						[sessionId]: avatarData,
					},
				}));
				return avatarData;
			},
			requiresUpdate: (sessionId: string, usernames: string[]): boolean => {
				const data = userAvatarStore.getState().avatarDataBySession[sessionId];
				if (!data) return true;
				const existingUsernames = Object.keys(data);

				if (!data) return true;
				return !haveSameUsernames(existingUsernames, usernames);
			},
		}),
		{
			name: "user-avatars", // Key for local storage
		},
	),
);

function haveSameUsernames(expectedUsernames: string[], usernames: string[]): boolean {
	if (expectedUsernames.length !== usernames.length) return false;

	const sortedExpectedUsernames = [...expectedUsernames].sort();
	const sortedUsernames = [...usernames].sort();

	return sortedExpectedUsernames.every((username, index) => username === sortedUsernames[index]);
}

const colors = [
	"#66c5cc",
	"#f6cf71",
	"#f89c74",
	"#dcb0f2",
	"#87c55f",
	"#9eb9f3",
	"#fe88b1",
	"#c9db74",
	"#8be0a4",
	"#b497e7",
	"#b3b3b3",
];

function createAvatarData(usernames: string[]): AvatarDataMap {
	const charsToIgnore = ["-", "_", " "];
	const avatarDataByUsername: AvatarDataMap = {};
	const usedInitials = new Set<string>();

	function getUniqueInitial(baseInitial: string, alternative: string): string {
		if (!usedInitials.has(baseInitial)) {
			usedInitials.add(baseInitial);
			return baseInitial;
		}

		if (!usedInitials.has(alternative)) {
			usedInitials.add(alternative);
			return alternative;
		}

		let index = 2;
		while (usedInitials.has(`${baseInitial}${index}`)) {
			index++;
		}

		const uniqueInitial = `${baseInitial}${index}`;
		usedInitials.add(uniqueInitial);
		return uniqueInitial;
	}

	let index = 0;
	for (const username of usernames.sort()) {
		const parts = username.split("").filter((c) => !charsToIgnore.includes(c));
		if (parts.length === 0) continue;

		const firstLetter = parts[0].toUpperCase();
		const secondLetter = parts[1]?.toUpperCase() ?? "X";
		const initial = getUniqueInitial(firstLetter, firstLetter + secondLetter);

		const color = colors[index % colors.length];
		avatarDataByUsername[username] = { username, initial, color };
		index++;
	}

	return avatarDataByUsername;
}
