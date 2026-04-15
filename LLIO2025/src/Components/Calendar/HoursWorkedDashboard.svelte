<script lang="ts">
  import TimeBank from './TimeBank.svelte';

  type Props = {
    hoursTotal: number;
    dateStart: string;
    dateEnd: string;
    textHoursWorked: string;
  };

  const { hoursTotal = 0, dateStart, dateEnd, textHoursWorked }: Props = $props();

  const formatDate = (date: string | Date) => {
    let dateObj: Date;

    if (date instanceof Date) {
      dateObj = date;
    } else {
      const parts = date.split('-');
      dateObj = new Date(Number(parts[0]), Number(parts[1]) - 1, Number(parts[2]), 12, 0, 0);
    }

    return new Intl.DateTimeFormat('fr-FR', {
      day: 'numeric',
      month: 'long',
    }).format(dateObj);
  };

  const areDatesEqual = () => {
    if (!dateStart || !dateEnd) return false;
    return dateStart === dateEnd;
  };
</script>

<div class="dashboard-card">
  <div class="bilan-container">
    <div class="header" data-testid="hours-worked-period">
      <h2>
        {#if areDatesEqual()}
          Bilan du {formatDate(dateStart)}
        {:else}
          Bilan du {formatDate(dateStart)} au {formatDate(dateEnd)}
        {/if}
      </h2>
    </div>

    <p class="summary-text" data-testid="hours-worked-summary">
      Vous avez travaillé <strong>{hoursTotal.toFixed(2)}</strong> heures {textHoursWorked}.
    </p>
  </div>

  <TimeBank />
</div>

<style>
  .dashboard-card {
    background: #f5f5f5;
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
  }

  .bilan-container {
    padding: 1rem;
  }

  .header {
    margin-bottom: 0.75rem;
  }

  .header h2 {
    margin: 0;
    font-size: 1.125rem;
    font-weight: 500;
    color: #111827;
  }

  .summary-text {
    margin: 0;
    line-height: 1.6;
    color: #1f2937;
  }
</style>
