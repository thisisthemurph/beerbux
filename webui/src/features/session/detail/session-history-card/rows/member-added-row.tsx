import {
	BaseHistoryEventRow,
	type HistoryEventRow,
} from "@/features/session/detail/session-history-card/rows/base-row.tsx";

interface MemberAddedRowProps extends HistoryEventRow {
	actorUsername: string;
	addedMemberUsername: string;
}

export function MemberAddedRow({
	actorUsername,
	actorAvatarData,
	addedMemberUsername,
}: MemberAddedRowProps) {
	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div>
				<span className="font-semibold">{actorUsername}</span> added{" "}
				<span className="font-semibold">{addedMemberUsername}</span> to the
				session
			</div>
		</BaseHistoryEventRow>
	);
}
