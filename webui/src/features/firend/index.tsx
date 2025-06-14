import useFriendsClient from "@/api/friends-client";
import { CreditScoreStatusBadge } from "@/components/credit-score/credit-score-status-badge.tsx";
import { CreditScore } from "@/components/credit-score/credit-score.tsx";
import { getCreditScoreStatus } from "@/components/credit-score/functions.ts";
import { PageError } from "@/components/page-error.tsx";
import { PageHeading } from "@/components/page-heading.tsx";
import { SessionListing } from "@/components/session-listing.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
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
			<PageHeading title={friend.name}>
				<CreditScoreStatusBadge status={getCreditScoreStatus(friend.account.creditScore)} />
			</PageHeading>

			<Card>
				<CardHeader>
					<CardTitle>Credit score</CardTitle>
					<CardDescription className="sr-only">Credit score indicator</CardDescription>
				</CardHeader>
				<CardContent>
					<CreditScore value={friend.account.creditScore} />
				</CardContent>
			</Card>

			<SessionListing
				title="Shared sessions"
				sessions={sessions ?? []}
				noSessionsMessage={`You do not have any sessions with @${friend.username}`}
				parentPath={`/friend/${friendId}`}
			/>
		</>
	);
}
