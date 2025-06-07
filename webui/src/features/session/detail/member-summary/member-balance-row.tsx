import { UserAvatar } from "@/components/user-avatar.tsx";
import { cn } from "@/lib/utils.ts";
import type { UserAvatarData } from "@/stores/user-avatar-store.ts";

type MemberBalanceRowProps = {
	username: string;
	amount: number;
	userAvatarData: UserAvatarData;
};

export function MemberBalanceRow({ username, amount, userAvatarData }: MemberBalanceRowProps) {
	const absAmount = amount < 0 ? amount * -1 : amount;

	return (
		<div className="flex items-center justify-between font-semibold py-2">
			<div className="flex items-center gap-4">
				<UserAvatar data={userAvatarData} />
				<span className="text-xl tracking-wider">{username}</span>
			</div>
			<span className={cn("flex justify-between items-end gap-2 text-2xl text-muted-foreground")}>
				{amount !== 0 ? <>{amount < 0 ? "is owed" : "owes"}</> : "all square"}
				{absAmount !== 0 && (
					<span className={cn("text-3xl", amount < 0 && "text-red-500", amount > 0 && "text-green-500")}>
						{absAmount}
					</span>
				)}
			</span>
		</div>
	);
}
