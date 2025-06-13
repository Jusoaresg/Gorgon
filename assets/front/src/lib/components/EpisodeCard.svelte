<script>
import { createEventDispatcher, onMount, tick } from 'svelte';

export let episode;
export let seasonNumber;

const dispatch = createEventDispatcher();

let status = episode.Tracking;
$: status = episode.Tracking;
let showDropdown = false;
let dropdownElement;
let buttonElement;
let portalTarget;

const statusOptions = ['wanted', 'skipped'];

function changeStatus(newStatus) {
	status = newStatus;
	showDropdown = false;
	dispatch('statusChange', { id: episode.ID, newStatus });
}

async function toggleDropdown() {
	showDropdown = !showDropdown;

	if (showDropdown && buttonElement) {
		await tick(); // Wait for DOM update
		positionDropdown();
	}
}

function positionDropdown() {
	if (!buttonElement || !dropdownElement) return;

	const buttonRect = buttonElement.getBoundingClientRect();
	const dropdownRect = dropdownElement.getBoundingClientRect();

	// Position dropdown below button
	dropdownElement.style.position = 'fixed';
	dropdownElement.style.top = `${buttonRect.bottom + 8}px`;
	dropdownElement.style.left = `${buttonRect.right - dropdownRect.width}px`;

	// Ensure dropdown stays within viewport
	const viewportWidth = window.innerWidth;
	const viewportHeight = window.innerHeight;

	if (buttonRect.right - dropdownRect.width < 0) {
		dropdownElement.style.left = `${buttonRect.left}px`;
	}

	if (buttonRect.bottom + dropdownRect.height + 8 > viewportHeight) {
		dropdownElement.style.top = `${buttonRect.top - dropdownRect.height - 8}px`;
	}
}

function handleClickOutside(event) {
	if (dropdownElement && !dropdownElement.contains(event.target) && 
		buttonElement && !buttonElement.contains(event.target)) {
		showDropdown = false;
	}
}

function handleEscapeKey(event) {
	if (event.key === 'Escape') {
		showDropdown = false;
	}
}

onMount(() => {
	// Create portal target
	portalTarget = document.createElement('div');
	portalTarget.style.position = 'absolute';
	portalTarget.style.top = '0';
	portalTarget.style.left = '0';
	portalTarget.style.zIndex = '9999';
	document.body.appendChild(portalTarget);

	document.addEventListener('click', handleClickOutside);
	document.addEventListener('keydown', handleEscapeKey);
	window.addEventListener('scroll', () => {
		if (showDropdown) positionDropdown();
	});
	window.addEventListener('resize', () => {
		if (showDropdown) positionDropdown();
	});

	return () => {
		document.removeEventListener('click', handleClickOutside);
		document.removeEventListener('keydown', handleEscapeKey);
		if (portalTarget && portalTarget.parentNode) {
			portalTarget.parentNode.removeChild(portalTarget);
		}
	};
});

function getStatusClass(status) {
	switch (status.toLowerCase()) {
		case 'downloaded': return 'status-downloaded';
		case 'missing': return 'status-missing';
		case 'wanted': return 'status-wanted';
		case 'skipped': return 'status-skipped';
		case 'snatched': return 'status-snatched';
		default: return 'status-unknown';
	}
}

function getStatusIcon(status) {
	switch (status.toLowerCase()) {
		case 'downloaded': return 'fas fa-check';
		case 'missing': return 'fas fa-times';
		case 'wanted': return 'fas fa-clock';
		case 'skipped': return 'fas fa-ban';
		case 'snatched': return 'fas fa-magnet';
		default: return 'fas fa-question';
	}
}

function formatDate(dateStr) {
	if (!dateStr) return '';
	const date = new Date(dateStr);
	return date.toLocaleDateString('pt-BR', {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	});
}
</script>

