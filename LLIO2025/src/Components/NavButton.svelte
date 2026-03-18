<script lang="ts">
  import { goto } from "$app/navigation";
  import { UserRole } from "$lib/types/enums";
  import { page } from '$app/state';

  type Props = {
      currentUserRole: UserRole;
    };
    
  const path = $derived(page.url.pathname);
  const { currentUserRole } : Props = $props();

</script>

    <div class="dashboard-item ">
      <div class="inline-flex rounded-md shadow-xs" role="group">

        <button
            {#if path === './calendar'}
                class="px-4 py-2 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer rounded-l-lg"
            {:else}    
            onclick={() => goto('./calendar')}
            type="button"
            class="px-4 py-2 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer rounded-l-lg"
            {/if}
        >
            Calendrier
        </button>

        <button 
            onclick={() => goto('./projects')}
            type="button" 
            class="py-2 px-4 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer"
        >
            Projets
        </button>
        
        {#if currentUserRole === UserRole.Admin}
            <button 
                onclick={() => goto('./administrator')}
                type="button" 
                class="py-2 px-4 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 rounded-r-lg hover:bg-[#014446] hover:text-white cursor-pointer"
            >
                Admin
            </button>
        {/if}
      </div>
    </div>

