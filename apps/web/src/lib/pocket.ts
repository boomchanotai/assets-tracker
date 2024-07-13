import { pockets } from './constants/pocket';

export const getPocketName = (targetId: string | null) => {
	if (!targetId) return '';

	switch (targetId) {
		case 'cashbox':
			return 'Cashbox';
		case 'trash':
			return 'ใช้จ่าย';
		default:
			return pockets.filter(({ id }) => id === targetId)[0].name;
	}
};
