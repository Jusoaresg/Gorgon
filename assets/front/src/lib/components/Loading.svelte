<script>
import { fade } from "svelte/transition";

export let text = '';
</script>

<div class="loading-state" in:fade={{ duration: 300 }}>
	<div class="loading-spinner">
		<div class="spinner-ring"></div>
		<div class="spinner-ring"></div>
		<div class="spinner-ring"></div>
	</div>
	<p class="loading-text">{text}</p>
</div>

<style>
.loading-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 4rem 2rem;
	text-align: center;
}

.loading-spinner {
	position: relative;
	width: 60px;
	height: 60px;
	margin-bottom: 1.5rem;
}

.spinner-ring {
	position: absolute;
	width: 100%;
	height: 100%;
	border: 3px solid transparent;
	border-top: 3px solid var(--highlight, #6366f1);
	border-radius: 50%;
	animation: spin 1s linear infinite;
}

.spinner-ring:nth-child(2) {
	width: 80%;
	height: 80%;
	top: 10%;
	left: 10%;
	border-top-color: var(--accent, #8b5cf6);
	animation-duration: 1.5s;
	animation-direction: reverse;
}

.spinner-ring:nth-child(3) {
	width: 60%;
	height: 60%;
	top: 20%;
	left: 20%;
	border-top-color: var(--text-secondary, #94a3b8);
	animation-duration: 2s;
}

@keyframes spin {
	0% {
		transform: rotate(0deg);
	}

	100% {
		transform: rotate(360deg);
	}
}

.loading-text {
	font-size: 1.125rem;
	color: var(--text-secondary, #94a3b8);
	margin: 0;
}
</style>
