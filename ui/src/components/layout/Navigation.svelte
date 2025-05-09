<script lang="ts">
  export let mobile = false;
  
  // SVG path data for icons
  const icons = {
    dashboard: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />',
    topics: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />',
    learners: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />',
    chat: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15h.01M12 8h.01M12 17h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />',
    settings: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />'
  };

  // Get the current hash route
  const currentHash = window.location.hash.slice(1) || 'dashboard';

  // Navigation items
  const navItems = [
    { href: '/dashboard', text: 'Dashboard', icon: icons.dashboard },
    { href: '/practice-topics', text: 'Practice Topics', icon: icons.topics },
    { href: '/learners', text: 'Learners', icon: icons.learners },
    { href: '/chat', text: 'Chat', icon: icons.chat },
    { href: '/settings', text: 'Settings', icon: icons.settings }
  ];

  import NavItem from './NavItem.svelte';
</script>

<nav class={mobile ? "mt-5 px-2 space-y-1" : "mt-5 flex-1 px-2 bg-white space-y-1"}>
  {#each navItems as item}
    <NavItem 
      href={item.href} 
      text={item.text} 
      icon={item.icon} 
      active={currentHash === item.href.substring(1) || (item.href === '/dashboard' && currentHash === '')} 
      mobile={mobile} 
    />
  {/each}
</nav> 