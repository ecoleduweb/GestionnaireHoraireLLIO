<script lang="ts">
  import { formatDateHoursWorked, areDatesEqual } from '../../utils/date';

  let { dateStart, dateEnd, displayedHoursTotal, displayedTextHoursWorked, onOpenConfig } = $props();
</script>

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
    onclick={onOpenConfig}
  >
    Configurer
  </button>
</div>

<style>
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
</style>
