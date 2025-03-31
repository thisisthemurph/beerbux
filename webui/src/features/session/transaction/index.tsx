import useSessionClient from "@/api/session-client.ts";
import useTransactionClient from "@/api/transaction-client.ts";
import type { TransactionMemberAmounts } from "@/api/types.ts";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { TransactionForm } from "@/features/session/transaction/transaction-form.tsx";
import { useBackNavigation } from "@/hooks/use-back-navigation.ts";
import { tryCatch } from "@/lib/try-catch.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";

function TransactionPage() {
	const user = useUserStore((state) => state.user);
	const { getSession } = useSessionClient();
	const { sessionId } = useParams();
	const navigate = useNavigate();
	useBackNavigation(`/session/${sessionId}`);
	const [total, setTotal] = useState(0);
	const { createTransaction } = useTransactionClient();
	const queryClient = useQueryClient();

	async function handleNewTransaction(transaction: TransactionMemberAmounts) {
		if (!sessionId) return;
		const total = Object.values(transaction).reduce(
			(acc, value) => acc + value,
			0,
		);

		const { err } = await tryCatch(createTransaction(sessionId, transaction));
		if (err) {
			console.error(err);
			toast.error("There was an issue creating the transaction.");
			return;
		}

		toast.success("Transaction created:", {
			description: (
				<p>
					A transaction of <span className="font-semibold">${total}</span> has
					been created.
				</p>
			),
		});

		await queryClient.invalidateQueries({ queryKey: ["session", sessionId] });
		navigate(`/session/${sessionId}`);
	}

	const { data: session, isPending: sessionIsPending } = useQuery({
		queryKey: ["session", sessionId],
		queryFn: () => {
			if (!sessionId) {
				return null;
			}
			return getSession(sessionId);
		},
	});

	if (sessionIsPending) {
		return <p>Loading</p>;
	}

	if (!session) {
		return <p>There has been an issue loading the session.</p>;
	}

	return (
		<div>
			<h1>{session.name}</h1>

			<Card>
				<CardHeader>
					<section className="flex justify-between items-center">
						<CardTitle className="tracking-wider">Transaction</CardTitle>
						<p className="font-semibold text-muted-foreground tracking-wider">
							${total}
						</p>
					</section>
				</CardHeader>
				<CardContent>
					<TransactionForm
						members={session.members.filter((m) => m.id !== user?.id)}
						onTransactionCreate={handleNewTransaction}
						onTotalChanged={(t) => setTotal(t)}
					/>
				</CardContent>
			</Card>
		</div>
	);
}

export default TransactionPage;
