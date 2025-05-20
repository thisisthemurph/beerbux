export type HistoryEventType =
	| "transaction_created"
	| "member_added"
	| "member_removed"
	| "member_left"
	| "session_closed"
	| "session_opened"
	| "promoted_to_admin"
	| "demoted_from_admin";

export type SessionHistory = {
	sessionId: string;
	events: SessionHistoryEvent[] | null;
};

interface BaseSessionHistoryEvent {
	id: number;
	eventType: HistoryEventType;
	eventData: unknown;
	memberId: string;
	createdAt: string;
}

export interface TransactionCreatedSessionHistoryEvent extends BaseSessionHistoryEvent {
	eventType: "transaction_created";
	eventData: {
		transactionId: string;
		lines: {
			memberId: string;
			amount: number;
		}[];
	};
}

export interface MemberPromotedToAdminEvent extends BaseSessionHistoryEvent {
	eventType: "promoted_to_admin";
	eventData: {
		memberId: string;
	};
}

export interface MemberDemotedFromAdminEvent extends BaseSessionHistoryEvent {
	eventType: "demoted_from_admin";
	eventData: {
		memberId: string;
	};
}

export interface MemberAddedSessionHistoryEvent extends BaseSessionHistoryEvent {
	eventType: "member_added";
	eventData: {
		memberId: string;
	};
}

export interface MemberRemovedSessionHistoryEvent extends BaseSessionHistoryEvent {
	eventType: "member_removed";
	eventData: {
		memberId: string;
	};
}

export interface NoEventDataSessionHistoryEvent extends BaseSessionHistoryEvent {
	eventType: "member_left" | "session_closed" | "session_opened";
	eventData: undefined;
}

export type SessionHistoryEvent =
	| TransactionCreatedSessionHistoryEvent
	| MemberAddedSessionHistoryEvent
	| MemberRemovedSessionHistoryEvent
	| NoEventDataSessionHistoryEvent
	| MemberPromotedToAdminEvent
	| MemberDemotedFromAdminEvent;
