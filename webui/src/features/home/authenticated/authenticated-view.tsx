import type { User } from "@/api/types.ts";
import useUserClient from "@/api/user-client.ts";
import {
	PrimaryActionCard,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
} from "@/components/primary-action-card";
import { SessionListing } from "@/components/session-listing.tsx";
import { UserCard } from "@/features/home/authenticated/user-card.tsx";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { SquareChevronRight } from "lucide-react";
import { Suspense } from "react";
import { Link } from "react-router";

type AuthenticatedViewProps = {
	user: User;
};

export function AuthenticatedView({ user }: AuthenticatedViewProps) {
	const { getSessions, getBalance } = useUserClient();

	const { data: sessions } = useSuspenseQuery({
		queryKey: ["sessions"],
		queryFn: () => getSessions(3),
	});

	const { data: balance } = useQuery({
		queryKey: ["balance", user.id],
		queryFn: () => getBalance(user.id),
		placeholderData: { credit: 0, debit: 0, net: 0 },
	});

	return (
		<>
			<UserCard {...user} netBalance={balance?.net ?? 0} />

			<PrimaryActionCard>
				<PrimaryActionCardContent>
					<PrimaryActionCardLinkItem
						to="/session/create"
						text="Start new session"
						icon={<SquareChevronRight className="text-green-300 w-8 h-8" />}
					/>
				</PrimaryActionCardContent>
			</PrimaryActionCard>

			<Suspense fallback={<SessionListing.Skeleton />}>
				<SessionListing sessions={sessions}>
					{sessions && <AllSessionsLink />}
				</SessionListing>
			</Suspense>
		</>
	);
}

function AllSessionsLink() {
	return (
		<Link to="/sessions" className="text-blue-400">
			All sessions
		</Link>
	);
}
