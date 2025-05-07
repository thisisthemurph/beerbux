import type { SessionMember, SessionTransaction } from "@/api/types/session.ts";
import {
	Drawer,
	DrawerContent,
	DrawerDescription,
	DrawerFooter,
	DrawerHeader,
	DrawerTitle,
} from "@/components/ui/drawer.tsx";
import { UserAvatar } from "@/components/user-avatar.tsx";
import { useUserAvatarDataBySession } from "@/hooks/user-avatar-data.ts";
import { pluralize } from "@/lib/strings.ts";
import { cn } from "@/lib/utils.ts";
import type { AvatarData } from "@/stores/user-avatar-store.ts";
import type { DrawerToggleProps } from "@/types.ts";
import { ChevronDown, ChevronUp } from "lucide-react";

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
	const bought = transactions.filter((t) => t.creatorId === member.id);
	const received = transactions.filter((t) => t.members.some((m) => m.userId === member.id));
	const transactionTotalsByMember = calculateTransactionTotalsByMember(member.id, transactions);
	const avatarData = useUserAvatarDataBySession(sessionId);

	return (
		<Drawer open={open} onOpenChange={onOpenChange}>
			<DrawerContent>
				<DrawerHeader className="px-6">
					<DrawerTitle>{member.username}</DrawerTitle>
					<DrawerDescription className="text-lg">
						{member.username} has bought {bought.length} {pluralize(bought.length, "round", "rounds")} and has
						had {received.length} {pluralize(received.length, "round", "rounds")} bought for them.
					</DrawerDescription>
				</DrawerHeader>
				<DrawerFooter className="px-6">
					{otherMembers.map((m) => (
						<MemberBalanceRow
							key={m.id}
							username={m.username}
							amount={transactionTotalsByMember.get(m.id) ?? 0}
							avatarData={avatarData[m.username] ?? {}}
						/>
					))}
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	);
}

type MemberBalanceRowProps = {
	username: string;
	amount: number;
	avatarData: AvatarData;
};

function MemberBalanceRow({ avatarData, username, amount }: MemberBalanceRowProps) {
	const absAmount = amount < 0 ? amount * -1 : amount;
	const symbol = amount < 0 ? <ChevronDown /> : amount > 0 ? <ChevronUp /> : "";

	return (
		<div className="flex items-center justify-between font-semibold py-2">
			<div className="flex items-center gap-4">
				<UserAvatar data={avatarData} />
				<span className="tracking-wide">{username}</span>
			</div>
			<span
				className={cn(
					"flex items-center text-lg text-muted-foreground",
					amount < 0 && "text-red-500",
					amount > 0 && "text-green-500",
				)}
			>
				<span>{symbol}</span>
				{absAmount}
			</span>
		</div>
	);
}
