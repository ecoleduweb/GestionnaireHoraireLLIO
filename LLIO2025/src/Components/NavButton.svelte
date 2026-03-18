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

  const baseClass = "px-4 py-2 text-sm transition-colors font-semibold"

  const bouttonCalendarClass =  $derived(isCalendarActive ? `${baseClass} bg-[#014446] text-white rounded-l-lg` : `${baseClass} bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer rounded-l-lg`)
  const bouttonProjectsClass =  $derived( currentUserRole === UserRole.Admin
                                            ? (
                                                isProjectsActive
                                                  ? `${baseClass} bg-[#014446] text-white`
                                                  : `${baseClass} bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer`
                                              )
                                            : (
                                                isProjectsActive
                                                  ? `${baseClass} bg-[#014446] text-white rounded-r-lg`
                                                  : `${baseClass} bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer rounded-r-lg`
                                              ));
  const bouttonAdminClass =  $derived(isAdminActive ? `${baseClass} bg-[#014446] text-white rounded-r-lg` : `${baseClass} bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer rounded-r-lg`)

</script>
      <div class="inline-flex rounded-md shadow-xs" role="group">
        <button
            onclick={() => goto('./calendar')}
            type="button"
            class={bouttonCalendarClass}
        >
            Calendrier
        </button>
        <button 
            onclick={() => goto('./projects')}
            type="button" 
            class={bouttonProjectsClass}
        >
            Projets
        </button>
        {#if currentUserRole === UserRole.Admin}
            <button 
                onclick={() => goto('./administrator')}
                type="button" 
                class={bouttonAdminClass}
            >
                Admin
            </button>
        {/if}
      </div>

