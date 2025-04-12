import { CardTitle } from "@/components/ui/card.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { OverviewCard } from "@/features/session/detail/overview-card.tsx";
import { SessionMenu } from "@/features/session/detail/session-menu.tsx";

export function SessionDetailSkeleton() {
	return (
		<>
			<div className="flex justify-between items-center mb-8">
				<h1 className="mb-0">Loading session</h1>
				<SessionMenu.Skeleton />
			</div>
			<OverviewCard.Skeleton />
			<Skeleton className="h-[125px] rounded-xl p-6">
				<p className="sr-only">Loading</p>
			</Skeleton>
			<Skeleton className="h-58 rounded-xl p-6">
				<CardTitle className="animate-pulse">Members</CardTitle>
			</Skeleton>
			<Skeleton className="h-58 rounded-xl p-6">
				<CardTitle className="animate-pulse">Rounds</CardTitle>
			</Skeleton>
		</>
	);
}
