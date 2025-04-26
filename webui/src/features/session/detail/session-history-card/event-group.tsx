import type { SessionHistoryEvent } from "@/api/types/session-history.ts";
import type { SessionMember } from "@/api/types/session.ts";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { GroupLabel } from "@/features/session/detail/session-history-card/group-label.tsx";
import { HistoryEventRow } from "@/features/session/detail/session-history-card/rows/history-event-row.tsx";

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
			{events.map((event) => (
				<HistoryEventRow
					key={event.id}
					event={event}
					members={members}
					avatarData={avatarData}
				/>
			))}
		</>
	);
}
