import type { Session } from "@/api/types.ts";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { withBackLinkOverride } from "@/hooks/use-back-navigation.ts";
import { ShieldOff } from "lucide-react";
import type { ReactNode } from "react";
import { Link } from "react-router";

type SessionListingProps = {
	title?: string;
	sessions: Session[];
	children?: ReactNode;
	parentPath?: string;
};

function SessionListing({
	title,
	sessions,
	children,
	parentPath,
}: SessionListingProps) {
	return (
		<Card>
			<CardHeader>
				<section className="flex justify-between items-center">
					<CardTitle>{title ?? "Your sessions"}</CardTitle>
					<p className="text-muted-foreground">
						{sessions.length > 0 ? sessions.length : "No"} sessions
					</p>
				</section>
			</CardHeader>
			<CardContent className="px-0">
				<section className="flex flex-col">
					{sessions.length === 0 && <NoSessionsIndicator />}
					{sessions.map((session) => {
						const url = withBackLinkOverride(
							`/session/${session.id}`,
							parentPath,
						);

						return (
							<SessionListingItem
								key={session.id}
								session={session}
								url={url}
							/>
						);
					})}
				</section>
			</CardContent>
			{children && <CardFooter>{children}</CardFooter>}
		</Card>
	);
}

function SessionListingItem({
	session,
	url,
}: { session: Session; url: string }) {
	return (
		<Link
			to={url}
			key={session.id}
			className="group hover:bg-muted transition-colors"
		>
			<div className="flex items-center gap-6 py-4 px-6">
				<Avatar className="w-10 h-10">
					<AvatarFallback className="group-hover:bg-card transition-colors">
						{getAvatarText(session.name)}
					</AvatarFallback>
				</Avatar>
				<div className="flex justify-between items-center w-full">
					<p>{session.name}</p>
					{!session.isActive && <InactiveIcon />}
				</div>
				<p className="text-xl text-muted-foreground font-semibold">
					${session.total}
				</p>
			</div>
		</Link>
	);
}

function InactiveIcon() {
	return (
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger>
					<ShieldOff className="text-muted-foreground" />
				</TooltipTrigger>
				<TooltipContent>
					<p>Inactive session</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>
	);
}

function NoSessionsIndicator() {
	return (
		<p className="text-center py-8 font-semibold text-lg tracking-wide">
			You don't have any sessions yet.
			<br /> Create one to get started!
		</p>
	);
}

function SkeletonCard() {
	return (
		<Skeleton className="h-[125px] rounded-xl p-6">
			<section className="flex justify-between items-center">
				<CardTitle className="text-muted-foreground">Your sessions</CardTitle>
				<p className="text-muted-foreground">loading</p>
			</section>
		</Skeleton>
	);
}

function getAvatarText(name: string) {
	if (name.length === 0) return "S";
	if (name.split(" ").length > 1) {
		const [first, last] = name.split(" ");
		return (first[0] + last[0]).toUpperCase();
	}
	return name[0].toUpperCase();
}

SessionListing.Skeleton = SkeletonCard;
export { SessionListing };
