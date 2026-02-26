<script lang="ts">
  import { UserApiService } from '../../services/UserApiService';
  import type { HoursConfig } from '$lib/types/type';
  import HoursWorkedDashboardSummary from './HoursWorkedDashboardSummary.svelte';
  import HoursWorkedConfigModal from './HoursWorkedConfigModal.svelte';

  const { hoursTotal, dateStart, dateEnd, textHoursWorked } = $props();
  let displayedHoursTotal = $state<number>(hoursTotal);
  let displayedTextHoursWorked = $state<string>(textHoursWorked);

  let showModal = $state(false);
  let isSaving = $state(false);
  let errorMsg = $state<string | null>(null);

  let config = $state<HoursConfig>({
    startDate: '',
    offset: 0,
    hoursWorked: 0,
  });

  // Keep the displayed values in sync with parent props when the calendar view/date range changes.
  $effect(() => {
    displayedHoursTotal = hoursTotal;
    displayedTextHoursWorked = textHoursWorked;
  });

  $effect(() => {
    const saved = localStorage.getItem('hoursConfig');
    if (saved) {
      try {
        config = JSON.parse(saved);
      } catch (err) {
        console.warn('Failed to parse hoursConfig from localStorage', err);
      }
    }
  });

  // on mount, try to load server-calculated time in bank
  $effect(() => {
    (async () => {
      try {
        const bank = await UserApiService.getTimeInBank();
        if (bank && typeof bank.timeInBank === 'number') {
          displayedHoursTotal = bank.timeInBank;
        }
        if (bank && typeof (bank as any).textHoursWorked === 'string') {
          displayedTextHoursWorked = (bank as any).textHoursWorked;
        }
      } catch (err) {
        console.warn('Time bank fetch failed', err);
      }
    })();
  });

  const handleSave = async () => {
    errorMsg = null;
    isSaving = true;

    try {
      localStorage.setItem('hoursConfig', JSON.stringify(config));

      await UserApiService.saveTimeBankConfig({
        startDate: config.startDate,
        hoursPerWeek: Number(config.hoursWorked),
        offset: Number(config.offset),
      });

      const bank = await UserApiService.getTimeInBank();

      if (bank && typeof bank.timeInBank === 'number') {
        displayedHoursTotal = bank.timeInBank;
      }

      if (bank && typeof (bank as any).textHoursWorked === 'string') {
        displayedTextHoursWorked = (bank as any).textHoursWorked;
      }

      showModal = false;
    } catch (err) {
      errorMsg = "Impossible de calculer pour l'instant. Vérifie l'API /user/time-bank et /user/time-bank/config.";
      console.error(err);
    } finally {
      isSaving = false;
    }
  };

  const openConfigModal = () => {
    showModal = true;
  };

  const closeConfigModal = () => {
    showModal = false;
  };

  const updateConfig = (patch: Partial<HoursConfig>) => {
    config = { ...config, ...patch };
  };
</script>

<div class="bilan-container">
  <HoursWorkedDashboardSummary
    {dateStart}
    {dateEnd}
    {displayedHoursTotal}
    {displayedTextHoursWorked}
    onOpenConfig={openConfigModal}
  />
</div>

<HoursWorkedConfigModal
  {showModal}
  {isSaving}
  {errorMsg}
  {config}
  onClose={closeConfigModal}
  onSave={handleSave}
  onUpdateConfig={updateConfig}
/>

<style>
  .bilan-container {
    padding: 1rem;
    border: 1px solid #ddd;
    border-radius: 5px;
    position: relative;
  }
</style>
