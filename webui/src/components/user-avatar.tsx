import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";

type UserAvatarVariant = "default" | "prominent";

export function UserAvatar({
	data,
	tooltip,
	variant,
}: { data: AvatarData; tooltip?: string; variant?: UserAvatarVariant }) {
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

function InnerUserAvatar({
	data,
	variant,
}: { data: AvatarData; variant: UserAvatarVariant }) {
	return (
		<Avatar className="size-10">
			<AvatarFallback
				className="border-2 border-dotted"
				style={{
					backgroundColor: data.color,
					borderColor:
						data.color && variant === "prominent"
							? darkenHex(data.color, 40)
							: data.color,
				}}
			>
				{data.initial ?? "?"}
			</AvatarFallback>
		</Avatar>
	);
}

function darkenHex(hex: string, percent: number) {
	const hexValue = hex.replace(/^#/, "");

	let r = Number.parseInt(hexValue.substring(0, 2), 16);
	let g = Number.parseInt(hexValue.substring(2, 4), 16);
	let b = Number.parseInt(hexValue.substring(4, 6), 16);

	r = Math.max(0, Math.floor(r * (1 - percent / 100)));
	g = Math.max(0, Math.floor(g * (1 - percent / 100)));
	b = Math.max(0, Math.floor(b * (1 - percent / 100)));

	const toHex = (c: number) => c.toString(16).padStart(2, "0");
	return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
}
