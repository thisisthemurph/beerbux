import useUserClient from "@/api/userClient.ts";
import { SessionListing } from "@/components/SessionListing";
import { useQuery } from "@tanstack/react-query";

function SessionListingPage() {
	const { getSessions } = useUserClient();

	const { data: sessions, isLoading: sessionsLoading } = useQuery({
		queryKey: ["all-sessions"],
		queryFn: () => getSessions(),
	});

	return (
		<>
			<h1>Your sessions</h1>
			{sessionsLoading ? (
				<SessionListing.Skeleton />
			) : (
				<SessionListing sessions={sessions ?? []} />
			)}
		</>
	);
}

export default SessionListingPage;
