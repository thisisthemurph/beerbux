import { ValidationError } from "@/api/api-fetch.ts";
import useAuthClient from "@/api/auth-client.ts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { SignupForm, type SignupFormValues } from "@/features/auth/signup/signup-form.tsx";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";
import { PageHeading } from "@/components/page-heading.tsx";

function SignupPage() {
	const navigate = useNavigate();
	const { signup } = useAuthClient();

	async function handleSignup({ name, username, password, verificationPassword }: SignupFormValues) {
		try {
			await signup(name, username, password, verificationPassword);
			navigate("/login?signup=true");
		} catch (err) {
			handleSignupError(err);
		}
	}

	function handleSignupError(err: unknown) {
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
			<PageHeading title="Add member" />
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
		</>
	);
}

export default SignupPage;
