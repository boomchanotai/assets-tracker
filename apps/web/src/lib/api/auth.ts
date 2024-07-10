export const signin = async ({ email, password }: { email: string; password: string }) => {
	const reponse = await fetch(`${import.meta.env.VITE_BASE_URL}/auth/login`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ email, password })
	});

	if (!reponse.ok) {
		throw new Error('Login failed');
	}

	return reponse.json();
};

export const signup = async ({
	email,
	name,
	password
}: {
	email: string;
	name: string;
	password: string;
}) => {
	const reponse = await fetch(`${import.meta.env.VITE_BASE_URL}/auth/register`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ email, password, name })
	});

	if (!reponse.ok) {
		throw new Error('Registration failed');
	}

	return reponse.json();
};
