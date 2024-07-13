<script lang="ts">
	import scb from '$lib/images/scb.png';
	import kbank from '$lib/images/kbank.png';
	import bbl from '$lib/images/bbl.jpg';
	import Icon from '@iconify/svelte';

	import Container from '$lib/components/Container.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Select from '$lib/components/ui/select';

	import { authStore } from '@/store/auth';

	import { banks } from '$lib/constants/bank';
	import type { Account, SelectOptions } from '@/types';
	import NewAccount from './dialog/NewAccount.svelte';
	import { accountStore } from '@/store/account';

	export let accounts: Account[];

	let selectedAccount: SelectOptions | undefined;

	if (accounts.length === 0) {
		selectedAccount = undefined;
		accountStore.clear();
	} else {
		selectedAccount = {
			value: accounts[0].id,
			label: accounts[0].name
		};
		accountStore.set(selectedAccount?.value);
	}

	function getAccount(accountId: string) {
		return accounts.filter(({ id }) => id === accountId)[0];
	}

	function getLogo(bankId: string) {
		switch (bankId) {
			case 'scb':
				return scb;
			case 'kbank':
				return kbank;
			case 'bbl':
				return bbl;
			default:
				return '';
		}
	}

	function getBankName(bankId: string) {
		if (bankId === '') return '';

		const b = banks.filter(({ value }) => value === bankId);
		return b.length > 0 ? b[0].label : '';
	}
</script>

<Container class="py-4 space-y-4">
	<div class="flex flex-row justify-between gap-6">
		<h1 class="text-xl font-semibold">Cloud Pocket</h1>
		<div>
			<Button on:click={() => authStore.clear()}><Icon icon="ph:sign-out-bold" /></Button>
		</div>
	</div>
	<div class="flex items-center gap-4">
		<div>
			<Select.Root bind:selected={selectedAccount}>
				<Select.Trigger
					class="flex flex-row gap-6 items-center border border-black px-4 py-2 rounded-lg hover:bg-black/10 transition-colors duration-150 h-auto w-64"
				>
					{#if selectedAccount !== undefined}
						<div class="flex flex-row gap-2 items-center">
							<div>
								<img
									src={getLogo(getAccount(selectedAccount.value).bank)}
									alt="scb"
									class="size-6 rounded-full"
								/>
							</div>
							<div class="flex flex-col justify-start items-start">
								<p class="text-sm font-medium">{getAccount(selectedAccount.value).name}</p>
								<p class="text-xs">{getBankName(getAccount(selectedAccount.value).bank)}</p>
							</div>
						</div>
					{:else}
						<div class="flex flex-row gap-2 items-center">
							<div class="rounded-full bg-gray-200 size-6"></div>
							<div class="flex flex-col justify-start items-start">
								<p class="text-sm font-medium">Select Account</p>
								<p class="text-xs">เลือกบัญชี</p>
							</div>
						</div>
					{/if}
				</Select.Trigger>
				<Select.Content class="max-h-64 overflow-y-auto">
					{#each accounts as account}
						<Select.Item value={account.id}>
							<div class="flex flex-row gap-2 items-center">
								<div><img src={getLogo(account.bank)} alt="scb" class="size-6 rounded-full" /></div>
								<div class="flex flex-col justify-start items-start">
									<p class="text-sm font-medium">{account.name}</p>
									<p class="text-xs">{getBankName(account.bank)}</p>
								</div>
							</div>
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
		<div>
			<NewAccount />
		</div>
	</div>
</Container>
