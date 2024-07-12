import { getAccounts } from '@/api/account';
import type { Account } from '@/types';
import { createQuery } from '@tanstack/svelte-query';

interface getAccountsResponse {
	result: Account[];
}
export const useAccounts = () => {
	return createQuery<getAccountsResponse>({
		queryKey: ['accounts'],
		queryFn: getAccounts
	});
};
