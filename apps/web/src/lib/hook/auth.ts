import { goto } from '$app/navigation';
import { signin, signup } from '@/api/auth';
import { authStore } from '@/store/auth';
import { createMutation } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';

export const useLoginMutation = () => {
	return createMutation({
		mutationFn: signin,
		onSuccess(data) {
			authStore.set(data.result);
			toast.success('Login successful');
			goto('/account');
		}
	});
};

export const useRegisterMutation = () => {
	return createMutation({
		mutationFn: signup,
		onSuccess() {
			toast.success('Register successful');
		}
	});
};
