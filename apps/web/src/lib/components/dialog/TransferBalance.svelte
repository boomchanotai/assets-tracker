<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Icon from '@iconify/svelte';
	import type { Pocket } from '$lib/types';

	export let openBalance = false;
	export let setOpenBalance: (state: boolean) => void;

	export let fromPocket: string | null;
	export let toPocket: string | null;
	export let pockets: Pocket[];

	const getPocketName = (targetId: string | null) => {
		if (!targetId) return '';

		switch (targetId) {
			case 'cashbox':
				return 'Cashbox';
			case 'trash':
				return 'Out';
			default:
				return pockets.filter(({ id }) => id === targetId)[0].name;
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
		<div class="flex flex-col justify-center items-center gap-4">
			<Input id="balance" type="number" placeholder="Amount" />
			<div>
				<Button type="submit">Save changes</Button>
			</div>
		</div>
	</Dialog.Content>
</Dialog.Root>
