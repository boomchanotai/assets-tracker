<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input/index.js';
	import Button from '@/components/ui/button/button.svelte';
	import { useDepositMutation } from '@/hook/mutation/account';

	export let accountId: string;
	let balance: number = 0;

	const depositMutation = useDepositMutation({ accountId });
	const handleSubmit = (event: SubmitEvent) => {
		event.preventDefault();

		$depositMutation.mutate({ id: accountId, amount: balance });
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button>Add Balance</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header class="mb-4">
			<Dialog.Title>Add Balance</Dialog.Title>
		</Dialog.Header>
		<form on:submit={handleSubmit} class="grid gap-4 py-4">
			<Input bind:value={balance} id="balance" type="number" placeholder="Amount" />
			<Button type="submit">Save changes</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
