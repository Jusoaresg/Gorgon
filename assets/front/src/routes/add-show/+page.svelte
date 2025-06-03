<script>
import { goto } from '$app/navigation';

import { PUBLIC_API_BASE_URL } from '$env/static/public'
import SearchForm from '$lib/components/SearchForm.svelte'

let shows = [];
let currentShow = null;

async function fetchShows(search) {
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
	} catch {
		console.error("Error searching for shows:", err);
	}
};

async function addShow(showID) {
	const show = shows.find(s => s.Show.id === showID);
	if (!show) return;

	try {
		// const res = await fetch(`${PUBLIC_API_BASE_URL}/database/show`, {
		// 	method: 'POST',
		// 	headers: { 'Content-Type': 'application/json' },
		// 	body: JSON.stringify({ id: Number(showID) })
		// });
		//
		// if(!res.ok) {
		// 	console.error('Failed to add show', await res.text());
		// 	return
		// }
		//
		// show.IsAdded = true;
		// shows = [...shows];
		//
		// currentShow = show;

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
	<h1>Add Show</h1>
	<SearchForm placeholder="Enter show name..." hasButton=true on:search={e => fetchShows(e.detail)} />
	<div id="search-results" class="results-container">
		{#if currentShow != null}
			<div class="popup">
				<h2>Editando: {currentShow.Show.name}</h2>
				<p>{@html currentShow.Show.summary}</p>
				<button on:click={() => currentShow = null}>Fechar</button>
			</div>
		{/if}

		{#if shows.length > 0 }
			{#each shows as show}
				<div class="show-result">
					<div class="show-image">
						<img src={ show.Show.image.medium } alt={ show.Show.name }/>
					</div>
					<div class="show-info">
						<h3 class="show-title">{ show.Show.name }</h3>
						{@html show.Show.summary }
					</div>
					<button class="add-button" class:added={show.IsAdded} on:click={() => addShow(show.Show.id)} disabled={show.IsAdded}>
						{#if show.IsAdded}
							Added
						{:else}
							Add
						{/if}
					</button>
				</div>

			{/each}
		{:else}
			<div class="empty-state">
				<p>Search for a show to add it to your collection</p>
			</div>
		{/if}
	</div>

</div>
