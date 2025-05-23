import { ValidationError } from "@/api/api-fetch.ts";
import useAuthClient from "@/api/auth-client.ts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { LoginForm, type LoginFormValues } from "@/features/auth/login/login-form.tsx";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useEffect, useRef } from "react";
import { Link, useNavigate, useSearchParams } from "react-router";
import { toast } from "sonner";
import { PageHeading } from "@/components/page-heading.tsx";

function LoginPage() {
	const [searchParams] = useSearchParams();
	const navigatedAfterSignup = searchParams.get("signup") === "true";
	const hasShownToast = useRef(false);
	const setUser = useUserStore((state) => state.setUser);
	const navigate = useNavigate();
	const { login } = useAuthClient();

	useEffect(() => {
		if (navigatedAfterSignup && !hasShownToast.current) {
			toast.success("Signed up successfully! You can now login.");
			hasShownToast.current = true;
		}
	}, [navigatedAfterSignup]);

	async function handleLogin({ username, password }: LoginFormValues) {
		const { data: user, err } = await tryCatch(login(username, password));
		if (err) {
			handleLoginError(err);
			return;
		}

		setUser(user);
		navigate("/");
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
		<>
			<PageHeading title="Login" />
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
		</>
	);
}

export default LoginPage;
