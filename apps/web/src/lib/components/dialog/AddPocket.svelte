<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input/index.js';
	import Button from '@/components/ui/button/button.svelte';
	import { useCreatePocketMutation } from '@/hook/mutation/pocket';
	import Icon from '@iconify/svelte';

	export let accountId: string;
	let name: string;

	const addPocketMutation = useCreatePocketMutation({ accountId });
	const handleSubmit = (event: SubmitEvent) => {
		event.preventDefault();

		$addPocketMutation.mutate({ accountId: accountId, name: name });
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button class="gap-2">
			<Icon icon="ph:plus-bold" /> Add Pocket
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header class="mb-4">
			<Dialog.Title>Add Pocket</Dialog.Title>
		</Dialog.Header>
		<form on:submit={handleSubmit} class="grid gap-4 py-4">
			<Input bind:value={name} id="balance" type="text" placeholder="Name" />
			<Button type="submit">Save changes</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
