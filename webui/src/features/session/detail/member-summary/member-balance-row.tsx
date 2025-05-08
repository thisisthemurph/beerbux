import { UserAvatar } from "@/components/user-avatar.tsx";
import { cn } from "@/lib/utils.ts";
import type { UserAvatarData } from "@/stores/user-avatar-store.ts";
import { Check, ChevronDown, ChevronUp } from "lucide-react";

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
			<span className={cn("flex justify-between items-center text-xl text-muted-foreground w-1/3")}>
				<span>
					{amount < 0 ? (
						<ChevronDown className="size-8 text-red-500" />
					) : amount > 0 ? (
						<ChevronUp className="size-8 text-green-500" />
					) : (
						<Check className="size-8 text-green-500" />
					)}
				</span>
				<span>
					{amount !== 0 && <span className="font-normal">{amount < 0 ? "owes" : "is owed"} </span>}
					<span className={cn("", amount < 0 && "text-red-500", amount > 0 && "text-green-500")}>
						{absAmount === 0 ? "all square" : absAmount}
					</span>
				</span>
			</span>
		</div>
	);
}
