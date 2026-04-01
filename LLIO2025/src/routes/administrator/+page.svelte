<script lang="ts">
  import type { UserInfo } from "../../Models";
  import { UserApiService } from "../../services/UserApiService";
  import UsersModal from '../../Components/Administrator/UsersModal.svelte';
  import { LogOut } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import NavButton from "../../Components/NavButton.svelte";
 import { ReportApiService } from '../../services/ReportApiService';
  
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
      <NavButton currentUserRole = {currentUser?.role} />
    </div>
    <div class="dashboard-item-button flex-col justify-center-safe content-center">
        <button 
          onclick={() => showUsersModal = !showUsersModal }
          id="user-button"
          type="button" 
          class="w-full mt-4 py-2 px-4 text-sm font-medium transition-colors bg-[#e6f0f0] text-[#005e61] rounded-md hover:bg-[#d0e6e6] flex items-center justify-center cursor-pointer"
          >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
          </svg>
          Gestion des utilisateurs
        </button>
        <button 
          onclick={() => ReportApiService.getReportExcel()}
          id="user-button"
          type="button" 
          class="w-full mt-4 py-2 px-4 text-sm font-medium transition-colors bg-[#e6f0f0] text-[#005e61] rounded-md hover:bg-[#d0e6e6] flex items-center justify-center cursor-pointer"
          >
          <svg class="w-4 h-4 mr-2" viewBox="0 0 192 192" xmlns="http://www.w3.org/2000/svg" fill="none">
            <g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
            <g id="SVGRepo_iconCarrier">
              <path d="M56 30c0-1.662 1.338-3 3-3h108c1.662 0 3 1.338 3 3v132c0 1.662-1.338 3-3 3H59c-1.662 0-3-1.338-3-3v-32m0-68V30" style="fill-opacity:.402658;stroke:#005e61;stroke-width:12;stroke-linecap:round;paint-order:stroke fill markers"></path>
              <rect width="68" height="68" x="-58.1" y="40.3" rx="3" style="fill:none;fill-opacity:.402658;stroke:#005e61;stroke-width:12;stroke-linecap:round;stroke-linejoin:miter;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" transform="translate(80.1 21.7)"></rect>
              <path d="M138.79 164.725V27.175M56.175 58.792H170M170 96H90.328M169 133.21H56.175M44.5 82l23 28m0-28-23 28" style="fill:none;stroke:#005e61;stroke-width:12;stroke-linecap:round;stroke-linejoin:round;stroke-dasharray:none;stroke-opacity:1"></path>
            </g>
          </svg>         
          Générer un export Excel
        </button>
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