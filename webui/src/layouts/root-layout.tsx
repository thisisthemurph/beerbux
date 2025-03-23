import NavigationDrawer from "@/components/navigation-drawer.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useUserStore } from "@/stores/user-store.tsx";
import { Link, Outlet } from "react-router";
import { Toaster } from "sonner";

function RootLayout() {
	const user = useUserStore((state) => state.user);

	return (
		<div>
			<nav className="flex justify-between items-center px-4 py-5 mb-4">
				<Link to="/" className="text-xl tracking-wider">
					<span className="font-semibold">Beer</span>
					<span className="text-slate-600 dark:text-gray-400">bux</span>
				</Link>
				<div className="flex items-center gap-4">
					{!user && <LoginButton />}
					<NavigationDrawer />
				</div>
			</nav>
			<main className="p-4">
				<Outlet />
				<Toaster />
			</main>
		</div>
	);
}

function LoginButton() {
	return (
		<Button size="sm" variant="secondary" asChild className="rounded-full">
			<Link to="/login" className="text-sm tracking-wider">
				Login
			</Link>
		</Button>
	);
}

export default RootLayout;
