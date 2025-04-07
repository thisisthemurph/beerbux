import { useEffect } from "react";

type SessionTransactionCreatedMessage = {
	creatorId: string;
	sessionId: string;
	transactionId: string;
	total: number;
};

type Options = {
	sessionId: string;
	userId: string;
	onEventReceived: (msg: SessionTransactionCreatedMessage) => void;
};

const BASE_URL = import.meta.env.VITE_SSE_BASE_URL;

export function useSessionTransactionCreatedEventSource({
	sessionId,
	userId,
	onEventReceived,
}: Options) {
	const url = `${BASE_URL}/session?session_id=${sessionId}&user_id=${userId}`;

	useEffect(() => {
		const eventSource = new EventSource(url, { withCredentials: true });

		eventSource.onerror = (event) => {
			console.error("SSE error:", event);
			eventSource.close();
		};

		eventSource.addEventListener("session.transaction.created", (event) =>
			onEventReceived(
				JSON.parse(event.data) as SessionTransactionCreatedMessage,
			),
		);

		return () => eventSource.close();
	}, [url, onEventReceived]);
}
