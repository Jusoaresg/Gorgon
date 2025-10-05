<script>

import IndexerSection from '$lib/components/settings/IndexerSection.svelte';
import KeywordSection from '$lib/components/settings/KeywordSection.svelte';
import TorrentSection from '$lib/components/settings/TorrentSection.svelte';
import FolderSection from  '$lib/components/settings/FolderSection.svelte';
import { onMount } from 'svelte';
import { fade, fly } from 'svelte/transition';
import { toast } from 'svelte-sonner';

import '$lib/styles/components/configurations/configurations.css';

// Configuration state with empty initial values
let config = {

	// Prowlarr settings
	prowlarrApiKey: "",
	prowlarrHost: "",
	prowlarrPort: "",

	// qBittorrent settings
	qBittorrentHost: "",
	qBittorrentPort: "",
	qBittorrentUsername: "",
	qBittorrentPassword: "",
	qBittorrentDownloadFolder: "",

	// Folder settings
	defaultShowInfoFolder: "",
	showsFolder: "",

	// Keyword Settings
	preferredKeywords: [],
	blacklistedKeywords: [],
	requiredKeywords: []
};

let activeSection = 'prowlarr';
let hasChanges = false;
let saveStatus = '';

// Sections configuration
const sections = [
	{ id: 'prowlarr', title: 'Prowlarr', icon: '🔍' },
	{ id: 'qbittorrent', title: 'qBittorrent', icon: '⬇️' },
	{ id: 'folders', title: 'Folders', icon: '📁' },
	{ id: 'keywords', title: 'Keywords', icon: '🏷️' },
	{ id: 'about', title: 'About', icon: 'ℹ️' }
];

async function updateKeywords(event) {
	const all = event.detail.all;

	config = {
		...config,
		preferredKeywords: [...all.preferredKeywords],
		blacklistedKeywords: [...all.blacklistedKeywords],
		requiredKeywords: [...all.requiredKeywords]
	};

	hasChanges = true;
	saveStatus = '';
}

async function updateTorrentClient(event) {
	const all = event.detail;

	config = {
		...config,
		qBittorrentHost: all.config.host,
		qBittorrentPort: all.config.port,
		qBittorrentUsername: all.config.username,
		qBittorrentPassword: all.config.password,
		qBittorrentDownloadFolder: all.config.downloadFolder
	};

	hasChanges = true;
	saveStatus = '';
}

async function updateIndexer(event) {
	const all = event.detail;

	config = {
		...config,
		prowlarrHost: all.config.host,
		prowlarrPort: all.config.port,
		prowlarrApiKey: all.config.key,
	};

	hasChanges = true;
	saveStatus = '';
}

async function updateFolders(event) {
	const all = event.detail;

	config = {
		...config,
		defaultShowInfoFolder: all.showInfoFolder,
		showsFolder: all.showsFolder
	};

	hasChanges = true;
	saveStatus = '';
}

function switchSection(sectionId) {
	activeSection = sectionId;
}

async function saveConfig() {
	try {
		saveStatus = 'saving';

		const response = await fetch('/api/v1/app/config', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(config)
		});

		hasChanges = false;
		saveStatus = 'saved';
		setTimeout(() => saveStatus = '', 3000);
		toast.success("Configuration saved successfully!");
	} catch (error) {
		console.error('Error saving config:', error);
		saveStatus = 'error';
		setTimeout(() => saveStatus = '', 3000);
		toast.error("Failed to save configuration. Please try again.");
	}
}

function resetConfig() {
	if (confirm('Are you sure you want to reset all settings to default values?')) {
		// Reset to empty values
		config = {
			prowlarrApiKey: "",
			prowlarrHost: "",
			prowlarrPort: "",
			qBittorrentHost: "",
			qBittorrentPort: "",
			qBittorrentUsername: "",
			qBittorrentPassword: "",
			qBittorrentDownloadFolder: "",
			defaultShowInfoFolder: "",
			showsFolder: "",
			preferredKeywords: [],
			blacklistedKeywords: [],
			requiredKeywords: []
		};
		hasChanges = true;
		toast.info("Configuration has been reset to defaults.");
	}
}



onMount(async () => {
	try {
		//Load current configuration
		const res = await fetch("/api/v1/app/config");
		const json = await res.json();
		if (json.data) {
			config = { 
				...config, 
				...json.data,
				// Ensure arrays exist
				preferredKeywords: json.data.preferredKeywords || [],
				blacklistedKeywords: json.data.blacklistedKeywords || [],
				requiredKeywords: json.data.requiredKeywords || []
			};
		}
	} catch (error) {
		console.error('Error loading config:', error);
		toast.error("Failed to load configuration.");
	}
});
</script>

