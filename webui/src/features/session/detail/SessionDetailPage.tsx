import useSessionClient from "@/api/sessionClient.ts";
import { MemberDetailsCard } from "@/features/session/detail/MemberDetailsCard.tsx";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

export default function SessionDetailPage() {
	const { getSession } = useSessionClient();
	const { sessionId } = useParams();

	const { data: session, isLoading: sessionLoading } = useQuery({
		queryKey: ["session", sessionId],
		queryFn: () => {
			if (!sessionId) {
				return null;
			}
			return getSession(sessionId);
		},
	});

	if (sessionLoading) {
		return <p>Loading</p>;
	}

	if (!session) {
		return <p>There has been an issue loading the session.</p>;
	}

	return (
		<div>
			<h1>{session.name}</h1>
			<MemberDetailsCard members={session.members} />
		</div>
	);
}
