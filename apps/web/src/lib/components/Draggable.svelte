<script lang="ts">
	import { dragstart, dragover, dragenter, dragleave, drop } from '@/dragndrop';
	import type { Pocket } from '@/types';
	import TransferBalance from './dialog/TransferBalance.svelte';

	export let pockets: Pocket[];
	export let accountId: string;
	export let draggableId: string;
	export let draggable: boolean;

	let fromPocket: string;
	let toPocket: string;

	let openBalance = false;
	function setOpenBalance(state: boolean) {
		openBalance = state;
	}

	function move(from: string, to: string) {
		fromPocket = from;
		toPocket = to;
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

	<TransferBalance {accountId} {fromPocket} {toPocket} {pockets} {openBalance} {setOpenBalance} />
</div>
