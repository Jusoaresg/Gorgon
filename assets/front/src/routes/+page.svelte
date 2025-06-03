<script>
import SearchForm from '$lib/components/SearchForm.svelte';
import { goto } from '$app/navigation';
import { fade, fly, scale } from 'svelte/transition';
import { onMount } from 'svelte';

let allShows = [];
let shows = [];

function navigateToShow(showId) {
	goto(`/show/${showId}`);
}

function handleSearch(event) {
	const query = event.detail.toLowerCase();

	shows = allShows.filter(show => show.Name.toLowerCase().includes(query));
}

onMount(async () => {
	const res = await fetch("/api/v1/database/show")
	const json = await res.json();

	allShows = json.data
	shows = json.data
})
</script>

<div class="home-container">
	<h1>My Shows</h1>
	<div class="categories">
		<SearchForm placeholder="Filter the show" on:search={handleSearch} live=true />
		<!-- Add extra filters for release data and etc... -->
	</div>

	{#if shows.length > 0}
		<div class="shows-grid">
			{#each shows as show}
				<div 
					aria-hidden=true 
					class="show-grid-item" 
					on:click={() => navigateToShow(show.ID) } 
					style="cursor: pointer;" 
					transition:fly={{ y: 20, duration: 200 }}
				>
					<input type="hidden" value={ show.ID.toString() } name="id"/>
					<div class="show-poster">
						{#if show.Image.Medium != ""}
							<img src={show.Image.Medium}  alt={show.Name}/>
						{/if}
						<div class="show-overlay">
						</div>
					</div>
					<div class="show-grid-info">
						<h3 class="show-grid-title">{show.Name}</h3>
					</div>
				</div>

			{/each}
		</div>
	{:else}
		<div class="empty-state">
			<p>You haven't added any shows yet. <a href="/add-show">Add your first show</a>.</p>
		</div>
	{/if}
</div>
