<script lang="ts">
  import type { CalendarService } from '../../services/calendar.service';
  import { UserApiService } from '../../services/UserApiService';
  import { ProjectApiService } from '../../services/ProjectApiService';
  import { CalendarService as CS } from '../../services/calendar.service';
  import { onMount, tick } from 'svelte';
  import ActivityModal from '../../Components/Calendar/ActivityModal.svelte';
  import DashboardLeftPane from '../../Components/Calendar/DashboardLeftPane.svelte';
  import { ActivityApiService } from '../../services/ActivityApiService';
  import type { Activity, UserInfo, Project, DetailedProject } from '../../Models/index.ts';
  // Importez le fichier CSS
  import '../../style/modern-calendar.css';
  import { getDateOrDefault, formatDate } from '../../utils/date';
  // Importer FullCalendar en français
  import frLocale from '@fullcalendar/core/locales/fr';
  import { formatViewTitle } from '../../utils/date';
  import { Plus, Calendar, ChevronLeft, ChevronRight, LogOut } from 'lucide-svelte';
  import { goto } from '$app/navigation';

  let calendarEl = $state<HTMLElement | null>(null);
  let calendarService = $state<CalendarService | null>(null);
  let showModal = $state(false);
  let selectedDate: { start: Date; end: Date } | null = $state(null);
  let editMode = $state(false);
  let editActivity = $state(null);
  let activeView = $state('timeGridWeek');
  let currentViewTitle = $state('');
  let isLoading = $state(false);
  let dateStart = $state(null);
  let dateEnd = $state(null);
  let textHoursWorked = $state('');
  let totalHours = $state(0);
  let headerFormat: { weekday: string; day?: string; month?: string };

  const timeRanges = [
    { label: 'Heures de bureau', start: '06:00:00', end: '19:00:00', default: true },
    { label: 'Toute la journée', start: '00:00:00', end: '24:00:00' },
  ];

  let activeTimeRange = $state(timeRanges.find((range) => range.default));

  let currentUser = $state<UserInfo | null>(null);
  let projects = $state<Project[]>([]);
  let detailedProjects = $state<DetailedProject[]>([]);

  // Fonction pour attribuer une couleur à chaque événement
  const getEventClassName = (eventInfo: any) => {
    const eventTypes = [
      'event-blue',
      'event-green',
      'event-teal',
      'event-dark-teal',
      'event-light-teal',
    ];

    if (eventInfo.event.extendedProps.categoryId) {
      const categoryId = parseInt(eventInfo.event.extendedProps.categoryId);
      return eventTypes[categoryId % eventTypes.length];
    }

    const hash = eventInfo.event.title.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
    return eventTypes[hash % eventTypes.length];
  }

  // Fonction pour mettre à jour le titre de la période courante
  const updateViewTitle = () => {
    if (calendarService?.calendar) {
      const dateAPI = calendarService.calendar.getDate();
      const viewType = calendarService.calendar.view.type;
      currentViewTitle = formatViewTitle(viewType, dateAPI);
    }
  }

  // Fonction pour charger toutes les activités
  const loadActivities = async () =>{
    let day, diff;
    try {
      switch (activeView) {
        case 'dayGridMonth':
          dateStart = calendarService.calendar.getDate();
          dateStart.setDate(1); // Premier jour du mois
          dateEnd = new Date(dateStart.getFullYear(), dateStart.getMonth() + 1, 0); // Dernier jour du mois
          textHoursWorked = 'ce mois-ci';
          headerFormat = { weekday: 'short'};
          break;

        case 'timeGridWeek':
          dateStart = calendarService.calendar.getDate();
          day = dateStart.getDay();
          diff = dateStart.getDate() - day + (day === 0 ? -6 : 1); // Ajuster lorsque jour = dimanche
          dateStart.setDate(diff);
          dateEnd = new Date(dateStart);
          dateEnd.setDate(dateStart.getDate() + 6);
          textHoursWorked = 'cette semaine';
          headerFormat = { weekday: 'short', day: '2-digit', month: '2-digit' };
          break;
        case 'timeGridDay':
          dateStart = calendarService.calendar.getDate();
          dateEnd = new Date(dateStart);
          textHoursWorked = "aujourd'hui";
          break;
        default:
          dateStart = calendarService.calendar.getDate();
          day = dateStart.getDay();
          diff = dateStart.getDate() - day + (day === 0 ? -6 : 1); // Ajuster lorsque jour = dimanche
          dateStart.setDate(diff);

          dateEnd = new Date(dateStart);
          dateEnd.setDate(dateStart.getDate() + 6);
          headerFormat = { weekday: 'short', day: '2-digit', month: '2-digit' };
          break;
      }

      dateStart = formatDate(dateStart);
      dateEnd = formatDate(dateEnd);
      let activities = await ActivityApiService.getAllActivitesFromRange(dateStart, dateEnd);

      // Utiliser la méthode du service pour ajouter les activités au calendrier
      if (activities && calendarService) {
        calendarService.loadActivities(activities);
        totalHours = calendarService.getTotalHours();
      }
    } catch (error) {
      alert('Une erreur est survenue lors du chargement des activités.');
    }
  }

  const loadProjects = async () =>{
    try {
      isLoading = true;
      projects = await ProjectApiService.getProjects();
      detailedProjects = await ProjectApiService.getCurrentUserProjects();
    } catch (err) {
      console.error('Erreur lors de la récupération des projets:', err);
      alert('Une erreur est survenue lors de la récupération des projets.');
      projects = [];
      detailedProjects = [];
    } finally {
      isLoading = false;
    }
  }

  onMount(async () => {
    isLoading = true;
    if (calendarEl) {
      calendarService = new CS();

      // Configuration personnalisée pour FullCalendar
      const calendarOptions = {
        initialView: activeView,
        locale: frLocale, // Utiliser la locale française
        firstDay: 1, // 1 = lundi (standard français)
        buttonText: {
          today: "Aujourd'hui",
          month: 'Mois',
          week: 'Semaine',
          day: 'Jour',
        },
        slotDuration: '00:30:00', // Durée de chaque intervalle de temps
        allDaySlot: false,
        slotMinTime: activeTimeRange.start,
        slotMaxTime: activeTimeRange.end,
        nowIndicator: true,

        // Gestion du drag
        editable: true,
        eventDrop: handleEventDropOrResize,
        eventResize: handleEventDropOrResize,

        height: 'auto',
        contentHeight: 'auto', // Hauteur automatique
        slotHeight: 25, // Hauteur réduite des slots (plus compact)
        expandRows: true,
        dayHeaderFormat: headerFormat,
        eventClassNames: getEventClassName,
        eventTimeFormat: {
          hour: '2-digit',
          minute: '2-digit',
          meridiem: false,
          hour12: false,
        },
        slotLabelFormat: {
          hour: 'numeric',
          minute: '2-digit',
          hour12: false,
        },
        datesSet: () => {
          // Appelé à chaque changement de dates ou de vue
          updateViewTitle();
        },
      };

      calendarService.onDateSelect = (info) => {
        editMode = false;
        editActivity = null;
        selectedDate = info;
        showModal = true;
      };

      calendarService.onEventClick = (info) => {
        editMode = true;
        editActivity = {
          id: info.event.extendedProps.id,
          name: info.event.extendedProps.name,
          description: info.event.extendedProps.description,
          userId: info.event.extendedProps.userId,
          projectId: info.event.extendedProps.projectId,
          categoryId: info.event.extendedProps.categoryId,
          startDate: info.event.start,
          endDate: info.event.end,
        };
        selectedDate = {
          start: info.event.start,
          end: info.event.end,
        };
        showModal = true;
      };

      // Charger les informations utilisateur
      try {
        currentUser = await UserApiService.getUserInfo();
      } catch (error) {
        console.error('Erreur lors du chargement des informations utilisateur:', error);
      }

      // Initialiser avec les options personnalisées
      calendarService.initialize(calendarEl, calendarOptions);
      calendarService.render();

      // Mettre à jour le titre initial
      updateViewTitle();

      await loadActivities();
      await loadProjects();
    }
    isLoading = false;
  });

  const setView = (viewName: string) => {
    if (calendarService) {
      calendarService.setView(viewName);
      activeView = viewName;
      updateViewTitle();
      loadActivities();
    }
  }

  const handleActivitySubmit = async (activityData: Activity) =>{
    calendarService.addEvent({
      id: activityData.id.toString(),
      title: activityData.projectName,
      start: activityData.startDate,
      end: activityData.endDate,
      extendedProps: { ...activityData },
    });
    totalHours = calendarService.getTotalHours();
    detailedProjects = await ProjectApiService.getCurrentUserProjects();
  }

  const handleActivityUpdate = async (activity: Activity) =>{
    if (!calendarService?.calendar) return;

    try {
      // Utiliser la fonction pour valider les dates
      const now = new Date();
      activity.startDate = getDateOrDefault(activity.startDate, now);

      // si la date est invalide, définir par défaut 30 minutes après le début
      const defaultEndDate = new Date(activity.startDate.getTime() + 30 * 60000);
      activity.endDate = getDateOrDefault(activity.endDate, defaultEndDate);

      const updatedActivity = await ActivityApiService.updateActivity(activity);
      calendarService.updateEvent(updatedActivity);
      totalHours = calendarService.getTotalHours();
      detailedProjects = await ProjectApiService.getCurrentUserProjects();
    } catch (error) {
      console.error("Erreur lors de la mise à jour de l'activité", error);

      alert("Une erreur est survenue lors de la mise à jour de l'activité.");

      throw error;
    }
  }
  const handleActivityDelete = async (activity: Activity) =>{
    if (!calendarService?.calendar || !activity.id) return;
    try {
      await ActivityApiService.deleteActivity(activity.id);
      calendarService.deleteActivity(activity.id.toString());
      totalHours = calendarService.getTotalHours();
      detailedProjects = await ProjectApiService.getCurrentUserProjects();
    } catch (error) {
      console.error("Erreur lors de la suppression de l'activité", error);
      throw error;
    }
  }

  // Fonction pour gérer le déplacement et le redimmensionnement d'une tâche
  const handleEventDropOrResize = async (info) =>{
    try {
      const activity = calendarService.eventToActivity(info);

      const updatedActivity = await ActivityApiService.updateActivity(activity);

      calendarService.updateEvent(updatedActivity);
      totalHours = calendarService.getTotalHours();
      detailedProjects = await ProjectApiService.getCurrentUserProjects();

    } catch (error) {
      console.error("Erreur lors de la mise à jour de l'activité", error);
      alert("Une erreur est survenue lors de la mise à jour de l'activité.");
      info.revert();
    }
  }


  const handleNewActivity = () => {
    editMode = false;
    editActivity = null;
    selectedDate = {
      start: new Date(),
      end: new Date(new Date().getTime() + 30 * 60000), // 30 minutes par défaut
    };
    showModal = true;
  }

  const prevPeriod = () => {
    calendarService?.prev();
    updateViewTitle();
    loadActivities();
  }

  const nextPeriod = () => {
    calendarService?.next();
    updateViewTitle();
    loadActivities();
  }

  const goToday = () => {
    calendarService?.today();
    updateViewTitle();
    loadActivities();
  }

  const setTimeRange = async (range) => {
  activeTimeRange = range;

  const cal = calendarService?.calendar;
  if (!cal) return;

  cal.batchRendering(() => {
    cal.setOption('slotMinTime', range.start);
    cal.setOption('slotMaxTime', range.end === '24:00:00' ? '23:59:59' : range.end);
  });

  await tick();
  requestAnimationFrame(() => {
    cal.updateSize();
  });
};


  const currentDate = new Date();
  const formattedDate = currentDate.toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
