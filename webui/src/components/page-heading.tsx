import type { ReactNode } from "react";

export function PageHeading({
	title,
	children,
}: { title: string; children?: ReactNode }) {
	return (
		<div className="flex justify-between items-center mb-8">
			<h1 className="mb-0">{title}</h1>
			{children}
		</div>
	);
}
