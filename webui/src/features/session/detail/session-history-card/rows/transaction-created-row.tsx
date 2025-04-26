import type { TransactionCreatedSessionHistoryEvent } from "@/api/types/session-history.ts";
import { cn } from "@/lib/utils.ts";
import type { SessionMember } from "@/api/types/session.ts";
import { BaseHistoryEventRow, type HistoryEventRow } from "./base-row";

interface TransactionCreatedRowProps
	extends TransactionCreatedSessionHistoryEvent,
		HistoryEventRow {
	actor: SessionMember | undefined;
	members: SessionMember[];
}

export function TransactionCreatedRow({
	actor,
	actorAvatarData,
	members,
	...event
}: TransactionCreatedRowProps) {
	function stringifyMemberNames(usernames: string[]): string {
		if (usernames.length === members.length - 1) return "everyone";
		if (usernames.length === 1) return usernames[0];
		if (usernames.length === 2) return `${usernames[0]} and ${usernames[1]}`;
		return `${usernames.slice(0, -1).join(", ")}, and ${usernames.slice(-1)[0]}`;
	}

	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div className="grid grid-cols-5 grid-rows-2 w-full">
				<p
					className={cn(
						"col-span-4 font-semibold tracking-wider",
						actor?.isDeleted && "line-through",
					)}
				>
					{actor?.username ?? "unknown"}
				</p>
				<div className="row-span-2 flex items-center justify-end">
					<p className="font-semibold">
						${event.eventData.lines.reduce((sum, v) => sum + v.amount, 0)}
					</p>
				</div>
				<p className="col-span-4 text-muted-foreground">
					{stringifyMemberNames(
						event.eventData.lines.map(({ memberId }) => {
							return (
								members.find((m) => m.id === memberId)?.username ?? "unknown"
							);
						}),
					)}
				</p>
			</div>
		</BaseHistoryEventRow>
	);
}
