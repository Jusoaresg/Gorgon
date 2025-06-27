<script>
import { createEventDispatcher } from 'svelte';
import { fade, fly, scale } from 'svelte/transition';

const dispatch = createEventDispatcher();

// Props
export let preferredKeywords = [];
export let blacklistedKeywords = [];
export let requiredKeywords = [];

$: currentKeywords = {
	preferredKeywords,
	blacklistedKeywords,
	requiredKeywords
};

// Internal state
let keywordInputs = {
	preferredKeywords: '',
	blacklistedKeywords: '',
	requiredKeywords: ''
};

let inputFocused = {
	preferredKeywords: false,
	blacklistedKeywords: false,
	requiredKeywords: false
};

let suggestions = {
	preferredKeywords: ['1080p', '720p', 'BluRay', 'WEB-DL', 'HEVC', 'x264', 'x265'],
	blacklistedKeywords: ['CAM', 'TS', 'TC', 'HDCAM', 'HDTS', 'DVDSCR', 'WORKPRINT'],
	requiredKeywords: ['English', 'Subtitles', 'Multi-Subs']
};

let showSuggestions = {
	preferredKeywords: false,
	blacklistedKeywords: false,
	requiredKeywords: false
};

const keywordTypes = [
	{
		id: 'preferredKeywords',
		title: 'Preferred Keywords',
		icon: 'fas fa-star',
		description: 'Keywords that increase priority when searching for torrents',
		color: 'preferred',
		placeholder: 'e.g., 1080p, BluRay, HEVC...'
	},
	{
		id: 'requiredKeywords',
		title: 'Required Keywords',
		icon: 'fas fa-exclamation-circle',
		description: 'Keywords that must be present in torrent names',
		color: 'required',
		placeholder: 'e.g., English, Complete...'
	},
	{
		id: 'blacklistedKeywords',
		title: 'Blacklisted Keywords',
		icon: 'fas fa-ban',
		description: 'Keywords that will exclude torrents from search results',
		color: 'blacklisted',
		placeholder: 'e.g., CAM, TS, HDCAM...'
	}
];

function getCurrentKeywords(type) {
	switch(type) {
		case 'preferredKeywords': return preferredKeywords;
		case 'blacklistedKeywords': return blacklistedKeywords;
		case 'requiredKeywords': return requiredKeywords;
		default: return [];
	}
}

function addKeyword(type) {
	const value = keywordInputs[type].trim();
	const currentKeywords = getCurrentKeywords(type);
	
	if (value && !currentKeywords.includes(value)) {
		const newKeywords = [...currentKeywords, value];
		updateKeywords(type, newKeywords);
		keywordInputs[type] = '';
		showSuggestions[type] = false;
		return true;
	}
	return false;
}

function addMultipleKeywords(type, keywords) {
	const keywordsArray = keywords.split(',').map(k => k.trim()).filter(k => k);
	let currentKeywords = getCurrentKeywords(type);
	let addedCount = 0;
	
	keywordsArray.forEach(keyword => {
		if (keyword && !currentKeywords.includes(keyword)) {
			currentKeywords = [...currentKeywords, keyword]
			addedCount++;
		}
	});
	
	if (addedCount > 0) {
		updateKeywords(type, currentKeywords);
		keywordInputs[type] = '';
		showSuggestions[type] = false;
		return addedCount;
	}
	return 0;
}

function removeKeyword(type, keyword) {
	const currentKeywords = getCurrentKeywords(type);
	const newKeywords = currentKeywords.filter(k => k !== keyword);
	updateKeywords(type, newKeywords);
}

function addSuggestion(type, suggestion) {
	const currentKeywords = getCurrentKeywords(type);
	if (!currentKeywords.includes(suggestion)) {
		const newKeywords = [...currentKeywords, suggestion];
		updateKeywords(type, newKeywords);
	}
}

function updateKeywords(type, newKeywords) {
	switch(type) {
		case 'preferredKeywords':
			preferredKeywords = [...newKeywords];
			break;
		case 'blacklistedKeywords':
			blacklistedKeywords = [...newKeywords];
			break;
		case 'requiredKeywords':
			requiredKeywords = [...newKeywords];
			break;
	}
	
	// Dispatch change event
	dispatch('change', {
		type,
		keywords: newKeywords,
		all: {
			preferredKeywords,
			blacklistedKeywords,
			requiredKeywords
		}
	});
}