<svelte:head>
	<title>Configuration - My Shows</title>
</svelte:head>

<div class="config-page" in:fade={{ duration: 250, delay: 50 }}>
	<div class="config-container">
		<!-- Header -->
		<div class="config-header">
			<h1>Configuration</h1>
			<div class="config-meta">
				<span class="config-status">Settings Management</span>
				<span class="config-version">Version 0.1.0</span>
			</div>
			<div class="config-actions">
				{#if hasChanges}
					<button class="action-btn save-btn" on:click={saveConfig} disabled={saveStatus === 'saving'}>
						<i class="fas fa-save"></i>
						{#if saveStatus === 'saving'}
							Saving...
						{:else if saveStatus === 'saved'}
							✓ Saved
						{:else if saveStatus === 'error'}
							Error
						{:else}
							Save Changes
						{/if}
					</button>
				{/if}
				<button class="action-btn reset-btn" on:click={resetConfig}>
					<i class="fas fa-undo"></i> Reset to Default
				</button>
			</div>
		</div>

		<!-- Sections Navigation -->
		<div class="sections-navigation">
			{#each sections as section}
				<a href={"#section-" + section.id} class="section-pill" class:active={activeSection === section.id} on:click={() => switchSection(section.id)}>
					<span class="section-icon">{section.icon}</span>
					{section.title}
				</a>
			{/each}
		</div>

		<!-- Content Sections -->
		<div class="config-content">
			<!-- Prowlarr Settings -->
			{#if activeSection === 'prowlarr'}
				<IndexerSection
					prowlarrHost={config.prowlarrHost}
					prowlarrPort={config.prowlarrPort}
					prowlarrApiKey={config.prowlarrApiKey}
					on:change={updateIndexer}
				/>
			{/if}

			<!-- qBittorrent Settings -->
			{#if activeSection === 'qbittorrent'}
				<TorrentSection
					qBittorrentHost={config.qBittorrentHost}
					qBittorrentPort={config.qBittorrentPort}
					qBittorrentUsername={config.qBittorrentUsername}
					qBittorrentPassword={config.qBittorrentPassword}
					qBittorrentDownloadFolder={config.qBittorrentDownloadFolder}
					on:change={updateTorrentClient}
				/>
			{/if}

			<!-- Folders Settings -->
			{#if activeSection === 'folders'}
				<FolderSection
					defaultShowInfoFolder={config.defaultShowInfoFolder}
					showsFolder={config.showsFolder}
					on:change={updateFolders}
				/>
			{/if}

			<!-- Keywords Settings -->
			{#if activeSection === 'keywords'}
				<KeywordSection
					preferredKeywords={config.preferredKeywords}
					requiredKeywords={config.requiredKeywords}
					blacklistedKeywords={config.blacklistedKeywords}
					on:change={updateKeywords}
				/>
			{/if}

			<!-- About Section -->
			{#if activeSection === 'about'}
				<div id="section-about" class="config-section" in:fade={{ duration: 200 }}>
					<div class="section-header">
						<h2>
							<i class="fas fa-info-circle"></i>
							About Gorgon
						</h2>
					</div>

					<div class="about-content">
						<div class="app-info">
							<div class="app-logo">
								<i class="fas fa-dragon"></i>
							</div>
							<h3>Gorgon</h3>
							<p>Version 0.1.0</p>
							<p>A comprehensive media management system for organizing and downloading your favorite shows.</p>

							<div class="system-info">
								<h4>System Information</h4>
								<div class="info-grid">
									<div class="info-item">
										<span class="info-label">Platform:</span>
										<span class="info-value">Web Application</span>
									</div>
									<div class="info-item">
										<span class="info-label">Framework:</span>
										<span class="info-value">Svelte</span>
									</div>
									<div class="info-item">
										<span class="info-label">API Version:</span>
										<span class="info-value">v1.0</span>
									</div>
									<div class="info-item">
										<span class="info-label">Build Date:</span>
										<span class="info-value">2025-06-23</span>
									</div>
								</div>
							</div>

							<div class="links">
								<a href="https://github.com/Jusoaresg/Gorgon" target="_blank" class="link-btn">
									<i class="fas fa-book"></i> Documentation
								</a>
								<a href="https://github.com/Jusoaresg/Gorgon" target="_blank" class="link-btn">
									<i class="fab fa-github"></i> GitHub
								</a>
								<a href="https://github.com/Jusoaresg/Gorgon" target="_blank" class="link-btn">
									<i class="fas fa-question-circle"></i> Support
								</a>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
