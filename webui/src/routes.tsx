import RootLayout from "@/layouts/root-layout";
import HomePage from "@/pages/home";
import { Route, Routes } from "react-router";

function AppRoutes() {
	return (
		<Routes>
			<Route element={<RootLayout />}>
				<Route index element={<HomePage />} />
			</Route>
		</Routes>
	);
}

export default AppRoutes;
