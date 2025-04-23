import type { SessionHistoryEvent } from "@/api/types/session-history.ts";
import type { SessionMember } from "@/api/types/session.ts";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { GroupLabel } from "@/features/session/detail/session-history-card/group-label.tsx";
import { TransactionCreatedEvent } from "@/features/session/detail/session-history-card/transaction-created-event.tsx";

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
				const creator = members.find((m) => m.id === event.memberId);
				const creatorUsername = creator?.username ?? "unknown";
				const creatorAvatarData = avatarData[creatorUsername];

				if (event.eventType === "transaction_created") {
					return (
						<TransactionCreatedEvent
							creator={creator}
							creatorAvatarData={creatorAvatarData}
							members={members}
							{...event}
						/>
					);
				}
				return null;
			})}
		</>
	);
}
