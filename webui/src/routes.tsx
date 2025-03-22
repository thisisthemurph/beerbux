import LoginPage from "@/features/auth/login";
import SignupPage from "@/features/auth/signup";
import HomePage from "@/features/home";
import CreateSessionPage from "@/features/session/create";
import SessionDetailPage from "@/features/session/detail/SessionDetailPage.tsx";
import RootLayout from "@/layouts/root-layout";
import { Route, Routes } from "react-router";

function AppRoutes() {
	return (
		<Routes>
			<Route element={<RootLayout />}>
				<Route index element={<HomePage />} />
				<Route path="/login" element={<LoginPage />} />
				<Route path="/signup" element={<SignupPage />} />
				<Route path="/session/create" element={<CreateSessionPage />} />
				<Route path="/session/:sessionId" element={<SessionDetailPage />} />
			</Route>
		</Routes>
	);
}

export default AppRoutes;
