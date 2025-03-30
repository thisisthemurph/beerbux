import type { SessionMember, SessionTransaction } from "@/api/types.ts";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { UserAvatar } from "@/components/user-avatar.tsx";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";

type TransactionListingProps = {
	members: SessionMember[];
	transactions: SessionTransaction[];
	avatarData: Record<string, AvatarData>;
};

export function TransactionListing({
	transactions,
	members,
	avatarData,
}: TransactionListingProps) {
	return (
		<Card>
			<CardHeader>
				<section className="flex justify-between">
					<CardTitle>Transactions</CardTitle>
					<p>{transactions.length}</p>
				</section>
			</CardHeader>
			<CardContent className="px-0">
				{transactions.length === 0 && <NoTransactionsMessage />}
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
			</CardContent>
		</Card>
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
			<p className="col-span-4 font-semibold tracking-wider">
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

function NoTransactionsMessage() {
	return (
		<p className="p-6 text-muted-foreground text-center font-semibold w-[90%] mx-auto">
			You don't have any transactions in this session yet. Once you have some,
			they will show up here.
		</p>
	);
}
