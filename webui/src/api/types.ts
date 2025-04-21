export type User = {
	id: string;
	name: string;
	username: string;
};

export type UserBalance = {
	credit: number;
	debit: number;
	net: number;
};

export type Session = {
	id: string;
	name: string;
	total: number;
	isActive: boolean;
	members: SessionMember[];
	transactions: SessionTransaction[];
};

export type SessionMember = {
	id: string;
	name: string;
	username: string;
	isCreator: boolean;
	isAdmin: boolean;
	isDeleted: boolean;
	transactionSummary: {
		credit: number;
		debit: number;
	};
};

export type SessionTransaction = {
	id: string;
	creatorId: string;
	total: number;
	createdAt: string;
	members: SessionTransactionMember[];
};

export type SessionTransactionMember = {
	userId: string;
	name: string;
	username: string;
	amount: number;
};

// TransactionMemberAmounts is a Record where the key is the member ID and the value is the amount in the transaction.
export type TransactionMemberAmounts = Record<string, number>;
