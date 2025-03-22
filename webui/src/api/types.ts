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
