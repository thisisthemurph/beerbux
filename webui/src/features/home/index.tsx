import { AuthenticatedView } from "@/features/home/authenticated/AuthenticatedView.tsx";
import { useUser } from "@/hooks/useUser.tsx";

function HomePage() {
	const { user } = useUser();

	return (
		<>
			<div>
				{user ? <AuthenticatedView user={user} /> : <p>Welcome to Beerbux</p>}
			</div>
		</>
	);
}

export default HomePage;
