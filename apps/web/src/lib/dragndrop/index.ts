export function dragstart(e: DragEvent, id: string) {
	e.dataTransfer?.setData('id', id);
	navigator.vibrate(200);
}

export function dragover(e: DragEvent) {
	e.preventDefault();
	if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
}

export function dragenter(e: DragEvent) {
	e.preventDefault();
	if (e.target instanceof HTMLButtonElement) {
		e.target.classList.add('border-blue-500');
	}
}

export function dragleave(e: DragEvent) {
	e.preventDefault();
	if (e.target instanceof HTMLButtonElement) {
		e.target.classList.remove('border-blue-500');
	}
}

export function drop(
	e: DragEvent,
	id: string,
	setOpenBalance: (state: boolean) => void,
	move: (fromId: string, toId: string) => void
) {
	e.preventDefault();
	if (e.target instanceof HTMLButtonElement) {
		e.target.classList.remove('border-blue-500');
	}
	if (!e.dataTransfer) {
		console.error('No dataTransfer');
		return;
	}

	const fromId = e.dataTransfer.getData('id');
	const toId = id;

	if (fromId === toId) return;

	if (fromId === toId) return;

	move(fromId, toId);
	setOpenBalance(true);
}
