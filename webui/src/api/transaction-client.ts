import { apiFetch } from "@/api/api-fetch.ts";
import type { TransactionMemberAmounts } from "@/api/types.ts";

export default function useTransactionClient() {
	const createTransaction = async (
		sessionId: string,
		transaction: TransactionMemberAmounts,
	) => {
		return await apiFetch<void>(`/session/${sessionId}/transaction`, {
			method: "POST",
			body: JSON.stringify(transaction),
		});
	};

	return { createTransaction };
}
