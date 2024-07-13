<script lang="ts">
	import type { Currency } from '@/types';
	import Container from './Container.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input/index.js';
	import Button from '@/components/ui/button/button.svelte';
	import { useAccount } from '@/hook/queries/account';

	export let accountId: string;
	export let currency: Currency = 'THB';

	$: currentAccount = useAccount({ id: accountId });
</script>

<Container class="flex flex-row justify-between">
	<div>
		<h2 class="text-xs">Total</h2>
		<p class="font-semibold text-xl">$ {$currentAccount.data?.result.balance} {currency}</p>
	</div>
	<div>
		<Dialog.Root>
			<Dialog.Trigger>
				<Button>อัพเดตยอดเงิน</Button>
			</Dialog.Trigger>
			<Dialog.Content class="sm:max-w-[425px]">
				<Dialog.Header class="mb-4">
					<Dialog.Title>Update Balance</Dialog.Title>
				</Dialog.Header>
				<div>
					<Input id="balance" type="number" placeholder="Amount" />
				</div>
				<Dialog.Footer>
					<Button type="submit">Save changes</Button>
				</Dialog.Footer>
			</Dialog.Content>
		</Dialog.Root>
	</div>
</Container>
