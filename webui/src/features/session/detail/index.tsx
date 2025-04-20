import useSessionClient from "@/api/session-client.ts";
import useTransactionClient from "@/api/transaction-client.ts";
import type { TransactionMemberAmounts, User } from "@/api/types.ts";
import { PageError } from "@/components/page-error.tsx";
import { SessionDetailSkeleton } from "@/features/session/detail/skeleton.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
	const { getSession, updateSessionMemberAdmin } = useSessionClient();
	const { createTransaction } = useTransactionClient();
	const { addMemberToSession } = useSessionClient();
	useBackNavigation("/");

	const url = `${SSE_BASE_URL}/session?session_id=${sessionId}&user_id=${user.id}`;

	const sessionQuery = useQuery({
		queryKey: ["session", sessionId],
		queryFn: async () => {
			return await getSession(sessionId);
		},
	});

	const currentMember = sessionQuery.data?.members.find(
		(m) => m.id === user.id,
	);

	const updateMemberAdminStateMutation = useMutation({
		mutationFn: async (data: {
			sessionId: string;
			memberId: string;
			newAdminState: boolean;
		}) => {
			const { err } = await tryCatch(
				updateSessionMemberAdmin(
					data.sessionId,
					data.memberId,
					data.newAdminState,
				),
			);
			if (err) {
				toast.error(
					data.newAdminState
						? "There has been an issue setting the member as an admin."
						: "There has been an issue removing the member's admin status.",
				);
				return;
			}

			await queryClient.invalidateQueries({
				queryKey: ["session", data.sessionId],
			});
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
			toast.error("There was an issue creating the transaction.");
			return;
		}

		toast.success("You bought a round!");
	}

	async function handleAddMember(username: string) {
		if (!sessionId) return;
		const {err} = await tryCatch(addMemberToSession(sessionId, username));
		if (err) {
			toast.error("There was an issue adding the member.", {
				description: err.message,
			});
			return;
		}

		await queryClient.invalidateQueries({queryKey: ["session", sessionId]});
		toast.success(`${username} has been added to the session.`);
	}

	if (sessionQuery.isError) {
		return (
			<PageError
				message={
					sessionQuery.error.message.includes("not found")
						? "The session could not be found, please ensure you have the correct session."
						: (sessionQuery.error?.message ??
							"There has been an unexpected error fetching the session.")
				}
			/>
		);
	}

	if (sessionQuery.isPending) {
		return <SessionDetailSkeleton />;
	}

	if (!currentMember) {
		return <PageError message="You are not a member of this session." />;
	}

	return (
		<SessionDetailContent
			session={sessionQuery.data}
			user={user}
			handleNewTransaction={handleNewTransaction}
			handleAddMember={handleAddMember}
			onMemberAdminStateUpdate={(sessionId, memberId, newAdminState) =>
				updateMemberAdminStateMutation.mutate({
					sessionId,
					memberId,
					newAdminState,
				})
			}
		/>
	);
}
