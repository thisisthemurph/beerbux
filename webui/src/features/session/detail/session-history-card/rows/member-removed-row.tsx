import {
	BaseHistoryEventRow,
	type HistoryEventRow,
} from "@/features/session/detail/session-history-card/rows/base-row.tsx";

interface MemberRemovedEventProps extends HistoryEventRow {
	actorUsername: string;
	removedMemberUsername: string;
}

export function MemberRemovedRow({
	actorUsername,
	actorAvatarData,
	removedMemberUsername,
}: MemberRemovedEventProps) {
	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div>
				<span className="font-semibold">{actorUsername}</span> removed{" "}
				<span className="font-semibold">{removedMemberUsername}</span> from the
				session
			</div>
		</BaseHistoryEventRow>
	);
}
