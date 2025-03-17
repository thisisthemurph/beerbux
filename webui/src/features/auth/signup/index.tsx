import useAuthClient from "@/api/authClient.ts";
import { isValidationErrorResponse } from "@/api/types.ts";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import {
	SignupForm,
	type SignupFormValues,
} from "@/features/auth/signup/SignupForm.tsx";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";

function SignupPage() {
	const navigate = useNavigate();
	const authClient = useAuthClient();

	function handleSignup(values: SignupFormValues) {
		authClient
			.signup(
				values.name,
				values.username,
				values.password,
				values.verificationPassword,
			)
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

				navigate("/login?signup=true");
			})
			.catch((error) => {
				toast.error("An error occurred", { description: error.message });
			});
	}

	return (
		<div>
			<h1>Signup</h1>
			<Card>
				<CardHeader>
					<CardTitle>Sign up</CardTitle>
					<CardDescription>Sign up to the app</CardDescription>
				</CardHeader>
				<CardContent>
					<SignupForm onSubmit={handleSignup} />
				</CardContent>
				<CardFooter>
					<Link to="/login" className="text-sm">
						Already have an account?
					</Link>
				</CardFooter>
			</Card>
		</div>
	);
}

export default SignupPage;
