import LoginPage from "@/features/auth/login";
import SignupPage from "@/features/auth/signup";
import HomePage from "@/features/home";
import AddMemberPage from "@/features/session/add_member";
import CreateSessionPage from "@/features/session/create";
import SessionDetailPage from "@/features/session/detail";
import SessionListingPage from "@/features/session/listing";
import TransactionPage from "@/features/session/transaction";
import RootLayout from "@/layouts/root-layout";
import { Route, Routes } from "react-router";

function AppRoutes() {
	return (
		<Routes>
			<Route element={<RootLayout />}>
				<Route index element={<HomePage />} />
				<Route path="/login" element={<LoginPage />} />
				<Route path="/signup" element={<SignupPage />} />
				<Route path="/sessions" element={<SessionListingPage />} />
				<Route path="/session/create" element={<CreateSessionPage />} />
				<Route path="/session/:sessionId" element={<SessionDetailPage />} />
				<Route path="/session/:sessionId/member" element={<AddMemberPage />} />
				<Route
					path="/session/:sessionId/transaction"
					element={<TransactionPage />}
				/>
			</Route>
		</Routes>
	);
}

export default AppRoutes;
