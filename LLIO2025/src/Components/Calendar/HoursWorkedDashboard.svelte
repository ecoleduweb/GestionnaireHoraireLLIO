<script lang="ts">
  import { formatDateHoursWorked, areDatesEqual } from '../../utils/date';
  import { X } from 'lucide-svelte';
  import { UserApiService } from '../../services/UserApiService';
  import type { HoursConfig } from '$lib/types/type';
  

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
      // Step 1: Save config locally
      localStorage.setItem('hoursConfig', JSON.stringify(config));

      // Step 2: Send config to server (POST /user/time-bank/config)
      await UserApiService.saveTimeBankConfig({
        startDate: config.startDate,
        hoursPerWeek: Number(config.hoursWorked),
        offset: Number(config.offset),
      });

      // Step 3: Fetch calculated time-in-bank (GET /user/time-bank)
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
</script>

<div class="bilan-container">
  <div class="header">
    <h2>
      {#if areDatesEqual(dateStart, dateEnd)}
        Bilan du {formatDateHoursWorked(dateStart)}
      {:else}
        Bilan du {formatDateHoursWorked(dateStart)} au {formatDateHoursWorked(dateEnd)}
      {/if}
    </h2>
  </div>

  <span>
    Vous avez travaillé <strong>{displayedHoursTotal}</strong> heures {displayedTextHoursWorked}.
  </span>

  <div class="config-section">
    <button
      type="button"
      class="config-btn"
      title="Configurer les heures en banque"
      onclick={() => (showModal = true)}
    >
      Configurer
    </button>
  </div>
</div>

{#if showModal}
  <div
    class="modal-overlay"
    role="button"
    tabindex="0"
    aria-label="Fermer la fenêtre"
    onclick={() => (showModal = false)}
    onkeydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') showModal = false;
      if (e.key === 'Escape') showModal = false;
    }}
  >
    <div
      class="modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="hours-modal-title"
      tabindex="-1"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => {
        if (e.key === 'Escape') showModal = false;
      }}
    >
      <div class="modal-header">
        <h3 id="hours-modal-title">Configuration des heures en banque</h3>
        <button type="button" class="close-btn" aria-label="Fermer" onclick={() => (showModal = false)}>
          <X size={18} />
        </button>
      </div>

      <div class="modal-body">
        {#if errorMsg}
          <div class="error">{errorMsg}</div>
        {/if}

        <div class="form-group">
          <label for="startDate">Début de la période</label>
          <input id="startDate" type="date" bind:value={config.startDate} class="form-input" />
        </div>

        <div class="form-group">
          <label for="hoursWorked">Nombre d'heures par semaine</label>
          <input
            id="hoursWorked"
            type="number"
            bind:value={config.hoursWorked}
            min="0"
            step="0.5"
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label for="offset">Décalage (offset)</label>
          <input id="offset" type="number" bind:value={config.offset} step="0.5" class="form-input" />
        </div>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn-cancel" onclick={() => (showModal = false)} disabled={isSaving}>
          Annuler
        </button>

        <button type="button" class="btn-save" onclick={handleSave} disabled={isSaving}>
          {#if isSaving}
            Calcul...
          {:else}
            Enregistrer
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .bilan-container {
    padding: 1rem;
    border: 1px solid #ddd;
    border-radius: 5px;
    position: relative;
  }

  .header {
    margin-bottom: 1rem;
  }

  .header h2 {
    margin: 0;
  }

  .config-section {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid #eee;
    display: flex;
    justify-content: flex-start;
  }

  .config-btn {
    background: #015e61;
    color: white;
    border: none;
    cursor: pointer;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    transition: all 0.2s;
    font-weight: 500;
    font-size: 0.95rem;
  }

  .config-btn:hover {
    background: #014446;
  }

  /* Modal styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: white;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.15);
    max-width: 400px;
    width: 90%;
    max-height: 85vh;
    overflow-y: auto;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #eee;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.1rem;
  }

  .close-btn {
    background: none;
    border: none;
    cursor: pointer;
    color: #666;
    padding: 0;
    display: flex;
    align-items: center;
  }

  .close-btn:hover {
    color: #000;
  }

  .modal-body {
    padding: 1.5rem;
  }

  .error {
    margin-bottom: 1rem;
    padding: 0.75rem;
    border: 1px solid #f3c2c2;
    background: #fff3f3;
    border-radius: 6px;
    font-size: 0.9rem;
    color: #d32f2f;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group:last-child {
    margin-bottom: 0;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    font-size: 0.9rem;
  }

  .form-input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    font-family: inherit;
    box-sizing: border-box;
  }

  .form-input:focus {
    outline: none;
    border-color: #015e61;
    box-shadow: 0 0 0 3px rgba(1, 94, 97, 0.1);
  }

  .modal-footer {
    display: flex;
    gap: 1rem;
    padding: 1.5rem;
    border-top: 1px solid #eee;
    justify-content: flex-end;
  }

  .btn-cancel,
  .btn-save {
    padding: 0.7rem 1.4rem;
    border: none;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.9rem;
  }

  .btn-cancel {
    background: #f0f0f0;
    color: #333;
  }

  .btn-cancel:hover {
    background: #e0e0e0;
  }

  .btn-save {
    background: #015e61;
    color: white;
  }

  .btn-save:hover {
    background: #014446;
  }

  .btn-cancel:disabled,
  .btn-save:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
</style>
