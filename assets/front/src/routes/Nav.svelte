<script>
  import { page } from '$app/stores';
  
  // Navigation items with icons
  const navItems = [
    { href: '/', label: 'Shows', icon: '📺' },
    { href: '/add-show', label: 'Add Show', icon: '➕' },
    { href: '/config', label: 'Config', icon: '⚙️' }
  ];
  
  // Check if current route is active
  function isActive(href, currentPath) {
    return currentPath === href || (href !== '/' && currentPath.startsWith(href));
  }
</script>

<nav class="sidebar">
  <div class="sidebar-header">
    <h1>Gorgon</h1>
    <div class="brand-accent"></div>
  </div>
  
  <div class="sidebar-controls">
    {#each navItems as item}
      <a 
        href={item.href} 
        class="nav-link"
        class:active={isActive(item.href, $page.url.pathname)}
      >
        <span class="nav-icon">{item.icon}</span>
        <span class="nav-label">{item.label}</span>
      </a>
    {/each}
  </div>
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
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 12px;
    position: relative;
    font-size: 14px;
    border: 1px solid transparent;
  }
  
  .nav-link:hover {
    background-color: var(--hover-bg);
    border-color: var(--border-color);
    transform: translateX(2px);
  }
  
  .nav-link.active {
    background-color: var(--card-bg);
    border-color: var(--highlight);
    color: var(--text-color);
    box-shadow: 0 2px 8px rgba(125, 0, 163, 0.2);
  }
  
  .nav-link.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 20px;
    background-color: var(--highlight);
    border-radius: 0 2px 2px 0;
  }
  
  .nav-icon {
    font-size: 16px;
    width: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .nav-label {
    flex: 1;
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
