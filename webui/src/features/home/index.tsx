import { AuthenticatedView } from "@/features/home/authenticated/AuthenticatedView.tsx";
import { useUserStore } from "@/stores/userStore.tsx";

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
