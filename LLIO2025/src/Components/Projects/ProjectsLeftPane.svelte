<script lang="ts">
  import { goto } from '$app/navigation';
  import { quintOut } from 'svelte/easing';
  import { slide } from 'svelte/transition';
  import ProjectItem from './ProjectPaneItem.svelte';
  import { Plus } from 'lucide-svelte';
  import ProjectModal from './ProjectModal.svelte';
  import type { Project, UserInfo } from '../../Models';
  import { UserRole } from '$lib/types/enums';

  type Props = {
    currentUser: UserInfo;
    filteredComponant : Project[];
    onProjectsRefresh: () => void;
  };

  let { filteredComponant = [], currentUser, onProjectsRefresh }: Props = $props();
  let isArchivedVisible = $state(false);
  let showModal = $state(false);
  let projectToEdit = $state<Project | null>(null);
    
  const handleNewProject = () =>{
    projectToEdit = null;
    showModal = true;
  }

  const handleEditProject = (project) =>{
    projectToEdit = filteredComponant.find((x) => x.id === project.id);
    showModal = true;
  }

  const handleCloseModal = () =>{
    showModal = false;
    projectToEdit = null;
  }
</script>


<div class="dashboard-panel">
  <!-- En-tête du dashboard -->
  <div class="dashboard-header">Tableau de bord</div>

  <!-- Contenu du dashboard -->
  <div class="dashboard-content">
    <!-- Éléments du dashboard -->
    <div class="dashboard-item">
      <div class="inline-flex rounded-md shadow-xs" role="group">
        <button
          onclick={() => goto('./calendar')}
          type="button"
          class="py-2 px-4 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 rounded-l-lg hover:bg-[#014446] hover:text-white cursor-pointer"
        >
          Calendrier
        </button>
        <button
          type="button"
          class="px-4 py-2 text-sm transition-colors font-semibold bg-[#014446] text-white rounded-r-lg"
        >
          Projets
        </button>
      </div>

      {#if currentUser.role == UserRole.Admin || currentUser.role == UserRole.ProjectManager} 
        <button
          type="button"
          id="new-project-button"
          title="Créer un nouveau projet"
          onclick={handleNewProject}
          class="ml-12 px-3 py-2 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 rounded-lg hover:bg-[#014446] hover:text-white cursor-pointer"
        >
          <Plus class="h-4 w-4" />
        </button>
      {/if}
    </div>

    <div class="overflow-y-auto max-h-[calc(100vh-150px)]">
      {#each filteredComponant.filter((x) => !x.isArchived) as project}
        <ProjectItem {project} {currentUser} onEdit={handleEditProject} />
      {/each}

      <!-- Projets archivés -->
      {#if filteredComponant.some((x) => x.isArchived)}
        <div>
          <button
            class="w-full p-4 flex items-center justify-between text-gray-600 hover:bg-gray-50 cursor-pointer"
            onclick={() => (isArchivedVisible = !isArchivedVisible)}
          >
            <span class="font-medium"
              >Projets archivés ({filteredComponant.filter((x) => x.isArchived).length})</span
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
              {#each filteredComponant.filter((x) => x.isArchived) as project}
                <ProjectItem {project} {currentUser} onEdit={handleEditProject} />
              {/each}
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
</div>

{#if showModal}
<ProjectModal
  projectToEdit={projectToEdit}
  onSuccess={onProjectsRefresh}
  onClose={handleCloseModal}
/>
{/if}

<style>
  .dashboard-panel {
    width: 300px;
    height: 100vh;
    background-color: white;
    border-right: 1px solid #e5e7eb;
    box-shadow: 2px 0 5px rgba(0, 0, 0, 0.05);
    overflow-y: auto;
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 30;
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
    padding: 0px;
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