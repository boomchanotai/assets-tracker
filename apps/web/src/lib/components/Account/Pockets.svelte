<script lang="ts">
	import Container from '$lib/components/Container.svelte';
	import Pocket from '$lib/components/Pocket.svelte';
	import Trash from '$lib/components/Trash.svelte';
	import type { Pocket as PocketType } from '$lib/types';
	import AddPocket from '../dialog/AddPocket.svelte';

	export let accountId: string;
	export let cashboxPocket;
	export let pockets: PocketType[];

	const realPockets = pockets.filter((pocket) => pocket.type !== 'CASHBOX');
</script>

<Container class="space-y-8">
	<div class="grid grid-cols-12 gap-4">
		<div class="col-span-9 h-full">
			<Pocket
				{accountId}
				{pockets}
				id={cashboxPocket?.id}
				name={cashboxPocket?.name}
				amount={cashboxPocket?.balance ? parseFloat(cashboxPocket?.balance) : 0}
			/>
		</div>
		<div class="col-span-3">
			<Trash {accountId} {pockets} />
		</div>
	</div>
	<div class="flex flex-row justify-between items-center">
		<div class="font-semibold text-lg">Pockets</div>
		<AddPocket {accountId} />
	</div>
	{#if !pockets}
		<p class="text-center text-gray-500">No pockets found</p>
	{:else}
		<div class="grid grid-cols-2 gap-8">
			{#each realPockets as pocket}
				<Pocket
					{accountId}
					{pockets}
					class="aspect-square"
					id={pocket.id}
					name={pocket.name}
					amount={parseFloat(pocket.balance)}
				/>
			{/each}
		</div>
	{/if}
</Container>
