import type { SessionHistoryEvent } from "@/api/types/session-history.ts";
import { format, parse } from "date-fns";

type GroupedEventRecords = Record<string, SessionHistoryEvent[]>;
const DATE_FMT_LONG = "EEEE do MMMM, yyyy";

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

type GroupedEvents = {
	// Contains the most recent date label to be displayed first
	firstLabel: string;
	// Contains the events of the most recent label
	firstEvents: SessionHistoryEvent[];
	// Contains all the events grouped by the label
	events: GroupedEventRecords;
	// Contains the labels sorted by date
	sortedLabels: string[];
}

export function useGroupedEvents(eventsToGroup: SessionHistoryEvent[]): GroupedEvents {
	const events = groupEventsByDate(eventsToGroup);
	const sortedLabels = Object.keys(events).sort((a, b) => {
		return (
			parse(b, DATE_FMT_LONG, new Date()).getTime() -
			parse(a, DATE_FMT_LONG, new Date()).getTime()
		);
	});

	const firstLabel = sortedLabels[0];
	const firstEvents = events[firstLabel];

	return {
		firstLabel,
		firstEvents,
		events,
		sortedLabels,
	};
}
