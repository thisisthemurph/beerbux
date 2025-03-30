import useSessionClient from "@/api/session-client.ts";
import {
	PrimaryActionCard,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
	PrimaryActionCardSeparator,
} from "@/components/primary-action-card";
import { Badge } from "@/components/ui/badge.tsx";
import { MemberDetailsCard } from "@/features/session/detail/member-details-card.tsx";
import { TransactionListing } from "@/features/session/detail/transaction-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useUserAvatarData } from "@/hooks/user-avatar-data.ts";
import { useQuery } from "@tanstack/react-query";
import { Beer, SquarePlus } from "lucide-react";
import { useParams } from "react-router";

export default function SessionDetailPage() {
	const { getSession } = useSessionClient();
	const { sessionId } = useParams();
	useBackNavigation("/");

	const { data: session, isPending: sessionIsPending } = useQuery({
		queryKey: ["session", sessionId],
		queryFn: () => {
			if (!sessionId) {
				return null;
			}
			return getSession(sessionId);
		},
	});

	const { avatarData } = useUserAvatarData(
		session?.members.map((m) => m.username) ?? [],
	);

	if (sessionIsPending) {
		return <p>Loading</p>;
	}

	if (!session) {
		return <p>There has been an issue loading the session.</p>;
	}

	return (
		<div className="space-y-6">
			<div className="flex justify-between m-0">
				<h1>{session.name}</h1>
				{session.total > 0 && (
					<Badge
						variant="secondary"
						className="flex justify-between gap-4 mb-8 px-6 text-lg font-normal text-muted-foreground"
					>
						<span className="">total:</span>
						<span className="font-semibold">${session.total}</span>
					</Badge>
				)}
			</div>
			{session.isActive && (
				<PrimaryActionCard>
					<PrimaryActionCardContent>
						<PrimaryActionCardLinkItem
							to={`/session/${sessionId}/member`}
							text="Add a member"
							icon={<SquarePlus className="text-green-300 w-8 h-8" />}
						/>
						{session.members.length > 1 && (
							<>
								<PrimaryActionCardSeparator />
								<PrimaryActionCardLinkItem
									text="Buy a beer"
									icon={<Beer className="text-green-300 w-8 h-8" />}
									to={`/session/${sessionId}/transaction`}
								/>
							</>
						)}
					</PrimaryActionCardContent>
				</PrimaryActionCard>
			)}
			<MemberDetailsCard members={session.members} avatarData={avatarData} />

			<TransactionListing
				transactions={session.transactions}
				members={session.members}
				avatarData={avatarData}
			/>

			{/*

			I want a listing of the transactions here.
			In reality, I want a timeline of all events that happened in the session; including user added etc.
			- Who bought the beer
			- Who received the beer (everyone vs specific people)
			- How much was the transaction?

			This should be achieved by the session service listening for ledger.updated
			events and adding the data to the sessions database.

		*/}
		</div>
	);
}
