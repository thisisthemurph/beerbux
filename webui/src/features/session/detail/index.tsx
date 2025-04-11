import useSessionClient from "@/api/session-client.ts";
import useTransactionClient from "@/api/transaction-client.ts";
import type { TransactionMemberAmounts } from "@/api/types.ts";
import {
	PrimaryActionCard,
	PrimaryActionCardButtonItem,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
} from "@/components/primary-action-card";
import { CreateTransactionDrawer } from "@/features/session/detail/create-transaction/create-transaction-drawer.tsx";
import { useSessionTransactionCreatedEventSource } from "@/features/session/detail/hooks/use-session-transaction-created-event-source.ts";
import { MemberDetailsCard } from "@/features/session/detail/member-details-card.tsx";
import { OverviewCard } from "@/features/session/detail/overview-card.tsx";
import { SessionMenu } from "@/features/session/detail/session-menu.tsx";
import { TransactionListing } from "@/features/session/detail/transaction-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useUserAvatarData } from "@/hooks/user-avatar-data.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Beer, SquarePlus } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";
import { toast } from "sonner";

export default function SessionDetailPage() {
	const queryClient = useQueryClient();
	const user = useUserStore((state) => state.user);
	const { getSession } = useSessionClient();
	const { sessionId } = useParams();
	const { createTransaction } = useTransactionClient();
	const [createDrawerOpen, setCreateDrawerOpen] = useState(false);
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

	async function handleNewTransaction(transaction: TransactionMemberAmounts) {
		if (!sessionId) return;
		const total = Object.values(transaction).reduce(
			(acc, value) => acc + value,
			0,
		);

		const { err } = await tryCatch(createTransaction(sessionId, transaction));
		if (err) {
			console.error(err);
			toast.error("There was an issue creating the transaction.");
			return;
		}

		toast.success("Transaction created:", {
			description: (
				<p>
					A transaction of <span className="font-semibold">${total}</span> has
					been created.
				</p>
			),
		});

		setCreateDrawerOpen(false);
	}

	const currentSessionMember = session?.members.find((m) => m.id === user?.id);
	const otherSessionMembers = (session?.members ?? [])
		.filter((m) => m.id !== user?.id)
		.sort((a, b) => a.name.localeCompare(b.name));

	if (sessionIsPending) {
		return <p>Loading</p>;
	}

	if (!session || !currentSessionMember) {
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
							<PrimaryActionCardButtonItem
								text="Buy a round"
								icon={<Beer className="text-green-400 w-8 h-8" />}
								onClick={() => setCreateDrawerOpen(true)}
							/>
						)}
					</PrimaryActionCardContent>
				</PrimaryActionCard>
			)}

			<MemberDetailsCard
				members={[currentSessionMember, ...otherSessionMembers]}
				avatarData={avatarData}
			/>

			<TransactionListing
				transactions={session.transactions}
				members={session.members}
				avatarData={avatarData}
			/>

			<CreateTransactionDrawer
				members={otherSessionMembers}
				onTransactionCreate={handleNewTransaction}
				open={createDrawerOpen}
				onOpenChange={(open) => setCreateDrawerOpen(open)}
			/>
		</div>
	);
}
