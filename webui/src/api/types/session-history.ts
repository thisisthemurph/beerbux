export type HistoryEventType = "transaction_created" | "member_added";

export type SessionHistory = {
	sessionId: string;
	events: SessionHistoryEvent[];
};

interface BaseSessionHistoryEvent {
	id: number;
	eventType: HistoryEventType;
	memberId: string;
	createdAt: string;
}

export interface TransactionCreatedSessionHistoryEvent
	extends BaseSessionHistoryEvent {
	eventType: "transaction_created";
	eventData: {
		transactionId: string;
		lines: {
			memberId: string;
			amount: number;
		}[];
	};
}

export interface MemberAddedSessionHistoryEvent
	extends BaseSessionHistoryEvent {
	eventType: "member_added";
	eventData: {
		memberId: string;
	};
}

export type SessionHistoryEvent =
	| TransactionCreatedSessionHistoryEvent
	| MemberAddedSessionHistoryEvent;
