<script lang="ts">
	import { dragstart, dragover, dragenter, dragleave, drop } from '@/dragndrop';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { pockets } from '@/constants/pocket';
	import Icon from '@iconify/svelte';

	export let draggableId: string;
	export let draggable: boolean;

	let openBalance = false;
	let fromPocket: string | null;
	let toPocket: string | null;

	function setOpenBalance(state: boolean) {
		openBalance = state;
	}

	function move(from: string, to: string) {
		fromPocket = from;
		toPocket = to;
	}

	function getPocketName(targetId: string | null) {
		if (!targetId) return '';

		switch (targetId) {
			case 'cashbox':
				return 'Cashbox';
			case 'trash':
				return 'ใช้จ่าย';
			default:
				const pks = pockets.filter(({ id }) => id === targetId);
				if (pks.length === 0) return '';
				return pks[0].name;
		}
	}
</script>

<div class="h-full w-full">
	<button
		{draggable}
		aria-grabbed={true}
		on:dragstart={(e) => dragstart(e, draggableId)}
		on:dragover={dragover}
		on:dragenter={dragenter}
		on:dragleave={dragleave}
		on:drop={(e) => drop(e, draggableId, setOpenBalance, move)}
		class={$$props.class}
	>
		<slot />
	</button>

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
</div>
