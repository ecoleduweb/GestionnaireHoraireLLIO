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
  import { ChevronDown, X, Plus } from 'lucide-svelte';
  import ConfirmationCreateCategory from './ConfirmModal.svelte';

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

  let categories: Category[] = []; // Toutes les catégories (utilisé en mode édition)
  let projectCategories: Category[] = $state([]); // Catégories spécifiques au projet
  let initialActivity = activityTemplate.generate();
  initialActivity.projectId = '' as unknown as number;
  initialActivity.categoryId = null;

  let isSubmitting = $state(false);
  let showCategoryConfirmModal = $state(false);
  let categoryToAdd = $state('');

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

  if (activityToEdit) {
    Object.assign(activity, activityToEdit);
    time.startHours = getHoursFromDate(activityToEdit.startDate);
    time.startMinutes = getMinutesFromDate(activityToEdit.startDate);
    time.endHours = getHoursFromDate(activityToEdit.endDate);
    time.endMinutes = getMinutesFromDate(activityToEdit.endDate);
  }

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

  // État pour le dropdown de catégories
  let categoryDropdownOpen = $state(false);
  let searchTerm = $state('');
  let isSearchFocused = $state(false);

  // Filtrer les catégories selon le terme de recherche
  const filteredCategories = $derived(
    searchTerm
      ? projectCategories.filter((category) =>
          category.name.toLowerCase().includes(searchTerm.toLowerCase())
        )
      : projectCategories
  );

  // Effet pour charger les catégories lors d'un changement de projet
  $effect(() => {
    if (activity.projectId) {
      loadCategoriesByProject(activity.projectId);
    } else {
      projectCategories = [];
      activity.categoryId = null;
    }
  });

  // Fonction pour charger les catégories spécifiques à un projet
  const loadCategoriesByProject = async (projectId) => {
    if (!projectId) {
      projectCategories = [];
      return;
    }
    try {
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
    }
  };

  // Fonction pour charger toutes les catégories
  const loadCategories = async () => {
    try {
      // En mode édition, charger les catégories du projet directement
      if (editMode && activityToEdit && activityToEdit.projectId) {
        await loadCategoriesByProject(activityToEdit.projectId);
      }
    } catch (error) {
      console.error('Erreur lors du chargement des catégories:', error);
    }
  };

  // Rafraîchir les catégories quand le modal s'ouvre
  $effect(() => {
    if (show) {
      if (activity.projectId) {
        loadCategoriesByProject(activity.projectId);
      } else {
        loadCategories();
      }
    }
  });

  // Fonction pour sélectionner une catégorie
  const selectCategory = (categoryId) => {
    // Convertir explicitement en nombre
    activity.categoryId = Number(categoryId);
    searchTerm = '';
    categoryDropdownOpen = false;
  };

  // Gérer la fermeture du dropdown si on clique ailleurs
  const handleOutsideClick = (event) => {
    if (categoryDropdownOpen && !event.target.closest('.category-dropdown-container')) {
      categoryDropdownOpen = false;
      searchTerm = ''; // Réinitialiser la recherche
    }
  };

  const handleSubmit = async () => {
    if (isSubmitting) return; // Empêche les soumissions multiples

    // S'assurer que les valeurs sont des nombres
    if (typeof activity.projectId === 'string') {
      activity.projectId = Number(activity.projectId);
    }

    if (typeof activity.categoryId === 'string') {
      activity.categoryId = Number(activity.categoryId);
    }

    isSubmitting = true;

    try {
      // Créer une nouvelle date de début avec les heures et minutes sélectionnées
      const updatedStartDate = createDateWithTime(
        activity.startDate,
        time.startHours,
        time.startMinutes
      );

      // Créer une nouvelle date de fin avec les heures et minutes sélectionnées
      const updatedEndDate = createDateWithTime(activity.endDate, time.endHours, time.endMinutes);

      // Mise à jour des dates de l'activité
      activity.startDate = updatedStartDate;
      activity.endDate = updatedEndDate;

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
      alert('Activité incomplète \nVeuillez remplir tous les champs obligatoires (suivis d\'une astérisque rouge) de l\'activité avant de l\'enregistrer.');
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
      } catch (error) {
        console.error('Erreur lors de la suppression', error);
      }
    }
  };

  const handleClose = () => {
    onClose();
  };

  const handleAddCategory = (e) => {
    e.stopPropagation();
    categoryToAdd = searchTerm;
    showCategoryConfirmModal = true;
  };

  // Fonction pour ajouter une nouvelle catégorie après confirmation
  const confirmAddCategory = async () => {
    if (!categoryToAdd.trim()) return;

    try {
      const selectedProjectId = activity.projectId;

      const newCategory = await CategoryApiService.createCategory(
        {
          name: categoryToAdd.trim(),
          description: '',
          billable: false,
        },
        selectedProjectId
      );

      // Ajouter la nouvelle catégorie à la liste des catégories du projet
      projectCategories = [...projectCategories, newCategory];

      // Sélectionner la nouvelle catégorie
      activity.categoryId = Number(newCategory.id);

      searchTerm = '';
      categoryDropdownOpen = false;
    } catch (error) {
      console.error("Erreur lors de l'ajout d'une catégorie", error);
      alert("Erreur lors de l'ajout de la catégorie");
    }
  };

  const getTruncatedDisplayText = (uniqueId, name, maxLength = 30) => {
    const separator = " | ";
    const availableForName = maxLength - uniqueId.length - separator.length;
    if (name === undefined || name === null || name.trim() === "") {
      return uniqueId; // Si le nom est vide, retourner uniquement l'uniqueId
    }
    
    if (availableForName <= 0) {
      // Si l'uniqueId est déjà trop long, on le tronque aussi
      return uniqueId.substring(0, maxLength - 3) + "...";
    }
    
    if (name.length <= availableForName) {
      return `${uniqueId}${separator}${name}`;
    }
    
    const truncatedName = name.substring(0, availableForName - 3) + "...";
    return `${uniqueId}${separator}${truncatedName}`;
  };

  $effect(() => {
  if (activity.projectId && projectCategories.length >= 1) {
    const currentCategoryExists = projectCategories.find(c => c.id === activity.categoryId);
    
    if (!activity.categoryId || !currentCategoryExists) {
      activity.categoryId = projectCategories[0].id;
    }
  }
  });

  const { form, errors } = validateActivityForm(handleSubmit, activity);
