import { createPocket, transfer, widthdraw } from '@/api/pocket';
import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';

export const useCreatePocketMutation = ({ accountId }: { accountId: string }) => {
	const client = useQueryClient();

	return createMutation({
		mutationFn: createPocket,
		onSuccess() {
			client.invalidateQueries({
				queryKey: ['account', accountId]
			});
			toast.success('Pocket created');
		}
	});
};

export const useTransferMutation = ({ accountId }: { accountId: string }) => {
	const client = useQueryClient();

	return createMutation({
		mutationFn: transfer,
		onSuccess() {
			client.invalidateQueries({
				queryKey: ['account', accountId]
			});
			toast.success('Transfer success');
		}
	});
};

export const useWithdrawMutation = ({ accountId }: { accountId: string }) => {
	const client = useQueryClient();

	return createMutation({
		mutationFn: widthdraw,
		onSuccess() {
			client.invalidateQueries({
				queryKey: ['account', accountId]
			});
			toast.success('Withdraw success');
		}
	});
};
