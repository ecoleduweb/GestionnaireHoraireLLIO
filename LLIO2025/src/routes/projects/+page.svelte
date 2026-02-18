<script lang="ts">
  import { onMount } from "svelte";
  import "../../style/app.css"
  import ProjectsLeftPane from "../../Components/Projects/ProjectsLeftPane.svelte";
  import ProjectComponent from "../../Components/Projects/ProjectComponent.svelte";
  import type { DetailedProject, Project, User, UserInfo } from '../../Models';
  import { ProjectApiService } from "../../services/ProjectApiService";
  import { UserApiService } from "../../services/UserApiService";
  import AddCoManagerModal from '../../Components/Projects/AddCoManagerModal.svelte';

  // État des projets
  let projects = $state<Project[]>([]);
  let isLoading = $state(true);
  let error = $state<string | null>(null);

  let currentUser = $state<UserInfo | null>(null);

  // État pour la modale d'ajout de responsable
  let showAddCoManagerModal = $state(false);
  let selectedUser = $state<string>("");
  let selectedUserNumber = $derived(Number.parseInt(selectedUser))
  let selectedProjectId = $state<number | null>(null);
  let selectedProject = $state<DetailedProject | null>(null);
  let users = $state<User[]>([]);
  let usersToDisplay = $derived<User[]>(
    users.filter(u =>
      selectedProject.managerId != u.id
      &&
      selectedProject.coLeads.findIndex(cl => cl.id === u.id) === -1
    )
  )

  const handleAddCoManagerModalOpen = async (projectId: number) => {
    selectedUser = "";
    selectedProjectId = projectId;
    await loadSelectedProject(projectId);
    await loadUsers();
    showAddCoManagerModal = true;
  }

  const loadProjects = async () => {
    try {
      isLoading = true;
      error = null;
      const response = await ProjectApiService.getDetailedProjects();
      projects = response;
    } catch (err) {
      console.error('Erreur lors de la récupération des projets:', err);
      error = "Impossible de charger les projets. Veuillez réessayer plus tard.";
      projects = [];
    } finally {
      isLoading = false;
    }
  }

  const loadUsers = async () => {
    try {
      isLoading = true;
      error = null;
      users = await UserApiService.getAllManagersAdmin();
    } catch (err) {
      console.error('Erreur lors de la récupération des utilisateurs:', err);
      users = [];
    } finally {
      isLoading = false;
    }
  }

  const loadSelectedProject = async (projectId: number) => {
    try {
      isLoading = true;
      error = null;
      const response = await ProjectApiService.getDetailedProjects();
      selectedProject = response.find(p => p.id === projectId);
    } catch (err) {
      console.error('Erreur lors de la récupération des projets détaillés :', err);
      error = "Impossible de charger les projets détaillés. Veuillez réessayer plus tard.";
      users = [];
    } finally {
      isLoading = false;
    }
  }

  const handleAddCoManager = async () => {
    try {
      isLoading = true;
      error= null;
      await ProjectApiService.addCoManagerToProject(selectedProjectId, selectedUserNumber);
    } catch (e) {
      console.error('Erreur lors de l\'ajout du co-chargé de projet :', e);
      alert(e.message);
    } finally {
      isLoading = false;
      showAddCoManagerModal = false;
      selectedProject = selectedProjectId = selectedUser = users = null;
      await loadProjects();
    }
  }

  onMount(async () => {
    try {
        currentUser = await UserApiService.getUserInfo();
      } catch (error) {
        console.error('Erreur lors du chargement des informations utilisateur:', error);
        alert('Impossible de charger les informations utilisateur.');
      }
    await loadProjects();
    await loadUsers();
  });
</script>

<div class="bg-gray-100">
  {#if currentUser}
  <ProjectsLeftPane {projects} {currentUser} onProjectsRefresh={loadProjects} />
  {/if}
  
  <div class="flex flex-col" style="padding-left: 300px;">
    <!-- Project Details -->
    <div class="p-4">
      <h1 class="text-2xl font-medium text-gray-800">Vos projets en cours</h1>
    </div>
    {#if isLoading}
    <div class="flex justify-center items-center h-screen">
      <p class="text-gray-500">Chargement des projets...</p>
    </div>
  {:else if error}
    <div class="flex justify-center items-center h-screen">
      <p class="text-red-500">{error}</p>
    </div>
  {:else}
    {#each projects as project}
      <ProjectComponent {project} onClickAddCoManager={() => handleAddCoManagerModalOpen(project.id)} />
    {/each}
  {/if}
  </div>
</div>

{#if showAddCoManagerModal}
  <AddCoManagerModal show={showAddCoManagerModal} users={usersToDisplay} bind:valueSelected={selectedUser} projectName={selectedProject.name} onAdd={handleAddCoManager} onCancel={() => showAddCoManagerModal = false} />
{/if}
