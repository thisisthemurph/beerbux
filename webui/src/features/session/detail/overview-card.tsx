import { Badge } from "@/components/ui/badge.tsx";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { ShieldOff } from "lucide-react";

type OverviewCardProps = {
	name: string;
	isActive: boolean;
	total: number;
	average: number;
};

export function OverviewCard({
	name,
	isActive,
	total,
	average,
}: OverviewCardProps) {
	return (
		<Card className="bg-primary">
			<CardHeader className="sr-only">
				<CardTitle>Overview card</CardTitle>
			</CardHeader>
			<CardContent>
				<div className="flex flex-col justify-between min-h-24">
					<section className="flex justify-between items-start">
						<p className="font-semibold tracking-wider font-mono">
							<span className="flex gap-2">
								{isActive ? (
									name
								) : (
									<>
										<ShieldOff />
										<span>Session closed</span>
									</>
								)}
							</span>
						</p>
						<div className="text-right">
							<p className="font-semibold">
								<span className="text-2xl">$</span>
								<span className="text-4xl">{total}</span>
							</p>
							<p className="text-xs tracking-wide font-mono font-semibold">
								total
							</p>
						</div>
					</section>

					<section className="flex gap-1">
						{average > 0 && (
							<Badge variant="secondary" className="text-sm">
								Average ${average.toFixed(2)}
							</Badge>
						)}
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
