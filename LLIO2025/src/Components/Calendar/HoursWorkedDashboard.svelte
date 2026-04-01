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

  let { hoursTotal }: Props = $props();

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

  .link {
    color: #2563eb;
    text-decoration: underline;
    cursor: pointer;
    background: none;
    border: none;
  }
</style>

<div class="card">
  <div class="section">
    <p>
      Vous avez
      {#if displayedHoursTotal === null || displayedHoursTotal === 0}
        <button class="link" on:click={openConfigModal}> configurer </button>
      {:else}
        <button class="link" on:click={openConfigModal} data-testid="total-hours">
          {displayedHoursTotal}
        </button>
      {/if}
      heures en banque.
    </p>
  </div>
</div>

{#if showModal}
  <HoursWorkedConfigModal
    onClose={closeConfigModal}
    onSave={handleSave}
    initialConfig={timeBankConfig}
  />
{/if}
