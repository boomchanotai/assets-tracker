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
		}

		subscribe((value) => localStorage.setItem('auth', JSON.stringify(value)));
	}

	return {
		subscribe,
		set: (auth: Auth) => update(() => auth)
	};
};

export const authStore = createAuthStore();
