import type { TransactionCreatedSessionHistoryEvent } from "@/api/types/session-history.ts";
import { UserAvatar } from "@/components/user-avatar.tsx";
import { cn } from "@/lib/utils.ts";
import type { SessionMember } from "@/api/types/session.ts";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { EventContainer } from "@/features/session/detail/session-history-card/event-container.tsx";

interface TransactionCreatedEventProps
	extends TransactionCreatedSessionHistoryEvent {
	creator: SessionMember | undefined;
	creatorAvatarData: AvatarData;
	members: SessionMember[];
}

export function TransactionCreatedEvent({
	creator,
	creatorAvatarData,
	members,
	...event
}: TransactionCreatedEventProps) {
	function stringifyMemberNames(usernames: string[]): string {
		if (usernames.length === members.length - 1) return "everyone";
		if (usernames.length === 1) return usernames[0];
		if (usernames.length === 2) return `${usernames[0]} and ${usernames[1]}`;
		return `${usernames.slice(0, -1).join(", ")}, and ${usernames.slice(-1)[0]}`;
	}

	return (
		<EventContainer>
			<UserAvatar data={creatorAvatarData} />
			<div className="grid grid-cols-5 grid-rows-2 w-full">
				<p
					className={cn(
						"col-span-4 font-semibold tracking-wider",
						creator?.isDeleted && "line-through",
					)}
				>
					{creator?.username ?? "unknown"}
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
		</EventContainer>
	);
}
