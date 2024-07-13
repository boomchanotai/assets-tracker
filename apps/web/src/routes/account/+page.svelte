<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Container from '@/components/Container.svelte';
	import Pocket from '@/components/Pocket.svelte';
	import Balance from '@/components/Account/Balance.svelte';
	import Trash from '@/components/Trash.svelte';
	import { pockets } from '@/constants/pocket';
	import { useAccounts } from '@/hook/queries/account';
	import HeaderSkeleton from '@/components/skeleton/HeaderSkeleton.svelte';
	import { accountStore } from '@/store/account';
	import Account from '@/components/Account.svelte';
	import AccountSkeleton from '@/components/skeleton/AccountSkeleton.svelte';

	$: accounts = useAccounts();

	let accountId: string | null = null;
	accountStore.subscribe((value) => {
		accountId = value;
	});
</script>

<div class="space-y-4">
	{#if $accounts.isFetching}
		<HeaderSkeleton />
	{:else}
		<Header accounts={$accounts.data?.result ?? []} />
	{/if}

	{#if accountId === null}
		<AccountSkeleton />
	{:else}
		<Account {accountId} />
	{/if}

	<Container class="space-y-8">
		<div class="grid grid-cols-12 gap-4">
			<div class="col-span-9 h-full">
				<Pocket id={'cashbox'} name={'Cashbox'} amount={13000} />
			</div>
			<div class="col-span-3">
				<Trash />
			</div>
		</div>
		<div class="grid grid-cols-2 gap-8">
			{#each pockets as pocket}
				<Pocket class="aspect-square" id={pocket.id} name={pocket.name} amount={pocket.amount} />
			{/each}
		</div>
	</Container>
</div>
