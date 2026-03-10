<script lang="ts">
  import type { UserInfo } from "../../Models";
  import { UserApiService } from "../../services/UserApiService";
  import UsersModal from '../../Components/Administrator/UsersModal.svelte';
  import { LogOut } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  
  let currentUser = $state<UserInfo | null>(null);
  let isLoading = $state(false);    
  let showUsersModal = $state(false);
   
  onMount(async () => {
    isLoading = true;    
    // Charger les informations utilisateur
    try {
      currentUser = await UserApiService.getUserInfo();
    } catch (error) {
      console.error('Erreur lors du chargement des informations utilisateur:', error);
    } finally {
      isLoading = false;
    }
  });

</script>
<div class="flex">
  <div class="dashboard-container">
  <!-- En-tête du dashboard -->
  <div class="dashboard-header">Tableau de bord</div>

  <!-- Contenu du dashboard -->
  <div class="dashboard-content">
    <!-- Contenu à venir -->
    <div class="dashboard-item ">
      <div class="inline-flex rounded-md shadow-xs" role="group">
        <button
            onclick={() => goto('./calendar')}
            type="button"
            class="px-4 py-2 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 hover:bg-[#014446] hover:text-white cursor-pointer rounded-l-lg"
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
        <button             
            type="button" 
            class="py-2 px-4 text-sm transition-colors font-semibold bg-[#014446] text-white rounded-r-lg "
        >
            Admin
        </button>
      </div>
    </div>
    <div class="dashboard-item-button flex-col justify-center-safe content-center">
      <div class="mt-4">
        <button 
          onclick={() => showUsersModal = !showUsersModal }
          id="user-button"
          type="button" 
          class="w-full py-2 px-4 text-sm font-medium transition-colors bg-[#e6f0f0] text-[#005e61] rounded-md hover:bg-[#d0e6e6] flex items-center justify-center cursor-pointer"
          >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
          </svg>
          Gestion des utilisateurs
        </button>
      </div>
    </div>
  </div>
  </div>    
  <div class="space-between-dashboard-calendar w-full h-full bg-white px-4 py-6">
    <div class="max-w-7xl mx-auto">
      <!-- Nouvelle section avec salutation -->
      <div class="flex justify-between items-center mb-6">
        <!-- Affichage nom d'utilisateur -->
        <h1 class="text-2xl font-bold text-gray-800 flex items-center gap-2">
          Espace administrateur,
          <span class="text-[#015e61] font-bold">
            {#if currentUser}
              {currentUser.firstName} {currentUser.lastName}
            {:else}
              <span class="inline-block w-24 h-6 bg-gray-200 animate-pulse rounded"></span>
            {/if}
          </span>
          <button
            class="ml-2 mt-1 p-1.5 rounded-full hover:bg-gray-100 text-gray-600 hover:text-[#015e61] transition-colors"
            title="Se déconnecter" 
            onclick={async () => {
              await UserApiService.logOut();
              goto("/");
            }} 
            >
            <LogOut class="w-5 h-5" />
          </button>
        </h1>
      </div>
    </div>
  </div>
</div>   
  {#if showUsersModal}
    <UsersModal/>
  {/if}
 

<style>
  /* Gestion espace dashboard */
  .space-between-dashboard-calendar {
    margin-left: 300px;
  }

  .dashboard-container {
    width: 300px;
    height: 100vh;
    background-color: white;
    border-right: 1px solid #e5e7eb;
    box-shadow: 2px 0 5px rgba(0, 0, 0, 0.05);
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 30;
    display: flex;
    flex-direction: column;
  }

  .dashboard-header {
    background-color: #005e61;
    color: white;
    padding: 16px;
    font-weight: 600;
    font-size: 1.25rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .dashboard-content {
    flex: 1;
    overflow-y: auto;
  }

  .dashboard-item {
    padding: 16px;
    border-bottom: 1px solid #f0f0f0;
    transition: background-color 0.2s;
  }

  .dashboard-item-button {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    border-bottom: 1px solid #f0f0f0;
    transition: background-color 0.2s;
  }

  .dashboard-item:hover {
    background-color: #f5f5f5;
  }
</style>