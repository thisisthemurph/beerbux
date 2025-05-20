import type { SessionHistoryEvent } from "@/api/types/session-history.ts";
import type { SessionMember } from "@/api/types/session.ts";
import {
	MemberAddedRow,
	MemberLeftRow,
	MemberRemovedRow,
	SessionClosedRow,
	SessionOpenedRow,
	TransactionCreatedRow,
} from "@/features/session/detail/session-history-card/rows";
import { MemberDemotedFromAdminRow } from "@/features/session/detail/session-history-card/rows/member-demoted-from-admin-row.tsx";
import { MemberPromotedToAdminRow } from "@/features/session/detail/session-history-card/rows/member-promoted-to-admin-row.tsx";
import type { UserAvatarData } from "@/stores/user-avatar-store.ts";

export function HistoryEventRow({
	event,
	members,
	avatarData,
}: {
	event: SessionHistoryEvent;
	members: SessionMember[];
	avatarData: Record<string, UserAvatarData>;
}) {
	const actor = members.find((m) => m.id === event.memberId);
	const actorUsername = actor?.username ?? "unknown";
	const actorAvatarData = avatarData[actorUsername];

	function getMemberUsername(memberId: string) {
		return members.find((m) => m.id === memberId)?.username ?? "someone";
	}

	if (actor === undefined) return null;

	switch (event.eventType) {
		case "transaction_created":
			return (
				<TransactionCreatedRow actor={actor} actorAvatarData={actorAvatarData} members={members} {...event} />
			);
		case "member_added": {
			return (
				<MemberAddedRow
					actorUsername={actorUsername}
					actorAvatarData={actorAvatarData}
					addedMemberUsername={getMemberUsername(event.eventData.memberId)}
				/>
			);
		}
		case "member_removed": {
			return (
				<MemberRemovedRow
					actorUsername={actorUsername}
					actorAvatarData={actorAvatarData}
					removedMemberUsername={getMemberUsername(event.eventData.memberId)}
				/>
			);
		}
		case "member_left":
			return <MemberLeftRow actorUsername={actorUsername} actorAvatarData={actorAvatarData} />;
		case "session_closed":
			return <SessionClosedRow actorUsername={actorUsername} actorAvatarData={actorAvatarData} />;
		case "session_opened":
			return <SessionOpenedRow actorUsername={actorUsername} actorAvatarData={actorAvatarData} />;
		case "promoted_to_admin":
			return (
				<MemberPromotedToAdminRow
					actorUsername={actorUsername}
					promotedMemberUsername={getMemberUsername(event.eventData.memberId)}
					actorAvatarData={actorAvatarData}
				/>
			);
		case "demoted_from_admin":
			return (
				<MemberDemotedFromAdminRow
					actorUsername={actorUsername}
					promotedMemberUsername={getMemberUsername(event.eventData.memberId)}
					actorAvatarData={actorAvatarData}
				/>
			);
		default:
			console.warn(`Unhandled event type: ${(event as { eventType: string }).eventType}`);
			return null;
	}
}
