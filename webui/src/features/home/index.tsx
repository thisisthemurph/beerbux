import { AuthenticatedView } from "@/features/home/authenticated/authenticated-view.tsx";
import { useUserStore } from "@/stores/user-store.tsx";

function HomePage() {
	const user = useUserStore((state) => state.user);
	const message: string = import.meta.env.VITE_TEST ?? "Default message";

	return user ? <AuthenticatedView user={user} /> : <p>Welcome to Beerbux: {message}</p>;
}

export default HomePage;
