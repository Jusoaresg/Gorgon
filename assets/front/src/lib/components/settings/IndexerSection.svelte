<script>
import { createEventDispatcher } from "svelte";
import '$lib/styles/components/configurations/configurations.css';
import { fade } from "svelte/transition";

export let indexer = "prowlarr";

export let prowlarrHost = '';
export let prowlarrPort = '';
export let prowlarrApiKey = '';

const dispatch = createEventDispatcher()

async function handleConfigChange() {
	dispatch('change', {
		indexer: indexer,
		config: {
			host: prowlarrHost,
			port: prowlarrPort,
			key: prowlarrApiKey
		}
	});
}

async function testConnection() {
	// dummy function
}
</script>

<div id="section-prowlarr" class="config-section" in:fade={{ duration: 200 }}>
	<div class="section-header">
		<h2>
			<i class="fas fa-search"></i>
			Prowlarr Configuration
		</h2>
		<div class="section-actions">
			<button class="action-btn test-btn" on:click={testConnection}>
				<i class="fas fa-plug"></i> Test Connection
			</button>
		</div>
	</div>

	<div class="settings-grid">
		<div class="setting-group">
			<label for="prowlarrHost">
				<i class="fas fa-server"></i>
				Prowlarr Host
			</label>
			<input 
				type="url" 
				id="prowlarrHost" 
				bind:value={prowlarrHost} 
				on:input={handleConfigChange}
				placeholder="http://localhost"
			>
		</div>

		<div class="setting-group">
			<label for="prowlarrPort">
				<i class="fas fa-ethernet"></i>
				Prowlarr Port
			</label>
			<input 
				type="text" 
				id="prowlarrPort" 
				bind:value={prowlarrPort} 
				on:input={handleConfigChange}
				placeholder="9696"
			>
		</div>

		<div class="setting-group full-width">
			<label for="prowlarrApiKey">
				<i class="fas fa-key"></i>
				Prowlarr API Key
			</label>
			<input 
				type="password" 
				id="prowlarrApiKey" 
				bind:value={prowlarrApiKey} 
				on:input={handleConfigChange}
				placeholder="Enter your Prowlarr API key"
			>
		</div>
	</div>
</div>
