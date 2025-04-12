import useSessionClient from "@/api/session-client.ts";
import useTransactionClient from "@/api/transaction-client.ts";
import type { Session, TransactionMemberAmounts, User } from "@/api/types.ts";
import { Container } from "@/components/container.tsx";
import { PageError } from "@/components/page-error.tsx";
import { SessionDetailSkeleton } from "@/features/session/detail/skeleton.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { useParams } from "react-router";
import { toast } from "sonner";
import {
	type SessionTransactionCreatedMessage,
	useSessionEventSource,
} from "./hooks/use-session-event-source.ts";
import { SessionDetailContent } from "@/features/session/detail/session-detail-content.tsx";

const SSE_BASE_URL = import.meta.env.VITE_SSE_BASE_URL;

export default function SessionDetailPage() {
	const user = useUserStore((state) => state.user) as User;
	const { sessionId } = useParams() as { sessionId: string };
	const queryClient = useQueryClient();
	const { getSession } = useSessionClient();
	const { createTransaction } = useTransactionClient();
	useBackNavigation("/");

	const url = `${SSE_BASE_URL}/session?session_id=${sessionId}&user_id=${user.id}`;

	const sessionQuery = useQuery({
		queryKey: ["session", sessionId],
		queryFn: async () => {
			return await getSession(sessionId);
		},
	});

	const handleTransactionCreatedMessage = useCallback(
		async (msg: SessionTransactionCreatedMessage) => {
			await queryClient.invalidateQueries({
				queryKey: ["session", msg.sessionId],
			});

			if (msg.creatorId === user.id) return;
			if (!sessionQuery.data) return;

			const creator = sessionQuery.data.members.find(
				(m) => m.id === msg.creatorId,
			);

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
		[queryClient, sessionQuery.data, user.id],
	);

	useSessionEventSource(url, handleTransactionCreatedMessage);

	async function handleNewTransaction(transaction: TransactionMemberAmounts) {
		if (!sessionId) return;
		const { err } = await tryCatch(createTransaction(sessionId, transaction));
		if (err) {
			console.error(err);
			toast.error("There was an issue creating the transaction.");
			return;
		}

		toast.success("You bought a round!");
	}

	if (sessionQuery.error) {
		return (
			<PageError
				message={
					sessionQuery.error.message.includes("not found")
						? "The session could not be found, please ensure you have the correct session."
						: sessionQuery.error.message
				}
			/>
		);
	}

	return (
		<Container
			isPending={sessionQuery.isPending}
			pendingComponent={<SessionDetailSkeleton />}
		>
			<SessionDetailContent
				session={sessionQuery.data as Session}
				user={user}
				handleNewTransaction={handleNewTransaction}
			/>
		</Container>
	);
}
