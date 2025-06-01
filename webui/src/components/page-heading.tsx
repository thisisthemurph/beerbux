import { cn } from "@/lib/utils.ts";
import type { ReactNode } from "react";

type PageHeadingProps = {
	title: string;
	children?: ReactNode;
	className?: string;
};

export function PageHeading({ title, children, className }: PageHeadingProps) {
	return (
		<div className={cn("flex justify-between items-center mb-8", className)}>
			<h1 className="mb-0">{title}</h1>
			{children}
		</div>
	);
}
