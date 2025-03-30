import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";

export function UserAvatar({
	data,
	tooltip,
}: { data: AvatarData; tooltip?: string }) {
	if (!tooltip) {
		return <InnerUserAvatar data={data} />;
	}

	return (
		<Tooltip>
			<TooltipTrigger>
				<InnerUserAvatar data={data} />
			</TooltipTrigger>
			<TooltipContent>
				<p>{tooltip}</p>
			</TooltipContent>
		</Tooltip>
	);
}

function InnerUserAvatar({ data }: { data: AvatarData }) {
	return (
		<Avatar className="size-10">
			<AvatarFallback
				style={{
					backgroundColor: data.color,
				}}
			>
				{data.initial ?? "?"}
			</AvatarFallback>
		</Avatar>
	);
}
