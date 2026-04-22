<script lang="ts">
  import { goto } from '$app/navigation';
  import { quintOut } from 'svelte/easing';
  import { slide } from 'svelte/transition';
  import DashboardProjectItem from '../Projects/DashboardPaneProjectItem.svelte';
  import type { DetailedProject, UserInfo } from '../../Models/index.ts';
  import { UserRole } from '../../lib/types/enums';
  import NavButton from '../NavButton.svelte';
  import HoursWorkedDashboard from './HoursWorkedDashboard.svelte';

  import TimeBank from './TimeBank.svelte';

  type Props = {
    detailedProjects: DetailedProject[];
    currentUser: UserInfo | null;
    totalHours: number;
    dateStart: string;
    dateEnd: string;
    textHoursWorked: string;
  };

  const {
    totalHours,
    detailedProjects = [],
    dateStart,
    dateEnd,
    textHoursWorked,
    currentUser,
  }: Props = $props();
  let isArchivedVisible = $state(false);
</script>



<div class="dashboard-container">
  <!-- En-tête du dashboard -->
  <div class="dashboard-header">Tableau de bord</div>

  <!-- Contenu du dashboard -->
  <div class="dashboard-content">
    <!-- Contenu à venir -->
    <div class="dashboard-item">
      <NavButton currentUserRole={currentUser.role} />
    </div>

    <!-- Projets en cours -->
    <div>
      {#each detailedProjects.filter((x) => !x.isArchived) as project}
        <DashboardProjectItem {project} />
      {/each}

      <!-- Projets archivés -->
      {#if detailedProjects.some((x) => x.isArchived)}
        <div>
          <button
            class="w-full p-4 flex items-center justify-between text-gray-600 hover:bg-gray-50 cursor-pointer"
            onclick={() => (isArchivedVisible = !isArchivedVisible)}
          >
            <span class="font-medium"
              >Projets archivés ({detailedProjects.filter((x) => x.isArchived).length})</span
            >
            <svg
              class="w-5 h-5 transform transition-transform {isArchivedVisible ? 'rotate-180' : ''}"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 9l-7 7-7-7"
              ></path>
            </svg>
          </button>

          {#if isArchivedVisible}
            <div transition:slide={{ duration: 300, easing: quintOut }}>
              {#each detailedProjects.filter((x) => x.isArchived) as project}
                <DashboardProjectItem {project} />
              {/each}
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
  
 <div class="dashboard-card">

  <HoursWorkedDashboard hoursTotal={totalHours} {dateStart} {dateEnd} {textHoursWorked} />
  <TimeBank {totalHours}/>
  </div>
</div>


<style>
 .dashboard-card {
    background: #f5f5f5;
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
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

  .dashboard-item:hover {
    background-color: #f5f5f5;
  }
</style>