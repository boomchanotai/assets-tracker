import { goto } from '$app/navigation';
import { writable } from 'svelte/store';

interface Auth {
	accessToken: string;
	refreshToken: string;
	exp: number;
}

const initialAuth = {
	accessToken: '',
	refreshToken: '',
	exp: 0
};

const createAuthStore = () => {
	const { subscribe, update } = writable<Auth>(initialAuth);

	if (typeof localStorage !== 'undefined') {
		const auth = localStorage.getItem('auth');
		if (auth) {
			update(() => JSON.parse(auth));
		} else {
			goto('/');
		}

		subscribe((value) => localStorage.setItem('auth', JSON.stringify(value)));
	}

	const set = (auth: Auth) => update(() => auth);

	const clear = () => {
		update(() => initialAuth);
		if (localStorage !== undefined) {
			localStorage.removeItem('auth');
			goto('/');
		}
	};

	return {
		subscribe,
		set,
		clear
	};
};

export const authStore = createAuthStore();
