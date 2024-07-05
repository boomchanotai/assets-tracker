<script lang="ts">
	import type { Currency } from '@/types';
	import { cn } from '@/utils';
	import type { DragEventHandler } from 'svelte/elements';

	type Props = {
		id: string;
		name: string;
		amount: number;
		currency?: Currency;
	};

	export let id: Props['id'];
	export let name: Props['name'];
	export let amount: Props['amount'];
	export let currency: Props['currency'] = 'THB';

	export function dragstart(e: DragEvent, id: string) {
		e.dataTransfer?.setData('id', id);
	}
	export function dragover(e: DragEvent) {
		e.preventDefault();
		if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
	}
	export function dragenter(e: DragEvent) {
		e.preventDefault();
		if (e.target instanceof HTMLButtonElement) {
			e.target.classList.add('bg-gray-100');
		}
	}
	export function dragleave(e: DragEvent) {
		e.preventDefault();
		if (e.target instanceof HTMLButtonElement) {
			e.target.classList.remove('bg-gray-100');
		}
	}
	export function drop(e: DragEvent) {
		e.preventDefault();
		if (e.target instanceof HTMLButtonElement) {
			e.target.classList.remove('bg-gray-100');
		}
		const fromId = e.dataTransfer?.getData('id');
		const toId = id;
		console.log(fromId, toId);
	}
</script>

<button
	draggable={true}
	aria-grabbed={true}
	class={cn(
		'w-full h-full flex flex-col justify-start items-start border border-black aspect-square rounded-lg p-4',
		$$props.class
	)}
	on:dragstart={(e) => dragstart(e, id)}
	on:dragover={dragover}
	on:dragenter={dragenter}
	on:dragleave={dragleave}
	on:drop={drop}
>
	<h3>{name}</h3>
	<p>$ {amount} {currency}</p>
</button>
