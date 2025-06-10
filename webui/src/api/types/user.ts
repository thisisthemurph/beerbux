export type UserAuthDetails = {
	id: string;
	name: string;
	username: string;
};

export type UserBalance = {
	credit: number;
	debit: number;
	net: number;
};

export type User = {
	id: string;
	username: string;
	email: string;
	name: string;
	createdAt: string;
	updatedAt: string;
	account: UserAccount;
};

type UserAccount = {
	debit: number;
	credit: number;
	creditScore: number;
};
