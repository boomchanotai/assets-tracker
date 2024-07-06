<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Container from '@/components/Container.svelte';
	import Pocket from '@/components/Pocket.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input/index.js';
	import Balance from '@/components/Balance.svelte';
	import Trash from '@/components/Trash.svelte';
	import { pockets } from '@/constants/pocket';
</script>

<div class="space-y-4">
	<Header />

	<Container class="flex flex-row justify-between">
		<Balance amount={14500} />
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
	</Container>

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
