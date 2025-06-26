<script>
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { fade, fly, scale } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  
  // Navigation items with icons
  const navItems = [
    { href: '/', label: 'Shows', icon: '📺' },
    { href: '/add-show', label: 'Add Show', icon: '➕' },
    { href: '/settings', label: 'Settings', icon: '⚙️' }
  ];
  
  let mounted = false;
  
  // Check if current route is active
  function isActive(href, currentPath) {
    return currentPath === href || (href !== '/' && currentPath.startsWith(href));
  }
  
  onMount(() => {
    mounted = true;
  });
</script>

<nav class="sidebar">
  <div class="sidebar-header" 
       in:fly={{ y: -20, duration: 300, delay: 50, easing: quintOut }}>
    <h1 in:scale={{ duration: 400, delay: 200, easing: quintOut }}>Gorgon</h1>
    <div class="brand-accent" 
         in:scale={{ duration: 300, delay: 350, easing: quintOut }}></div>
  </div>
  
  <div class="sidebar-controls">
    {#each navItems as item, i}
      <a 
        href={item.href} 
        class="nav-link"
        class:active={isActive(item.href, $page.url.pathname)}
        in:fly={{ x: -30, duration: 300, delay: 400 + i * 100, easing: quintOut }}
      >
        <span class="nav-icon" 
              in:scale={{ duration: 200, delay: 500 + i * 100 }}>{item.icon}</span>
        <span class="nav-label" 
              in:fade={{ duration: 200, delay: 550 + i * 100 }}>{item.label}</span>
        
        <!-- Active indicator animation -->
        {#if mounted && isActive(item.href, $page.url.pathname)}
          <div class="active-indicator" 
               in:scale={{ duration: 200, easing: quintOut }}></div>
        {/if}
      </a>
    {/each}
  </div>
  
  <!-- Animated background pulse -->
  <div class="bg-pulse" in:fade={{ duration: 500, delay: 800 }}></div>
</nav>

<style>
  .sidebar {
    background-color: var(--nav-bg);
    border-right: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 220px;
    box-sizing: border-box;
    position: relative;
    overflow: hidden;
  }
  
  .sidebar::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--highlight), transparent);
    opacity: 0.5;
    animation: shimmer 3s ease-in-out infinite;
  }
  
  @keyframes shimmer {
    0%, 100% { opacity: 0.3; }
    50% { opacity: 0.7; }
  }
  
  .sidebar-header {
    padding: 25px 20px 20px;
    border-bottom: 1px solid var(--border-color);
    position: relative;
  }
  
  .sidebar-header h1 {
    font-size: 24px;
    font-weight: 700;
    color: var(--text-color);
    margin: 0;
    text-align: center;
    letter-spacing: -0.5px;
  }
  
  .brand-accent {
    width: 40px;
    height: 3px;
    background: linear-gradient(90deg, var(--highlight), var(--button-bg));
    margin: 8px auto 0;
    border-radius: 2px;
    animation: glow 2s ease-in-out infinite alternate;
  }
  
  @keyframes glow {
    from { box-shadow: 0 0 5px rgba(125, 0, 163, 0.3); }
    to { box-shadow: 0 0 15px rgba(125, 0, 163, 0.6); }
  }
  
  .sidebar-controls {
    display: flex;
    flex-direction: column;
    gap: 5px;
    flex: 1;
    padding: 20px 15px;
  }
  
  .nav-link {
    color: var(--text-color);
    text-decoration: none;
    font-weight: 500;
    padding: 12px 15px;
    border-radius: 8px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    align-items: center;
    gap: 12px;
    position: relative;
    font-size: 14px;
    border: 1px solid transparent;
    overflow: hidden;
  }
  
  .nav-link::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent);
    transition: left 0.5s ease;
  }
  
  .nav-link:hover::before {
    left: 100%;
  }
  
  .nav-link:hover {
    background-color: var(--hover-bg);
    border-color: var(--border-color);
    transform: translateX(4px) scale(1.02);
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  }
  
  .nav-link:hover .nav-icon {
    transform: scale(1.1) rotate(5deg);
  }
  
  .nav-link.active {
    background-color: var(--card-bg);
    border-color: var(--highlight);
    color: var(--text-color);
    box-shadow: 0 2px 8px rgba(125, 0, 163, 0.2);
    transform: translateX(2px);
  }
  
  .nav-link.active .nav-icon {
    animation: bounce 0.6s ease-out;
  }
  
  @keyframes bounce {
    0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
    40% { transform: translateY(-4px); }
    60% { transform: translateY(-2px); }
  }
  
  .active-indicator {
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 20px;
    background-color: var(--highlight);
    border-radius: 0 2px 2px 0;
    animation: slideIn 0.3s ease-out;
  }
  
  @keyframes slideIn {
    from { 
      height: 0;
      opacity: 0;
    }
    to { 
      height: 20px;
      opacity: 1;
    }
  }
  
  .nav-icon {
    font-size: 16px;
    width: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .nav-label {
    flex: 1;
    transition: opacity 0.2s ease;
  }
  
  .bg-pulse {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, var(--highlight), transparent, var(--highlight));
    opacity: 0.3;
    animation: pulse 2s ease-in-out infinite;
  }
  
  @keyframes pulse {
    0%, 100% { 
      transform: scaleX(0.8);
      opacity: 0.3;
    }
    50% { 
      transform: scaleX(1);
      opacity: 0.6;
    }
  }
  
  /* Responsive */
  @media (max-width: 768px) {
    .sidebar {
      width: 60px;
    }
    
    .sidebar-header h1,
    .nav-label,
    .brand-accent {
      display: none;
    }
    
    .nav-link {
      justify-content: center;
      padding: 12px 8px;
    }
    
    .nav-link:hover {
      transform: translateY(-2px) scale(1.05);
    }
    
    .sidebar-controls {
      padding: 20px 8px;
    }
    
    .sidebar-header {
      padding: 20px 8px 15px;
    }
  }
  
  /* Focus states for accessibility */
  .nav-link:focus {
    outline: 2px solid var(--highlight);
    outline-offset: 2px;
  }
</style>
