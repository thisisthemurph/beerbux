export function getAvatarText(name: string, defaultValue = "UN"): string {
	if (!name) return defaultValue;
	const words = name.split(" ");
	if (words.length === 0) return defaultValue;
	if (words.length === 1) return words[0][0].toUpperCase();
	return `${words[0][0].toUpperCase()}${words[1][0].toUpperCase()}`;
}
