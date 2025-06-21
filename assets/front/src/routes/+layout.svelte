<script>
  import Nav from './Nav.svelte';
  import { onMount } from 'svelte';
  import { Toaster } from 'svelte-sonner';
  // Styles
  import '$lib/styles/style.css';

  let showScrollButton = false;
  let mainContent;

  onMount(() => {
    const handleScroll = () => {
      if (mainContent) {
        showScrollButton = mainContent.scrollTop > 300;
      }
    };

    if (mainContent) {
      mainContent.addEventListener('scroll', handleScroll);
      
      return () => {
        mainContent.removeEventListener('scroll', handleScroll);
      };
    }
  });

  const scrollToTop = () => {
    if (mainContent) {
      mainContent.scrollTo({
        top: 0,
        behavior: 'smooth'
      });
    }
  };
</script>

<div class="app-container">
  <Nav />
  <main class="main-content" bind:this={mainContent}>
    <slot />
    <Toaster
      theme="dark"
      position="top-center"
      richColors={true}
      closeButton
      expand
    />
    
    <!-- Scroll to top button -->
    <button 
      class="scroll-to-top" 
      class:visible={showScrollButton}
      on:click={scrollToTop}
      aria-label="Scroll to top"
    >
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M18 15l-6-6-6 6"/>
      </svg>
    </button>
  </main>
</div>

<style>
  .app-container {
    display: flex;
    height: 100vh;
  }
  
  .main-content {
    flex: 1;
    overflow-y: auto;
    position: relative;
  }
  
  .scroll-to-top {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    width: 3rem;
    height: 3rem;
    background: var(--highlight);
    color: white;
    border: 1px solid var(--border-color);
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    transition: all 0.3s ease;
    opacity: 0;
    visibility: hidden;
    transform: translateY(10px);
    z-index: 1000;
  }
  
  .scroll-to-top:hover {
    background: var(--button-bg);
    border-color: var(--button-bg);
    color: var(--bg-color);
    transform: translateY(-2px) scale(1.05);
  }
  
  .scroll-to-top.visible {
    opacity: 1;
    visibility: visible;
    transform: translateY(0);
  }
  
  .scroll-to-top:active {
    transform: translateY(0) scale(0.95);
  }
  
  .scroll-to-top:focus {
    outline: 2px solid var(--highlight);
    outline-offset: 2px;
  }
</style>
