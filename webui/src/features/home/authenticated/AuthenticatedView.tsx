import useUserClient from "@/api/userClient.ts";
import { PrimaryActions } from "@/components/PrimaryActions.tsx";
import {
	AllSessionsLink,
	SessionListing,
} from "@/components/SessionListing.tsx";
import { UserCard } from "@/features/home/authenticated/UserCard.tsx";
import type { User } from "@/stores/userStore.tsx";
import { useQuery } from "@tanstack/react-query";
import { SquareChevronRight } from "lucide-react";

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
			<PrimaryActions
				items={[
					{
						text: "Start new session",
						href: "/session/create",
						icon: <SquareChevronRight className="text-green-300 w-8 h-8" />,
					},
				]}
			/>
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
