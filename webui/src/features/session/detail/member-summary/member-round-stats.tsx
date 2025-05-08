import { pluralize } from "@/lib/strings.ts";

type MemberRoundStatsProps = {
	roundsBought: number;
	roundsReceived: number;
};

export function MemberRoundStats({ roundsBought, roundsReceived }: MemberRoundStatsProps) {
	const balance = roundsBought - roundsReceived;

	return (
		<section className="flex mb-4">
			<div className="flex flex-col items-center w-full">
				<span className="text-4xl font-semibold">{roundsBought}</span>
				<span className="text-muted-foreground tracking-wide">
					{pluralize(roundsBought, "round", "rounds")} bought
				</span>
			</div>
			<div className="flex flex-col items-center w-full">
				<span className="text-4xl font-semibold">{roundsReceived}</span>
				<span className="text-muted-foreground tracking-wide">
					{pluralize(roundsReceived, "round", "rounds")} received
				</span>
			</div>
			<div className="flex flex-col items-center w-full">
				<span className="text-4xl font-semibold">{balance}</span>
				<span className="text-muted-foreground tracking-wide">balance</span>
			</div>
		</section>
	);
}
