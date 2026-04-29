<script lang="ts">
  import type { Activity, User, Project, Category } from '../../Models';
  import '../../style/app.css';
  import { ChevronDown, X, Plus, Trash2 } from 'lucide-svelte';
  import ActivityEntryForm from "./ActivityEntryForm.svelte";
  import ConfirmationModal from "../Modal/ConfirmationModal.svelte";

  type Props = {
    show: boolean;
    projects: Project[];
    activityToEdit: Activity | null;
    selectedDate?: { start: Date; end: Date } | null;
    onClose: () => void;
    onDelete: (activity: Activity) => void;
    onSubmit: (activity: Activity) => void;
    onUpdate: (activity: Activity) => void;
  };

  let {
    show,
    projects,
    activityToEdit,
    selectedDate = null,
    onClose,
    onDelete,
    onSubmit,
    onUpdate,
  }: Props = $props();

  const editMode = activityToEdit !== null;

  let showCloseConfirmModal = $state(false);
  let isActivityDirty = $state(false);

  // Le e.stopPropagation empêche le clic de traverser la modale,
  // car sinon le cliquer sur le fond gris pourrait déclencher une action sur le calendrier
  // qui se trouve derrière comme ajouter une activité.
  // On a ajouté un if avec un (e) pour empêché que la méthode soit appeler autre part
  // sans raison du coup un événement mouse est obligatoire
  const handlePreventClosingIfDirty = (e: MouseEvent) => {
    if (e) e.stopPropagation();

    if (isActivityDirty) {
      showCloseConfirmModal = true;
    } else {
      onClose();
    }
  };

  const confirmClose = () => {
    showCloseConfirmModal = false;
    onClose();
  };
</script>

{#if show}
  <div class="fixed inset-0 z-40 flex justify-start">
    <div
            class="absolute inset-0 bg-gray-950/40 transition-opacity"
            onclick={handlePreventClosingIfDirty}
    ></div>

    <div
            class="w-full max-w-[300px] bg-white h-full overflow-y-auto relative flex flex-col z-50 animate-slideIn border-r border-gray-300 shadow-xl"
    >
      <!-- Header -->
      <div class="flex items-center justify-between bg-[#015e61] text-white px-6 py-4">
        <h2 class="text-xl font-medium">
          {editMode ? "Modifier l'activité" : 'Nouvelle activité'}
        </h2>
        <button type="button" class="text-white hover:text-gray-200" onclick={handlePreventClosingIfDirty}>
          <X />
        </button>
      </div>

      <!-- Form -->
      <div class="p-6 flex-grow">
        <ActivityEntryForm
                {projects}
                {activityToEdit}
                {selectedDate}
                bind:isDirty={isActivityDirty}
                {onClose}
                {onDelete}
                {onSubmit}
                {onUpdate}
        />
      </div>
    </div>
  </div>
{/if}

{#if showCloseConfirmModal}
  <ConfirmationModal
          modalTitle="Modifications non enregistrées"
          modalText="Vous avez des modifications non enregistrées. Voulez-vous vraiment quitter sans sauvegarder ?"
          confirmText="Oui, quitter"
          cancelText="Non, rester"
          errorText=""
          onSuccess={confirmClose}
          onClose={() => {
          showCloseConfirmModal = false;
        }}
  />
{/if}

<style>
  @keyframes slideIn {
    from { transform: translateX(-100%); }
    to { transform: translateX(0); }
  }

  .animate-slideIn {
    animation: slideIn 0.3s forwards;
  }
</style>
