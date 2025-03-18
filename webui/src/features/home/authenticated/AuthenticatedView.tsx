import { UserCard } from "@/features/home/authenticated/UserCard.tsx";
import type { User } from "@/stores/userStore.tsx";

type AuthenticatedViewProps = {
	user: User;
};

export function AuthenticatedView({ user }: AuthenticatedViewProps) {
	return (
		<div>
			<UserCard {...user} />
		</div>
	);
}
