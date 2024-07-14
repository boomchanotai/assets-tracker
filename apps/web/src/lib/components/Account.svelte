<script lang="ts">
	import { useAccount } from '@/hook/queries/account';
	import Balance from './Account/Balance.svelte';
	import BalanceSkeleton from './skeleton/BalanceSkeleton.svelte';
	import Pockets from './Account/Pockets.svelte';
	import PocketsSkeleton from './skeleton/PocketsSkeleton.svelte';

	export let accountId: string;

	$: currentAccount = useAccount({ id: accountId });
	$: cashboxPocket = $currentAccount.data?.result.pockets.find(
		(pocket) => pocket.type === 'CASHBOX'
	);
</script>

{#if $currentAccount.isFetching}
	<BalanceSkeleton />
	<PocketsSkeleton />
{:else}
	<Balance
		{accountId}
		amount={$currentAccount.data?.result.balance
			? parseFloat($currentAccount.data?.result.balance)
			: 0}
	/>
	{#if $currentAccount.data?.result.pockets}
		<Pockets {accountId} {cashboxPocket} pockets={$currentAccount.data?.result.pockets} />
	{/if}
{/if}
