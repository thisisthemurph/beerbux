import { useState } from "react";
import { useUserAvatarData } from "@/hooks/user-avatar-data.ts";
import { PageError } from "@/components/page-error.tsx";
import { SessionMenu } from "@/features/session/detail/session-menu.tsx";
import { OverviewCard } from "@/features/session/detail/overview-card.tsx";
import {
	PrimaryActionCard,
	PrimaryActionCardButtonItem,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
} from "@/components/primary-action-card.tsx";
import { Beer, SquarePlus } from "lucide-react";
import { MemberDetailsCard } from "@/features/session/detail/member-details-card.tsx";
import { TransactionListing } from "@/features/session/detail/transaction-listing.tsx";
import { CreateTransactionDrawer } from "@/features/session/detail/create-transaction/create-transaction-drawer.tsx";
import type { Session, TransactionMemberAmounts, User } from "@/api/types.ts";

type SessionDetailContentProps = {
	session: Session;
	user: User;
	handleNewTransaction: (
		transaction: TransactionMemberAmounts,
	) => Promise<void>;
};

export function SessionDetailContent({
	session,
	user,
	handleNewTransaction,
}: SessionDetailContentProps) {
	const [createDrawerOpen, setCreateDrawerOpen] = useState(false);
	const currentSessionMember = session?.members.find((m) => m.id === user?.id);
	const otherSessionMembers = (session?.members ?? [])
		.filter((m) => m.id !== user.id)
		.sort((a, b) => a.name.localeCompare(b.name));

	const { avatarData } = useUserAvatarData(
		session?.members.map((m) => m.username) ?? [],
	);

	if (!currentSessionMember) {
		return <PageError message="You are not a member of this session." />;
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
							to={`/session/${session.id}/member`}
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
				onTransactionCreate={async (data) => {
					await handleNewTransaction(data);
					setCreateDrawerOpen(false);
				}}
				open={createDrawerOpen}
				onOpenChange={(open) => setCreateDrawerOpen(open)}
			/>
		</>
	);
}
