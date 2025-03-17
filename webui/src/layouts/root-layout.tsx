import { Outlet } from "react-router";
import { Toaster } from "sonner";

function RootLayout() {
	return (
		<div>
			<nav>
				<ul className="flex justify-center space-x-4">
					<li>
						<a href="/">Home</a>
					</li>
					<li>
						<a href="/login">Login</a>
					</li>
				</ul>
			</nav>
			<main className="p-4">
				<Outlet />
				<Toaster />
			</main>
		</div>
	);
}

export default RootLayout;
