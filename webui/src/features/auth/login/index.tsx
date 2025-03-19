import { ValidationError } from "@/api/apiFetch.ts";
import useAuthClient from "@/api/authClient.ts";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import {
	LoginForm,
	type LoginFormValues,
} from "@/features/auth/login/LoginForm.tsx";
import { useUserStore } from "@/stores/userStore.tsx";
import { useEffect, useRef } from "react";
import { Link, useNavigate, useSearchParams } from "react-router";
import { toast } from "sonner";

function LoginPage() {
	const [searchParams] = useSearchParams();
	const navigatedAfterSignup = searchParams.get("signup") === "true";
	const hasShownToast = useRef(false);
	const { fetchUser } = useUserStore();
	const navigate = useNavigate();
	const { login } = useAuthClient();

	useEffect(() => {
		if (navigatedAfterSignup && !hasShownToast.current) {
			toast.success("Signed up successfully! You can now login.");
			hasShownToast.current = true;
		}
	}, [navigatedAfterSignup]);

	async function handleLogin({ username, password }: LoginFormValues) {
		try {
			await login(username, password);
			await fetchUser();
			navigate("/");
		} catch (err) {
			handleLoginError(err);
		}
	}

	function handleLoginError(err: unknown) {
		if (err instanceof ValidationError) {
			toast.error("There was an issue with the data you provided", {
				description: (
					<pre className="p-2 bg-foreground text-xs text-background rounded font-mono">
						{JSON.stringify(err.validationErrors.errors, null, 2)}
					</pre>
				),
			});
			return;
		}

		toast.error("An error occurred", {
			description: err instanceof Error ? err.message : "Unknown error",
		});
	}

	return (
		<div>
			<h1>Login</h1>
			<Card>
				<CardHeader>
					<CardTitle>Login</CardTitle>
					<CardDescription>Login to the app</CardDescription>
				</CardHeader>
				<CardContent>
					<LoginForm onSubmit={handleLogin} />
				</CardContent>
				<CardFooter>
					<Link to="/signup" className="text-sm">
						Don't have an account?
					</Link>
				</CardFooter>
			</Card>
		</div>
	);
}

export default LoginPage;
