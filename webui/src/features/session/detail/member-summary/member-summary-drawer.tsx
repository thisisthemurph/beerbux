import type { SessionMember, SessionTransaction } from "@/api/types/session.ts";
import {
	Drawer,
	DrawerContent,
	DrawerDescription,
	DrawerFooter,
	DrawerHeader,
	DrawerTitle,
} from "@/components/ui/drawer.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { MemberBalanceRow } from "@/features/session/detail/member-summary/member-balance-row.tsx";
import { MemberRoundStats } from "@/features/session/detail/member-summary/member-round-stats.tsx";
import { useUserAvatarDataBySession } from "@/hooks/user-avatar-data.ts";
import { pluralize } from "@/lib/strings.ts";
import type { DrawerToggleProps } from "@/types.ts";

interface MemberSummaryDrawerProps extends DrawerToggleProps {
	sessionId: string;
	member: SessionMember;
	members: SessionMember[];
	transactions: SessionTransaction[];
}

function calculateTransactionTotalsByMember(memberId: string, transactions: SessionTransaction[]) {
	const totalsByMember = new Map<string, number>();

	for (const t of transactions) {
		if (t.creatorId === memberId) {
			// Transaction is created by this member
			for (const m of t.members) {
				totalsByMember.set(m.userId, (totalsByMember.get(m.userId) ?? 0) + m.amount);
			}
		} else {
			// Transactions this member is involved in
			for (const m of t.members) {
				if (m.userId === memberId) {
					totalsByMember.set(t.creatorId, (totalsByMember.get(t.creatorId) ?? 0) - m.amount);
					break;
				}
			}
		}
	}

	return totalsByMember;
}

export function MemberSummaryDrawer({
	open,
	onOpenChange,
	sessionId,
	member,
	members,
	transactions,
}: MemberSummaryDrawerProps) {
	const otherMembers = members.filter((m) => m.id !== member.id);
	const roundsBought = transactions.filter((t) => t.creatorId === member.id).length;
	const roundsReceived = transactions.filter((t) => t.members.some((m) => m.userId === member.id)).length;
	const transactionTotalsByMember = calculateTransactionTotalsByMember(member.id, transactions);
	const avatarData = useUserAvatarDataBySession(sessionId);

	return (
		<Drawer open={open} onOpenChange={onOpenChange}>
			<DrawerContent>
				<DrawerHeader className="px-6">
					<DrawerTitle className="sr-only">{member.username}</DrawerTitle>
					<DrawerDescription className="sr-only">
						{member.username} has bought {roundsBought} {pluralize(roundsBought, "round", "rounds")} and has
						had {roundsReceived} {pluralize(roundsReceived, "round", "rounds")} bought for them.
					</DrawerDescription>
					<MemberRoundStats roundsBought={roundsBought} roundsReceived={roundsReceived} />
				</DrawerHeader>
				<Separator className="mb-4" />
				<DrawerFooter className="px-6">
					<h2>Beers owed to or by {member.username}</h2>
					{otherMembers.map((m) => (
						<MemberBalanceRow
							key={m.id}
							username={m.username}
							amount={transactionTotalsByMember.get(m.id) ?? 0}
							userAvatarData={avatarData[m.username] ?? {}}
						/>
					))}
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	);
}
