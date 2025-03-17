import useAuthClient from "@/api/authClient.ts";
import { isValidationErrorResponse } from "@/api/types";
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
import { useUser } from "@/hooks/useUser.tsx";
import { useEffect, useRef } from "react";
import { Link, useNavigate, useSearchParams } from "react-router";
import { toast } from "sonner";

function LoginPage() {
	const [searchParams] = useSearchParams();
	const navigatedAfterSignup = searchParams.get("signup") === "true";
	const hasShownToast = useRef(false);

	const { setUser } = useUser();
	const navigate = useNavigate();
	const authClient = useAuthClient();

	useEffect(() => {
		if (navigatedAfterSignup && !hasShownToast.current) {
			toast.success("Signed up successfully! You can now login.");
			hasShownToast.current = true;
		}
	}, [navigatedAfterSignup]);

	function handleLogin(values: LoginFormValues) {
		authClient
			.login(values.username, values.password)
			.then((response) => {
				if (isValidationErrorResponse(response)) {
					toast.error("There was an issue with the data you provided", {
						description: (
							<pre className="p-2 bg-foreground text-xs text-background rounded font-mono">
								{JSON.stringify(response.errors, null, 2)}
							</pre>
						),
					});

					return;
				}

				setUser({ id: response.id, username: response.username });
				navigate("/");
			})
			.catch((error) => {
				toast.error("An error occurred", { description: error.message });
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
