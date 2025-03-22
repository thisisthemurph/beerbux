import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import type { ReactNode } from "react";
import { Link } from "react-router";

type ActionLink = {
	text: string;
	icon: ReactNode;
	href: string;
};

type PrimaryActionsProps = {
	items: ActionLink[];
};

export function PrimaryActions({ items }: PrimaryActionsProps) {
	return (
		<Card>
			<CardHeader className="sr-only">
				<CardTitle>Actions</CardTitle>
			</CardHeader>
			<CardContent>
				<section>
					{items.map((x) => (
						<Link to={x.href} key={x.href} className="flex items-center gap-4">
							{x.icon}
							<span>{x.text}</span>
						</Link>
					))}
				</section>
			</CardContent>
		</Card>
	);
}
