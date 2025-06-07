import NotFoundPage from "@/features/NotFound.tsx";
import LoginPage from "@/features/auth/login";
import SignupPage from "@/features/auth/signup";
import DashboardPage from "@/features/dashboard";
import FriendDetailPage from "@/features/firend";
import HomePage from "@/features/home";
import SessionDetailPage from "@/features/session/detail";
import SessionListingPage from "@/features/session/listing";
import SettingsPage from "@/features/settings";
import RootLayout from "@/layouts/root-layout";
import { useUserStore } from "@/stores/user-store.tsx";
import type { JSX } from "react";
import { Route, Routes } from "react-router";

function AuthGuard({ page, alt }: { page: JSX.Element; alt?: JSX.Element }) {
	const user = useUserStore((state) => state.user);
	return user ? page : alt ? alt : <LoginPage />;
}

function AppRoutes() {
	return (
		<Routes>
			<Route element={<RootLayout />}>
				<Route index element={<AuthGuard page={<DashboardPage />} alt={<HomePage />} />} />
				<Route path="/login" element={<LoginPage />} />
				<Route path="/signup" element={<SignupPage />} />
				<Route path="/sessions" element={<AuthGuard page={<SessionListingPage />} />} />
				<Route path="/sessions" element={<AuthGuard page={<SessionListingPage />} />} />
				<Route path="/session/:sessionId" element={<AuthGuard page={<SessionDetailPage />} />} />
				<Route path="/settings" element={<AuthGuard page={<SettingsPage />} />} />
				<Route path="/friend/:friendId" element={<FriendDetailPage />} />
				<Route path="/friend/:friendId" element={<AuthGuard page={<FriendDetailPage />} />} />
				<Route path="*" element={<NotFoundPage />} />
			</Route>
		</Routes>
	);
}

export default AppRoutes;
