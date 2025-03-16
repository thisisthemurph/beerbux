import { Outlet } from "react-router";

function RootLayout() {
	return (
		<div>
			<main className="p-4">
				<Outlet />
			</main>
		</div>
	);
}

export default RootLayout;
