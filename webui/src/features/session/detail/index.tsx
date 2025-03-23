import useSessionClient from "@/api/session-client.ts";
import {
	PrimaryActionCard,
	PrimaryActionCardButtonItem,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
	PrimaryActionCardSeparator,
} from "@/components/primary-action-card";
import { MemberDetailsCard } from "@/features/session/detail/member-details-card.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
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

	if (sessionIsPending) {
		return <p>Loading</p>;
	}

	if (!session) {
		return <p>There has been an issue loading the session.</p>;
	}

	return (
		<div className="space-y-6">
			<h1>{session.name}</h1>
			{session.isActive && (
				<PrimaryActionCard>
					<PrimaryActionCardContent>
						<PrimaryActionCardLinkItem
							to={`/session/${sessionId}/member`}
							text="Add a member"
							icon={<SquarePlus className="text-green-300 w-8 h-8" />}
						/>
						<PrimaryActionCardSeparator />
						<PrimaryActionCardButtonItem
							text="Buy a beer"
							icon={<Beer className="text-green-300 w-8 h-8" />}
							onClick={() => console.log("buy a beer")}
						/>
					</PrimaryActionCardContent>
				</PrimaryActionCard>
			)}
			<MemberDetailsCard members={session.members} />
		</div>
	);
}
