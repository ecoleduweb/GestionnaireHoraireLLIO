<script lang="ts">
  import { onMount } from 'svelte';
  import { UserApiService } from '../../services/UserApiService';
  import type { TimeBalance, TimeBankConfig } from '../../Models/index';
  import HoursWorkedConfigModal from '../Calendar/HoursWorkedConfigModal.svelte';
  import { refreshTimeBankSignal } from '../../lib/refreshTimeBankSignal.svelte';
    
  let timeBalance = $state<TimeBalance>({
    isConfigured: false,
    displayedHoursTotal: null,
  });

  let showModal = $state(false);

  let config = $state<TimeBankConfig>({
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
      const timeBankConfig = await UserApiService.getTimeBankConfig();

      if (timeBankConfig) {
        Object.assign(config, timeBankConfig);
      }
    } catch (err) {
      console.error(err);
    }
  });

    // $effect agit comme onMount et est aussi mis à jour quand une variable state est mise à jour.
  $effect(() => {
    void refreshTimeBankSignal.tick; // permet de déclancher le rafraichissement lorsque refreshTimeBankSignal.tick est incrémenté (voir le fichier refreshTimeBankSignal.svelte.ts
   refreshTimeBankBalance();
  });

  const openConfigModal = () => {
    showModal = true;
  };

  const closeConfigModal = () => {
    showModal = false;
  };

  const handleSave = async (values: TimeBankConfig) => {
    Object.assign(config, values);
    await refreshTimeBankBalance();
    showModal = false;
  };
</script>

<div class="time-bank-section">
  {#if !timeBalance.isConfigured}
    <button class="link" on:click={openConfigModal}>Configurer votre banque d'heures</button>
  {:else}
    <p class="time-bank-text">
      Vous avez
      {' '}
      <button class="link" on:click={openConfigModal} data-testid="total-hours">{timeBalance.displayedHoursTotal ?? 0}</button>
      {' '}heures en banque.
    </p>
  {/if}
</div>

{#if showModal}
  <HoursWorkedConfigModal
    onClose={closeConfigModal}
    onSave={handleSave}
    initialConfig={config}
  />
{/if}

<style>
  .time-bank-section {
    padding: 1rem;
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

  .time-bank-text {
    margin: 0;
    line-height: 1.6;
    color: #1f2937;
  }
</style>
