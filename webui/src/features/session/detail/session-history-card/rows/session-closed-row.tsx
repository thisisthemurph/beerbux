import {
	BaseHistoryEventRow,
	type HistoryEventRow,
} from "@/features/session/detail/session-history-card/rows/base-row.tsx";

interface SessionClosedRowProps extends HistoryEventRow {
	actorUsername: string;
}

export function SessionClosedRow({ actorUsername, actorAvatarData }: SessionClosedRowProps) {
	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div>
				<span className="font-semibold">{actorUsername}</span> closed the session
			</div>
		</BaseHistoryEventRow>
	);
}
