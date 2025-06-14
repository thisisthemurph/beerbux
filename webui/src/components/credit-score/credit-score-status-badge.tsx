import { Badge } from "@/components/ui/badge.tsx";
import { cn } from "@/lib/utils.ts";
import {
	Scale as BalancedBrewerIcon,
	Trophy as RoundChampionIcon,
	Ghost as RoundDodgerIcon,
} from "lucide-react";
import type { CreditScoreStatus } from "./functions.ts";

function getStatusIcon(status: CreditScoreStatus) {
	switch (status) {
		case "Round Dodger":
			return <RoundDodgerIcon className="w-4 h-4" />;
		case "Balanced Brewer":
			return <BalancedBrewerIcon className="w-4 h-4" />;
		default:
			return <RoundChampionIcon className="w-4 h-4" />;
	}
}

export function CreditScoreStatusBadge({ status }: { status: CreditScoreStatus }) {
	return (
		<Badge
			className={cn(
				"flex gap-2 p-2 text-sm rounded-full shadow-xl border border-dashed",
				status === "Round Champion" && "bg-yellow-300 text-yellow-800 border-yellow-600",
				status === "Balanced Brewer" && "bg-green-300 text-green-800 border-green-600",
				status === "Round Dodger" && "bg-red-300 text-red-800 border-red-600",
			)}
		>
			<span>{status}</span>
			<span>{getStatusIcon(status)}</span>
		</Badge>
	);
}
