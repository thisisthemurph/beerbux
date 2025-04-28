import { Button } from "@/components/ui/button.tsx";
import { Info } from "lucide-react";

export function InformationButton({ onClick }: { onClick: () => void }) {
	return (
		<Button variant="ghost" size="icon" className="group translate-x-2 p-0" onClick={onClick}>
			<Info className="size-5 text-muted-foreground/50 group-hover:text-blue-400 transition-colors" />
		</Button>
	);
}
