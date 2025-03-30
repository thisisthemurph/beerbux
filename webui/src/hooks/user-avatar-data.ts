export type AvatarData = {
	username: string | undefined;
	initial: string | undefined;
	color: string | undefined;
};

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

export function useUserAvatarData(usernames: string[]) {
	const avatarData = createAvatarData(usernames);
	return { avatarData };
}

function createAvatarData(usernames: string[]): Record<string, AvatarData> {
	const charsToIgnore = ["-", "_", " "];
	const avatarDataByUsername: Record<string, AvatarData> = {};
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
