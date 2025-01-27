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

export const getAccount = async ({ id }: { id: string }) => {
	const response = await fetch(`${import.meta.env.VITE_BASE_URL}/account/${id}`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json',
			Authorization: 'Bearer ' + get(authStore).accessToken
		}
	});

	if (!response.ok) {
		throw new Error('Account not found');
	}

	return response.json();
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

export const deposit = async ({ id, amount }: { id: string; amount: number }) => {
	const response = await fetch(`${import.meta.env.VITE_BASE_URL}/account/${id}/deposit`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: 'Bearer ' + get(authStore).accessToken
		},
		body: JSON.stringify({ amount })
	});

	if (!response.ok) {
		throw new Error('Deposit failed');
	}

	return response.json();
};
