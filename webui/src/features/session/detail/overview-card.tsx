import { Badge } from "@/components/ui/badge.tsx";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";

type OverviewCardProps = {
	total: number;
};

export function OverviewCard({ total }: OverviewCardProps) {
	return (
		<Card className="bg-primary">
			<CardHeader className="sr-only">
				<CardTitle>Overview card</CardTitle>
			</CardHeader>
			<CardContent>
				<div className="flex gap-12 min-h-24">
					<section className="flex flex-col justify-end w-2/3">
						<Badge variant="secondary" className="text-sm">
							Average ${0}
						</Badge>
					</section>

					<section className="w-full text-right">
						<p className="font-semibold">
							<span className="text-2xl">$</span>
							<span className="text-4xl">{total}</span>
						</p>
						<p className="text-xs tracking-wide font-mono font-semibold">
							total
						</p>
					</section>
				</div>
			</CardContent>
		</Card>
	);
}

OverviewCard.Skeleton = function OverviewCardSkeleton() {
	return (
		<Card className="bg-primary animate-pulse">
			<CardHeader className="sr-only">
				<CardTitle>Overview</CardTitle>
			</CardHeader>
			<CardContent>
				<div className="flex gap-12 min-h-24">
					<section className="w-full text-right">
						<p className="font-semibold">
							<span className="text-2xl">$</span>
							<span className="text-4xl">0</span>
						</p>
						<p className="text-xs tracking-wide font-mono font-semibold">
							total
						</p>
					</section>
				</div>
			</CardContent>
		</Card>
	);
};
