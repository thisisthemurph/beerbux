import { cn } from "@/lib/utils.ts";

type Variant = "default" | "subtle";

type UsernameProps = {
	username: string | undefined;
	isDeleted?: boolean | undefined;
	variant?: Variant;
};

export function Username({ username, isDeleted, variant }: UsernameProps) {
	return (
		<span
			className={cn(
				"font-semibold tracking-wider",
				isDeleted && "line-through",
				variant === "subtle" && "text-muted-foreground",
			)}
		>
			{username ?? "unknown"}
		</span>
	);
}

type UsernameGroupProps = {
	usernames: string[];
	maxMembers: number;
};

export function UsernameGroup({ usernames, maxMembers }: UsernameGroupProps) {
	const len = usernames.length;

	if (len >= maxMembers)
		return <Username variant="subtle" username="everyone" />;
	if (len === 0) return <Username variant="subtle" username="unknown" />;
	if (len === 1) return <Username variant="subtle" username={usernames[0]} />;
	if (len === 2)
		return (
			<>
				<Username variant="subtle" username={usernames[0]} /> and{" "}
				<Username variant="subtle" username={usernames[1]} />
			</>
		);

	return (
		<>
			{usernames.slice(0, -1).map((username) => (
				<>
					<Username variant="subtle" key={username} username={username} />,
					{"\u00A0"}
				</>
			))}
			and {"\u00A0"}
			<Username variant="subtle" username={usernames.slice(-1)[0]} />
		</>
	);
}
