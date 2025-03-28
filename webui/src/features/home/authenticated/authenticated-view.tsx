import type { User } from "@/api/types.ts";
import useUserClient from "@/api/user-client.ts";
import {
	PrimaryActionCard,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
} from "@/components/primary-action-card";
import { SessionListing } from "@/components/session-listing.tsx";
import { UserCard } from "@/features/home/authenticated/user-card.tsx";
import { useQuery } from "@tanstack/react-query";
import { SquareChevronRight } from "lucide-react";
import { Link } from "react-router";

type AuthenticatedViewProps = {
	user: User;
};

export function AuthenticatedView({ user }: AuthenticatedViewProps) {
	const { getSessions } = useUserClient();

	const { data: sessions, isLoading: sessionsLoading } = useQuery({
		queryKey: ["sessions"],
		queryFn: () => getSessions(3),
	});

	return (
		<div className="space-y-6">
			<UserCard {...user} />
			<PrimaryActionCard>
				<PrimaryActionCardContent>
					<PrimaryActionCardLinkItem
						to="/session/create"
						text="Start new session"
						icon={<SquareChevronRight className="text-green-300 w-8 h-8" />}
					/>
				</PrimaryActionCardContent>
			</PrimaryActionCard>
			{sessionsLoading ? (
				<SessionListing.Skeleton />
			) : (
				<SessionListing sessions={sessions ?? []}>
					{sessions && <AllSessionsLink />}
				</SessionListing>
			)}
		</div>
	);
}

function AllSessionsLink() {
	return (
		<Link to="/sessions" className="text-blue-400">
			All sessions
		</Link>
	);
}
