<script lang="ts">
  import CrayonEdit from '../../assets/svg/crayon.svg.svelte';
  import Modal from '../Modal/BaseModal.svelte';

  let { employee, onRateUpdated = () => {} }: {
    employee: any;
    onRateUpdated?: (employee: any, newRate: number) => void;
  } = $props();

  let showRateModal = $state(false);
  let newRate = $state<number | string>('');

  function openModal(e: MouseEvent) {
    e.stopPropagation();
    newRate = employee.hourlyRate ?? '';
    showRateModal = true;
  }

  async function handleSave() {
    onRateUpdated(employee, Number(newRate));
    showRateModal = false;
  }
</script>

<button
  aria-label="Modifier le taux horaire"
  onclick={openModal}
  type="button"
  class="p-1 rounded-md hover:bg-gray-200 flex items-center justify-center"
>
  <CrayonEdit />
</button>

{#if showRateModal}
  <Modal
    modalTitle="Modifier le taux horaire"
    confirmText="Enregistrer"
    cancelText="Annuler"
    onClose={() => (showRateModal = false)}
    onSuccess={handleSave}
  >
    {#snippet children()}
      <div class="flex flex-col gap-3">
        <p class="text-sm text-gray-500">{employee.name}</p>
        <input
          type="number"
          bind:value={newRate}
          class="w-full border rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-[#005e61]"
          placeholder="Nouveau taux horaire"
        />
      </div>
    {/snippet}
  </Modal>
{/if}