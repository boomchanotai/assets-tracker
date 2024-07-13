<script lang="ts">
	import Header from '$lib/components/Header.svelte';
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
</div>
