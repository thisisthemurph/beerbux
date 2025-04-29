import {
	BaseHistoryEventRow,
	type HistoryEventRow,
} from "@/features/session/detail/session-history-card/rows/base-row.tsx";

interface SessionOpenedRowProps extends HistoryEventRow {
	actorUsername: string;
}

export function SessionOpenedRow({ actorUsername, actorAvatarData }: SessionOpenedRowProps) {
	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div>
				<span className="font-semibold">{actorUsername}</span> opened the session
			</div>
		</BaseHistoryEventRow>
	);
}
