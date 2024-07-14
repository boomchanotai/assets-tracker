import { authStore } from '@/store/auth';
import { get } from 'svelte/store';

export const transfer = async ({
	fromId,
	toId,
	amount
}: {
	fromId: string;
	toId: string;
	amount: number;
}) => {
	const response = await fetch(`${import.meta.env.VITE_BASE_URL}/pocket/${fromId}/transfer`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: 'Bearer ' + get(authStore).accessToken
		},
		body: JSON.stringify({ toPocketId: toId, amount })
	});

	if (!response.ok) {
		throw new Error('Transfer failed');
	}

	return response.json();
};

export const widthdraw = async ({ id, amount }: { id: string; amount: number }) => {
	const response = await fetch(`${import.meta.env.VITE_BASE_URL}/pocket/${id}/withdraw`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: 'Bearer ' + get(authStore).accessToken
		},
		body: JSON.stringify({ amount })
	});

	if (!response.ok) {
		throw new Error('Withdraw failed');
	}

	return response.json();
};
