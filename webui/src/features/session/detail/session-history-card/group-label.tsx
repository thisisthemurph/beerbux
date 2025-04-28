import { format, isThisYear, isToday, isYesterday, parse } from "date-fns";

const DATE_FMT_LONG = "EEEE do MMMM, yyyy";
const DATE_FMT_SHORT = "EEEE do MMMM";

export function GroupLabel({ text }: { text: string }) {
	const d = parse(text, DATE_FMT_LONG, new Date());
	const label = isToday(d)
		? "Today"
		: isYesterday(d)
			? "Yesterday"
			: format(d, isThisYear(d) ? DATE_FMT_SHORT : DATE_FMT_LONG);

	return <p className="px-6 py-4 text-muted-foreground font-semibold tracking-wide">{label}</p>;
}
