import { createAccount, deposit } from '@/api/account';
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

export const useDepositMutation = ({ accountId }: { accountId: string }) => {
	const client = useQueryClient();

	return createMutation({
		mutationFn: deposit,
		onSuccess() {
			client.invalidateQueries({
				queryKey: ['account', accountId]
			});
			toast.success('Deposit success');
		}
	});
};
