import useUserClient from "@/api/user-client.ts";
import { SessionListing } from "@/components/session-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useQuery } from "@tanstack/react-query";
import { PageHeading } from "@/components/page-heading.tsx";

function SessionListingPage() {
	useBackNavigation("/");
	const { getSessions } = useUserClient();

	const { data: sessions, isPending: sessionsPending } = useQuery({
		queryKey: ["all-sessions"],
		queryFn: () => getSessions(),
	});

	const activeSessions = sessions?.filter((s) => s.isActive) ?? [];
	const inactiveSessions = sessions?.filter((s) => !s.isActive) ?? [];

	return (
		<>
			<PageHeading title="Your sessions" />

			{sessionsPending ? (
				<SessionListing.Skeleton />
			) : (
				<SessionListing sessions={activeSessions ?? []} parentPath={"/sessions"} />
			)}

			{inactiveSessions.length > 0 && (
				<SessionListing title="Your inactive settings" sessions={inactiveSessions} />
			)}
		</>
	);
}

export default SessionListingPage;
