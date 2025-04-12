import { Badge } from "@/components/ui/badge.tsx";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import type { ReactNode } from "react";

type OverviewCardProps = {
	total: number;
};

export function OverviewCard({ total }: OverviewCardProps) {
	return (
		<Card className="bg-green-400">
			<CardHeader className="sr-only">
				<CardTitle>Overview</CardTitle>
			</CardHeader>
			<CardContent>
				<div className="flex gap-12 min-h-24">
					<section className="flex flex-col justify-end w-2/3">
						<Row
							title="Avg"
							content={<Badge variant="secondary">${0}</Badge>}
						/>
					</section>

					<section className="w-full text-right">
						<p className="text-white font-semibold">
							<span className="text-2xl">$</span>
							<span className="text-4xl">{total}</span>
						</p>
						<p className="text-xs text-white tracking-wide font-mono font-semibold">
							total
						</p>
					</section>
				</div>
			</CardContent>
		</Card>
	);
}

function Row({ title, content }: { title: string; content: ReactNode }) {
	return (
		<div className="flex items-end justify-between w-full">
			<p className="font-semibold tracking-wider text-green-900">{title}</p>
			{content}
		</div>
	);
}

OverviewCard.Skeleton = function OverviewCardSkeleton() {
	return (
		<Card className="bg-green-400 animate-pulse">
			<CardHeader className="sr-only">
				<CardTitle>Overview</CardTitle>
			</CardHeader>
			<CardContent>
				<div className="flex gap-12 min-h-24">
					<section className="flex flex-col justify-end w-2/3">
						<Row title="Avg" content={<Badge variant="secondary">$0</Badge>} />
					</section>

					<section className="w-full text-right">
						<p className="text-white font-semibold">
							<span className="text-2xl">$</span>
							<span className="text-4xl">0</span>
						</p>
						<p className="text-xs text-white tracking-wide font-mono font-semibold">
							total
						</p>
					</section>
				</div>
			</CardContent>
		</Card>
	);
};
