import { useEffect } from "react";

export interface SessionTransactionCreatedMessage {
	sessionId: string;
	creatorId: string;
	transactionId: string;
	total: number;
}

export function useSessionEventSource(
	url: string,
	handleEvent: (msg: SessionTransactionCreatedMessage) => Promise<void>,
) {
	useEffect(() => {
		const eventSource = new EventSource(url, { withCredentials: true });

		eventSource.onerror = (event) => {
			console.error("SSE error:", event);
			eventSource.close();
		};

		eventSource.addEventListener(
			"session.transaction.created",
			async (event) => {
				const data = JSON.parse(event.data) as SessionTransactionCreatedMessage;
				await handleEvent(data);
			},
		);

		return () => eventSource.close();
	}, [url, handleEvent]);
}
