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
		<Card>
			<CardHeader className="sr-only">
				<CardTitle>Actions</CardTitle>
			</CardHeader>
			{children}
		</Card>
	);
}

function PrimaryActionCardContent({ children }: { children: ReactNode }) {
	return <CardContent>{children}</CardContent>;
}

function PrimaryActionContent({
	text,
	icon,
}: { text: string; icon: ReactNode }) {
	return (
		<div className="flex items-center gap-4 cursor-pointer">
			{icon}
			<span className="py-4">{text}</span>
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
		<Link to={to} key={text}>
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
		<button type="button" onClick={onClick}>
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
