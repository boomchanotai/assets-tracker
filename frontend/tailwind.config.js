/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				poppins: ['Poppins', 'sans-serif'],
				'ibm-plex-sans-thai': ['IBM Plex Sans Thai', 'sans-serif']
			}
		}
	},
	plugins: []
};
