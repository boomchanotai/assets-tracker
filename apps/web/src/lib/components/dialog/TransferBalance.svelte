<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Icon from '@iconify/svelte';
	import type { Pocket } from '$lib/types';
	import { useTransferMutation, useWithdrawMutation } from '@/hook/mutation/pocket';

	export let accountId: string;
	export let openBalance = false;
	export let setOpenBalance: (state: boolean) => void;

	export let fromPocket: string;
	export let toPocket: string;
	export let pockets: Pocket[];

	const getPocketName = (targetId: string) => {
		if (!targetId) return '';

		switch (targetId) {
			case 'TRASH':
				return 'Out';
			default:
				return pockets.filter(({ id }) => id === targetId)[0].name;
		}
	};

	const transferMutation = useTransferMutation({ accountId: accountId });
	const withdrawMutation = useWithdrawMutation({ accountId: accountId });
	let amount: number = 0;
	const handleSubmit = (event: SubmitEvent) => {
		event.preventDefault();
		if (toPocket === 'TRASH') {
			$withdrawMutation.mutate({ id: fromPocket, amount });
		} else {
			$transferMutation.mutate({ fromId: fromPocket, toId: toPocket, amount });
		}
	};
</script>

<Dialog.Root open={openBalance} onOpenChange={setOpenBalance}>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title class="flex flex-row justify-center gap-4 mb-4">
				<span>{getPocketName(fromPocket)}</span>
				<Icon icon="ph:arrow-right-bold" />
				<span> {getPocketName(toPocket)}</span>
			</Dialog.Title>
		</Dialog.Header>
		<form on:submit={handleSubmit} class="grid gap-4 py-4">
			<Input bind:value={amount} id="balance" type="number" placeholder="Amount" />
			<Button type="submit">Save changes</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
