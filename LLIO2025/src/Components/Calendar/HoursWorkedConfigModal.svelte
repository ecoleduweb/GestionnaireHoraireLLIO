<script lang="ts">
  import { onMount } from 'svelte';
  import { X } from 'lucide-svelte';
  import { UserApiService } from '../../services/UserApiService';
  import type { TimeBankConfig } from '../../Models/index';
  import { validateTimeBankForm } from '../../Validation/TimeBank';

 type Props = {
  onClose: () => void;
  onSave: (values: TimeBankConfig) => void;
  initialConfig: TimeBankConfig;
};

let { onClose, onSave, initialConfig }: Props = $props();

  const config = $state<TimeBankConfig>({
    startDate: initialConfig?.startDate ?? '',
    hoursPerWeek: initialConfig?.hoursPerWeek ?? 0,
    offset: initialConfig?.offset ?? 0,
  });

  let isSubmitting = $state(false);
  let isLoading = $state(true);

  onMount(async () => {
    try {
      const data = await UserApiService.getTimeBankConfig();

      if (data) {
        config.startDate = data.startDate;
        config.hoursPerWeek = data.hoursPerWeek;
        config.offset = data.offset;
      }
    } catch (err) {
      console.error(err);
    } finally {
      isLoading = false;
    }
  });

  const handleSubmit = async (values: TimeBankConfig) => {
    if (isSubmitting) return;

    try {
      isSubmitting = true;

      await UserApiService.saveTimeBankConfig(values);

      onSave(values); 
      onClose();
    } catch (err) {
      console.error(err);
      alert("Erreur lors de la configuration de la banque d'heure");
    } finally {
      isSubmitting = false;
    }
  };

  const { form, errors } = validateTimeBankForm(handleSubmit, config);
</script>

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .modal {
    background: white;
    border-radius: 8px;
    width: 400px;
    max-width: 90%;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid #eee;
  }

  .modal-body {
    padding: 1rem;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  input {
    width: 100%;
    padding: 0.6rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  .modal-footer button {
    padding: 0.6rem 1.4rem;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    border: none;
  }

  .modal-footer button[type='button'] {
    background: #f3f4f6;
  }

  .modal-footer button[type='submit'] {
    background: #015e61;
    color: white;
  }

  .error-text {
    color: red;
    font-size: 0.8rem;
  }

  .close-btn {
    background: none;
    border: none;
    cursor: pointer;
  }
</style>

<div class="modal-overlay">
  <div class="modal">
    <div class="modal-header">
      <h3>Configuration des heures en banque</h3>
      <button class="close-btn" onclick={onClose}>
        <X size={18} />
      </button>
    </div>

    <div class="modal-body">
      {#if isLoading}
        <p>Chargement...</p>
      {:else}
        <form use:form onsubmit={(e) => e.preventDefault()}>
          <div class="form-group">
            <label>Début de la période</label>
            <input type="date" name="startDate" bind:value={config.startDate} />
          </div>

          <div class="form-group">
            <label>Heures par semaine</label>
            <input
              type="number"
              name="hoursPerWeek"
              bind:value={config.hoursPerWeek}
            />
          </div>

          <div class="form-group">
            <label>Heure en banque</label>
            <input type="number" name="offset" bind:value={config.offset} />
          </div>

          <div class="modal-footer">
            <button type="button" onclick={onClose}>Annuler</button>

            <button type="submit" disabled={isSubmitting}>
              Enregistrer
            </button>
          </div>
        </form>
      {/if}
    </div>
  </div>
</div>