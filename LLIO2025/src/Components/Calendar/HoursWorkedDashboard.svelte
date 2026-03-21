<script lang="ts">
  import { onMount } from 'svelte';
  import { UserApiService } from '../../services/UserApiService';
  import { formatDateHoursWorked, areDatesEqual } from '../../utils/date';
  import HoursWorkedConfigModal from './HoursWorkedConfigModal.svelte';

  type Props = {
    hoursTotal: number | null;
    dateStart: Date;
    dateEnd: Date;
    textHoursWorked: string;
  };

  let { hoursTotal, dateStart, dateEnd, textHoursWorked }: Props = $props();

  let displayedHoursTotal = $state<number | null>(hoursTotal ?? null);
  let showModal = $state(false);

  onMount(async () => {
    try {
      const bank = await UserApiService.getTimeInBank();
      displayedHoursTotal = bank.timeInBank;
    } catch (err) {
      console.error(err);
    }
  });

  const openConfigModal = () => {
    showModal = true;
  };

  const closeConfigModal = () => {
    showModal = false;
  };

  const handleSave = async () => {
    try {
      const bank = await UserApiService.getTimeInBank();
      displayedHoursTotal = bank.timeInBank;
      showModal = false;
    } catch (err) {
      alert("Erreur lors de la configuration");
    }
  };
</script>

<style>
  .card {
    background: #f5f5f5;
    padding: 2rem;
    border-radius: 12px;
    max-width: 500px;
  }

  .section {
    margin-bottom: 1rem;
  }

  .label {
    font-size: 0.75rem;
    letter-spacing: 1px;
    color: #999;
    margin-bottom: 0.5rem;
  }

  .link {
    color: #2563eb;
    text-decoration: underline;
    cursor: pointer;
    background: none;
    border: none;
    font: inherit;
    padding: 0;
  }
</style>

<div class="card">

 
  {#if displayedHoursTotal === null}

    <div class="section">
      <p>
        Vous avez 
        <button class="link" on:click={openConfigModal}>
          configurer
        </button>
        heures en banque.
      </p>
    </div>

  {:else}

    
    <div class="section">
      

      <p>
        Vous avez 
        <button
          class="link"
          on:click={openConfigModal}
          data-testid="total-hours"
          title="Modifier les heures"
        >
          {displayedHoursTotal}
        </button>
        heures en banque.
      </p>
    </div>

  {/if}

</div>

{#if showModal}
  <HoursWorkedConfigModal 
    onClose={closeConfigModal} 
    onSave={handleSave} 
  />
{/if}