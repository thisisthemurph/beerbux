import type { Session } from "@/api/userClient.ts";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { Link } from "react-router";

type SessionListingProps = {
	sessions: Session[];
};

function SessionListing({ sessions }: SessionListingProps) {
	return (
		<Card>
			<CardHeader>
				<section className="flex justify-between items-center">
					<CardTitle>Your sessions</CardTitle>
					<p className="text-muted-foreground">
						{sessions.length > 0 ? sessions.length : "No"} sessions
					</p>
				</section>
			</CardHeader>
			<CardContent>
				<section className="flex flex-col">
					{sessions.length === 0 && <NoSessionsIndicator />}
					{sessions.map((session, i) => (
						<Link to={`/session/${session.id}`} key={session.id + i.toString()}>
							<div className="flex items-center gap-6 py-6">
								<Avatar className="w-10 h-10">
									<AvatarFallback>{getAvatarText(session.name)}</AvatarFallback>
								</Avatar>
								<p>{session.name}</p>
							</div>
							{i < sessions.length - 1 && <Separator />}
						</Link>
					))}
				</section>
			</CardContent>
			{sessions.length > 0 && (
				<CardFooter>
					<Link to="/sessions" className="text-blue-400">
						All sessions
					</Link>
				</CardFooter>
			)}
		</Card>
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
		return first[0] + last[0];
	}
	return name[0];
}

SessionListing.Skeleton = SkeletonCard;
export { SessionListing };
