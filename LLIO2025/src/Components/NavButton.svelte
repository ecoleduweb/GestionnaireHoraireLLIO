<script lang="ts">
  import { goto } from "$app/navigation";
  import { UserRole } from "$lib/types/enums";
  import { page } from '$app/state';

  type Props = {
      currentUserRole: UserRole;
    };
    
  const path = $derived(page.url.pathname);
  const { currentUserRole } : Props = $props();

  const isCalendarActive = $derived(path === '/calendar');
  const isProjectsActive = $derived(path === '/projects');
  const isAdminActive = $derived(path === '/administrator');

  const bouttonCalendarClass =  $derived(isCalendarActive ? `btnActiveBaseClass` : `btnNotActiveBaseClass`)
  const bouttonProjectsClass =  $derived(isProjectsActive ? `btnActiveBaseClass` : `btnNotActiveBaseClass`);
  const bouttonAdminClass =  $derived(isAdminActive ? `btnActiveBaseClass` : `btnNotActiveBaseClass`)

</script>

<div class="inline-flex rounded-md shadow-xs" role="group">
  <button
      onclick={() => goto('./calendar')}
      type="button"
      class={`btnBaseClass  ${bouttonCalendarClass}`}
  >
      Calendrier
  </button>
  <button 
      onclick={() => goto('./projects')}
      type="button" 
      class={`btnBaseClass ${bouttonProjectsClass}`}
  >
      Projets
  </button>
  {#if currentUserRole === UserRole.Admin}
      <button 
          onclick={() => goto('./administrator')}
          type="button" 
          class={`btnBaseClass  ${bouttonAdminClass}`}
      >
          Admin
      </button>
  {/if}
</div>

<style>
@reference "tailwindcss";
.btnBaseClass
{
    @apply px-4 py-2 text-sm transition-colors font-semibold ;
}
.btnActiveBaseClass{
    @apply bg-[#014446] text-white ;
}
.btnNotActiveBaseClass{
    @apply bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer ;
}
:global(button:first-child){
    @apply rounded-l-lg ;
}
:global(button:last-child){
    @apply rounded-r-lg ;
}
</style>

