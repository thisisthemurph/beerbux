import AppRoutes from "@/routes.tsx";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import "./index.css";
import { ThemeProvider } from "@/components/theme-provider.tsx";
import { UserProvider } from "@/stores/userProvider.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";

const rootElement = document.getElementById("root");
if (!rootElement) throw new Error("Root element not found");

const queryClient = new QueryClient();

createRoot(rootElement).render(
	<StrictMode>
		<BrowserRouter>
			<UserProvider>
				<QueryClientProvider client={queryClient}>
					<ThemeProvider>
						<AppRoutes />
					</ThemeProvider>
				</QueryClientProvider>
			</UserProvider>
			<Toaster />
		</BrowserRouter>
	</StrictMode>,
);
