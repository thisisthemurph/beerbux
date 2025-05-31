import useFriendsClient from "@/api/friends-client.ts";
import useSessionClient from "@/api/session-client.ts";
import type { User } from "@/api/types/user.ts";
import useUserClient from "@/api/user-client.ts";
import {
	PrimaryActionCard,
	PrimaryActionCardButtonItem,
	PrimaryActionCardContent,
} from "@/components/primary-action-card";
import { SessionListing } from "@/components/session-listing.tsx";
import { CreateSessionDrawer } from "@/features/home/authenticated/create-session/create-session-drawer.tsx";
import { UserCard } from "@/features/home/authenticated/user-card.tsx";
import { tryCatch } from "@/lib/try-catch.ts";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { SquareChevronRight } from "lucide-react";
import { Suspense, useState } from "react";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";
import { FriendListing } from "./friend-listing";

type AuthenticatedViewProps = {
	user: User;
};

export function AuthenticatedView({ user }: AuthenticatedViewProps) {
	const { getSessions, getBalance } = useUserClient();
	const { createSession } = useSessionClient();
	const { getFriends } = useFriendsClient();
	const [createSessionOpen, setCreateSessionOpen] = useState(false);
	const navigate = useNavigate();

	const { data: sessions } = useSuspenseQuery({
		queryKey: ["sessions"],
		queryFn: () => getSessions(3),
	});

	const { data: balance } = useQuery({
		queryKey: ["balance", user.id],
		queryFn: () => getBalance(user.id),
		placeholderData: { credit: 0, debit: 0, net: 0 },
	});

	const { data: friends } = useQuery({
		queryKey: ["friends"],
		queryFn: () => getFriends(),
		placeholderData: [],
	});

	async function handleCreateSession({ name }: { name: string }) {
		const { data, err } = await tryCatch(createSession(name));
		if (err) {
			console.error(err);
			toast.error("Failed to create session");
			return;
		}

		navigate(`/session/${data.id}`);
		toast.error("Session created", {
			description: (
				<>
					<p>
						Your <span className="font-semibold underline">{name}</span> session has been created.
					</p>
					<p>Add members to get started!</p>
				</>
			),
		});
	}

	const userBalance = balance ?? { credit: 0, debit: 0, net: 0 };

	return (
		<>
			<UserCard
				{...user}
				credit={userBalance.credit}
				debit={userBalance.debit}
				netBalance={userBalance.net}
			/>

			<PrimaryActionCard>
				<PrimaryActionCardContent>
					<PrimaryActionCardButtonItem
						text="Start new session"
						icon={<SquareChevronRight className="text-primary w-8 h-8" />}
						onClick={() => setCreateSessionOpen(true)}
					/>
				</PrimaryActionCardContent>
			</PrimaryActionCard>

			<Suspense fallback={<SessionListing.Skeleton />}>
				<SessionListing sessions={sessions}>{sessions && <AllSessionsLink />}</SessionListing>
			</Suspense>

			<FriendListing friends={friends ?? []} />

			<CreateSessionDrawer
				open={createSessionOpen}
				onOpenChange={setCreateSessionOpen}
				onCreate={handleCreateSession}
			/>
		</>
	);
}

function AllSessionsLink() {
	return (
		<Link to="/sessions" className="text-blue-400">
			All sessions
		</Link>
	);
}
