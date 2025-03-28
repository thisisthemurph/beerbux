import { ValidationError } from "@/api/api-fetch.ts";
import useSessionClient from "@/api/session-client.ts";
import {
	NewSessionForm,
	type NewSessionFormValues,
} from "@/features/session/create/new-session-form.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useNavigate } from "react-router";
import { toast } from "sonner";

function CreateSessionPage() {
	useBackNavigation("/");
	const navigate = useNavigate();
	const { createSession } = useSessionClient();

	async function handleCreateSession({ name }: NewSessionFormValues) {
		try {
			await createSession(name);
			navigate("/");
		} catch (err) {
			handleCreateSessionError(err);
		}
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
			<h1>Start a session</h1>
			<NewSessionForm onCreate={handleCreateSession} />
		</>
	);
}

export default CreateSessionPage;
