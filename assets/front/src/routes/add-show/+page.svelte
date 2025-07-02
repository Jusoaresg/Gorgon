<script>
import { goto } from '$app/navigation';
import SearchForm from '$lib/components/SearchForm.svelte';
import Loading from '$lib/components/Loading.svelte';
import { fade, fly, scale, blur } from 'svelte/transition';
import { flip } from 'svelte/animate';
import { quintOut, elasticOut } from 'svelte/easing';
import { onMount } from 'svelte';

import '$lib/styles/pages/add-show.css'

let shows = [];
let currentShow = null;
let isLoading = false;
let searchQuery = '';
let mounted = false;

onMount(() => {
	mounted = true;
});

async function fetchShows(search) {
	if (!search.trim()) {
		shows = [];
		return;
	}

	searchQuery = search;
	isLoading = true;

	try {
		const resp = await fetch("/api/v1/tvmaze/search/name", {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ name: search }),
		});

		if (!resp.ok) {
			console.error("Error on shows request:", resp.status);
			return;
		}

		const data = await resp.json();
		shows = data.data || [];
	} catch (err) {
		console.error("Error searching for shows:", err);
	} finally {
		isLoading = false;
	}
}

async function addShow(showID) {
	const show = shows.find(s => s.Show.id === showID);
	if (!show) return;

	try {
		sessionStorage.setItem('showToEdit', JSON.stringify({
			Show: show.Show,
			isAdded: show.IsAdded
		}));
		goto(`/add-show/${showID}/edit`);
	} catch (e) {
		console.error('Error adding show', e);
	}
}

function openShowDetails(show) {
	currentShow = show;
}

function closeModal() {
	currentShow = null;
}

