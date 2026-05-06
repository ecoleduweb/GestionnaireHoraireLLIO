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
    startDate: '',
    hoursPerWeek: 0,
    offset: 0,
  });

  
  if (initialConfig) {
    Object.assign(config, initialConfig);
  }

  let isSubmitting = $state(false);
  let isLoading = $state(true);

  onMount(async () => {
    try {
      const data = await UserApiService.getTimeBankConfig();

      if (data) {
        Object.assign(config, data); 
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
      alert("Erreur lors de la configuration");
    } finally {
      isSubmitting = false;
    }
  };

  const { form, errors } = validateTimeBankForm(handleSubmit, config);
</script>

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
        
        <form use:form>

          <div class="form-group">
            <label for="startDate-input" >Début de la période</label>
            <input type="date" name="startDate" id="startDate-input" bind:value={config.startDate} />
            {#if $errors.startDate}
              <span class="error-text">{$errors.startDate}</span>
            {/if}
          </div>

          <div class="form-group">
            <label for="hoursPerWeek-input" >Heures par semaine</label>
            <input type="number" name="hoursPerWeek" id="hoursPerWeek-input" bind:value={config.hoursPerWeek} />
            {#if $errors.hoursPerWeek}
              <span class="error-text">{$errors.hoursPerWeek}</span>
            {/if}
          </div>

          <div class="form-group">
            <label for="offset-input" >Heure en banque</label>
            <input type="number" name="offset" id="offset-input" bind:value={config.offset} />
            {#if $errors.offset}
              <span class="error-text">{$errors.offset}</span>
            {/if}
          </div>

        
          <div class="modal-footer">
            <button
              type="button"
              class="btn-secondary"
              onclick={onClose}
            >
              Annuler
            </button>

            <button
              type="submit"
              class="btn-primary"
              disabled={isSubmitting}
            >
              {isSubmitting ? "En cours..." : "Enregistrer"}
            </button>
          </div>

        </form>
      {/if}
    </div>
  </div>
</div>

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

  .error-text {
    color: red;
    font-size: 0.8rem;
  }

  .close-btn {
    background: none;
    border: none;
    cursor: pointer;
  }

 
  .btn-secondary {
    padding: 0.6rem 1.4rem;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    border: 1px solid #e5e7eb;
    background: #f3f4f6;
    color: #374151;
    transition: all 0.2s ease;
  }

  .btn-secondary:hover {
    background: #e5e7eb;
    transform: translateY(-1px);
  }

  .btn-primary {
    padding: 0.6rem 1.4rem;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    border: none;
    background: #015e61;
    color: white;
    transition: all 0.2s ease;
  }

  .btn-primary:hover {
    background: #014446;
    transform: translateY(-1px);
  }

  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>