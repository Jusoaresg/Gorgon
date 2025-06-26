<script>
import { goto } from '$app/navigation';
import { PUBLIC_API_BASE_URL } from '$env/static/public'
import SearchForm from '$lib/components/SearchForm.svelte'
import { fade, fly, scale } from 'svelte/transition';
import { flip } from 'svelte/animate';
import { quintOut } from 'svelte/easing';

// Styles
import '$lib/styles/pages/add-show.css';
import '$lib/styles/components/search-form.css';

let shows = [];
let currentShow = null;
let isLoading = false;

async function fetchShows(search) {
	isLoading = true;
	try {
		const resp = await fetch(`${PUBLIC_API_BASE_URL}/tvmaze/search/name`, {
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
};

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
</script>

<div class="add-show-container">
	<h1 in:fade={{ duration: 600, delay: 100 }}>Add Show</h1>
	
	<div in:fly={{ y: 20, duration: 500, delay: 200 }}>
		<SearchForm placeholder="Enter show name..." hasButton=true on:search={e => fetchShows(e.detail)} />
	</div>
	
	<div id="search-results" class="results-container" in:fly={{ y: 30, duration: 500, delay: 300 }}>
		<!-- Loading state -->
		{#if isLoading}
			<div class="loading-state" in:fade={{ duration: 200 }}>
				<div class="loading-spinner"></div>
				<p>Searching for shows...</p>
			</div>
		{/if}

		<!-- Popup modal -->
		{#if currentShow != null}
			<div class="popup-backdrop" in:fade={{ duration: 200 }} out:fade={{ duration: 200 }}>
				<div class="popup" in:scale={{ duration: 300, easing: quintOut }} out:scale={{ duration: 200 }}>
					<h2>Editando: {currentShow.Show.name}</h2>
					<p>{@html currentShow.Show.summary}</p>
					<button on:click={() => currentShow = null}>Fechar</button>
				</div>
			</div>
		{/if}

		<!-- Results -->
		{#if shows.length > 0 && !isLoading}
			{#each shows as show, i (show.Show.id)}
				<div 
					class="show-result"
					in:fly={{ y: 30, duration: 400, delay: i * 50 }}
					out:fly={{ y: -20, duration: 200 }}
					animate:flip={{ duration: 300 }}
				>
					<div class="show-image">
						<img 
							src={show.Show.image?.medium} 
							alt={show.Show.name}
							in:fade={{ duration: 300, delay: 100 }}
						/>
					</div>
					<div class="show-info">
						<h3 class="show-title">{show.Show.name}</h3>
						{@html show.Show.summary}
					</div>
					<button 
						class="add-button" 
						class:added={show.IsAdded} 
						on:click={() => addShow(show.Show.id)} 
						disabled={show.IsAdded}
						in:scale={{ duration: 200, delay: 200 }}
					>
						{#if show.IsAdded}
							<span in:fade={{ duration: 200 }}>Added</span>
						{:else}
							<span in:fade={{ duration: 200 }}>Add</span>
						{/if}
					</button>
				</div>
			{/each}
		{:else if !isLoading}
			<div class="empty-state" in:fade={{ duration: 400, delay: 200 }}>
				<p>Search for a show to add it to your collection</p>
			</div>
		{/if}
	</div>
</div>
