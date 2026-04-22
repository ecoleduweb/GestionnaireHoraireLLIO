<script lang="ts">
  import type { CoLead, DetailedProject } from '../../Models';

  type Props = {
    show: boolean;
    project: DetailedProject;
    coManager: CoLead;
    onDelete: (userId: number, projectId: number) => void | Promise<void>;
    onCancel: () => void;
  };

  let {
    show,
    project,
    coManager,
    onDelete,
    onCancel,
  }: Props = $props();
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <div class="absolute inset-0 bg-gray-800/50 transition-opacity" onclick={onCancel}></div>

    <div class="bg-white rounded-lg shadow-xl p-8 w-full max-w-xl z-10">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Supprimer un co-chargé</h3>

      <p class="text-gray-600 mb-6">
        Voulez-vous supprimer {coManager.name} des co-chargés de projet de {project.name}?
      </p>

      <div class="flex justify-end gap-3">
        <button
          type="button"
          class="py-2 px-4 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 transition"
          onclick={onCancel}
        >
          Annuler
        </button>
        <button
          type="button"
          class="py-2 px-4 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700 transition"
          onclick={() => onDelete(coManager.id, project.id)}
        >
          Supprimer
        </button>
      </div>
    </div>
  </div>
{/if}
