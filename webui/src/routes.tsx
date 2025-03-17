import LoginPage from "@/features/auth/login";
import SignupPage from "@/features/auth/signup";
import HomePage from "@/features/home";
import RootLayout from "@/layouts/root-layout";
import { Route, Routes } from "react-router";

function AppRoutes() {
	return (
		<Routes>
			<Route element={<RootLayout />}>
				<Route index element={<HomePage />} />
				<Route path="/login" element={<LoginPage />} />
				<Route path="/signup" element={<SignupPage />} />
			</Route>
		</Routes>
	);
}

export default AppRoutes;
