import type * as React from "react";

export function EventContainer({ children }: { children: React.ReactNode }) {
	return (
		<div className="flex items-center gap-4 px-6 py-4 hover:bg-muted transition-colors">
			{children}
		</div>
	);
}
