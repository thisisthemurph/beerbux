import useFriendsClient from "@/api/friends-client";
import { PageError } from "@/components/page-error.tsx";
import { PageHeading } from "@/components/page-heading.tsx";
import { SessionListing } from "@/components/session-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

export default function FriendDetailPage() {
	const { getFriend, getJointSessions } = useFriendsClient();
	const { friendId } = useParams() as { friendId: string };
	useBackNavigation("/");

	const { data: friend, isPending: friendsIsPending } = useQuery({
		queryKey: ["friend", friendId],
		queryFn: () => getFriend(friendId),
	});

	const { data: sessions, isPending: sessionsIsPending } = useQuery({
		queryKey: ["friend-sessions", friendId],
		queryFn: () => getJointSessions(friendId),
		placeholderData: [],
	});

	if (friendsIsPending || sessionsIsPending) {
		return <p>Loading...</p>;
	}

	if (!friend) {
		return <PageError message="The specified friend could not be found" />;
	}

	return (
		<>
			<PageHeading title={friend.name} className="flex-col items-start">
				<p className="text-muted-foreground tracking-wider font-mono">@{friend.username}</p>
			</PageHeading>
			{sessionsIsPending && <p>Loading shared sessions...</p>}
			<SessionListing
				title={`Your sessions with ${friend.username}`}
				sessions={sessions ?? []}
				noSessionsMessage={`You do not have any sessions with @${friend.username}`}
				parentPath={`/friend/${friendId}`}
			/>
		</>
	);
}
