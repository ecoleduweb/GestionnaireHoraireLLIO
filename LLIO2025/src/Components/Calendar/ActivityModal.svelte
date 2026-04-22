<script lang="ts">
  import { validateActivityForm } from '../../Validation/Activity';
  import type { Activity, User, Project, Category } from '../../Models';
  import { activityTemplate } from '../../forms/activity';
  import { ActivityApiService } from '../../services/ActivityApiService';
  import { CategoryApiService } from '../../services/CategoryApiService';
  import {
    getHoursFromDate,
    getMinutesFromDate,
    createDateWithTime,
    initializeActivityDates,
    applyEndTime as applyEndTimeUtil,
  } from '../../utils/date';
  import '../../style/app.css';
  import { ChevronDown, X, Plus, Trash2 } from 'lucide-svelte';
  import SearchSelect from "../Global/SearchSelect.svelte";
  import ConfirmationModal from '../Modal/ConfirmationModal.svelte';
  import ActivityEntryForm from "./ActivityEntryForm.svelte";

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
</script>

{#if show}
  <div class="fixed inset-0 z-40 flex justify-start">
    <div
            class="absolute inset-0 bg-gray-950/40 transition-opacity"
            onclick={onClose}
    ></div>

    <div
            class="w-full max-w-[300px] bg-white h-full overflow-y-auto relative flex flex-col z-50 animate-slideIn border-r border-gray-300 shadow-xl"
    >
      <!-- Header -->
      <div class="flex items-center justify-between bg-[#015e61] text-white px-6 py-4">
        <h2 class="text-xl font-medium">
          {editMode ? "Modifier l'activité" : 'Nouvelle activité'}
        </h2>
        <button type="button" class="text-white hover:text-gray-200" onclick={onClose}>
          <X />
        </button>
      </div>

      <!-- Form -->
      <div class="p-6 flex-grow">
        <ActivityEntryForm
                {projects}
                {activityToEdit}
                {selectedDate}
                {onClose}
                {onDelete}
                {onSubmit}
                {onUpdate}
        />
      </div>
    </div>
  </div>
  <!-- à transferer dans le code de form
  {#if showCategoryConfirmModal}
    <ConfirmationModal
      modalTitle="Confirmer l'ajout"
      modalText={`Voulez-vous ajouter la catégorie "${categoryToAdd}" ?`}
      errorText="Erreur lors de la suppression du projet, il a soit une ou des activités liées à ce projet ou bien le projet est inexistant"
      onSuccess={() => {
        showCategoryConfirmModal = false;
        confirmAddCategory();
      }}
      onClose={() => {
        showCategoryConfirmModal = false;
      }}
    />
  {/if} -->
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
