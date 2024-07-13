<script lang="ts">
	import Container from '$lib/components/Container.svelte';
	import Pocket from '$lib/components/Pocket.svelte';
	import Trash from '$lib/components/Trash.svelte';

	export let cashboxPocket;
	export let pockets;
</script>

<Container class="space-y-8">
	<div class="grid grid-cols-12 gap-4">
		<div class="col-span-9 h-full">
			<Pocket
				id="cashbox"
				name="Cashbox"
				amount={cashboxPocket?.balance ? parseFloat(cashboxPocket?.balance) : 0}
			/>
		</div>
		<div class="col-span-3">
			<Trash />
		</div>
	</div>
	{#if !pockets}
		<p class="text-center text-gray-500">No pockets found</p>
	{:else}
		<div class="grid grid-cols-2 gap-8">
			{#each pockets as pocket}
				<Pocket
					class="aspect-square"
					id={pocket.id}
					name={pocket.name}
					amount={parseFloat(pocket.balance)}
				/>
			{/each}
		</div>
	{/if}
</Container>
