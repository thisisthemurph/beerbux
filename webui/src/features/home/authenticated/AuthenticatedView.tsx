import useUserClient from "@/api/userClient.ts";
import {
	AllSessionsLink,
	SessionListing,
} from "@/components/SessionListing.tsx";
import { PrimaryActions } from "@/features/home/authenticated/PrimaryActions.tsx";
import { UserCard } from "@/features/home/authenticated/UserCard.tsx";
import type { User } from "@/stores/userStore.tsx";
import { useQuery } from "@tanstack/react-query";

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
			<PrimaryActions />
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