<div class="episode-card" data-episode-id={episode.ID} data-status={status}>
	<div class={"episode-number " + getStatusClass(status)}>
		<span>{episode.Number}</span>
		<div class="absolute-number">
			{(seasonNumber - 1) * 10 + episode.Number}
		</div>
	</div>
	<div class="episode-content">
		<div class="episode-header">
			<h3 class="episode-title">
				{episode.Name}
			</h3>
			<div class="episode-actions">
				<button 
					aria-hidden="true" 
					class="episode-btn force-search-btn" 
					title="Force Search"
				>
					<i class="fas fa-search"></i>
				</button>
				<button 
					aria-hidden="true" 
					class="episode-btn manual-search-btn" 
					title="Manual Search"
				>
					<i class="fas fa-list"></i>
				</button>
				<button 
					bind:this={buttonElement}
					aria-hidden="true" 
					class={"episode-btn episode-status-toggle " + getStatusClass(status)} 
					title="Change Tracking Status"
					on:click={toggleDropdown}
				>
					<i class={getStatusIcon(status)}></i>
				</button>
			</div>
		</div>
		<div class="episode-info">
			{#if episode.AirStamp}
				<div class="episode-date mt-1">
					<span class="air-date">
						<i class="fas fa-calendar-alt"></i>
						{formatDate(episode.AirStamp)}
					</span>
				</div>
			{/if}
			{#if episode.Summary}
				<div class="episode-summary mt-2">
					{@html episode.Summary}
				</div>
			{/if}
		</div>
	</div>
</div>

<!-- Portal dropdown -->
{#if showDropdown && portalTarget}
	<div bind:this={dropdownElement} class="status-dropdown portal-dropdown">
		{#each statusOptions as option}
			<div
				class="dropdown-item"
				class:selected={option === status}
				data-status={option}
				on:click={() => changeStatus(option)}
				role="button"
				tabindex="0"
				on:keydown={(e) => e.key === 'Enter' && changeStatus(option)}
			>
				<i class={getStatusIcon(option)}></i>
				{option}
			</div>
		{/each}
	</div>
{/if}

<style>
:global(.portal-dropdown) {
	position: fixed !important;
	z-index: 9999 !important;
	background-color: var(--card-bg);
	border: 2px solid rgba(255, 255, 255, 0.1);
	border-radius: 8px;
	box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5);
	min-width: 140px;
	padding: 0.5rem 0;
	backdrop-filter: blur(10px);
	animation: dropdownFadeIn 0.15s ease-out;
}

:global(.portal-dropdown) {
	position: fixed !important;
	z-index: 9999 !important;
	background-color: var(--card-bg);
	border: 2px solid rgba(255, 255, 255, 0.1);
	border-radius: 8px;
	box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5);
	min-width: 140px;
	padding: 0.5rem 0;
	backdrop-filter: blur(10px);
	animation: portalDropdownFadeIn 0.15s ease-out;
	/* Ensure it appears above everything */
	pointer-events: auto;
}

:global(.portal-dropdown .dropdown-item) {
	display: flex;
	align-items: center;
	gap: 0.75rem;
	padding: 0.75rem 1rem;
	cursor: pointer;
	transition: all 0.2s ease;
	color: var(--text-color);
	font-size: 0.9rem;
	text-transform: capitalize;
	border: none;
	background: none;
	width: 100%;
	text-align: left;
	white-space: nowrap;
}

:global(.portal-dropdown .dropdown-item:hover) {
	background-color: rgba(255, 255, 255, 0.15);
	color: white;
}

:global(.portal-dropdown .dropdown-item.selected) {
	background-color: var(--highlight);
	color: white;
}

:global(.portal-dropdown .dropdown-item.selected:hover) {
	background-color: var(--highlight);
	opacity: 0.9;
}

:global(.portal-dropdown .dropdown-item i) {
	width: 14px;
	text-align: center;
	font-size: 0.85rem;
	flex-shrink: 0;
}

@keyframes portalDropdownFadeIn {
from {
	opacity: 0;
	transform: translateY(-8px) scale(0.95);
}
to {
	opacity: 1;
	transform: translateY(0) scale(1);
}
}

:global(.portal-dropdown .dropdown-item[data-status="downloaded"]:hover:not(.selected)) {
	background-color: rgba(39, 174, 96, 0.8);
	color: white;
}

:global(.portal-dropdown .dropdown-item[data-status="missing"]:hover:not(.selected)) {
	background-color: rgba(231, 76, 60, 0.8);
	color: white;
}

:global(.portal-dropdown .dropdown-item[data-status="wanted"]:hover:not(.selected)) {
	background-color: rgba(243, 156, 18, 0.8);
	color: white;
}

:global(.portal-dropdown .dropdown-item[data-status="skipped"]:hover:not(.selected)) {
	background-color: rgba(127, 140, 141, 0.8);
	color: white;
}

:global(.portal-dropdown .dropdown-item[data-status="snatched"]:hover:not(.selected)) {
	background-color: rgba(52, 152, 219, 0.8);
	color: white;
}

:global(.portal-dropdown .dropdown-item:focus) {
	outline: 2px solid var(--highlight);
	outline-offset: -2px;
	background-color: rgba(255, 255, 255, 0.1);
}

:global(.portal-dropdown) {
	z-index: 99999 !important;
}

:global(.portal-dropdown .dropdown-item:not(:last-child)) {
	border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

:global(.portal-dropdown .dropdown-item:first-child) {
	border-top-left-radius: 6px;
	border-top-right-radius: 6px;
}

:global(.portal-dropdown .dropdown-item:last-child) {
	border-bottom-left-radius: 6px;
	border-bottom-right-radius: 6px;
}

@media (max-width: 768px) {
	:global(.portal-dropdown) {
		min-width: 120px;
		font-size: 0.85rem;
	}

	:global(.portal-dropdown .dropdown-item) {
		padding: 0.6rem 0.8rem;
		gap: 0.6rem;
	}

	:global(.portal-dropdown .dropdown-item i) {
		width: 12px;
		font-size: 0.8rem;
	}
}

.episode-actions {
	position: relative;
	display: flex;
	gap: 0.5rem;
}

.episode-btn.episode-status-toggle:active {
	transform: scale(0.95);
}

/* Loading state for dropdown (optional) */
:global(.portal-dropdown.loading) {
	pointer-events: none;
	opacity: 0.7;
}

:global(.portal-dropdown.loading::after) {
	content: '';
	position: absolute;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	width: 20px;
	height: 20px;
	border: 2px solid transparent;
	border-top: 2px solid var(--highlight);
	border-radius: 50%;
	animation: spin 1s linear infinite;
}

@keyframes spin {
0% { transform: translate(-50%, -50%) rotate(0deg); }
100% { transform: translate(-50%, -50%) rotate(360deg); }
}
</style>
