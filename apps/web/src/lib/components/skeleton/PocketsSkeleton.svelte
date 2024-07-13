<script>
	import Container from '../Container.svelte';
	import PocketSkeleton from './PocketSkeleton.svelte';
	import TrashSkeleton from './TrashSkeleton.svelte';
</script>

<Container class="space-y-8">
	<div class="grid grid-cols-12 gap-4 items-center">
		<div class="col-span-9 h-full">
			<PocketSkeleton />
		</div>
		<div class="col-span-3">
			<TrashSkeleton />
		</div>
	</div>
	<div class="grid grid-cols-2 gap-8">
		<div class="aspect-square bg-gray-200 rounded-lg"></div>
		<div class="aspect-square bg-gray-200 rounded-lg"></div>
		<div class="aspect-square bg-gray-200 rounded-lg"></div>
		<div class="aspect-square bg-gray-200 rounded-lg"></div>
		<div class="aspect-square bg-gray-200 rounded-lg"></div>
	</div>
</Container>
