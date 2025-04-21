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
import { TransactionListing } from "@/features/session/detail/transaction-listing.tsx";
import { CreateTransactionDrawer } from "@/features/session/detail/create-transaction/create-transaction-drawer.tsx";
import type { Session, TransactionMemberAmounts, User } from "@/api/types.ts";
import { PageHeading } from "@/components/page-heading.tsx";
import { AddMemberDrawer } from "@/features/session/detail/add-member/add-member-drawer.tsx";

type SessionDetailContentProps = {
	session: Session;
	user: User;
	onMemberAdminStateUpdate: (
		sessionId: string,
		memberId: string,
		newAdminState: boolean,
	) => void;
	handleNewTransaction: (
		transaction: TransactionMemberAmounts,
	) => Promise<void>;
	handleAddMember: (username: string) => Promise<void>;
	onLeaveSession: () => void;
	onRemoveMember: (memberId: string) => void;
};

export function SessionDetailContent({
	session,
	user,
	onMemberAdminStateUpdate,
	handleNewTransaction,
	handleAddMember,
	onLeaveSession,
	onRemoveMember,
}: SessionDetailContentProps) {
	const [createDrawerOpen, setCreateDrawerOpen] = useState(false);
	const [addMemberDrawerOpen, setAddMemberDrawerOpen] = useState(false);
	const currentMember = session?.members.find((m) => m.id === user?.id);
	const otherSessionMembers = (session?.members ?? [])
		.filter((m) => m.id !== user.id)
		.sort((a, b) => a.name.localeCompare(b.name));

	const { avatarData } = useUserAvatarData(
		session?.members.map((m) => m.username) ?? [],
	);

	if (!currentMember) {
		return <PageError message="You are not a member of this session." />;
	}

	return (
		<>
			<PageHeading title={session.name}>
				<SessionMenu
					showAdminActions={currentMember.isAdmin}
					onLeave={onLeaveSession}
				/>
			</PageHeading>

			<OverviewCard total={session.total} />

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
				members={[
					currentMember,
					...otherSessionMembers.filter((m) => !m.isDeleted),
				]}
				avatarData={avatarData}
				onChangeMemberAdminState={(memberId, newAdminState) =>
					onMemberAdminStateUpdate(session.id, memberId, newAdminState)
				}
				onRemoveMember={onRemoveMember}
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
