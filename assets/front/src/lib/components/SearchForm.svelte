<script>
    import { createEventDispatcher } from "svelte";

	export let placeholder = "PlaceHolder";
	export let hasButton = false;
	export let live = false;

	let query = '';
	const dispatch = createEventDispatcher();

    function handleSubmit(event) {
	    event.preventDefault();
	    dispatch('search', query)
    }

    function handleInput(event) {
	    if (live) {
		    dispatch('search', query)
	    }
    }
</script>

<div class="search-container">
	<form class="slide-it" on:submit|preventDefault={handleSubmit}>
		<div class="input-group">
			<input
				type="text"
				bind:value={query}
				placeholder={ placeholder }
				autocomplete="off"
				on:input={handleInput}
			/>
			{#if hasButton}
			<button type="submit">
				<span >Search</span>
				<!-- <span >Searching...</span> -->
			</button>
			{/if}
		</div>
	</form>
</div>
