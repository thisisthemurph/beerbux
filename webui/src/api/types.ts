export type User = {
	id: string;
	name: string;
	username: string;
	netBalance: number;
};

export type Session = {
	id: string;
	name: string;
	isActive: boolean;
	members: SessionMember[];
};

export type SessionMember = {
	id: string;
	name: string;
	username: string;
};

// TransactionMemberAmounts is a Record where the key is the member ID and the value is the amount in the transaction.
export type TransactionMemberAmounts = Record<string, number>;
