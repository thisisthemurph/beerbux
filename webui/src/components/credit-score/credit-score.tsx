import { type CreditScoreStatus, getCreditScoreStatus } from "@/components/credit-score/functions.ts";

type CreditScoreProps = {
	value: number;
};

export function CreditScore({ value }: CreditScoreProps) {
	const clampedValue = Math.max(0, Math.min(100, value));
	const creditScoreStatus = getCreditScoreStatus(clampedValue);

	const size = 200;
	const strokeWidth = 24;
	const radius = (size - strokeWidth) / 2;
	const circumference = 2 * Math.PI * radius;
	const offset = circumference - (clampedValue / 100) * circumference - strokeWidth;

	const getColor = (status: CreditScoreStatus) => {
		switch (status) {
			case "Round Dodger":
				return "#f43f5e"; // red
			case "Balanced Brewer":
				return "#facc15"; // yellow
			case "Round Champion":
				return "#22c55e"; // green
			default:
				return "#e5e7eb"; // default gray
		}
	};

	return (
		<div className="relative flex flex-col items-center justify-center">
			<svg width={size} height={size} className="-rotate-90">
				<title>Credit score indicator</title>
				<circle
					cx={size / 2}
					cy={size / 2}
					r={radius}
					fill="none"
					stroke="#e5e7eb"
					strokeWidth={strokeWidth}
				/>
				<circle
					cx={size / 2}
					cy={size / 2}
					r={radius}
					fill="none"
					stroke={getColor(creditScoreStatus)}
					strokeWidth={strokeWidth}
					strokeDasharray={circumference}
					strokeDashoffset={offset}
					strokeLinecap="round"
				/>
			</svg>

			<div className="absolute text-center">
				<div className="text-4xl font-semibold mt-2">{Math.round(clampedValue * 9)}</div>
				<div className="text-muted-foreground text-sm tracking-wide">{creditScoreStatus}</div>
			</div>
		</div>
	);
}
