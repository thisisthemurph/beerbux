import useSessionClient from "@/api/sessionClient.ts";
import { PrimaryActions } from "@/components/PrimaryActions.tsx";
import { MemberDetailsCard } from "@/features/session/detail/MemberDetailsCard.tsx";
import { useQuery } from "@tanstack/react-query";
import { SquarePlus } from "lucide-react";
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
		<div className="space-y-6">
			<h1>{session.name}</h1>
			{session.isActive && (
				<PrimaryActions
					items={[
						{
							text: "Add a member",
							href: `/session/${sessionId}/member`,
							icon: <SquarePlus className="text-green-300 w-8 h-8" />,
						},
					]}
				/>
			)}
			<MemberDetailsCard members={session.members} />
		</div>
	);
}
