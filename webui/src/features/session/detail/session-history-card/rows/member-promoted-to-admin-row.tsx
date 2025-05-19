import {
	BaseHistoryEventRow,
	type HistoryEventRow,
} from "@/features/session/detail/session-history-card/rows/base-row.tsx";

interface MemberPromotedToAdminRowProps extends HistoryEventRow {
	actorUsername: string;
	promotedMemberUsername: string;
}

export function MemberPromotedToAdminRow({
	actorUsername,
	actorAvatarData,
	promotedMemberUsername,
}: MemberPromotedToAdminRowProps) {
	return (
		<BaseHistoryEventRow actorAvatarData={actorAvatarData}>
			<div>
				<span className="font-semibold">{actorUsername}</span> promoted{" "}
				<span className="font-semibold">{promotedMemberUsername}</span> to admin
			</div>
		</BaseHistoryEventRow>
	);
}
