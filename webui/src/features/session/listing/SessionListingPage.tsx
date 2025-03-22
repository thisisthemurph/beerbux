import useUserClient from "@/api/userClient.ts";
import { SessionListing } from "@/components/SessionListing";
import { useQuery } from "@tanstack/react-query";

function SessionListingPage() {
	const { getSessions } = useUserClient();

	const { data: sessions, isLoading: sessionsLoading } = useQuery({
		queryKey: ["all-sessions"],
		queryFn: () => getSessions(),
	});

	const activeSessions = sessions?.filter((s) => s.isActive) ?? [];
	const inactiveSessions = sessions?.filter((s) => !s.isActive) ?? [];

	return (
		<section className="space-y-6">
			<h1>Your sessions</h1>
			{sessionsLoading ? (
				<SessionListing.Skeleton />
			) : (
				<SessionListing sessions={activeSessions ?? []} />
			)}

			{inactiveSessions && (
				<SessionListing
					title="Your inactive settings"
					sessions={inactiveSessions}
				/>
			)}
		</section>
	);
}

export default SessionListingPage;
