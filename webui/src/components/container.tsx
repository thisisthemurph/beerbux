import type { ReactNode } from "react";

interface PendingComponentProps {
	isPending: boolean;
	pendingComponent: ReactNode;
	children: ReactNode;
}

interface NonPendingComponentProps {
	isPending?: never;
	pendingComponent?: never;
	children: ReactNode;
}

type ContainerProps = PendingComponentProps | NonPendingComponentProps;

export function Container({
	children,
	isPending,
	pendingComponent,
}: ContainerProps) {
	return (
		<section className="space-y-6">
			{isPending ? pendingComponent : children}
		</section>
	);
}
