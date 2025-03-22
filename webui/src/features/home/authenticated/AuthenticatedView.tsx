import useUserClient from "@/api/userClient.ts";
import { PrimaryActions } from "@/features/home/authenticated/PrimaryActions.tsx";
import { SessionListing } from "@/features/home/authenticated/SessionListing.tsx";
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
		queryFn: getSessions,
	});

	return (
		<div className="space-y-6">
			<UserCard {...user} />
			<PrimaryActions />
			{sessionsLoading ? (
				<SessionListing.Skeleton />
			) : (
				<SessionListing sessions={sessions ?? []} />
			)}
		</div>
	);
}
