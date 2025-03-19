import AppRoutes from "@/routes.tsx";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import "./index.css";
import { ThemeProvider } from "@/components/theme-provider.tsx";
import { UserProvider } from "@/stores/userProvider.tsx";
import { Toaster } from "sonner";

const rootElement = document.getElementById("root");
if (!rootElement) throw new Error("Root element not found");

createRoot(rootElement).render(
	<StrictMode>
		<BrowserRouter>
			<UserProvider>
				<ThemeProvider>
					<AppRoutes />
				</ThemeProvider>
			</UserProvider>
			<Toaster />
		</BrowserRouter>
	</StrictMode>,
);
