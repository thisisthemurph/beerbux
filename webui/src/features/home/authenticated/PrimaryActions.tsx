import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { SquareChevronRight } from "lucide-react";
import { Link } from "react-router";

export function PrimaryActions() {
	return (
		<Card>
			<CardHeader className="sr-only">
				<CardTitle>Actions</CardTitle>
			</CardHeader>
			<CardContent>
				<section>
					<Link to="/sessions/create" className="flex items-center gap-4">
						<SquareChevronRight className="text-green-300 w-8 h-8" />
						<span>Start new session</span>
					</Link>
				</section>
			</CardContent>
		</Card>
	);
}
