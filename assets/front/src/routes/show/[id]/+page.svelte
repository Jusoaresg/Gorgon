<script>
import EpisodeCard from '$lib/components/EpisodeCard.svelte'
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { onMount } from 'svelte';
import { page } from '$app/stores';
import { goto } from '$app/navigation';
import { tick } from 'svelte';
import { toast } from 'svelte-sonner';

// Styles
import '$lib/styles/pages/show.css';
import '$lib/styles/components/episode-card.css';

	let episodeDropdown = false;
	let dropdownPos = { top: 0, left: 0 };
	let dropdownEpisodeId = null;
	let dropdownOptions = [];
	let dropdownStatus = null;

	function handleDropdownToggle(e) {
		dropdownEpisodeId = e.detail.id;
		dropdownPos = e.detail.position;
		dropdownOptions = e.detail.options;
		dropdownStatus = e.detail.currentStatus;
		episodeDropdown = true;
	}

let show = null;
let seasons = [];
let episodes = [];

// Season fold/unfold state
let collapsedSeasons = new Set();

let socket;

onMount(async () => {
	const showId = $page.params.id;

	try {
		const [showRes, seasonsRes, episodesRes] = await Promise.all([
			fetch(`/api/v1/database/show/${showId}`),
			fetch(`/api/v1/database/show/season/${showId}`),
			fetch(`/api/v1/database/show/episode/${showId}`)
		])

		if (!showRes.ok || !seasonsRes.ok || !episodesRes.ok) {
			console.error("One or more requests failed", {
				show: showRes.status,
				seasons: seasonsRes.status,
				episodes: episodesRes.status
			});
			return;
		}

		const [showJson, seasonsJson, episodesJson] = await Promise.all([
			showRes.json(),
			seasonsRes.json(),
			episodesRes.json(),
		])

		show = showJson.data;
		seasons = seasonsJson.data;
		episodes = episodesJson.data;
	}
	catch (error) {
		console.error("Unexpected error during show loading", error)
	}

	socket = new WebSocket('ws://localhost:8080/api/v1/ws')

	socket.onopen = () => {
		console.log("WebSocket connected")
	};

socket.onmessage = async (event) => {
	const data = JSON.parse(event.data);

	if (data.type === "EpisodeTrackingUpdated") {
		// recria todos os episódios com novo objeto para disparar reatividade
		episodes = episodes.map(episode =>
			episode.ID === data.episodeID
				? { ...episode, Tracking: data.tracking }
				:  episode
		);
	}
}
})

function deleteShow() {
	toast('Are you sure you want to delete this show?', {
		action: {
			label: 'Delete',
			onClick: () => {

	fetch(`${PUBLIC_API_BASE_URL}/database/show`, {
		method: "DELETE",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({ id: show.ID })
	}).then(() => {
			goto('/');
			// alert("Show deleted.");
			toast.success("Show deleted.");
		});
			}
		},
		cancel: {
			label: 'Cancel',
		}
	})

}

async function deleteEpisode(event) {
	const episodeID = event.detail.id

	try {
		const res = await fetch(`${PUBLIC_API_BASE_URL}/database/show/episode/${episodeID}`, {
			method: 'DELETE',
			headers: {
				"Content-Type": "application/json",
			}
		});

		if(!res.ok) {
			console.error("Failed to delete downloaded episode")
		}
	} catch(error) {
			console.error("Error trying to delete downloaded episode")
	}
}

// async function changeTrackingStatus(event) {
// 	const { id, newStatus } = event.detail
//
// 	try {
// 		const res = await fetch(`${PUBLIC_API_BASE_URL}/database/show/episode/status`, {
// 			method: 'POST',
// 			headers: {
// 				"Content-Type": "application/json"
// 			},
// 			body: JSON.stringify({ episode_id: id, tracking: newStatus })
// 		});
//
// 		if(!res.ok) {
// 			console.error("Failed to update episode status", res.status)
// 		}
//
// 		const data = await res.json()
// 		console.log("New Episode status", data)
// 	} catch(error) {
// 		console.error("Error trying to change episode status", error)
// 	}
// }

function toggleSeason(seasonNumber) {
	if (collapsedSeasons.has(seasonNumber)) {
		collapsedSeasons.delete(seasonNumber);
	} else {
		collapsedSeasons.add(seasonNumber);
	}
	collapsedSeasons = new Set(collapsedSeasons); // Trigger reactivity
}

function bulkyEdit() {
	// Placeholder function - no functionality yet
	console.log("Bulky Edit clicked");
}
</script>

<svelte:head>
	<link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet" />
</svelte:head>

{#if show}
	<div class="show-page">
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
				<button class="action-btn bulky-edit-btn" on:click={bulkyEdit}>
					<i class="fas fa-edit"></i> Bulky Edit
				</button>
				<button class="action-btn delete-show-btn" on:click={deleteShow}>
					<i class="fas fa-trash"></i> Delete Show
				</button>
			</div>
		</div>

		<div class="show-overview">
			<div class="show-poster">
				{#if show.ImageMedium}
					<img src={show.ImageMedium} alt={`${show.Name} Poster`} />
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
			{#each seasons as season}
				<a href={"#season-" + season.Number} class="season-pill">
					Season {season.Number}
				</a>
			{/each}
		</div>

		<!-- Temporadas em ordem decrescente -->
		{#each [...seasons].reverse() as season}
			<div id={"season-" + season.Number} class="season-block">
				<div class="season-header">
					<h2>
						<button class="season-toggle-btn" on:click={() => toggleSeason(season.Number)}>
							<i class="fas {collapsedSeasons.has(season.Number) ? 'fa-chevron-right' : 'fa-chevron-down'}"></i>
							Season {season.Number}
						</button>
					</h2>
					<div class="season-actions">
						<button class="action-btn season-search-btn" data-season={season.Number}>
							<i class="fas fa-search"></i> Search Season
						</button>
					</div>
				</div>

				{#if !collapsedSeasons.has(season.Number)}
					<div class="episodes-container">
						{#each [...episodes.filter(e => e.Season === season.Number)].reverse() as episode}
							<EpisodeCard {episode} seasonNumber={season.Number} on:deleteEpisode={deleteEpisode} />
									<!-- on:statusChange={changeTrackingStatus} /> -->
						{/each}
					</div>
				{/if}
			</div>
		{/each}
	</div>
</div>
{/if}
