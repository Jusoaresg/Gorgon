<script>
import EpisodeCard from '$lib/components/EpisodeCard.svelte'
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { onMount } from 'svelte';
import { page } from '$app/stores';
import { goto } from '$app/navigation';

let show = null;

onMount(async () => {
	const showId = $page.params.id;

	const res = await fetch(`/api/v1/database/show/${showId}`)
	const json = await res.json()
	show = json.data;
})



function deleteShow() {
	if (!confirm("Are you sure?")) return;

	fetch(`${PUBLIC_API_BASE_URL}/database/show`, {
		method: "DELETE",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({ id: show.ID })
	}).then(() => {
			goto('/');
			alert("Show deleted.");
		});
}

async function changeTrackingStatus(event) {
	const { id, newStatus } = event.detail

	try {
		const res = await fetch(`${PUBLIC_API_BASE_URL}/database/show/episode/status`, {
			method: 'POST',
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({ episode_id: id, tracking: newStatus })
		});

		if(!res.ok) {
			console.error("Failed to update episode status", res.status)
		}

		const data = await res.json()
		console.log("New Episode status", data)
	} catch(error) {
		console.error("Error trying to change episode status", error)
	}
} 
</script>

<svelte:head>
	<link href="/css/pages/show.css" rel="stylesheet" />
	<link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet" />
</svelte:head>

{#if show}
	<div class="show-container">
		<div class="show-header">
			<h1>{show.Name}</h1>
			<div class="show-meta">
				<span class="show-status">{show.Status}</span>
				{#if show.Premiered}
					<span class="show-date">{show.Premiered} - {show.Ended}</span>
				{/if}
				{#if show.Language}
					<span class="show-language">{show.Language}</span>
				{/if}
			</div>
			<div class="show-actions">
				<button class="action-btn refresh-btn">
					<i class="fas fa-sync-alt"></i> Refresh Show Info
				</button>
				<button class="action-btn search-all-btn">
					<i class="fas fa-search"></i> Search All Missing
				</button>
				<button class="action-btn delete-show-btn" on:click={deleteShow}>
					<i class="fas fa-trash"></i> Delete Show
				</button>
			</div>
		</div>

		<div class="show-overview">
			<div class="show-poster">
				{#if show.Image?.Medium}
					<img src={show.Image.Medium} alt={`${show.Name} Poster`} />
				{:else}
					<div class="placeholder-poster">
						<span>{show.Name.charAt(0)}</span>
					</div>
				{/if}
			</div>
			<div class="show-summary">
				{@html show.Summary}
			</div>
		</div>

		<div class="seasons-navigation">
			{#each show.Seasons as season}
				<a href={"#season-" + season.Number} class="season-pill">
					Season {season.Number}
				</a>
			{/each}
		</div>

		<!-- Temporadas em ordem decrescente -->
		{#each [...show.Seasons].reverse() as season}
			<div id={"season-" + season.Number} class="season-block">
				<div class="season-header">
					<h2>Season {season.Number}</h2>
					<div class="season-actions">
						<button class="action-btn season-search-btn" data-season={season.Number}>
							<i class="fas fa-search"></i> Search Season
						</button>
					</div>
				</div>

				<div class="episodes-container">
					{#each show.Episodes.filter(e => e.Season === season.Number) as episode}
						<EpisodeCard {episode} seasonNumber={season.Number} on:statusChange={changeTrackingStatus} />
					{/each}
				</div>
			</div>
		{/each}
	</div>

{/if}