function handleKeywordInput(event, type) {
	const value = event.target.value;
	keywordInputs[type] = value;
	
	// Show suggestions when typing
	if (value.length > 0) {
		showSuggestions[type] = true;
	} else {
		showSuggestions[type] = false;
	}
}

function handleKeywordKeydown(event, type) {
	if (event.key === 'Enter') {
		event.preventDefault();
		const input = keywordInputs[type];
		
		// Check if input contains commas (multiple keywords)
		if (input.includes(',')) {
			const addedCount = addMultipleKeywords(type, input);
			if (addedCount > 0) {
				console.log(`Added ${addedCount} keywords`);
			}
		} else {
			addKeyword(type);
		}
	} else if (event.key === 'Escape') {
		showSuggestions[type] = false;
		event.target.blur();
	}
}

function handleInputFocus(type) {
	inputFocused[type] = true;
	if (keywordInputs[type].length > 0) {
		showSuggestions[type] = true;
	}
}

function handleInputBlur(type) {
	inputFocused[type] = false;
	// Delay hiding suggestions to allow clicking on them
	setTimeout(() => {
		if (!inputFocused[type]) {
			showSuggestions[type] = false;
		}
	}, 200);
}

function clearAllKeywords(type) {
	const typeName = type.replace('Keywords', '').toLowerCase();
	if (confirm(`Are you sure you want to remove all ${typeName} keywords?`)) {
		updateKeywords(type, []);
	}
}

function getFilteredSuggestions(type) {
	const input = keywordInputs[type].toLowerCase();
	const currentKeywords = getCurrentKeywords(type);
	return suggestions[type].filter(suggestion => 
		suggestion.toLowerCase().includes(input) && 
		!currentKeywords.includes(suggestion)
	);
}
</script>

<svelte:head>
	<link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet" />
</svelte:head>

