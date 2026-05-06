<script lang="ts">
  import { projectTemplate } from '../../forms/project';
  import { ProjectApiService } from '../../services/ProjectApiService';
  import { UserApiService } from '../../services/UserApiService';
  import type { Project, ProjectBase, User } from '../../Models/index';
  import { X } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { validateProjectForm } from '../../Validation/Project';

  type Props = {
    projectToEdit: Project | null;
    onClose: () => void;
    onSuccess: () => void;
  };

  let { projectToEdit, onClose, onSuccess }: Props = $props();

  const project = $state<ProjectBase>(projectTemplate.generate());

  $effect(() => {
    if (projectToEdit) {
      Object.assign(project, projectToEdit);
      project.estimatedHours = projectToEdit.totalTimeEstimated || 0;
    }
  });

  const editMode = $derived(projectToEdit !== null);

  let isSubmitting = $state(false);
  let isLoadingManagers = $state(true);
  let isLoadingProject = $state(false);
  let managers = $state<User[]>([]);

  onMount(async () => {
    try {
      isLoadingManagers = true;
      managers = await UserApiService.getAllManagersAdmin();
    } catch (err) {
      console.error('Failed to load managers:', err);
      alert('Impossible de charger les chargés de projet.');
    } finally {
      isLoadingManagers = false;
    }
  });

  const handleClose = () => {
    onClose();
  };

  const handleSubmit = async () => {
    if (isSubmitting) return;
    try {
      isSubmitting = true;
      if (editMode) {
        await ProjectApiService.updateProject(project);
      } else {
        await ProjectApiService.createProject(project);
      }
      onSuccess();
      onClose();
    } catch (err) {
      console.error('Erreur lors de la soumission du projet', err);
      alert('Erreur lors de la soumission du projet');
    } finally {
      isSubmitting = false;
    }
  };
  const { form, errors } = validateProjectForm(handleSubmit, project);
</script>

  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2 class="modal-title">{editMode ? 'Modifier le projet' : 'Créer un nouveau projet'}</h2>
        <button type="button" class="text-black hover:text-gray-600" onclick={handleClose}>
          <X />
        </button>
      </div>

      <div class="modal-content">
        {#if isLoadingProject}
          <div class="py-2 px-4 bg-gray-100 rounded">Chargement du projet...</div>
        {:else}
          <form
            class="flex flex-col h-full"
            use:form
            onsubmit={(e) => {
              e.preventDefault();
            }}
          >
            <div class="form-group">
              <label for="project-uniqueId">Identifiant unique (numéro Airtable)*</label>
              <label>
                <input id="project-uniqueId" name="uniqueId" type="text" bind:value={project.uniqueId} />
              </label>
              {#if $errors.uniqueId}
                <span class="text-red-500 text-sm">{$errors.uniqueId}</span>
              {/if}
            </div>

            <div class="form-group">
              <label for="project-manager">Chargé de projet*</label>
              {#if isLoadingManagers}
                <div class="py-2 px-4 bg-gray-100 rounded">Chargement des managers...</div>
              {:else}
                <select id="project-manager" name="managerId" bind:value={project.managerId}>
                  <option value="">-- Sélectionner un manager --</option>
                  {#each managers as manager}
                    <option value={manager.id}>
                      {`${manager.firstName} ${manager.lastName}`}
                    </option>
                  {/each}
                </select>
              {/if}
              {#if $errors.managerId}
                <span class="text-red-500 text-sm">{$errors.managerId}</span>
              {/if}
            </div>

            <div class="form-group">
              <label for="project-name">Nom du projet</label>
              <input
                id="project-name"
                name="name"
                bind:value={project.name}
                
              >
              {#if $errors.name}
                <span class="text-red-500 text-sm">{$errors.name}</span>
              {/if}
            </div>

            <div>
              <label for="estimated-hours">Heures estimées</label>
              <input
                id="estimated-hours"
                name="estimatedHours"
                type="number"
                min="0"
                step="1"
                bind:value={project.estimatedHours}
              >
              {#if $errors.estimatedHours}
                <span class="text-red-500 text-sm">{$errors.estimatedHours}</span>
              {/if}
            </div>

            <div class="modal-footer">
              {#if editMode}
                <button
                  type="button"
                  class="py-3 px-6 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 hover:-translate-y-0.5 hover:shadow-sm active:translate-y-0 transition border border-gray-200"
                  onclick={handleClose}
                >
                  Retour
                </button>
                <button
                  type="submit"
                  class="py-3 px-6 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition disabled:opacity-50"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? 'En cours...' : 'Modifier'}
                </button>
              {:else}
                <button
                  type="button"
                  class="py-3 px-6 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 hover:-translate-y-0.5 hover:shadow-sm active:translate-y-0 transition border border-gray-200"
                  onclick={handleClose}
                >
                  Annuler
                </button>
                <button
                  type="submit"
                  class="py-3 px-6 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition disabled:opacity-50"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? 'En cours...' : 'Soumettre'}
                </button>
              {/if}
            </div>
          </form>
        {/if}
      </div>
    </div>
  </div>


<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal {
    background-color: white;
    border-radius: 4px;
    width: 400px;
    max-width: 90%;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }

  .modal-header {
    padding: 12px 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #eee;
  }

  .modal-title {
    font-size: 18px;
    margin: 0;
    color: #333;
  }

  .modal-content {
    padding: 24px;
  }

  .form-group {
    margin-bottom: 16px;
  }

  label {
    display: block;
    margin-bottom: 8px;
    color: #666;
    font-size: 14px;
  }

  input,
  select {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #ccc;
    border-radius: 4px;
    box-sizing: border-box;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
  }

</style>
