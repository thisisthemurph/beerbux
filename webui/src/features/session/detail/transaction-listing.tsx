import type { SessionMember, SessionTransaction } from "@/api/types.ts";
import { Button } from "@/components/ui/button";
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
} from "@/components/ui/collapsible";
import { UserAvatar } from "@/components/user-avatar.tsx";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { format, isThisYear, isToday, isYesterday, parse } from "date-fns";
import { ChevronDown, ChevronUp } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils.ts";

type TransactionListingProps = {
	members: SessionMember[];
	transactions: SessionTransaction[];
	avatarData: Record<string, AvatarData>;
};

type GroupedTransactions = Record<string, SessionTransaction[]>;
const DATE_FMT_LONG = "EEEE do MMMM, yyyy";
const DATE_FMT_SHORT = "EEEE do MMMM";

export function TransactionListing({
	transactions,
	members,
	avatarData,
}: TransactionListingProps) {
	const groupedTransactions = groupTransactionsByDate(transactions);
	const [isOpen, setIsOpen] = useState(false);

	const sortedDateKeys = Object.keys(groupedTransactions).sort((a, b) => {
		return (
			parse(b, DATE_FMT_LONG, new Date()).getTime() -
			parse(a, DATE_FMT_LONG, new Date()).getTime()
		);
	});

	const firstDateLabel = sortedDateKeys[0];
	const firstTransactionGroup = groupedTransactions[firstDateLabel];
	const showCollapsibleTrigger =
		sortedDateKeys.length > 1 ||
		(firstTransactionGroup && firstTransactionGroup.length > 5);

	return (
		<Collapsible open={isOpen} onOpenChange={setIsOpen}>
			<Card>
				<CardHeader>
					<section className="flex justify-between items-center">
						<CardTitle>Rounds</CardTitle>
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
					{members.length <= 1 ? (
						<NoMembers />
					) : (
						transactions.length === 0 && <NoTransactionsMessage />
					)}

					{firstTransactionGroup && firstTransactionGroup.length > 0 && (
						<TransactionGroup
							transactions={
								isOpen
									? firstTransactionGroup
									: firstTransactionGroup.slice(0, 5)
							}
							dateLabel={firstDateLabel}
							members={members}
							avatarData={avatarData}
						/>
					)}
					<CollapsibleContent>
						{sortedDateKeys.slice(1).map((dateLabel) => (
							<TransactionGroup
								key={dateLabel}
								transactions={groupedTransactions[dateLabel]}
								dateLabel={dateLabel}
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

function TransactionGroup({
	dateLabel,
	transactions,
	members,
	avatarData,
}: {
	transactions: SessionTransaction[];
	dateLabel: string;
	members: SessionMember[];
	avatarData: Record<string, AvatarData>;
}) {
	return (
		<>
			<GroupLabel date={dateLabel} />
			{transactions.map((t) => {
				const creator = members.find((m) => m.id === t.creatorId);
				const creatorUsername = creator?.username ?? "unknown";
				const creatorAvatarData = avatarData[creatorUsername];

				return (
					<div
						key={t.id}
						className="flex items-center gap-4 px-6 py-4 hover:bg-muted transition-colors"
					>
						<UserAvatar data={creatorAvatarData} />
						<TransactionText
							creator={creator}
							transaction={t}
							members={members}
						/>
					</div>
				);
			})}
		</>
	);
}

type TransactionTextProps = {
	creator: SessionMember | undefined;
	transaction: SessionTransaction;
	members: SessionMember[];
};

function TransactionText({
	creator,
	transaction,
	members,
}: TransactionTextProps) {
	function stringifyMemberNames(usernames: string[]): string {
		if (usernames.length === members.length - 1) {
			return "everyone";
		}

		if (usernames.length === 1) {
			return usernames[0];
		}

		if (usernames.length === 2) {
			return `${usernames[0]} and ${usernames[1]}`;
		}

		return `${usernames.slice(0, -1).join(", ")}, and ${usernames.slice(-1)[0]}`;
	}

	return (
		<div className="grid grid-cols-5 grid-rows-2 w-full">
			<p
				className={cn(
					"col-span-4 font-semibold tracking-wider",
					creator?.isDeleted && "line-through",
				)}
			>
				{creator?.username ?? "unknown"}
			</p>
			<div className="row-span-2 flex items-center justify-end">
				<p className="font-semibold">${transaction.total}</p>
			</div>
			<p className="col-span-4 text-muted-foreground">
				{stringifyMemberNames(transaction.members.map((m) => m.username))}
			</p>
		</div>
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

function NoTransactionsMessage() {
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

function GroupLabel({ date }: { date: string }) {
	const d = parse(date, DATE_FMT_LONG, new Date());
	const label = isToday(d)
		? "Today"
		: isYesterday(d)
			? "Yesterday"
			: format(d, isThisYear(d) ? DATE_FMT_SHORT : DATE_FMT_LONG);

	return (
		<p className="px-6 py-4 text-muted-foreground font-semibold tracking-wide">
			{label}
		</p>
	);
}

function groupTransactionsByDate(
	transactions: SessionTransaction[],
): GroupedTransactions {
	const groupedTransactions = transactions.reduce((acc, transaction) => {
		const formattedDate = format(
			new Date(transaction.createdAt),
			DATE_FMT_LONG,
		);

		if (!acc[formattedDate]) {
			acc[formattedDate] = [];
		}

		acc[formattedDate].push(transaction);

		return acc;
	}, {} as GroupedTransactions);

	for (const date in groupedTransactions) {
		groupedTransactions[date].sort((a, b) => {
			return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
		});
	}

	return groupedTransactions;
}
