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
import { format, parse } from "date-fns";
import { Button } from "@/components/ui/button.tsx";
import { ChevronDown, ChevronUp } from "lucide-react";
import { EventGroup } from "./event-group";

type GroupedEventRecords = Record<string, SessionHistoryEvent[]>;
const DATE_FMT_LONG = "EEEE do MMMM, yyyy";

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
	const groupedEvents = groupEventsByDate(events);
	const sortedGroupLabels = Object.keys(groupedEvents).sort((a, b) => {
		return (
			parse(b, DATE_FMT_LONG, new Date()).getTime() -
			parse(a, DATE_FMT_LONG, new Date()).getTime()
		);
	});

	const firstGroupLabel = sortedGroupLabels[0];
	const firstGroupEvents = groupedEvents[firstGroupLabel];
	const showCollapsibleTrigger =
		sortedGroupLabels.length > 1 ||
		(firstGroupEvents && firstGroupEvents.length > 5);

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

					{firstGroupEvents.length > 0 && (
						<EventGroup
							label={firstGroupLabel}
							events={isOpen ? firstGroupEvents : firstGroupEvents.slice(0, 5)}
							members={members}
							avatarData={avatarData}
						/>
					)}

					<CollapsibleContent>
						{sortedGroupLabels.slice(1).map((label) => (
							<EventGroup
								key={label}
								label={label}
								events={groupedEvents[label]}
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

function groupEventsByDate(events: SessionHistoryEvent[]): GroupedEventRecords {
	const groupedEvents = events.reduce((acc, transaction) => {
		const formattedDate = format(
			new Date(transaction.createdAt),
			DATE_FMT_LONG,
		);

		if (!acc[formattedDate]) {
			acc[formattedDate] = [];
		}

		acc[formattedDate].push(transaction);
		return acc;
	}, {} as GroupedEventRecords);

	for (const date in groupedEvents) {
		groupedEvents[date].sort((a, b) => {
			return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
		});
	}

	return groupedEvents;
}
