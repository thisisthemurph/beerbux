import type { Friend } from "@/api/types/friend";
import { getAvatarText } from "@/components/avatar";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { pluralize } from "@/lib/strings";
import { Link } from "react-router";

type FriendListingProps = {
	friends: Friend[];
};

function FriendListing({ friends }: FriendListingProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Your drinking buddies</CardTitle>
				<CardDescription>A list of your friends from shared sessions.</CardDescription>
			</CardHeader>
			<CardContent className="px-0">
				<section className="flex flex-col">
					{friends.length === 0 && <NoFriendsMessage />}
					{friends.map((f) => (
						<FriendListingItem friend={f} key={f.id} />
					))}
				</section>
			</CardContent>
		</Card>
	);
}

function FriendListingItem({ friend }: { friend: Friend }) {
	return (
		<Link to={`friend/${friend.id}`} className="group hover:bg-muted transition-colors">
			<div className="flex items-start gap-6 py-4 px-6">
				<Avatar className="w-10 h-10">
					<AvatarFallback className="group-hover:bg-card transition-colors">
						{getAvatarText(friend.username)}
					</AvatarFallback>
				</Avatar>
				<div className="flex flex-col gap-1 w-full">
					<p className="text-lg tracking-wide font-semibold">{friend.name}</p>
					<p className="mt-0 tracking-wider font-mono text-muted-foreground text-sm">@{friend.username}</p>
				</div>
				<Badge
					variant="secondary"
					className="text-xs text-muted-foreground font-semibold group-hover:bg-card transition-colors"
					title={`${friend.sharedSessionCount} shared ${pluralize(friend.sharedSessionCount, "session", "sessions")} with @${friend.username}`}
				>
					{friend.sharedSessionCount} {pluralize(friend.sharedSessionCount, "session", "sessions")}
				</Badge>
			</div>
		</Link>
	);
}

function NoFriendsMessage() {
	return (
		<p className="text-center py-8 font-semibold text-lg tracking-wide">
			You don't have any friends. <br /> Create a session and add friends to make some!
		</p>
	);
}

export { FriendListing };
