<script lang="ts">
  import { onMount } from 'svelte';
  import { UserApiService } from '../../services/UserApiService';
  import type { TimeBankConfig } from '../../Models/index';
  import HoursWorkedConfigModal from './HoursWorkedConfigModal.svelte';

  type Props = {
    hoursTotal: number | null;
    dateStart: Date;
    dateEnd: Date;
    textHoursWorked: string;
  };

  interface TimeBalance {
    isConfigured: boolean;
    displayedHoursTotal: number | null;
  }

  let { hoursTotal, textHoursWorked }: Props = $props();

  let timeBalance = $state<TimeBalance>({
    isConfigured: false,
    displayedHoursTotal: null,
  });
  let showModal = $state(false);

  let timeBankConfig = $state<TimeBankConfig>({
    startDate: '',
    hoursPerWeek: 0,
    offset: 0,
  });

  const refreshTimeBankBalance = async () => {
    try {
      const balance = await UserApiService.getTimeInBank();
      timeBalance.isConfigured = balance.isConfigured;
      timeBalance.displayedHoursTotal = balance.timeInBank ?? null;
    } catch (err) {
      console.error(err);
    }
  };

  onMount(async () => {
    try {
      const config = await UserApiService.getTimeBankConfig();

      if (config) {
        Object.assign(timeBankConfig, config);
      }

      await refreshTimeBankBalance();
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

  const handleSave = async (values: TimeBankConfig) => {
    Object.assign(timeBankConfig, values);
    await refreshTimeBankBalance();
    showModal = false;
  };
</script>

<div class="card">
  <div class="section">
    {#if !timeBalance.isConfigured}
      <button class="link" on:click={openConfigModal}>Configurer votre banque d'heures</button>
    {:else}
      <p>
        Vous avez
        {' '}
        <button class="link" on:click={openConfigModal} data-testid="total-hours">{timeBalance.displayedHoursTotal ?? 0}</button>
        {' '}heures en banque.
      </p>
    {/if}
  </div>
</div>

{#if showModal}
  <HoursWorkedConfigModal
    onClose={closeConfigModal}
    onSave={handleSave}
    initialConfig={timeBankConfig}
  />
{/if}

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
    display: inline;
    padding: 0;
  }
</style>
