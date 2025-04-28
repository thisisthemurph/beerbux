import {
	BaseHistoryEventRow,
	type HistoryEventRow,
} from "@/features/session/detail/session-history-card/rows/base-row.tsx";

interface MemberLeftRowProps extends HistoryEventRow {
	actorUsername: string;
}

export function MemberLeftRow({ actorUsername, actorAvatarData }: MemberLeftRowProps) {
	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div>
				<span className="font-semibold">{actorUsername}</span> left the session
			</div>
		</BaseHistoryEventRow>
	);
}
