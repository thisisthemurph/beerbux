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
