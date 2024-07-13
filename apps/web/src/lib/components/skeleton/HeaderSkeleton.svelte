<script lang="ts">
	import Container from '../Container.svelte';
</script>

<Container class="py-4 space-y-4 animate-pulse">
	<div class="flex flex-row justify-between gap-6">
		<div class="bg-gray-200 rounded-full w-48 h-6"></div>
		<div class="bg-gray-200 w-12 h-10 rounded"></div>
	</div>
	<div class="flex items-center gap-4">
		<div
			class="flex flex-row gap-6 items-center border border-gray-200 px-4 py-2 rounded-lg hover:bg-black/10 transition-colors duration-150 h-auto w-64"
		>
			<div class="flex flex-row gap-2 items-center">
				<div class="rounded-full bg-gray-200 size-6"></div>
				<div class="flex flex-col justify-start items-start gap-2">
					<div class="bg-gray-200 rounded-full w-32 h-3.5"></div>
					<div class="bg-gray-200 rounded-full w-24 h-3.5"></div>
				</div>
			</div>
		</div>
		<div class="bg-gray-200 w-12 h-10 rounded"></div>
	</div>
</Container>
