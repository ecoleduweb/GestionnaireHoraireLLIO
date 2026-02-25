<script lang="ts">
  import { X } from 'lucide-svelte';

  let { showModal, isSaving, errorMsg, config, onClose, onSave, onUpdateConfig } = $props();

  const updateStartDate = (event: Event) => {
    const input = event.currentTarget as HTMLInputElement;
    onUpdateConfig({ startDate: input.value });
  };

  const updateHoursWorked = (event: Event) => {
    const input = event.currentTarget as HTMLInputElement;
    onUpdateConfig({ hoursWorked: Number(input.value) });
  };

  const updateOffset = (event: Event) => {
    const input = event.currentTarget as HTMLInputElement;
    onUpdateConfig({ offset: Number(input.value) });
  };

  const handleOverlayKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ' || event.key === 'Escape') {
      onClose();
    }
  };

  const handleModalKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape') {
      onClose();
    }
  };
</script>

{#if showModal}
  <div
    class="modal-overlay"
    role="button"
    tabindex="0"
    aria-label="Fermer la fen�tre"
    onclick={onClose}
    onkeydown={handleOverlayKeydown}
  >
    <div
      class="modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="hours-modal-title"
      tabindex="-1"
      onclick={(e) => e.stopPropagation()}
      onkeydown={handleModalKeydown}
    >
      <div class="modal-header">
        <h3 id="hours-modal-title">Configuration des heures en banque</h3>
        <button type="button" class="close-btn" aria-label="Fermer" onclick={onClose}>
          <X size={18} />
        </button>
      </div>

      <div class="modal-body">
        {#if errorMsg}
          <div class="error">{errorMsg}</div>
        {/if}

        <div class="form-group">
          <label for="startDate">Début de la période</label>
          <input
            id="startDate"
            type="date"
            value={config.startDate}
            oninput={updateStartDate}
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label for="hoursWorked">Nombre d'heures par semaine</label>
          <input
            id="hoursWorked"
            type="number"
            value={config.hoursWorked}
            oninput={updateHoursWorked}
            min="0"
            step="0.5"
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label for="offset">Décalage (offset)</label>
          <input
            id="offset"
            type="number"
            value={config.offset}
            oninput={updateOffset}
            step="0.5"
            class="form-input"
          />
        </div>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn-cancel" onclick={onClose} disabled={isSaving}>
          Annuler
        </button>

        <button type="button" class="btn-save" onclick={onSave} disabled={isSaving}>
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
