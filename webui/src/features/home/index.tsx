import { AuthenticatedView } from "@/features/home/authenticated/authenticated-view.tsx";
import { useUserStore } from "@/stores/user-store.tsx";

function HomePage() {
	const user = useUserStore((state) => state.user);

	return (
		<>
			<div>
				{user ? <AuthenticatedView user={user} /> : <p>Welcome to Beerbux</p>}
			</div>
		</>
	);
}

export default HomePage;
