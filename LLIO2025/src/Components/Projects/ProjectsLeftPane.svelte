<script lang="ts">
  import { goto } from '$app/navigation';
  import { quintOut } from 'svelte/easing';
  import { slide } from 'svelte/transition';
  import ProjectItem from './ProjectPaneItem.svelte';
  import { Plus } from 'lucide-svelte';
  import ProjectModal from './ProjectModal.svelte';
  import type { Project, UserInfo } from '../../Models';
  import { UserRole } from '$lib/types/enums';
  import { ProjectApiService } from '../../services/ProjectApiService';
  import ConfirmationModal from '../ConfirmationModal.svelte';

  type Props = {
    currentUser: UserInfo;
    projects : Project[];
    onProjectsRefresh: () => void;
  };

  let { projects = [], currentUser, onProjectsRefresh }: Props = $props();
  let isArchivedVisible = $state(false);
  let showModal = $state(false);
  let showModalDelete = $state(false);
  let projectToEdit = $state<Project | null>(null);
  let projectToDelete = $state<Project | null>(null);
    
  const handleNewProject = () =>{
    projectToEdit = null;
    showModal = true;
  }

  const handleEditProject = (project) =>{
    projectToEdit = projects.find((x) => x.id === project.id);
    showModal = true;
  }

  const handleDeleteProject = (project) => {
    projectToDelete = projects.find((x) => x.id === project.id);
    showModalDelete = true;
  }

  const handleCloseModal = () =>{
    showModal = false;
    showModalDelete = false;
    projectToEdit = null;
    projectToDelete = null;
  }

  const handleSuccessDelete = async () => {
    if (projectToDelete?.id != null) {
      await ProjectApiService.deleteProject(projectToDelete.id);
    }
    onProjectsRefresh();
  }
</script>


<div class="dashboard-panel">
  <!-- En-tête du dashboard -->
  <div class="dashboard-header">Tableau de bord</div>

  <!-- Contenu du dashboard -->
  <div class="dashboard-content">
    <!-- Éléments du dashboard -->
    <div class="dashboard-item flex-col">
      <div class="inline-flex rounded-md shadow-xs" role="group">
        <button
          onclick={() => goto('./calendar')}
          type="button"
          class="py-2 px-4 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 rounded-l-lg hover:bg-[#014446] hover:text-white cursor-pointer"
        >
          Calendrier
        </button>
        {#if currentUser.role === UserRole.Admin}
          <button
            type="button" 
            class="py-2 px-4 text-sm transition-colors font-semibold bg-[#014446] text-white"
          >
            Projets
          </button>
          <button 
            onclick={() => goto('./administrator')}
            type="button" 
            class="py-2 px-4 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 rounded-r-lg hover:bg-[#014446] hover:text-white cursor-pointer"
          >
            Admin
          </button>
        {:else}  
          <button 
            type="button" 
            class="py-2 px-4 text-sm transition-colors font-semibold bg-[#014446] text-white rounded-r-lg "
          >
            Projets
          </button>
        {/if}
      </div>

      {#if currentUser.role === UserRole.Admin || currentUser.role === UserRole.ProjectManager} 
        <button
          type="button"
          id="new-project-button"
          title="Créer un nouveau projet"
          onclick={handleNewProject}
          class="mt-4 px-3 py-2 flex items-center gap-2 text-sm transition-colors font-semibold bg-gray-200 text-gray-900 rounded-lg hover:bg-[#014446] hover:text-white cursor-pointer"
        >
          <Plus class="h-4 w-4" /> Créer un nouveau projet 
        </button>
      {/if}
    </div>

    <div class="overflow-y-auto max-h-[calc(100vh-150px)]">
      {#each projects.filter((x) => !x.isArchived) as project}
        <ProjectItem {project} {currentUser} onEdit={handleEditProject} onDelete={handleDeleteProject} />
      {/each}

      <!-- Projets archivés -->
      {#if projects.some((x) => x.isArchived)}
        <div>
          <button
            class="w-full p-4 flex items-center justify-between text-gray-600 hover:bg-gray-50 cursor-pointer"
            onclick={() => (isArchivedVisible = !isArchivedVisible)}
          >
            <span class="font-medium"
              >Projets archivés ({projects.filter((x) => x.isArchived).length})</span
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
              {#each projects.filter((x) => x.isArchived) as project}
                <ProjectItem {project} {currentUser} onEdit={handleEditProject} onDelete={handleDeleteProject} />
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

{#if showModalDelete}
<ConfirmationModal
  modalTitle="Supprimer un projet"
  modalText="Voulez-vous vraiment supprimer le projet {projectToDelete.name} ?"
  errorText="Erreur lors de la suppression du projet, il a soit une ou des activités liées à ce projet ou bien le projet est inexistant"
  onSuccess={handleSuccessDelete}
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