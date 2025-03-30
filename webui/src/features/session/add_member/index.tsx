import useSessionClient from "@/api/session-client.ts";
import { AddMemberForm } from "@/features/session/add_member/add-member-form.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";

function AddMemberPage() {
	const { sessionId } = useParams();
	useBackNavigation(`/session/${sessionId}`);
	const navigate = useNavigate();
	const { addMemberToSession } = useSessionClient();

	async function handleAddMember(username: string) {
		if (!sessionId) return;
		try {
			await addMemberToSession(sessionId, username);
			navigate(`/session/${sessionId}`);
		} catch (err) {
			const message =
				err instanceof Error
					? (err.message ?? "An unknown error occurred")
					: "An unknown error occurred";

			toast.error("Could not add member to session", {
				description: message,
			});
		}
	}

	return (
		<>
			<h1>Add member</h1>
			<AddMemberForm onAdd={handleAddMember} />
		</>
	);
}

export default AddMemberPage;
