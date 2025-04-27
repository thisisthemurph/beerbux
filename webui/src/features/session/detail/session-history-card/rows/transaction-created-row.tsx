import type { TransactionCreatedSessionHistoryEvent } from "@/api/types/session-history.ts";
import { cn } from "@/lib/utils.ts";
import type { SessionMember } from "@/api/types/session.ts";
import { BaseHistoryEventRow, type HistoryEventRow } from "./base-row";
import { Username, UsernameGroup } from "@/components/username.tsx";

interface TransactionCreatedRowProps
	extends TransactionCreatedSessionHistoryEvent,
		HistoryEventRow {
	actor: SessionMember;
	members: SessionMember[];
}

export function TransactionCreatedRow({
	actor,
	actorAvatarData,
	members,
	...event
}: TransactionCreatedRowProps) {
	const totalAmount = event.eventData.lines.reduce(
		(sum, v) => sum + v.amount,
		0,
	);
	const memberUsernames = event.eventData.lines.map(
		({ memberId }) =>
			members.find((m) => m.id === memberId)?.username ?? "unknown",
	);

	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div className="flex justify-between gap-4 w-full">
				<p
					className={cn(
						"tracking-wider",
						actor?.isDeleted && "line-through",
					)}
				>
					<Username {...actor} /> bought a round for{" "}
					<UsernameGroup
						maxMembers={members.length-1}
						usernames={memberUsernames}
					/>
				</p>
				<p className="flex items-center justify-end font-semibold">
					${totalAmount}
				</p>
			</div>
		</BaseHistoryEventRow>
	);
}
