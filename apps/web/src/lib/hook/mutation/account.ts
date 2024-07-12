import { createAccount } from '@/api/account';
import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';

export const useCreateAccountMutation = () => {
	const client = useQueryClient();

	return createMutation({
		mutationFn: createAccount,
		onSuccess() {
			client.invalidateQueries({
				queryKey: ['accounts']
			});
			toast.success('Account created');
		}
	});
};
