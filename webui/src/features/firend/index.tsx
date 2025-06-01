import useFriendsClient from "@/api/friends-client";
import { PageError } from "@/components/page-error.tsx";
import { PageHeading } from "@/components/page-heading.tsx";
import { SessionListing } from "@/components/session-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

export default function FriendDetailPage() {
	const { getFriends, getJointSessions } = useFriendsClient();
	const { friendId } = useParams() as { friendId: string };
	useBackNavigation("/");

	// This query will already be cached, so we can determine the selected friend without needing to refetch.
	const { data: friends, isPending: friendsIsPending } = useQuery({
		queryKey: ["friends"],
		queryFn: () => getFriends(),
		placeholderData: [],
	});

	const { data: sessions, isPending: sessionsIsPending } = useQuery({
		queryKey: ["friend-sessions", friendId],
		queryFn: () => getJointSessions(friendId),
		placeholderData: [],
	});

	if (friendsIsPending || sessionsIsPending) {
		return <p>Loading...</p>;
	}

	const friend = (friends ?? []).find((f) => f.id === friendId);
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
