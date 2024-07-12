import { authStore } from '@/store/auth';
import { get } from 'svelte/store';

export const getAccounts = async () => {
	const reponse = await fetch(`${import.meta.env.VITE_BASE_URL}/account`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json',
			Authorization: 'Bearer ' + get(authStore).accessToken
		}
	});

	if (!reponse.ok) {
		throw new Error('Account not found');
	}

	return reponse.json();
};

export const createAccount = async ({
	type,
	name,
	bank
}: {
	type: string;
	name: string;
	bank: string;
}) => {
	const response = await fetch(`${import.meta.env.VITE_BASE_URL}/account`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: 'Bearer ' + get(authStore).accessToken
		},
		body: JSON.stringify({ type, name, bank })
	});

	if (!response.ok) {
		throw new Error('Account creation failed');
	}

	return response.json();
};
