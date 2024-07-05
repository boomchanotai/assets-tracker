<script lang="ts">
	import { dragstart, dragover, dragenter, dragleave, drop } from '@/dragndrop';
	export let draggableId: string;
	export let draggable: boolean;
</script>

<button
	{draggable}
	aria-grabbed={true}
	on:dragstart={(e) => dragstart(e, draggableId)}
	on:dragover={dragover}
	on:dragenter={dragenter}
	on:dragleave={dragleave}
	on:drop={(e) => drop(e, draggableId)}
	class={$$props.class}
>
	<slot />
</button>
