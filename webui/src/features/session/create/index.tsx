import { ValidationError } from "@/api/api-fetch.ts";
import useSessionClient from "@/api/session-client.ts";
import {
	NewSessionForm,
	type NewSessionFormValues,
} from "@/features/session/create/new-session-form.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import { PageHeading } from "@/components/page-heading.tsx";

function CreateSessionPage() {
	useBackNavigation("/");
	const navigate = useNavigate();
	const { createSession } = useSessionClient();

	async function handleCreateSession({ name }: NewSessionFormValues) {
		const { err } = await tryCatch(createSession(name));
		if (err) {
			handleCreateSessionError(err);
			return;
		}

		navigate("/");
	}

	function handleCreateSessionError(err: unknown) {
		if (err instanceof ValidationError) {
			toast.error("The session could not be created", {
				description: err.validationErrors.errors.name,
			});
		}
	}

	return (
		<>
			<PageHeading title="Start a session" />
			<NewSessionForm onCreate={handleCreateSession} />
		</>
	);
}

export default CreateSessionPage;
