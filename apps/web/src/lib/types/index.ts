export type Currency = 'THB' | 'USD' | 'JPY' | 'EUR' | 'CNY' | 'BTC' | 'ETH' | 'BNB';

export type SelectOptions = {
	label: string;
	value: string;
};

export type Account = {
	id: string;
	userId: string;
	type: string;
	name: string;
	bank: string;
	balance: string;
	createdAt: string;
	updatedAt: string;
};