</script>

<div class="flex">
  <!-- Dashboard toujours visible à gauche -->
  {#if isLoading}
    <div class="fixed left-0 top-0 w-[300px] h-full bg-gray-100 animate-pulse"></div>
  {:else if currentUser}
    <DashboardLeftPane
      {detailedProjects}
      {currentUser}
      {dateStart}
      {dateEnd}
      {totalHours}
      {textHoursWorked}
    />
  {/if}

  <!-- Contenu principal (calendrier) avec marge pour s'adapter au dashboard -->
  <div class="space-between-dashboard-calendar w-full h-full bg-white px-4 py-6">
    <div class="max-w-7xl mx-auto">
      <!-- Nouvelle section avec salutation et boutons d'heures -->
      <div class="flex justify-between items-center mb-6">
        <!-- Affichage nom d'utilisateur -->
        <h1 class="text-2xl font-bold text-gray-800 flex items-center gap-2">
          Bonjour,
          <span class="text-[#015e61] font-bold">
            {#if currentUser}
              {currentUser.firstName} {currentUser.lastName}
            {:else}
              <span class="inline-block w-24 h-6 bg-gray-200 animate-pulse rounded"></span>
            {/if}
          </span>
          <button
            class="ml-2 mt-1 p-1.5 rounded-full hover:bg-gray-100 text-gray-600 hover:text-[#015e61] transition-colors"
            title="Se déconnecter" 
            onclick={async () => {
              await UserApiService.logOut();
              goto("/");
            }} 
            >
            <LogOut class="w-5 h-5" />
          </button>
        </h1>

        <!-- Boutons pour changer les heures -->
        <div class="flex items-center">
          {#each timeRanges as range, index}
            <button
              class="py-2 px-4 text-sm transition-colors font-semibold
              {index === 0
                ? 'rounded-l-lg rounded-r-none'
                : index === timeRanges.length - 1
                  ? 'rounded-r-lg rounded-l-none'
                  : 'rounded-none border-x border-white/20'}
              {activeTimeRange.label === range.label
                ? 'bg-[#015e61] text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}
              "
              onclick={() => setTimeRange(range)}
            >
              {range.label}
            </button>
          {/each}
        </div>
      </div>

      <!-- Informations de date et contrôles du calendrier -->
      <div class="flex flex-col md:flex-row justify-between items-center mb-6 gap-4">
        <!-- Titre avec icône -->
        <div class="flex items-center">
          <Calendar class="mr-2" />
          <p class="text-md text-gray-600 text-xl font-semibold">Aujourd'hui, {formattedDate}</p>
        </div>

        <!-- Boutons de vue alignés au centre -->
        <div class="flex items-center bg-gray-100 p-1 rounded-lg">
          <button
            class="px-5 py-2 rounded-lg transition-colors {activeView === 'timeGridDay'
              ? 'bg-white text-[#015e61] font-medium'
              : 'text-gray-500 hover:bg-white hover:text-[#015e61]'}"
            onclick={() => setView('timeGridDay')}
          >
            Jour
          </button>
          <button
            class="px-5 py-2 rounded-lg transition-colors {activeView === 'timeGridWeek'
              ? 'bg-white text-[#015e61] font-medium'
              : 'text-gray-500 hover:bg-white hover:text-[#015e61]'}"
            onclick={() => setView('timeGridWeek')}
          >
            Semaine
          </button>
          <button
            class="px-5 py-2 rounded-lg transition-colors {activeView === 'dayGridMonth'
              ? 'bg-white text-[#015e61] font-medium'
              : 'text-gray-500 hover:bg-white hover:text-[#015e61]'}"
            onclick={() => setView('dayGridMonth')}
          >
            Mois
          </button>
        </div>

        <!-- Bouton New Activity -->
        <button
          onclick={handleNewActivity}
          class="bg-[#015e61] hover:bg-[#014446] text-white py-2 px-6 rounded-xl flex items-center gap-2 font-semibold transition-colors"
        >
          <Plus class="h-5 w-5" />
          Nouvelle activité
        </button>
      </div>

      <!-- Contrôles de navigation -->
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center space-x-3">
          <button
            onclick={prevPeriod}
            class="p-2 rounded-lg bg-gray-100 hover:bg-gray-200 transition-colors"
          >
            <ChevronLeft class="w-6 h-6 text-gray-600" />
          </button>
          <button
            onclick={nextPeriod}
            class="p-2 rounded-lg bg-gray-100 hover:bg-gray-200 transition-colors"
          >
            <ChevronRight class="w-6 h-6 text-gray-600" />
          </button>
          <button
            onclick={goToday}
            class="px-4 py-2 rounded-lg bg-gray-100 hover:bg-gray-200 text-gray-700 transition-colors"
          >
            Aujourd'hui
          </button>
        </div>
        <div class="text-lg font-medium text-gray-700">
          <!-- Titre dynamique de la période courante -->
          <span>{currentViewTitle}</span>
        </div>
        <div class="w-28">
          <!-- Élément vide pour maintenir l'alignement -->
        </div>
      </div>

      <!-- Calendrier -->
      <div class="border border-gray-200 rounded-lg overflow-hidden">
        <div bind:this={calendarEl} class="w-full"></div>
      </div>
    </div>
  </div>
</div>

<!-- Modal d'activité qui s'affiche par-dessus tout le reste -->
{#if showModal}
  <ActivityModal
    show={showModal}
    {projects}
    activityToEdit={editActivity}
    {selectedDate}
    onDelete={handleActivityDelete}
    onSubmit={handleActivitySubmit}
    onUpdate={handleActivityUpdate}
    onClose={() => {
      showModal = false;
    }}
  />
{/if}

<style>
  :global(.fc .fc-timegrid-slot) {
    height: 25px !important;
    min-height: 25px !important;
    max-height: 25px !important;
  }

  :global(.fc-timegrid-event) {
    min-height: 20px !important;
    max-height: none !important;
  }

  :global(.fc-timegrid-slot-label) {
    vertical-align: top !important;
    padding-top: 2px !important;
  }

  :global(.fc-theme-standard td),
  :global(.fc-theme-standard th) {
    padding: 1px !important;
  }

  :global(.fc .fc-daygrid-day-frame) {
    padding: 2px !important;
  }

  /* Couleur personnalisée pour l'indicateur de l'heure actuelle */
  :global(.fc .fc-timegrid-now-indicator-line) {
    border-color: #015e61 !important;
  }

  :global(.fc .fc-timegrid-now-indicator-arrow) {
    border-color: #015e61 !important;
    background-color: #015e61 !important;
  }

  /* Gestion espace dashboard */
  .space-between-dashboard-calendar {
    margin-left: 300px;
  }
</style>
