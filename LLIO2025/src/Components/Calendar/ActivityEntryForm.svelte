<script lang="ts">
  import { validateActivityForm } from '../../Validation/Activity';
  import type { Activity, Project, Category } from '../../Models';
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
  import { ChevronDown, Plus, Trash2 } from 'lucide-svelte';
  import SearchSelect from '../Global/SearchSelect.svelte';
  import ConfirmationModal from '../Modal/ConfirmationModal.svelte';

  type Props = {
    projects: Project[];
    activityToEdit?: Activity | null;
    activityToImport?: Activity | null;
    selectedDate?: { start: Date; end: Date } | null;
    isDirty?: boolean;
    onClose: () => void;
    onDelete: (activity: Activity) => void;
    onSubmit: (activity: Activity) => void;
    onUpdate: (activity: Activity) => void;
  };

  let {
    projects,
    activityToEdit,
    activityToImport = $bindable(null),
    selectedDate = null,
    isDirty = $bindable(false),
    onClose,
    onDelete,
    onSubmit,
    onUpdate,
  }: Props = $props();
  let editMode = $derived(!!activityToEdit);

  let projectCategories: Category[] = $state([]);
  let initialActivity = activityTemplate.generate();
  initialActivity.projectId = '' as unknown as number;
  initialActivity.categoryId = null;

  let isLoading = $state(false);
  let isSubmitting = $state(false);
  let showCategoryConfirmModal = $state(false);
  let categoryToAdd = $state('');
  let isProjectSelectFocused = $state(true);

  if (selectedDate && selectedDate.start) {
    const { startDate, endDate } = initializeActivityDates(selectedDate.start);
    initialActivity.startDate = startDate;
    initialActivity.endDate = endDate ? new Date(selectedDate.end) : endDate;
  }

  const activity = $state<Activity>(initialActivity);

  const time = $state({
    startHours: getHoursFromDate(activity.startDate),
    startMinutes: getMinutesFromDate(activity.startDate),
    endHours: getHoursFromDate(activity.endDate),
    endMinutes: getMinutesFromDate(activity.endDate),
  });

  let initialSnapshot = $state({
    projectId: activity.projectId,
    categoryId: activity.categoryId,
    name: activity.name || '',
    description: activity.description || '',
    sh: time.startHours,
    sm: time.startMinutes,
    eh: time.endHours,
    em: time.endMinutes,
  });


  const _isDirty = $derived(
          activity.projectId !== initialSnapshot.projectId ||
          activity.categoryId !== initialSnapshot.categoryId ||
          (activity.name || '') !== initialSnapshot.name ||
          (activity.description || '') !== initialSnapshot.description ||
          time.startHours !== initialSnapshot.sh ||
          time.startMinutes !== initialSnapshot.sm ||
          time.endHours !== initialSnapshot.eh ||
          time.endMinutes !== initialSnapshot.em
  );

  $effect(() => { isDirty = _isDirty; });

  const {
    time: { hours, minutes },
  } = activityTemplate;

  const applyEndTime = () => {
    const result = applyEndTimeUtil(
            time.startHours,
            time.startMinutes,
            time.endHours,
            time.endMinutes
    );
    time.endHours = result.endHours;
    time.endMinutes = result.endMinutes;
  };

  let categoryDropdownOpen = $state(false);
  let searchTerm = $state('');
  let isSearchFocused = $state(false);

  const filteredCategories = $derived(
          searchTerm
                  ? projectCategories.filter((category) =>
                          category.name.toLowerCase().includes(searchTerm.toLowerCase())
                  )
                  : projectCategories
  );

  const loadCategoriesByProject = async (projectId) => {
    if (!projectId) {
      projectCategories = [];
      return;
    }
    try {
      isLoading = true;
      projectCategories = await CategoryApiService.getCategoriesByProject(projectId);

      if (activity.categoryId) {
        const categoryExists = projectCategories.some((c) => c.id === activity.categoryId);

        if (!categoryExists && editMode) {
          activity.categoryId = projectCategories[0].id;
        } else if (!categoryExists) {
          activity.categoryId = null;
        }
      }
    } catch (error) {
      console.error(`Erreur lors du chargement des catégories pour le projet ${projectId}:`, error);
      projectCategories = [];
    } finally {
      isLoading = false;
    }
  };

  $effect(() => {
    if (activity.projectId) {
      loadCategoriesByProject(activity.projectId);
    } else {
      projectCategories = [];
      activity.categoryId = null;
    }
  });

  const selectCategory = (categoryId) => {
    activity.categoryId = Number(categoryId);
    searchTerm = '';
    categoryDropdownOpen = false;
  };

  const handleOutsideClick = (event) => {
    if (categoryDropdownOpen && !event.target.closest('.category-dropdown-container')) {
      categoryDropdownOpen = false;
      searchTerm = '';
    }
  };

  const handleSubmit = async () => {
    if (isSubmitting) return;

    if (typeof activity.projectId === 'string') {
      activity.projectId = Number(activity.projectId);
    }
    if (typeof activity.categoryId === 'string') {
      activity.categoryId = Number(activity.categoryId);
    }

    isSubmitting = true;

    try {
      activity.startDate = createDateWithTime(activity.startDate, time.startHours, time.startMinutes);
      activity.endDate = createDateWithTime(activity.endDate, time.endHours, time.endMinutes);

      if (editMode) {
        const updatedActivity = await ActivityApiService.updateActivity(activity);
        onUpdate(updatedActivity);
      } else {
        const newActivity = await ActivityApiService.createActivity(activity);
        onSubmit(newActivity);
      }

      onClose();
    } catch (error) {
      console.error('Erreur lors de la soumission', error);
      alert("Activité incomplète \nVeuillez remplir tous les champs obligatoires (suivis d'une astérisque rouge) de l'activité avant de l'enregistrer.");
    } finally {
      isSubmitting = false;
    }
  };

  const handleDelete = async () => {
    if (!activity.id) return;

    if (confirm('Supprimer cette tâche ?')) {
      try {
        await ActivityApiService.deleteActivity(activity.id);
        onDelete(activity);
        onClose();
      } catch (error) {
        console.error('Erreur lors de la suppression', error);
      }
    }
  };

  const deleteCategoryFromList = (category: Category) => {
    projectCategories = projectCategories.filter((cat) => cat !== category);
  };

  const handleDeleteCategory = async (category: Category) => {
    if (!category) return;

    if (confirm('Supprimer la catégorie ' + category.name + ' ?')) {
      try {
        await CategoryApiService.deleteCategory(category.id);
        deleteCategoryFromList(category);
      } catch (error) {
        alert('Erreur lors de la suppression de la catégorie');
        console.error('Erreur lors de la suppression de la catégorie', error);
      }
    }
  };

  const handleAddCategory = (e) => {
    e.stopPropagation();
    categoryToAdd = searchTerm;
    showCategoryConfirmModal = true;
  };

  const confirmAddCategory = async () => {
    if (!categoryToAdd.trim()) return;

    try {
      const newCategory = await CategoryApiService.createCategory(
              { name: categoryToAdd.trim(), description: '', billable: false },
              activity.projectId
      );

      projectCategories = [...projectCategories, newCategory];
      activity.categoryId = Number(newCategory.id);
      searchTerm = '';
      categoryDropdownOpen = false;
    } catch (error) {
      console.error("Erreur lors de l'ajout d'une catégorie", error);
      alert("Erreur lors de l'ajout de la catégorie");
    }
  };

  const getDisplayText = (uniqueId, name) => {
    const separator = ' | ';
    if (name === undefined || name === null || name.trim() === '') {
      return uniqueId;
    }
    return `${uniqueId}${separator}${name}`;
  };

  $effect(() => {
    if (activity.projectId && projectCategories.length >= 1) {
      const currentCategoryExists = projectCategories.find((c) => c.id === activity.categoryId);
      if (!activity.categoryId || !currentCategoryExists) {
        activity.categoryId = projectCategories[0].id;
      }
    }
  });

  const { form, errors, setFields } = validateActivityForm(handleSubmit, activity);

  function applyActivity(source: Activity | null | undefined) {
    const fresh = activityTemplate.generate();
    fresh.projectId = '' as unknown as number;
    fresh.categoryId = null;

    if (selectedDate?.start) {
      const { startDate, endDate } = initializeActivityDates(selectedDate.start);
      fresh.startDate = startDate;
      fresh.endDate = endDate ? new Date(selectedDate.end) : endDate;
    }

    const merged = Object.assign({}, fresh, source ?? {});

    Object.assign(activity, merged);

    setFields('name', merged.name ?? '');
    setFields('description', merged.description ?? '');
    setFields('projectId', merged.projectId);
    setFields('categoryId', merged.categoryId);

    time.startHours = getHoursFromDate(merged.startDate);
    time.startMinutes = getMinutesFromDate(merged.startDate);
    time.endHours = getHoursFromDate(merged.endDate);
    time.endMinutes = getMinutesFromDate(merged.endDate);

    initialSnapshot.projectId = merged.projectId;
    initialSnapshot.categoryId = merged.categoryId;
    initialSnapshot.name = merged.name || '';
    initialSnapshot.description = merged.description || '';
    initialSnapshot.sh = time.startHours;
    initialSnapshot.sm = time.startMinutes;
    initialSnapshot.eh = time.endHours;
    initialSnapshot.em = time.endMinutes;
  }

  let lastEditId = $state<number | null>(null);
  let lastImportId = $state<number | null>(null);
  let wasReset = $state(false);

  $effect(() => {
    if (activityToEdit) {
      if (activityToEdit.id !== lastEditId) {
        lastEditId = activityToEdit.id;
        wasReset = false;
        applyActivity(activityToEdit);
      }
      return;
    }

    if (activityToImport) {
      if (activityToImport.id !== lastImportId) {
        lastImportId = activityToImport.id;
        wasReset = false;
        applyActivity(activityToImport);
      }
      return;
    }

    if (!wasReset) {
      wasReset = true;
      lastEditId = null;
      lastImportId = null;
      applyActivity(null);
    }
  });
