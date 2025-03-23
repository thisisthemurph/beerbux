import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import type { ReactNode } from "react";
import { Link } from "react-router";

type LinkAction = {
	text: string;
	icon: ReactNode;
	href: string;
	onClick?: undefined;
};

type ButtonAction = {
	text: string;
	icon: ReactNode;
	href?: undefined;
	onClick: () => void;
};

type ActionItem = LinkAction | ButtonAction;

type PrimaryActionsProps = {
	items: ActionItem[];
};

export function PrimaryActions({ items }: PrimaryActionsProps) {
	return (
		<Card>
			<CardHeader className="sr-only">
				<CardTitle>Actions</CardTitle>
			</CardHeader>
			<CardContent>
				<section>
					{items.map((item) => {
						if (item.href) {
							return (
								<Link
									to={item.href}
									key={item.text}
									className="flex items-center gap-4 cursor-pointer"
								>
									{item.icon}
									<span>{item.text}</span>
								</Link>
							);
						}

						return (
							<button
								type="button"
								onClick={item.onClick}
								key={item.text}
								className="flex items-center gap-4 cursor-pointer"
							>
								{item.icon}
								<span>{item.text}</span>
							</button>
						);
					})}
				</section>
			</CardContent>
		</Card>
	);
}
