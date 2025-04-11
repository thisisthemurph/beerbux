import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import type { ReactNode } from "react";
import { Link } from "react-router";

function PrimaryActionCard({ children }: { children: ReactNode }) {
	return (
		<Card className="py-0">
			<CardHeader className="sr-only">
				<CardTitle>Actions</CardTitle>
			</CardHeader>
			{children}
		</Card>
	);
}

function PrimaryActionCardContent({ children }: { children: ReactNode }) {
	return <CardContent className="px-0">{children}</CardContent>;
}

function PrimaryActionContent({
	text,
	icon,
}: { text: string; icon: ReactNode }) {
	return (
		<div className="flex items-center gap-4 px-6 cursor-pointer">
			{icon}
			<span className="">{text}</span>
		</div>
	);
}

type PrimaryActionCardLinkItemProps = {
	to: string;
	text: string;
	icon: ReactNode;
};

function PrimaryActionCardLinkItem({
	to,
	text,
	icon,
}: PrimaryActionCardLinkItemProps) {
	return (
		<Link
			to={to}
			key={text}
			className="block w-full pb-4 pt-4 hover:bg-muted first:rounded-t-lg last:rounded-b-lg transition-colors"
		>
			<PrimaryActionContent text={text} icon={icon} />
		</Link>
	);
}

type PrimaryActionCardButtonItemProps = {
	text: string;
	icon: ReactNode;
	onClick: () => void;
};

function PrimaryActionCardButtonItem({
	text,
	icon,
	onClick,
}: PrimaryActionCardButtonItemProps) {
	return (
		<button
			type="button"
			onClick={onClick}
			className="w-full pb-4 pt-4 hover:bg-muted first:rounded-t-lg last:rounded-b-lg transition-colors cursor-pointer"
		>
			<PrimaryActionContent text={text} icon={icon} />
		</button>
	);
}

function PrimaryActionCardSeparator() {
	return <Separator />;
}

export {
	PrimaryActionCard,
	PrimaryActionCardContent,
	PrimaryActionCardLinkItem,
	PrimaryActionCardButtonItem,
	PrimaryActionCardSeparator,
};
