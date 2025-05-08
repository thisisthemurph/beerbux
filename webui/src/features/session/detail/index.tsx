import useSessionClient from "@/api/session-client.ts";
import useTransactionClient from "@/api/transaction-client.ts";
import type { TransactionMemberAmounts } from "@/api/types/transaction.ts";
import type { User } from "@/api/types/user.ts";
import { PageError } from "@/components/page-error.tsx";
import { SessionDetailContent } from "@/features/session/detail/session-detail-content.tsx";
import { SessionDetailSkeleton } from "@/features/session/detail/skeleton.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import {
	type SessionTransactionCreatedMessage,
	useSessionEventSource,
} from "./hooks/use-session-event-source.ts";

const SSE_BASE_URL = import.meta.env.VITE_SSE_BASE_URL;

export default function SessionDetailPage() {
	const user = useUserStore((state) => state.user) as User;
	const { sessionId } = useParams() as { sessionId: string };
	const queryClient = useQueryClient();
	const {
		getSession,
		updateSessionMemberAdmin,
		updateSessionActiveState,
		addMemberToSession,
		leaveSession,
		removeMemberFromSession,
		getSessionHistory,
	} = useSessionClient();
	const { createTransaction } = useTransactionClient();
	const navigate = useNavigate();
	useBackNavigation("/");

	const url = `${SSE_BASE_URL}/session?session_id=${sessionId}&user_id=${user.id}`;

	const sessionQuery = useQuery({
		queryKey: ["session", sessionId],
		queryFn: async () => {
			return await getSession(sessionId);
		},
	});

	const sessionHistoryQuery = useQuery({
		queryKey: ["session-history", sessionId],
		queryFn: async () => {
			return await getSessionHistory(sessionId);
		},
	});

	const currentMember = sessionQuery.data?.members.find((m) => m.id === user.id);

	async function handleLeaveSession() {
		if (!sessionId) {
			toast.error("Could not determine the ID of the current session.");
			return;
		}

		const { err } = await tryCatch(leaveSession(sessionId));
		if (err) {
			toast.error("There was an issue leaving the session.", {
				description: err.message,
			});
			return;
		}

		toast.success("You have left the session.");
		navigate("/");
	}

	async function handleRemoveMemberFromSession(memberId: string) {
		if (!sessionId) {
			toast.error("Could not determine the ID of the current session.");
			return;
		}

		const { err } = await tryCatch(removeMemberFromSession(sessionId, memberId));
		const member = sessionQuery.data?.members.find((m) => m.id === memberId);

		if (err) {
			toast.error(`There was an issue removing ${member?.username ?? "the member"} from the session.`, {
				description: err.message,
			});
			return;
		}

		toast.success(`${member?.username ?? "The member"} has been removed from the session.`);
		await queryClient.invalidateQueries({
			queryKey: ["session", sessionId],
		});
	}

	async function handleUpdateSessionActiveState(sessionId: string, newActiveState: boolean) {
		const { err } = await tryCatch(updateSessionActiveState(sessionId, newActiveState));
		if (err) {
			toast.error(
				newActiveState
					? "There was an issue re-opening the session."
					: "There was an issue closing the session.",
			);
			return;
		}

		toast.success(newActiveState ? "The session has been re-opened." : "The session has been closed.");
		await queryClient.invalidateQueries({
			queryKey: ["session", sessionId],
		});
	}

	const updateMemberAdminStateMutation = useMutation({
		mutationFn: async (data: {
			sessionId: string;
			memberId: string;
			newAdminState: boolean;
		}) => {
			const { err } = await tryCatch(
				updateSessionMemberAdmin(data.sessionId, data.memberId, data.newAdminState),
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

			const creator = sessionQuery.data.members.find((m) => m.id === msg.creatorId);

			toast.info("New round bought", {
				description: (
					<p>
						{creator?.username ?? "Someone"} bought a{" "}
						<span className="font-semibold">${Math.round(msg.total * 100) / 100}</span> round.
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
		const { err } = await tryCatch(addMemberToSession(sessionId, username));
		if (err) {
			toast.error("There was an issue adding the member.", {
				description: err.message,
			});
			return;
		}

		await queryClient.invalidateQueries({ queryKey: ["session", sessionId] });
		toast.success(`${username} has been added to the session.`);
	}

	if (sessionQuery.isError) {
		return (
			<PageError
				message={
					sessionQuery.error.message.includes("not found")
						? "The session could not be found, please ensure you have the correct session."
						: (sessionQuery.error?.message ?? "There has been an unexpected error fetching the session.")
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
			history={sessionHistoryQuery.data}
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
			onLeaveSession={handleLeaveSession}
			onChangeSessionActiveState={() =>
				handleUpdateSessionActiveState(sessionQuery.data.id, !sessionQuery.data.isActive)
			}
			onRemoveMember={handleRemoveMemberFromSession}
		/>
	);
}
