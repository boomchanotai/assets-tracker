import { writable } from 'svelte/store';

const createAccountStore = () => {
	const { subscribe, update } = writable<string | null>(null);

	return {
		subscribe,
		set: (accountId: string) => update(() => accountId),
		clear: () => update(() => null)
	};
};

export const accountStore = createAccountStore();
