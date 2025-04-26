import type { SessionHistoryEvent } from "@/api/types/session-history.ts";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible.tsx";
import { useState } from "react";
import type { SessionMember } from "@/api/types/session.ts";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { Button } from "@/components/ui/button.tsx";
import { ChevronDown, ChevronUp } from "lucide-react";
import { EventGroup } from "./event-group";
import { useGroupedEvents } from "@/features/session/detail/hooks/use-grouped-events.ts";

type SessionHistoryCardProps = {
	events: SessionHistoryEvent[];
	members: SessionMember[];
	avatarData: Record<string, AvatarData>;
};

export function SessionHistoryCard({
	events,
	members,
	avatarData,
}: SessionHistoryCardProps) {
	const [isOpen, setIsOpen] = useState(false);
	const grouped = useGroupedEvents(events);
	const showCollapsibleTrigger =
		grouped.sortedLabels.length > 1 ||
		(grouped.firstEvents && grouped.firstEvents.length > 5);

	return (
		<Collapsible open={isOpen} onOpenChange={setIsOpen}>
			<Card>
				<CardHeader>
					<section className="flex justify-between items-center">
						<CardTitle>History</CardTitle>
						{showCollapsibleTrigger && (
							<CollapsibleTrigger asChild>
								<Button variant="secondary">
									<span>{isOpen ? "See less" : "See more"}</span>
									{isOpen ? <ChevronUp /> : <ChevronDown />}
								</Button>
							</CollapsibleTrigger>
						)}
					</section>
				</CardHeader>
				<CardContent className="px-0">
					{members.length <= 1 && <NoMembers />}
					{members.length > 1 && events.length === 0 && <NoEventsMessage />}

					{grouped.firstEvents.length > 0 && (
						<EventGroup
							label={grouped.firstLabel}
							events={
								isOpen ? grouped.firstEvents : grouped.firstEvents.slice(0, 5)
							}
							members={members}
							avatarData={avatarData}
						/>
					)}

					<CollapsibleContent>
						{grouped.sortedLabels.slice(1).map((label) => (
							<EventGroup
								key={label}
								label={label}
								events={grouped.events[label]}
								members={members}
								avatarData={avatarData}
							/>
						))}
					</CollapsibleContent>
				</CardContent>
			</Card>
		</Collapsible>
	);
}

function NoMembers() {
	return (
		<div className="p-6 text-muted-foreground text-center  w-[90%] mx-auto tracking-wide">
			<p className="pb-4 font-semibold">Hey, Billy no mates!</p>
			<p>Add some friends to your session to get the ball rolling.</p>
		</div>
	);
}

function NoEventsMessage() {
	return (
		<div className="p-6 text-muted-foreground text-center  w-[90%] mx-auto tracking-wide">
			<p className="pb-4 font-semibold">Well this is a bit depressing!</p>
			<p>
				It looks like nobody's bought a round yet. Once someone gets one in, it
				will be shown here.
			</p>
		</div>
	);
}
