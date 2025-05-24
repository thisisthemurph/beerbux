import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils.ts";
import type { UserAvatarData } from "@/stores/user-avatar-store.ts";

type UserAvatarVariant = "default" | "prominent";

export function UserAvatar({
	data,
	tooltip,
	variant,
}: { data: UserAvatarData; tooltip?: string; variant?: UserAvatarVariant }) {
	if (!tooltip) {
		return <InnerUserAvatar data={data} variant={variant ?? "default"} />;
	}

	return (
		<Tooltip>
			<TooltipTrigger>
				<InnerUserAvatar data={data} variant={variant ?? "default"} />
			</TooltipTrigger>
			<TooltipContent>
				<p>{tooltip}</p>
			</TooltipContent>
		</Tooltip>
	);
}

function InnerUserAvatar({ data, variant }: { data: UserAvatarData; variant: UserAvatarVariant }) {
	return (
		<Avatar
			className={cn(
				"size-10 transition-all duration-300 ease-in-out",
				variant === "prominent" &&
					"ring-4 ring-offset-4 dark:ring-offset-0 ring-primary/75 dark:ring-primary hover:ring-8 hover:ring-offset-0",
			)}
		>
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
