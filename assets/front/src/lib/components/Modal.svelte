<script>
	import { createEventDispatcher, onDestroy, onMount } from "svelte";
	
	export let show = false;
	export let body = true;
	export let footer = true;
	
	// Size configuration options
	export let size = "medium"; // "small", "medium", "large", "xl", "full", "custom"
	export let width = null; // Custom width (e.g., "800px", "90vw")
	export let height = null; // Custom height (e.g., "600px", "80vh")
	export let maxWidth = null; // Custom max-width
	export let maxHeight = null; // Custom max-height
	
	const dispatch = createEventDispatcher();
	
	let isClosing = false;
	let modalVisible = false;
	
	// Predefined size configurations
	const sizeConfigs = {
		small: { maxWidth: "400px", maxHeight: "60vh" },
		medium: { maxWidth: "700px", maxHeight: "70vh" },
		large: { maxWidth: "900px", maxHeight: "80vh" },
		xl: { maxWidth: "1200px", maxHeight: "85vh" },
		full: { maxWidth: "95vw", maxHeight: "95vh" },
		custom: {} // Will use the custom width/height props
	};
	
	// Compute modal styles based on size configuration
	$: modalStyles = (() => {
		const config = sizeConfigs[size] || sizeConfigs.medium;
		const styles = {};
		
		// Apply predefined config
		if (config.maxWidth) styles['max-width'] = config.maxWidth;
		if (config.maxHeight) styles['max-height'] = config.maxHeight;
		
		// Override with custom props if provided
		if (width) styles.width = width;
		if (height) styles.height = height;
		if (maxWidth) styles['max-width'] = maxWidth;
		if (maxHeight) styles['max-height'] = maxHeight;
		
		// Convert to CSS string
		return Object.entries(styles)
			.map(([key, value]) => `${key}: ${value}`)
			.join('; ');
	})();
	
	// Watch for show prop changes
	$: if (show && !modalVisible) {
		modalVisible = true;
		isClosing = false;
	} else if (!show && modalVisible && !isClosing) {
		startCloseAnimation();
	}
	
	const startCloseAnimation = () => {
		isClosing = true;
		// Wait for animation to complete before hiding modal
		setTimeout(() => {
			modalVisible = false;
			isClosing = false;
		}, 200); // Match the animation duration
	};
	
	const close = () => {
		dispatch("close");
	};
	
	const handleKeydown = (event) => {
		if (event.key === 'Escape') {
			close();
		}
	};
	
	export function clickOutside(node, callback) {
		function handleClick(event) {
			if (!node.contains(event.target)) {
				callback();
			}
		}
		document.addEventListener('pointerdown', handleClick, true);
		return {
			destroy() {
				document.removeEventListener('pointerdown', handleClick, true);
			}
		};
	}

	$: if (show) {
		document.addEventListener('keydown', handleKeydown);
	} else {
		document.removeEventListener('keydown', handleKeydown);
	}

	onDestroy(() => {
		document.removeEventListener('keydown', handleKeydown);
	})
</script>

{#if modalVisible}
	<div 
		aria-hidden=true 
		class="modal-overlay" 
		class:closing={isClosing}
		on:click={close}
	>
		<div 
			aria-hidden=true 
			class="modal-container" 
			class:closing={isClosing}
			style={modalStyles}
			use:clickOutside={close} 
			on:click|stopPropagation
		>
			<div class="modal-header">
				<h2 class="modal-title">
					<slot name="title"/>
				</h2>
				<button aria-label="close" class="modal-close-btn" on:click={close} title="Close">
					<i class="fas fa-times"></i>
				</button>
			</div>
			
			{#if body}
			<div class="modal-body">
				<slot name="body"/>
			</div>
			{/if}
			
			{#if footer}
			<div class="modal-footer">
				<slot name="footer"/>
			</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.6);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
		animation: fadeIn 0.15s ease-out;
	}
	
	.modal-overlay.closing {
		animation: fadeOut 0.2s ease-in forwards;
	}
	
	.modal-container {
		background-color: var(--card-bg);
		border-radius: 12px;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
		width: 100%;
		/* max-width and max-height now set via inline styles */
		overflow: hidden;
		animation: slideIn 0.2s ease-out;
		border: 2px solid rgba(255, 255, 255, 0.1);
		display: flex;
		flex-direction: column;
	}
	
	.modal-container.closing {
		animation: slideOut 0.2s ease-in forwards;
	}
	
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1.5rem 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.02));
		flex-shrink: 0;
	}
	
	.modal-title {
		margin: 0;
		color: var(--text-color);
		font-size: 1.4rem;
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}
	
	.modal-close-btn {
		width: 36px;
		height: 36px;
		border-radius: 6px;
		background-color: rgba(255, 255, 255, 0.1);
		color: var(--text-color);
		border: none;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s ease;
	}
	
	.modal-close-btn:hover {
		background-color: rgba(231, 76, 60, 0.8);
		color: white;
		transform: scale(1.05);
	}
	
	.modal-body {
		padding: 2rem;
		overflow-y: auto;
		flex: 1;
		min-height: 0; /* Important for flexbox scrolling */
	}
	
	.modal-footer {
		padding: 1.5rem 2rem;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.02), rgba(255, 255, 255, 0.05));
		flex-shrink: 0;
	}
	
	/* Opening animations */
	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}
	
	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translateY(-20px) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}
	
	/* Closing animations */
	@keyframes fadeOut {
		from { opacity: 1; }
		to { opacity: 0; }
	}
	
	@keyframes slideOut {
		from {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
		to {
			opacity: 0;
			transform: translateY(-20px) scale(0.95);
		}
	}
	
	/* Responsive design */
	@media (max-width: 768px) {
		.modal-container {
			margin: 1rem;
			max-height: 95vh !important;
		}
		
		.modal-header,
		.modal-body,
		.modal-footer {
			padding: 1rem 1.5rem;
		}
	}
	
	/* Custom scrollbar */
	.modal-body::-webkit-scrollbar{
		width: 6px;
	}
	
	.modal-body::-webkit-scrollbar-track{
		background: rgba(255, 255, 255, 0.05);
		border-radius: 3px;
	}
	
	.modal-body::-webkit-scrollbar-thumb{
		background: rgba(255, 255, 255, 0.2);
		border-radius: 3px;
	}
	
	.modal-body::-webkit-scrollbar-thumb:hover{
		background: rgba(255, 255, 255, 0.3);
	}
</style>
