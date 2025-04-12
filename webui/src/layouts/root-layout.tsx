import { Navigation } from "@/components/navigation.tsx";
import { Outlet } from "react-router";
import { Toaster } from "sonner";

function RootLayout() {
	return (
		<div>
			<Navigation />
			<main className="p-4 space-y-4">
				<Outlet />
				<Toaster />
			</main>
		</div>
	);
}

export default RootLayout;
