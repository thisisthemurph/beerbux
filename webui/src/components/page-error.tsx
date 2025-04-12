import { BeerOff } from "lucide-react";

export function PageError({ message }: { message?: string }) {
	return (
		<div className="flex flex-col gap-8 mt-18 h-full w-full items-center justify-center">
			<BeerOff className="size-36 text-muted-foreground" />
			<p className="px-4 text-2xl text-center tracking-wide">
				{message ?? "There has been an issue."}
			</p>
		</div>
	);
}
