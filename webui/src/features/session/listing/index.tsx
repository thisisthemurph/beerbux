import useUserClient from "@/api/user-client.ts";
import { SessionListing } from "@/components/session-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useQuery } from "@tanstack/react-query";

function SessionListingPage() {
	useBackNavigation("/");
	const { getSessions } = useUserClient();

	const { data: sessions, isPending: sessionsPending } = useQuery({
		queryKey: ["all-sessions"],
		queryFn: () => getSessions(),
	});

	const activeSessions = sessions?.filter((s) => s.isActive) ?? [];
	const inactiveSessions = sessions?.filter((s) => !s.isActive) ?? [];
	//http://localhost:5173/sessions/1533f3c3-c3a0-4c57-8c35-ff7a9aa963fa?bl=%2Fsessions
	return (
		<section className="space-y-6">
			<h1>Your sessions</h1>
			{sessionsPending ? (
				<SessionListing.Skeleton />
			) : (
				<SessionListing
					sessions={activeSessions ?? []}
					parentPath={"/sessions"}
				/>
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
