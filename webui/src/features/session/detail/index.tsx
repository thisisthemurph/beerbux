import useSessionClient from "@/api/session-client.ts";
import useTransactionClient from "@/api/transaction-client.ts";
import type { TransactionMemberAmounts, User } from "@/api/types.ts";
import { PageError } from "@/components/page-error.tsx";
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
import { SessionDetailSkeleton } from "@/features/session/detail/skeleton.tsx";
import { TransactionListing } from "@/features/session/detail/transaction-listing.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useUserAvatarData } from "@/hooks/user-avatar-data.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Beer, SquarePlus } from "lucide-react";
import { Suspense, useState } from "react";
import type * as React from "react";
import { useParams } from "react-router";
import { toast } from "sonner";

export default function SessionDetailPage() {
	const user = useUserStore((state) => state.user);
	const { sessionId } = useParams();
	const [createDrawerOpen, setCreateDrawerOpen] = useState(false);
	useBackNavigation("/");

	if (!sessionId) throw new Error("Session id not found");
	if (!user) throw new Error("User not found");

	return (
		<div className="space-y-6">
			<Suspense fallback={<SessionDetailSkeleton />}>
				<SessionDetailsContent
					sessionId={sessionId}
					user={user}
					createDrawerOpen={createDrawerOpen}
					setCreateDrawerOpen={setCreateDrawerOpen}
				/>
			</Suspense>
		</div>
	);
}

function SessionDetailsContent({
	sessionId,
	user,
	createDrawerOpen,
	setCreateDrawerOpen,
}: {
	sessionId: string;
	user: User;
	createDrawerOpen: boolean;
	setCreateDrawerOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) {
	const queryClient = useQueryClient();
	const { getSession } = useSessionClient();
	const { createTransaction } = useTransactionClient();

	const {
		data: session,
		isLoading,
		error,
	} = useQuery({
		queryKey: ["session", sessionId],
		queryFn: async () => {
			console.log("fetching the session");
			return await getSession(sessionId);
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

	if (isLoading) {
		return <SessionDetailSkeleton />;
	}

	if (error) {
		return (
			<PageError
				message={
					error.message.includes("not found")
						? "The session could not be found, please ensure you have the correct session."
						: error.message
				}
			/>
		);
	}

	if (!session || !currentSessionMember) {
		return <p>There has been an issue loading the session.</p>;
	}

	return (
		<>
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
		</>
	);
}
