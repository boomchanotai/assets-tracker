<script lang="ts">
	import { useAccount } from '@/hook/queries/account';
	import Balance from './Account/Balance.svelte';
	import BalanceSkeleton from './skeleton/BalanceSkeleton.svelte';

	export let accountId: string;

	$: currentAccount = useAccount({ id: accountId });
</script>

{#if $currentAccount.isFetching}
	<BalanceSkeleton />
{:else}
	<Balance
		{accountId}
		amount={$currentAccount.data?.result.balance
			? parseFloat($currentAccount.data?.result.balance)
			: 0}
	/>
{/if}
