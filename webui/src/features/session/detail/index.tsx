import useSessionClient from "@/api/session-client.ts";
import {
	PrimaryActionCard,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
	PrimaryActionCardSeparator,
} from "@/components/primary-action-card";
import { useSessionTransactionCreatedEventSource } from "@/features/session/detail/hooks/use-session-transaction-created-event-source.ts";
import { MemberDetailsCard } from "@/features/session/detail/member-details-card.tsx";
import { OverviewCard } from "@/features/session/detail/overview-card.tsx";
import { SessionMenu } from "@/features/session/detail/session-menu.tsx";
import { TransactionListing } from "@/features/session/detail/transaction-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useUserAvatarData } from "@/hooks/user-avatar-data.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Beer, SquarePlus } from "lucide-react";
import { useParams } from "react-router";
import { toast } from "sonner";

export default function SessionDetailPage() {
	const queryClient = useQueryClient();
	const user = useUserStore((state) => state.user);
	const { getSession } = useSessionClient();
	const { sessionId } = useParams();
	useBackNavigation("/");

	if (!sessionId) throw Error("sessionId is required");

	const { data: session, isPending: sessionIsPending } = useQuery({
		queryKey: ["session", sessionId],
		queryFn: () => {
			return getSession(sessionId);
		},
	});

	useSessionTransactionCreatedEventSource({
		sessionId: sessionId,
		userId: user?.id ?? "",
		onEventReceived: async (msg) => {
			await queryClient.invalidateQueries({
				queryKey: ["session", msg.sessionId],
			});

			if (msg.creatorId === user?.id) return;

			const creator = session?.members.find((m) => m.id === msg.creatorId);

			toast.info("New round bought", {
				description: (
					<p>
						{creator?.username ?? "Someone"} bought a{" "}
						<span className="font-semibold">
							${Math.round(msg.total * 100) / 100}
						</span>{" "}
						round.
					</p>
				),
			});
		},
	});

	const { avatarData } = useUserAvatarData(
		session?.members.map((m) => m.username) ?? [],
	);

	if (sessionIsPending) {
		return <p>Loading</p>;
	}

	if (!session) {
		return <p>There has been an issue loading the session.</p>;
	}

	return (
		<div className="space-y-6">
			<div className="flex justify-between items-center mb-8">
				<h1 className="mb-0">{session.name}</h1>
				<SessionMenu />
			</div>
			<OverviewCard total={session.total} />
			{session.isActive && (
				<PrimaryActionCard>
					<PrimaryActionCardContent>
						<PrimaryActionCardLinkItem
							to={`/session/${sessionId}/member`}
							text="Add a member"
							icon={<SquarePlus className="text-green-400 w-8 h-8" />}
						/>
						{session.members.length > 1 && (
							<>
								<PrimaryActionCardSeparator />
								<PrimaryActionCardLinkItem
									text="Buy a round"
									icon={<Beer className="text-green-400 w-8 h-8" />}
									to={`/session/${sessionId}/transaction`}
								/>
							</>
						)}
					</PrimaryActionCardContent>
				</PrimaryActionCard>
			)}

			<MemberDetailsCard members={session.members} avatarData={avatarData} />

			<TransactionListing
				transactions={session.transactions}
				members={session.members}
				avatarData={avatarData}
			/>
		</div>
	);
}
