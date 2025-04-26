import type { SessionHistoryEvent } from "@/api/types/session-history.ts";
import type { SessionMember } from "@/api/types/session.ts";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { GroupLabel } from "@/features/session/detail/session-history-card/group-label.tsx";
import { TransactionCreatedRow } from "@/features/session/detail/session-history-card/rows/transaction-created-row.tsx";
import { MemberRemovedRow } from "@/features/session/detail/session-history-card/rows/member-removed-row.tsx";
import { MemberLeftRow } from "@/features/session/detail/session-history-card/rows/member-left-row.tsx";
import { MemberAddedRow } from "@/features/session/detail/session-history-card/rows/member-added-row.tsx";

type EventGroupProps = {
	label: string;
	events: SessionHistoryEvent[];
	members: SessionMember[];
	avatarData: Record<string, AvatarData>;
};

export function EventGroup({
	label,
	events,
	members,
	avatarData,
}: EventGroupProps) {
	return (
		<>
			<GroupLabel text={label} />
			{events.map((event) => {
				const actor = members.find((m) => m.id === event.memberId);
				const actorUsername = actor?.username ?? "unknown";
				const actorAvatarData = avatarData[actorUsername];

				switch (event.eventType) {
					case "transaction_created":
						return (
							<TransactionCreatedRow
								key={event.id}
								actor={actor}
								actorAvatarData={actorAvatarData}
								members={members}
								{...event}
							/>
						);
					case "member_added": {
						return (
							<MemberAddedRow
								key={event.id}
								actorUsername={actorUsername}
								actorAvatarData={actorAvatarData}
								addedMemberUsername={
									members.find((m) => m.id === event.eventData.memberId)
										?.username ?? "someone"
								}
							/>
						);
					}
					case "member_removed": {
						return (
							<MemberRemovedRow
								key={event.id}
								actorUsername={actorUsername}
								actorAvatarData={actorAvatarData}
								removedMemberUsername={
									members.find((m) => m.id === event.eventData.memberId)
										?.username ?? "someone"
								}
							/>
						);
					}
					case "member_left":
						return (
							<MemberLeftRow
								key={event.id}
								actorUsername={actorUsername}
								actorAvatarData={actorAvatarData}
							/>
						);
					default:
						console.warn(
							`Unhandled event type: ${(event as { eventType: string }).eventType}`,
						);
						return null;
				}
			})}
		</>
	);
}