<div class="keywords-manager" in:fade={{ duration: 200 }}>
	<div class="section-header">
		<h2>
			<i class="fas fa-tags"></i>
			Keywords Management
		</h2>
		<div class="section-actions">
			<button class="action-btn info-btn" title="Keywords help improve torrent search accuracy">
				<i class="fas fa-info-circle"></i> Help
			</button>
		</div>
	</div>
	
	<div class="keywords-overview">
		<div class="keywords-stats">
			<div class="stat-card preferred">
				<div class="stat-number">{preferredKeywords.length}</div>
				<div class="stat-label">Preferred</div>
			</div>
			<div class="stat-card required">
				<div class="stat-number">{requiredKeywords.length}</div>
				<div class="stat-label">Required</div>
			</div>
			<div class="stat-card blacklisted">
				<div class="stat-number">{blacklistedKeywords.length}</div>
				<div class="stat-label">Blacklisted</div>
			</div>
		</div>
	</div>

	{#each keywordTypes as keywordType}
		<div class="keyword-section enhanced" transition:fly={{ y: 20, duration: 300 }}>
			<div class="keyword-header">
				<div class="keyword-title">
					<h3>
						<i class={keywordType.icon}></i>
						{keywordType.title}
					</h3>
					<p class="keyword-description">{keywordType.description}</p>
				</div>
				<div class="keyword-actions">
					{#if getCurrentKeywords(keywordType.id).length > 0}
						<button 
							aria-label="clear-all"
							class="clear-all-btn" 
							on:click={() => clearAllKeywords(keywordType.id)}
							title="Clear all keywords"
						>
							<i class="fas fa-trash"></i>
						</button>
					{/if}
				</div>
			</div>
			
			<div class="keyword-input-container">
				<div class="keyword-input-wrapper" class:focused={inputFocused[keywordType.id]}>
					<input 
						type="text" 
						bind:value={keywordInputs[keywordType.id]}
						on:input={(e) => handleKeywordInput(e, keywordType.id)}
						on:keydown={(e) => handleKeywordKeydown(e, keywordType.id)}
						on:focus={() => handleInputFocus(keywordType.id)}
						on:blur={() => handleInputBlur(keywordType.id)}
						placeholder={keywordType.placeholder}
						autocomplete="off"
					>
					<button 
						aria-label="add-keyword"
						class="add-keyword-btn {keywordType.color}" 
						on:click={() => {
							if (keywordInputs[keywordType.id].includes(',')) {
								addMultipleKeywords(keywordType.id, keywordInputs[keywordType.id]);
							} else {
								addKeyword(keywordType.id);
							}
						}}
						disabled={!keywordInputs[keywordType.id].trim()}
					>
						<i class="fas fa-plus"></i>
					</button>
				</div>
				
				<div class="input-help">
					<span class="help-text">
						<i class="fas fa-lightbulb"></i>
						Tip: Use commas to add multiple keywords at once
					</span>
				</div>

				<!-- Suggestions dropdown -->
				{#if showSuggestions[keywordType.id] && getFilteredSuggestions(keywordType.id).length > 0}
					<div class="suggestions-dropdown" transition:scale={{ duration: 200 }}>
						<div class="suggestions-header">Suggestions</div>
						{#each getFilteredSuggestions(keywordType.id) as suggestion}
							<button 
								class="suggestion-item"
								on:click={() => addSuggestion(keywordType.id, suggestion)}
							>
								<i class="fas fa-plus-circle"></i>
								{suggestion}
							</button>
						{/each}
					</div>
				{/if}
			</div>
			
			<div class="keywords-list">
				{#if currentKeywords[keywordType.id].length === 0}
					<div class="empty-state">
						<i class={keywordType.icon}></i>
						<p>No {keywordType.title.toLowerCase()} added yet</p>
						<p class="empty-subtitle">Add keywords to improve your search results</p>
					</div>
				{:else}
					<!-- {#each getCurrentKeywords(keywordType.id) as keyword, index (keyword)} -->
					{#each currentKeywords[keywordType.id] as keyword, index (keyword)}
						<div 
							class="keyword-tag {keywordType.color}" 
							transition:scale={{ duration: 200, delay: index * 50 }}
						>
							<span class="keyword-text">{keyword}</span>
							<button 
								aria-label="remove-keyword"
								class="remove-keyword-btn"
								on:click={() => removeKeyword(keywordType.id, keyword)}
								title="Remove keyword"
							>
								<i class="fas fa-times"></i>
							</button>
						</div>
					{/each}
				{/if}
			</div>
		</div>
	{/each}
</div>

<style>
:root {
	--preferred-color: #00ff7f;
	--required-color: #7d00a3;
	--blacklisted-color: #ff6b6b;
	--focus-color: #4a9eff;
}

.keywords-manager {
	background-color: var(--nav-bg);
	border-radius: 12px;
	padding: 2rem;
	box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
	color: var(--text-color);
}

.section-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 2rem;
}

.section-header h2 {
	margin: 0;
	color: var(--text-color);
	font-size: 1.5rem;
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.section-header h2 i {
	color: var(--highlight);
}

.section-actions {
	display: flex;
	gap: 0.5rem;
}

.action-btn {
	background-color: var(--card-bg);
	color: var(--text-color);
	border: none;
	border-radius: 8px;
	padding: 0.7rem 1.2rem;
	font-size: 0.95rem;
	display: flex;
	align-items: center;
	gap: 0.5rem;
	cursor: pointer;
	transition: all 0.3s ease;
	font-weight: 500;
}

.action-btn:hover {
	background-color: var(--highlight);
	color: white;
	transform: translateY(-2px);
}

.info-btn:hover {
	background-color: #3498db !important;
}

/* Enhanced Keywords Section */
.keywords-overview {
	margin-bottom: 2rem;
	padding: 1.5rem;
	background: linear-gradient(135deg, var(--card-bg) 0%, rgba(125, 0, 163, 0.1) 100%);
	border-radius: 12px;
	border: 1px solid var(--border-color);
}

.keywords-stats {
	display: flex;
	gap: 1rem;
	justify-content: center;
	flex-wrap: wrap;
}

.stat-card {
	background-color: var(--nav-bg);
	padding: 1rem 1.5rem;
	border-radius: 8px;
	text-align: center;
	min-width: 100px;
	border: 2px solid transparent;
	transition: all 0.3s ease;
}

.stat-card:hover {
	transform: translateY(-2px);
}

.stat-card.preferred {
	border-color: var(--preferred-color);
}

.stat-card.required {
	border-color: var(--required-color);
}

.stat-card.blacklisted {
	border-color: var(--blacklisted-color);
}

.stat-number {
	font-size: 1.5rem;
	font-weight: bold;
	margin-bottom: 0.25rem;
}

.stat-card.preferred .stat-number {
	color: var(--preferred-color);
}

.stat-card.required .stat-number {
	color: var(--required-color);
}

.stat-card.blacklisted .stat-number {
	color: var(--blacklisted-color);
}

.stat-label {
	font-size: 0.9rem;
	color: var(--text-muted);
}

.keyword-section.enhanced {
	margin-bottom: 2.5rem;
	padding: 2rem;
	background-color: var(--card-bg);
	border-radius: 12px;
	border: 1px solid var(--border-color);
	position: relative;
	overflow: visible;
}

.keyword-header {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	margin-bottom: 1.5rem;
}

.keyword-title h3 {
	margin: 0 0 0.5rem 0;
	color: var(--text-color);
	display: flex;
	align-items: center;
	gap: 0.5rem;
	font-size: 1.2rem;
}

.keyword-title h3 i {
	color: var(--highlight);
}

.keyword-description {
	margin: 0;
	color: var(--text-muted);
	font-size: 0.9rem;
	line-height: 1.4;
}

.keyword-actions {
	display: flex;
	gap: 0.5rem;
}

.clear-all-btn {
	background-color: rgba(255, 107, 107, 0.2);
	color: #ff6b6b;
	border: 1px solid #ff6b6b;
	border-radius: 6px;
	padding: 0.5rem;
	cursor: pointer;
	transition: all 0.3s ease;
}

.clear-all-btn:hover {
	background-color: #ff6b6b;
	color: white;
}

.keyword-input-container {
	position: relative;
	margin-bottom: 1.5rem;
}

.keyword-input-wrapper {
	display: flex;
	gap: 0.5rem;
	background-color: var(--input-bg);
	border: 2px solid var(--border-color);
	border-radius: 8px;
	padding: 0.5rem;
	transition: all 0.3s ease;
}

.keyword-input-wrapper.focused {
	border-color: var(--focus-color);
	box-shadow: 0 0 0 3px rgba(74, 158, 255, 0.2);
}

.keyword-input-wrapper input {
	flex: 1;
	background: none;
	border: none;
	color: var(--text-color);
	font-size: 0.95rem;
	padding: 0.5rem;
	outline: none;
}

.keyword-input-wrapper input::placeholder {
	color: var(--text-muted);
}

.add-keyword-btn {
	background-color: var(--button-bg);
	color: var(--bg-color);
	border: none;
	border-radius: 6px;
	padding: 0.5rem 1rem;
	cursor: pointer;
	font-size: 0.9rem;
	font-weight: 600;
	transition: all 0.3s ease;
	display: flex;
	align-items: center;
	gap: 0.5rem;
	min-width: 80px;
	justify-content: center;
}

.add-keyword-btn:hover:not(:disabled) {
	background-color: var(--button-hover);
	transform: scale(1.05);
}

.add-keyword-btn:disabled {
	opacity: 0.5;
	cursor: not-allowed;
	transform: none;
}

.add-keyword-btn.preferred {
	background-color: var(--preferred-color);
}

.add-keyword-btn.required {
	background-color: var(--required-color);
}

.add-keyword-btn.blacklisted {
	background-color: var(--blacklisted-color);
}

.input-help {
	margin-top: 0.5rem;
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.help-text {
	font-size: 0.8rem;
	color: var(--text-muted);
	display: flex;
	align-items: center;
	gap: 0.25rem;
}

.help-text i {
	color: #ffd700;
}

.suggestions-dropdown {
	position: absolute;
	top: 100%;
	left: 0;
	right: 0;
	background-color: var(--nav-bg);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	z-index: 1000;
	max-height: 200px;
	overflow-y: auto;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	margin-top: 0.5rem;
}

.suggestions-header {
	padding: 0.75rem 1rem;
	background-color: var(--card-bg);
	color: var(--text-muted);
	font-size: 0.8rem;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.5px;
	border-bottom: 1px solid var(--border-color);
}

.suggestion-item {
	width: 100%;
	padding: 0.75rem 1rem;
	background: none;
	border: none;
	color: var(--text-color);
	text-align: left;
	cursor: pointer;
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.suggestion-item:hover {
	background-color: var(--hover-bg);
}

.suggestion-item i {
	color: var(--highlight);
	font-size: 0.8rem;
}

.keywords-list {
	display: flex;
	flex-wrap: wrap;
	gap: 0.75rem;
	min-height: 60px;
	align-items: flex-start;
	align-content: flex-start;
}

.empty-state {
	width: 100%;
	text-align: center;
	padding: 2rem;
	color: var(--text-muted);
	background-color: var(--input-bg);
	border-radius: 8px;
	border: 2px dashed var(--border-color);
}

.empty-state i {
	font-size: 2rem;
	margin-bottom: 0.5rem;
	opacity: 0.5;
}

.empty-state p {
	margin: 0.5rem 0;
}

.empty-subtitle {
	font-size: 0.8rem;
	opacity: 0.7;
}

.keyword-tag {
	background-color: var(--nav-bg);
	border: 1px solid var(--border-color);
	border-radius: 20px;
	padding: 0.5rem 1rem;
	display: flex;
	align-items: center;
	gap: 0.5rem;
	font-size: 0.9rem;
	transition: all 0.3s ease;
	position: relative;
	overflow: hidden;
}

.keyword-tag::before {
	content: '';
	position: absolute;
	left: 0;
	top: 0;
	bottom: 0;
	width: 3px;
	background-color: var(--highlight);
}

.keyword-tag.preferred::before {
	background-color: var(--preferred-color);
}

.keyword-tag.required::before {
	background-color: var(--required-color);
}

.keyword-tag.blacklisted::before {
	background-color: var(--blacklisted-color);
}

.keyword-tag:hover {
	transform: translateY(-2px);
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.keyword-text {
	color: var(--text-color);
	font-weight: 500;
}

.remove-keyword-btn {
	background: none;
	border: none;
	color: var(--text-muted);
	cursor: pointer;
	padding: 0.2rem;
	border-radius: 50%;
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: all 0.2s ease;
	font-size: 0.7rem;
}

.remove-keyword-btn:hover {
	background-color: #ff6b6b;
	color: white;
	transform: scale(1.1);
}

.quick-actions {
	margin-top: 2rem;
	padding: 1.5rem;
	background: linear-gradient(135deg, var(--card-bg) 0%, rgba(0, 255, 127, 0.05) 100%);
	border-radius: 12px;
	border: 1px solid var(--border-color);
}

.quick-actions-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
	gap: 1rem;
}

.quick-action-btn {
	background-color: var(--nav-bg);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	padding: 1rem;
	color: var(--text-color);
	cursor: pointer;
	transition: all 0.3s ease;
	display: flex;
	align-items: center;
	gap: 0.75rem;
	font-size: 0.9rem;
	text-align: left;
}

.quick-action-btn:hover {
	background-color: var(--highlight);
	color: white;
	transform: translateY(-2px);
	box-shadow: 0 4px 12px rgba(125, 0, 163, 0.3);
}

.quick-action-btn i {
	font-size: 1.1rem;
	opacity: 0.8;
}

/* Responsive Design */
@media (max-width: 768px) {
	.keywords-manager {
		padding: 1rem;
	}
	
	.keyword-section.enhanced {
		padding: 1.5rem;
	}
	
	.keyword-header {
		flex-direction: column;
		gap: 1rem;
		align-items: flex-start;
	}
	
	.keywords-stats {
		justify-content: center;
	}
	
	.stat-card {
		min-width: 80px;
		padding: 0.75rem 1rem;
	}
	
	.quick-actions-grid {
		grid-template-columns: 1fr;
	}
	
	.keyword-input-wrapper {
		flex-direction: column;
		gap: 0.75rem;
	}
	
	.add-keyword-btn {
		align-self: stretch;
	}
}

@media (max-width: 480px) {
	.keywords-stats {
		flex-direction: column;
		align-items: center;
	}
	
	.stat-card {
		width: 100%;
		max-width: 150px;
	}
}

/* Scrollbar Styling */
.suggestions-dropdown::-webkit-scrollbar {
	width: 6px;
}

.suggestions-dropdown::-webkit-scrollbar-track {
	background: var(--card-bg);
}

.suggestions-dropdown::-webkit-scrollbar-thumb {
	background: var(--border-color);
	border-radius: 3px;
}

.suggestions-dropdown::-webkit-scrollbar-thumb:hover {
	background: var(--text-muted);
}

/* Focus states for accessibility */
.action-btn:focus,
.add-keyword-btn:focus,
.remove-keyword-btn:focus,
.suggestion-item:focus,
.quick-action-btn:focus,
.clear-all-btn:focus {
	outline: 2px solid var(--focus-color);
	outline-offset: 2px;
}

/* Animation keyframes */
@keyframes pulse {
	0%, 100% {
		transform: scale(1);
	}
	50% {
		transform: scale(1.05);
	}
}

.keyword-tag:hover {
	animation: pulse 0.3s ease-in-out;
}
</style>
