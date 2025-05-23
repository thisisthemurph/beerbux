import { UserAvatar } from "@/components/user-avatar.tsx";
import type { UserAvatarData } from "@/stores/user-avatar-store.ts";
import type { ReactNode } from "react";

export interface HistoryEventRow {
	actorAvatarData: UserAvatarData;
}

interface EventHistoryBaseRowProps extends HistoryEventRow {
	children: ReactNode;
}

export function BaseHistoryEventRow({ actorAvatarData, children }: EventHistoryBaseRowProps) {
	return (
		<div className="flex items-center gap-4 px-6 py-4 tracking-wide hover:bg-muted transition-colors">
			<UserAvatar data={actorAvatarData} />
			{children}
		</div>
	);
}