// Handle escape key for modal
function handleKeydown(event) {
	if (event.key === 'Escape' && currentShow) {
		closeModal();
	}
}
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="add-show">
<div class="add-show-container">
	{#if mounted}
		<header in:fade={{ duration: 800, delay: 100 }}>
			<h1 class="page-title">
				<span class="title-gradient">Discover</span>
				<span class="title-accent">Your Next Show</span>
			</h1>
			<p class="page-subtitle">
				Search through thousands of TV shows and add them to your list
			</p>
		</header>

		<div class="search-section" in:fly={{ y: 30, duration: 600, delay: 300 }}>
			<SearchForm 
				placeholder="Search for shows" 
				hasButton={true} 
				on:search={e => fetchShows(e.detail)} 
			/>
		</div>

		<main class="results-section" in:fly={{ y: 40, duration: 600, delay: 500 }}>
			{#if isLoading}
				<Loading text="Searching for {searchQuery}..."/>
			{:else if shows.length > 0}
				<div class="results-header" in:fade={{ duration: 400 }}>
					<h2>Found {shows.length} show{shows.length !== 1 ? 's' : ''}</h2>
				</div>

				<div class="shows-grid">
					{#each shows as show, i (show.Show.id)}
						<article 
							class="add-show-card"
							class:added={show.IsAdded}
							in:fly={{ y: 50, duration: 500, delay: i * 100 }}
							out:fly={{ y: -30, duration: 300 }}
							animate:flip={{ duration: 400, easing: quintOut }}
						>
							<div class="show-image-container">
								{#if show.Show.image?.medium}
									<img 
										src={show.Show.image.medium} 
										alt={show.Show.name}
										class="show-image"
										in:fade={{ duration: 400, delay: 200 }}
										loading="lazy"
									/>
								{:else}
									<div class="image-placeholder" in:fade={{ duration: 400 }}>
										<svg viewBox="0 0 24 24" fill="currentColor">
											<path d="M21 3H3c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h18c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM5 17l3.5-4.5 2.5 3.01L14.5 11l4.5 6H5z"/>
										</svg>
									</div>
								{/if}

								{#if show.IsAdded}
									<div class="added-badge" in:scale={{ duration: 300, easing: elasticOut }}>
										<svg viewBox="0 0 24 24" fill="currentColor">
											<path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
										</svg>
									</div>
								{/if}
							</div>

							<div class="show-content">
								<h3 class="show-title">{show.Show.name}</h3>

								{#if show.Show.genres?.length}
									<div class="show-genres">
										{#each show.Show.genres.slice(0, 3) as genre}
											<span class="genre-tag">{genre}</span>
										{/each}
									</div>
								{/if}

								{#if show.Show.rating?.average}
									<div class="show-rating">
										<span class="rating-star">★</span>
										<span class="rating-value">{show.Show.rating.average}/10</span>
									</div>
								{/if}

								{#if show.Show.summary}
									<div class="show-summary">
										{@html show.Show.summary.length > 150 
											? show.Show.summary.substring(0, 150) + '...' 
											: show.Show.summary}
									</div>
								{/if}
							</div>

							<div class="show-actions">
								<button 
									class="details-button"
									on:click={() => openShowDetails(show)}
								>
									View Details
								</button>

								<button 
									class="add-button" 
									class:added={show.IsAdded} 
									on:click={() => addShow(show.Show.id)} 
									disabled={show.IsAdded}
									in:scale={{ duration: 300, delay: 300 }}
								>
									{#if show.IsAdded}
										<span in:fade={{ duration: 200 }}>
											<svg viewBox="0 0 24 24" fill="currentColor">
												<path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
											</svg>
											Added
										</span>
									{:else}
										<span in:fade={{ duration: 200 }}>
											<svg viewBox="0 0 24 24" fill="currentColor">
												<path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
											</svg>
											Add Show
										</span>
									{/if}
								</button>
							</div>
						</article>
					{/each}
				</div>
			{:else if searchQuery}
				<div class="empty-state" in:fade={{ duration: 500, delay: 200 }}>
					<div class="empty-icon">
						<svg viewBox="0 0 24 24" fill="currentColor">
							<path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
						</svg>
					</div>
					<h3>No shows found</h3>
					<p>Try searching with different keywords or check your spelling</p>
				</div>
			{:else}
				<div class="welcome-state" in:fade={{ duration: 500, delay: 200 }}>
					<div class="welcome-icon">
						<svg viewBox="0 0 24 24" fill="currentColor">
							<path d="M21 3H3c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h18c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM5 17l3.5-4.5 2.5 3.01L14.5 11l4.5 6H5z"/>
						</svg>
					</div>
					<h3>Ready to discover?</h3>
					<p>Start by searching for your favorite TV shows above</p>
				</div>
			{/if}
		</main>
	{/if}
</div>

<!-- Modal -->
{#if currentShow}
	<div 
		aria-hidden=true 
		class="modal-backdrop" 
		in:fade={{ duration: 200 }} 
		out:fade={{ duration: 200 }}
		on:click={closeModal}
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-title"
	>
		<div 
			aria-hidden=true
			class="modal-content" 
			in:scale={{ duration: 300, easing: quintOut }} 
			out:scale={{ duration: 200 }}
			on:click|stopPropagation
		>
			<div class="modal-header">
				<h2 id="modal-title">{currentShow.Show.name}</h2>
				<button class="close-button" on:click={closeModal} aria-label="Close modal">
					<svg viewBox="0 0 24 24" fill="currentColor">
						<path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
					</svg>
				</button>
			</div>

			<div class="modal-body">
				{#if currentShow.Show.image?.original}
					<img src={currentShow.Show.image.original} alt={currentShow.Show.name} class="modal-image" />
				{/if}

				<div class="modal-info">
					{#if currentShow.Show.genres?.length}
						<div class="modal-genres">
							{#each currentShow.Show.genres as genre}
								<span class="genre-tag">{genre}</span>
							{/each}
						</div>
					{/if}

					{#if currentShow.Show.rating?.average}
						<div class="modal-rating">
							<span class="rating-star">★</span>
							<span class="rating-value">{currentShow.Show.rating.average}/10</span>
						</div>
					{/if}

					{#if currentShow.Show.summary}
						<div class="modal-summary">
							{@html currentShow.Show.summary}
						</div>
					{/if}
				</div>
			</div>

			<div class="modal-actions">
				<button class="modal-add-button" on:click={() => addShow(currentShow.Show.id)} disabled={currentShow.IsAdded}>
					{currentShow.IsAdded ? 'Already Added' : 'Add to List'}
				</button>
			</div>
		</div>
	</div>
{/if}
	</div>
