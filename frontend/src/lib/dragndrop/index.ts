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

export function drop(e: DragEvent, id: string) {
	e.preventDefault();
	if (e.target instanceof HTMLButtonElement) {
		e.target.classList.remove('bg-gray-100');
	}
	const fromId = e.dataTransfer?.getData('id');
	const toId = id;
	console.log(fromId, toId);
}
