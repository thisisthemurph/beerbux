import { useState } from "react";
import { useUserAvatarData } from "@/hooks/user-avatar-data.ts";
import { PageError } from "@/components/page-error.tsx";
import { SessionMenu } from "@/features/session/detail/session-menu.tsx";
import { OverviewCard } from "@/features/session/detail/overview-card.tsx";
import {
	PrimaryActionCard,
	PrimaryActionCardButtonItem,
	PrimaryActionCardContent,
} from "@/components/primary-action-card.tsx";
import { Beer, SquarePlus } from "lucide-react";
import { MemberDetailsCard } from "@/features/session/detail/member-details-card.tsx";
import { CreateTransactionDrawer } from "@/features/session/detail/create-transaction/create-transaction-drawer.tsx";
import type { Session, SessionTransaction } from "@/api/types/session.ts";
import type { TransactionMemberAmounts } from "@/api/types/transaction.ts";
import type { User } from "@/api/types/user.ts";
import { PageHeading } from "@/components/page-heading.tsx";
import { AddMemberDrawer } from "@/features/session/detail/add-member/add-member-drawer.tsx";
import type { SessionHistory } from "@/api/types/session-history.ts";
import { SessionHistoryCard } from "@/features/session/detail/session-history-card/session-history-card.tsx";

type SessionDetailContentProps = {
	session: Session;
	history: SessionHistory | undefined;
	user: User;
	onMemberAdminStateUpdate: (sessionId: string, memberId: string, newAdminState: boolean) => void;
	handleNewTransaction: (transaction: TransactionMemberAmounts) => Promise<void>;
	handleAddMember: (username: string) => Promise<void>;
	onLeaveSession: () => void;
	onChangeSessionActiveState: () => void;
	onRemoveMember: (memberId: string) => void;
};

function calculateAverage(transactions: SessionTransaction[]) {
	if (transactions.length === 0) return 0;
	const total = transactions.reduce((acc, transaction) => {
		return acc + transaction.total;
	}, 0);

	return total / transactions.length;
}

export function SessionDetailContent({
	session,
	history,
	user,
	onMemberAdminStateUpdate,
	handleNewTransaction,
	handleAddMember,
	onLeaveSession,
	onChangeSessionActiveState,
	onRemoveMember,
}: SessionDetailContentProps) {
	const [createDrawerOpen, setCreateDrawerOpen] = useState(false);
	const [addMemberDrawerOpen, setAddMemberDrawerOpen] = useState(false);
	const currentMember = session?.members.find((m) => m.id === user?.id);
	const otherSessionMembers = (session?.members ?? [])
		.filter((m) => m.id !== user.id)
		.sort((a, b) => a.name.localeCompare(b.name));

	const { avatarData } = useUserAvatarData(session?.members.map((m) => m.username) ?? []);

	if (!currentMember) {
		return <PageError message="You are not a member of this session." />;
	}

	return (
		<>
			<PageHeading title={session.name}>
				<SessionMenu
					{...session}
					showAdminActions={currentMember.isAdmin}
					onLeave={onLeaveSession}
					onChangeActiveState={onChangeSessionActiveState}
				/>
			</PageHeading>

			<OverviewCard {...session} average={calculateAverage(session.transactions)} />

			{session.isActive && (
				<PrimaryActionCard>
					<PrimaryActionCardContent>
						{currentMember.isAdmin && (
							<PrimaryActionCardButtonItem
								text="Add a member"
								icon={<SquarePlus className="text-primary w-8 h-8" />}
								onClick={() => setAddMemberDrawerOpen(true)}
							/>
						)}
						{session.members.length > 1 && (
							<PrimaryActionCardButtonItem
								text="Buy a round"
								icon={<Beer className="text-primary w-8 h-8" />}
								onClick={() => setCreateDrawerOpen(true)}
							/>
						)}
					</PrimaryActionCardContent>
				</PrimaryActionCard>
			)}

			<MemberDetailsCard
				showMemberDropdownMenu={currentMember.isAdmin}
				members={[currentMember, ...otherSessionMembers.filter((m) => !m.isDeleted)]}
				avatarData={avatarData}
				onChangeMemberAdminState={(memberId, newAdminState) =>
					onMemberAdminStateUpdate(session.id, memberId, newAdminState)
				}
				onRemoveMember={onRemoveMember}
			/>

			<SessionHistoryCard events={history?.events ?? []} members={session.members} avatarData={avatarData} />

			<CreateTransactionDrawer
				members={otherSessionMembers}
				onTransactionCreate={async (data) => {
					await handleNewTransaction(data);
					setCreateDrawerOpen(false);
				}}
				open={createDrawerOpen}
				onOpenChange={(open) => setCreateDrawerOpen(open)}
			/>

			<AddMemberDrawer
				onMemberAdd={async (username) => {
					await handleAddMember(username);
					setAddMemberDrawerOpen(false);
				}}
				open={addMemberDrawerOpen}
				onOpenChange={(open) => setAddMemberDrawerOpen(open)}
			/>
		</>
	);
}
