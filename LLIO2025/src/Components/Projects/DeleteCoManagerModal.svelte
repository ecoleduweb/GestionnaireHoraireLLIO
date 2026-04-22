<script lang="ts">
  import { Trash2 } from 'lucide-svelte';
  import type { CoLead } from '../../Models';
  import { ProjectApiService } from '../../services/ProjectApiService';
  import ConfirmationModal from '../Modal/ConfirmationModal.svelte';

  let {
    projectId,
    coManager,
    onClickDeleteCoManager = () => {},
  }: {
    projectId: number;
    coManager: CoLead;
    onClickDeleteCoManager?: (projectId: number, coManager: CoLead) => void;
  } = $props();

  let showConfirmationModal = $state(false);

  const handleDeleteCoManager = async () => {
    await ProjectApiService.deleteCoManagerFromProject(projectId, coManager.id);
    onClickDeleteCoManager(projectId, coManager);
  };
</script>

<div class="text-sm wrap-normal">
  {coManager.name}
  <button
    class="p-1 text-gray-500 hover:text-red-700 hover:bg-red-50 rounded-full transition-colors"
    onclick={() => showConfirmationModal = true}
    aria-label="Supprimer le co-chargé {coManager.name}"
  >
    <Trash2 size={14} />
  </button>
</div>

{#if showConfirmationModal}
  <ConfirmationModal
    modalTitle="Supprimer un co-chargé"
    modalText="Voulez-vous vraiment supprimer {coManager.name} des co-chargés de ce projet ?"
    confirmText="Supprimer"
    errorText="Erreur lors de la suppression du co-chargé"
    onClose={() => showConfirmationModal = false}
    onSuccess={handleDeleteCoManager}
  />
{/if}