</script>

<svelte:window onclick={handleOutsideClick} />

{#if show}
  <!-- Structure principale avec Tailwind -->
  <div class="fixed inset-0 z-40 flex justify-start">
    <!-- Overlay semi-transparent avec opacité à 40% comme dans l'original -->
    <div
      class="absolute inset-0 bg-gray bg-opacity-40 transition-opacity"
      onclick={handleClose}
    ></div>

    <!-- Panneau latéral avec bordure et ombre à gauche pour délimiter -->
    <div
      class="w-full max-w-[300px] bg-white h-full overflow-y-auto relative flex flex-col z-50 animate-slideIn border-r border-gray-300 shadow-xl"
    >
      <!-- En-tête avec titre et bouton de fermeture -->
      <div class="flex items-center justify-between bg-[#015e61] text-white px-6 py-4">
        <h2 class="text-xl font-medium">
          {editMode ? "Modifier l'activité" : 'Nouvelle activité'}
        </h2>
        <button type="button" class="text-white hover:text-gray-200" onclick={handleClose}>
          <X />
        </button>
      </div>

      <!-- Contenu du formulaire - espace vertical ajusté -->
      <div class="p-6 flex-grow">
        <form
          class="flex flex-col h-full"
          use:form
          onsubmit={(e) => {
            e.preventDefault();
          }}
        >
          <!-- Champs de formulaire avec espacement vertical uniforme -->
          <div class="space-y-6">
            <!-- Projet -->
            <div>
              <label for="activity-project" class="block text-gray-700 font-medium mb-2">
                Projet
                <span class="text-red-500">*</span>
              </label>
              <div class="select-container">
                <select
                  id="activity-project"
                  name="projectId"
                  bind:value={activity.projectId}
                  required
                  class="form-select w-full"
                >
                  <option value="" disabled selected hidden>Sélectionner un projet...</option>
                  {#each projects as project}
                    <option value={project.id} title={project.name}>{getTruncatedDisplayText(project.uniqueId, project.name)}</option>
                  {/each}
                </select>
                <div class="select-icon">
                  <ChevronDown size={18} />
                </div>
                {#if $errors.projectId}
                  <span class="text-red-500 text-sm">{$errors.projectId}</span>
                {/if}
              </div>
            </div>

            <!-- Sélecteurs d'heure côte à côte -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">
                Heure de début
                <span class="text-red-500">*</span>
              </label>
              <div class="flex gap-3">
                <div class="select-container form-select-flex">
                  <select bind:value={time.startHours} class="form-select w-full">
                    {#each hours as hour}
                      <option value={hour}>{hour} h</option>
                    {/each}
                  </select>
                  <div class="select-icon">
                    <ChevronDown size={18} />
                  </div>
                </div>
                <div class="select-container form-select-flex">
                  <select bind:value={time.startMinutes} class="form-select w-full">
                    {#each minutes as minute}
                      <option value={minute}>{minute} min</option>
                    {/each}
                  </select>
                  <div class="select-icon">
                    <ChevronDown size={18} />
                  </div>
                </div>
              </div>
            </div>

            <!-- Sélecteurs d'heure de fin -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">
                Heure de fin
                <span class="text-red-500">*</span>
              </label>
              <div class="flex gap-3">
                <div class="select-container form-select-flex">
                  <select
                    bind:value={time.endHours}
                    onchange={applyEndTime}
                    class="form-select w-full"
                  >
                    {#each hours as hour}
                      <option value={hour}>{hour} h</option>
                    {/each}
                  </select>
                  <div class="select-icon">
                    <ChevronDown size={18} />
                  </div>
                </div>
                <div class="select-container form-select-flex">
                  <select
                    bind:value={time.endMinutes}
                    onchange={applyEndTime}
                    class="form-select w-full"
                  >
                    {#each minutes as minute}
                      <option value={minute}>{minute} min</option>
                    {/each}
                  </select>
                  <div class="select-icon">
                    <ChevronDown size={18} />
                  </div>
                </div>
              </div>
            </div>

            <!-- Champ Nom -->
            <div>
              <label for="activity-name" class="block text-gray-700 font-medium mb-2">
                Nom
                <span class="text-gray-400">(optionnel)</span>
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

            <!-- Champ Description -->
            <div>
              <label for="activity-description" class="block text-gray-700 font-medium mb-2">
                Description
                <span class="text-gray-400">(optionnel)</span>
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

            <!-- Catégorie avec dropdown et recherche intégrée -->
            <div>
              <label for="activity-category-search" class="block text-gray-700 font-medium mb-2">
                Catégorie
                <span class="text-red-500">*</span>
              </label>
              <div class="category-dropdown-container">
                <!-- Sélecteur personnalisé avec champ de recherche intégré -->
                <div
                  class="custom-select"
                  class:disabled-select={!activity.projectId}
                  onclick={(e) => {
                    if (!activity.projectId) {
                      e.preventDefault();
                      return;
                    }
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
                    onfocus={() => {
                      if (activity.projectId) {
                        categoryDropdownOpen = true;
                        isSearchFocused = true;
                      }
                    }}
                    onblur={() => {
                      isSearchFocused = false;
                    }}
                    onclick={(e) => e.stopPropagation()}
                  />

                  <div class="arrow-icon">
                    <ChevronDown size={18} class="text-gray-600" />
                  </div>
                </div>

                {#if categoryDropdownOpen && activity.projectId}
                  <div class="dropdown-content">
                    <!-- Liste des catégories -->
                    <div class="category-list">
                      {#if filteredCategories.length === 0}
                        <div class="no-results">
                          {searchTerm
                            ? 'Aucun résultat trouvé'
                            : 'Aucune catégorie disponible pour ce projet'}
                        </div>
                      {:else}
                        {#each filteredCategories as category}
                          <div
                            class="category-item"
                            class:selected={category.id === activity.categoryId}
                            onclick={(e) => {
                              e.stopPropagation();
                              selectCategory(category.id);
                            }}
                          >
                            {category.name}
                          </div>
                        {/each}
                      {/if}
                    </div>

                    <!-- Option "Ajouter une nouvelle catégorie" uniquement si un projet est sélectionné -->
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
          </div>

          <!-- Séparateur et boutons d'action -->
          <div class="mt-auto">
            <!-- Ligne de séparation -->
            <div class="border-t border-gray-200 my-6"></div>

            <!-- Actions en bas du formulaire -->
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
                  disabled={isSubmitting}
                >
                  {isSubmitting ? 'En cours...' : 'Modifier'}
                </button>
              {:else}
                <button
                  type="button"
                  class="py-3 px-6 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 hover:-translate-y-0.5 hover:shadow-sm active:translate-y-0 transition border border-gray-200"
                  onclick={handleClose}
                >
                  Annuler
                </button>
                <button
                  type="submit"
                  class="py-3 px-6 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition disabled:opacity-50"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? 'En cours...' : 'Créer'}
                </button>
              {/if}
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
  <ConfirmationCreateCategory
    show={showCategoryConfirmModal}
    title="Confirmer l'ajout"
    message={`Voulez-vous ajouter la catégorie "<strong>${categoryToAdd}</strong>" ?`}
    onConfirm={() => {
      showCategoryConfirmModal = false;
      confirmAddCategory();
    }}
    onCancel={() => {
      showCategoryConfirmModal = false;
    }}
  />
{/if}


<style>
  /* Animation pour le panneau latéral - ne peut pas être fait en Tailwind standard */
  @keyframes slideIn {
    from {
      transform: translateX(-100%);
    }
    to {
      transform: translateX(0);
    }
  }

  .animate-slideIn {
    animation: slideIn 0.3s forwards;
  }

  /* Classe condensée pour tous les selects */
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

  .form-select-flex {
    flex: 1;
  }

  .form-select::-ms-expand {
    display: none;
  }

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

  .select-container {
    position: relative;
    width: 100%;
  }

  .select-icon {
    position: absolute;
    right: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    pointer-events: none;
    color: #606060;
  }

  .fixed {
    z-index: 40; /* Plus élevé que le dashboard */
  }

  /* Styles pour le dropdown de catégories */
  .category-dropdown-container {
    position: relative;
    width: 100%;
  }

  .custom-select {
    width: 100%;
    height: 42px; /* Hauteur fixe pour correspondre aux autres champs */
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
    padding-right: 2.5rem; /* Espace pour la flèche */
    border: none;
    background: transparent;
    font-size: inherit;
    color: inherit;
    outline: none;
    cursor: text;
  }

  .search-input::placeholder {
    color: #9ca3af;
  }

  .select-value {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    padding: 0 1rem;
    padding-right: 2.5rem; /* Espace pour la flèche */
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
    pointer-events: none; /* permet les clics à travers */
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

  .category-list {
    max-height: 200px;
    overflow-y: auto;
  }

  .category-item {
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: background-color 0.15s;
  }

  .category-item:hover {
    background-color: #f3f4f6;
  }

  .category-item.selected {
    background-color: #e5e7eb;
  }

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

  .add-category:hover {
    background-color: #f3f4f6;
  }

  .no-results {
    padding: 0.75rem 1rem;
    color: #6b7280;
    text-align: center;
  }

  .disabled-select {
    background-color: #f3f4f6;
    cursor: not-allowed;
    opacity: 0.7;
  }
</style>
