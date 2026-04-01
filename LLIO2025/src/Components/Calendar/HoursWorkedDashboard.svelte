<script lang="ts">
  import { onMount } from 'svelte';
  import { UserApiService } from '../../services/UserApiService';
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

  
  let timeBankConfig = $state({
    startDate: '',
    hoursPerWeek: 0,
    offset: 0,
  });

  onMount(async () => {
  try {
    const config = await UserApiService.getTimeBankConfig();

    Object.assign(timeBankConfig, config);

    displayedHoursTotal = config.offset;
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

 
  const handleSave = (values) => {
    Object.assign(timeBankConfig, values); 
    displayedHoursTotal = values.offset; 
    showModal = false;
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

  {#if displayedHoursTotal === null || displayedHoursTotal === 0}

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
    initialConfig={timeBankConfig}
  />
{/if}