</script>

<svelte:window onclick={handleOutsideClick} />
<form
        class="flex flex-col"
        use:form
        onsubmit={(e) => e.preventDefault()}
>
  <!-- Action buttons -->
  <div class="mt-auto">
    <div class="flex justify-center gap-5">
      {#if editMode}
        <button
                type="button"
                class="py-3 px-6 bg-red-500 text-white rounded-lg font-medium hover:bg-red-600 hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition"
                onclick={handleDelete}
        >
          Supprimer
        </button>
        <button
                type="submit"
                class="py-3 px-6 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition disabled:opacity-50"
                disabled={isSubmitting || isLoading}
        >
          {isSubmitting ? 'En cours...' : 'Modifier'}
        </button>
      {:else}
        <button
                type="button"
                class="py-3 px-6 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 hover:-translate-y-0.5 hover:shadow-sm active:translate-y-0 transition border border-gray-200"
                onclick={onClose}
        >
          { activityToImport ? 'Ne pas importer' : 'Annuler' }
        </button>
        <button
                type="submit"
                class="py-3 px-6 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition disabled:opacity-50"
                disabled={isSubmitting || isLoading}
        >
          {isSubmitting ? 'En cours...' : (activityToImport ? 'Importer' : 'Créer')}
        </button>
      {/if}
    </div>

    <div class="border-t border-gray-200 my-6"></div>
  </div>

  <!-- Form fields -->
  <div class="space-y-6">
    <!-- Projet -->
    <div>
      <label for="activity-project" class="block text-gray-700 font-medium mb-2">
        Projet <span class="text-red-500">*</span>
      </label>
      <SearchSelect
              id="activity-project"
              name="projectId"
              items={projects.map((value) => ({
          value: value.id,
          label: getDisplayText(value.uniqueId, value.name),
        }))}
              bind:selectedValue={activity.projectId}
              placeholder="Sélectionner un projet"
              setFields={setFields}
              bind:focused={isProjectSelectFocused}
              onSubmit={() => document.querySelector('form')?.requestSubmit()}
              required
      />
      {#if $errors.projectId}
        <span class="text-red-500 text-sm">{$errors.projectId}</span>
      {/if}
    </div>

    <!-- Heure de début -->
    <div>
      <label class="block text-gray-700 font-medium mb-2">
        Heure de début <span class="text-red-500">*</span>
      </label>
      <div class="flex gap-3">
        <div class="select-container form-select-flex">
          <select bind:value={time.startHours} class="form-select w-full">
            {#each hours as hour}
              <option value={hour}>{hour} h</option>
            {/each}
          </select>
          <div class="select-icon"><ChevronDown size={18} /></div>
        </div>
        <div class="select-container form-select-flex">
          <select bind:value={time.startMinutes} class="form-select w-full">
            {#each minutes as minute}
              <option value={minute}>{minute} min</option>
            {/each}
          </select>
          <div class="select-icon"><ChevronDown size={18} /></div>
        </div>
      </div>
    </div>

    <!-- Heure de fin -->
    <div>
      <label class="block text-gray-700 font-medium mb-2">
        Heure de fin <span class="text-red-500">*</span>
      </label>
      <div class="flex gap-3">
        <div class="select-container form-select-flex">
          <select bind:value={time.endHours} onchange={applyEndTime} class="form-select w-full">
            {#each hours as hour}
              <option value={hour}>{hour} h</option>
            {/each}
          </select>
          <div class="select-icon"><ChevronDown size={18} /></div>
        </div>
        <div class="select-container form-select-flex">
          <select bind:value={time.endMinutes} onchange={applyEndTime} class="form-select w-full">
            {#each minutes as minute}
              <option value={minute}>{minute} min</option>
            {/each}
          </select>
          <div class="select-icon"><ChevronDown size={18} /></div>
        </div>
      </div>
    </div>

    <!-- Catégorie -->
    <div>
      <label for="activity-category-search" class="block text-gray-700 font-medium mb-2">
        Catégorie <span class="text-red-500">*</span>
      </label>
      <div class="category-dropdown-container">
        <div
                class="custom-select"
                class:disabled-select={!activity.projectId}
                onclick={(e) => {
            if (!activity.projectId) { e.preventDefault(); return; }
            categoryDropdownOpen = !categoryDropdownOpen;
            if (categoryDropdownOpen) {
              document.getElementById('activity-category-search')?.focus();
            }
          }}
        >
          {#if !categoryDropdownOpen || !isSearchFocused}
            <span class="select-value">
              {!activity.projectId
                      ? 'Sélectionner une catégorie'
                      : projectCategories.find((c) => c.id === activity.categoryId)?.name ||
                      'Sélectionner une catégorie'}
            </span>
          {/if}

          <input
                  id="activity-category-search"
                  type="text"
                  class="search-input"
                  placeholder="Rechercher une catégorie..."
                  bind:value={searchTerm}
                  disabled={!activity.projectId}
                  style="opacity: {categoryDropdownOpen ? '1' : '0'}"
                  onfocus={() => { if (activity.projectId) { categoryDropdownOpen = true; isSearchFocused = true; } }}
                  onblur={() => { isSearchFocused = false; }}
                  onclick={(e) => e.stopPropagation()}
          />

          <div class="arrow-icon">
            <ChevronDown size={18} class="text-gray-600" />
          </div>
        </div>

        {#if categoryDropdownOpen && activity.projectId}
          <div class="dropdown-content">
            <div class="category-list">
              {#if filteredCategories.length === 0}
                <div class="no-results">
                  {searchTerm ? 'Aucun résultat trouvé' : 'Aucune catégorie disponible pour ce projet'}
                </div>
              {:else}
                {#each filteredCategories as category}
                  <div
                          class="category-item"
                          class:selected={category.id === activity.categoryId}
                          onclick={(e) => { e.stopPropagation(); selectCategory(category.id); }}
                  >
                    {category.name}
                    {#if category.name !== 'Par défaut'}
                      <button
                              class="justify-end p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
                              onclick={() => handleDeleteCategory(category)}
                              aria-label="Supprimer la catégorie"
                      >
                        <Trash2 size={16} />
                      </button>
                    {/if}
                  </div>
                {/each}
              {/if}
            </div>

            {#if activity.projectId && searchTerm && !filteredCategories.some((cat) => cat.name.toLowerCase() === searchTerm.toLowerCase())}
              <div class="add-category" onclick={handleAddCategory}>
                <Plus size={16} />
                Ajouter "{searchTerm}"
              </div>
            {/if}
          </div>
        {/if}

        {#if $errors.categoryId}
          <span class="text-red-500 text-sm mt-1">{$errors.categoryId}</span>
        {/if}
      </div>
    </div>

    <!-- Nom -->
    <div>
      <label for="activity-name" class="block text-gray-700 font-medium mb-2">
        Nom <span class="text-gray-400">(optionnel)</span>
      </label>
      <input
              id="activity-name"
              name="name"
              type="text"
              bind:value={activity.name}
              placeholder="Nom de l'activité..."
              class="form-input"
      />
      {#if $errors.name}
        <span class="text-red-500 text-sm">{$errors.name}</span>
      {/if}
    </div>

    <!-- Description -->
    <div>
      <label for="activity-description" class="block text-gray-700 font-medium mb-2">
        Description <span class="text-gray-400">(optionnel)</span>
      </label>
      <textarea
              id="activity-description"
              name="description"
              bind:value={activity.description}
              placeholder="Description de l'activité..."
              rows="3"
              class="form-input"
      ></textarea>
      {#if $errors.description}
        <span class="text-red-500 text-sm">{$errors.description}</span>
      {/if}
    </div>
  </div>
</form>
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
{/if}

<style>
  .form-select {
    appearance: none;
    -webkit-appearance: none;
    -moz-appearance: none;
    background-color: white;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    padding: 0.75rem 1rem;
    padding-right: 2.5rem;
    color: #4b5563;
    transition: all 0.2s;
  }

  .form-select:focus {
    outline: none;
    border-color: #015e61;
    box-shadow: 0 0 0 3px rgba(1, 94, 97, 0.2);
  }

  .form-select-flex { flex: 1; }
  .form-select::-ms-expand { display: none; }

  .form-input {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    transition: all 0.2s;
    background-color: white;
    color: #4b5563;
  }

  .form-input:focus {
    outline: none;
    border-color: #015e61;
    box-shadow: 0 0 0 3px rgba(1, 94, 97, 0.2);
  }

  .select-container { position: relative; width: 100%; }

  .select-icon {
    position: absolute;
    right: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    pointer-events: none;
    color: #606060;
  }

  .category-dropdown-container { position: relative; width: 100%; }

  .custom-select {
    width: 100%;
    height: 42px;
    background-color: white;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    color: #4b5563;
    position: relative;
    cursor: pointer;
  }

  .custom-select:focus,
  .custom-select:focus-within {
    outline: none;
    border-color: #015e61;
    box-shadow: 0 0 0 3px rgba(1, 94, 97, 0.2);
  }

  .search-input {
    width: 100%;
    height: 100%;
    padding: 0 1rem;
    padding-right: 2.5rem;
    border: none;
    background: transparent;
    font-size: inherit;
    color: inherit;
    outline: none;
    cursor: text;
  }

  .search-input::placeholder { color: #9ca3af; }

  .select-value {
    position: absolute;
    top: 0; left: 0;
    width: 100%; height: 100%;
    padding: 0 1rem;
    padding-right: 2.5rem;
    display: flex;
    align-items: center;
    pointer-events: none;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .arrow-icon {
    position: absolute;
    right: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    pointer-events: none;
    cursor: pointer;
  }

  .dropdown-content {
    position: absolute;
    width: 100%;
    max-height: 300px;
    background-color: white;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    margin-top: 0.25rem;
    z-index: 50;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  .category-list { max-height: 200px; overflow-y: auto; }

  .category-item {
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: background-color 0.15s;
  }

  .category-item:hover { background-color: #f3f4f6; }
  .category-item.selected { background-color: #e5e7eb; }

  .add-category {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    border-top: 1px solid #e5e7eb;
    color: #015e61;
    cursor: pointer;
    font-weight: 500;
    background-color: #f9fafb;
  }

  .add-category:hover { background-color: #f3f4f6; }
  .no-results { padding: 0.75rem 1rem; color: #6b7280; text-align: center; }

  .disabled-select {
    background-color: #f3f4f6;
    cursor: not-allowed;
    opacity: 0.7;
  }
</style>
