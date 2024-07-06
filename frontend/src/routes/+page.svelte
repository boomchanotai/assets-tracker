<script lang="ts">
	import Container from '@/components/Container.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
</script>

<div class="space-y-4">
	<Container class="h-svh flex flex-col justify-center items-center">
		<Card.Root class="w-[350px]">
			<Card.Header>
				<Card.Title class="text-center">Sign in</Card.Title>
				<Card.Description class="text-center">Sign in to your account to continue</Card.Description>
			</Card.Header>
			<Card.Content>
				<form>
					<div class="grid w-full items-center gap-4">
						<div class="flex flex-col space-y-1.5">
							<Label for="name">Email</Label>
							<Input id="email" type="email" placeholder="Email | อีเมล" />
						</div>
						<div class="flex flex-col space-y-1.5">
							<Label for="name">Password</Label>
							<Input id="password" type="password" placeholder="Password | พาสเวิร์ด" />
						</div>
					</div>
				</form>
			</Card.Content>
			<Card.Footer class="flex justify-center">
				<Button>Sign in</Button>
			</Card.Footer>
		</Card.Root>
	</Container>
</div>
