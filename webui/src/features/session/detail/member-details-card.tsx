import type { SessionMember } from "@/api/types.ts";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { UserAvatar } from "@/components/user-avatar.tsx";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { cn } from "@/lib/utils";
import { useUserStore } from "@/stores/user-store.tsx";

type MemberDetailsCardProps = {
	members: SessionMember[];
	avatarData: Record<string, AvatarData>;
};

export function MemberDetailsCard({
	members,
	avatarData,
}: MemberDetailsCardProps) {
	const user = useUserStore((state) => state.user);
	const currentMember = members.find((member) => member.id === user?.id);

	return (
		<Card>
			<CardHeader>
				<CardTitle>Members</CardTitle>
			</CardHeader>
			<CardContent>
				{currentMember && (
					<div className="flex items-center gap-4">
						<UserAvatar
							data={avatarData[currentMember.username]}
							tooltip={currentMember.name}
						/>
						<div className="flex justify-between items-center w-full">
							<p className="py-6 font-semibold">{currentMember.username}</p>
							<Balance {...currentMember.transactionSummary} />
						</div>
					</div>
				)}
				{members
					.filter((m) => m.id !== currentMember?.id)
					.sort((a, b) => a.name.localeCompare(b.name))
					.map(({ id, name, username, transactionSummary }) => (
						<div key={id} className="flex items-center gap-4">
							<UserAvatar data={avatarData[username]} tooltip={name} />
							<div className="flex justify-between items-center w-full">
								<p className="py-6 font-semibold">{username}</p>
								<Balance {...transactionSummary} />
							</div>
						</div>
					))}
			</CardContent>
		</Card>
	);
}

function Balance({ credit, debit }: { credit: number; debit: number }) {
	const value = credit - debit;

	return (
		<p
			className={cn(
				"font-semibold text-muted-foreground",
				value > 0 && "text-green-600",
				value < 0 && "text-red-600",
			)}
		>
			${value < 0 ? value * -1 : value}
		</p>
	);
}
