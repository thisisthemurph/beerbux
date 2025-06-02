export interface Session {
	id: string;
	name: string;
	isActive: boolean;
	members: SessionMember[];
}

export interface SessionWithTransactions extends Session {
	transactions: SessionTransaction[];
	total: number;
}

export type SessionMember = {
	id: string;
	name: string;
	username: string;
	isCreator: boolean;
	isAdmin: boolean;
	isDeleted: boolean;
	transactionSummary: MemberTransactionSummary;
};

type MemberTransactionSummary = {
	credit: number;
	debit: number;
};

export type SessionTransaction = {
	id: string;
	userId: string;
	total: number;
	createdAt: string;
	lines: SessionTransactionLine[];
};

export type SessionTransactionLine = {
	userId: string;
	name: string;
	username: string;
	amount: number;
};
