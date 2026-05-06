import { isThisWeek } from 'date-fns';

// Extraire les heures et minutes d'un objet Date
export const getHoursFromDate = (date: Date | null | undefined): string => {
  if (!date || !(date instanceof Date) || isNaN(date.getTime())) {
    return '12';
  }
  return date.getHours().toString().padStart(2, '0');
};

export const getMinutesFromDate = (date: Date | null | undefined): string => {
  if (!date || !(date instanceof Date) || isNaN(date.getTime())) {
    return '00';
  }
  return date.getMinutes().toString().padStart(2, '0');
};

// Créer une date avec année/mois/jour spécifiques et heures/minutes données
export const createDateWithTime = (
  baseDate: Date,
  hours: string | number,
  minutes: string | number
): Date => {
  const result = new Date(baseDate);
  result.setHours(Number(hours));
  result.setMinutes(Number(minutes));
  result.setSeconds(0);
  result.setMilliseconds(0);
  return result;
};

// Initialiser les dates de début et fin pour une nouvelle tâche
export const initializeActivityDates = (targetDate?: Date): { startDate: Date; endDate: Date } => {
  const date = targetDate ? new Date(targetDate) : new Date();
  // Réinitialiser les secondes et millisecondes
  date.setSeconds(0);
  date.setMilliseconds(0);

  // Créer des copies pour éviter les références partagées
  const startDate = new Date(date);
  // Par défaut, l'activité dure 1 heure
  const endDate = new Date(date);
  endDate.setHours(endDate.getHours() + 1);

  return { startDate, endDate };
};

// Valider que l'heure de fin est après l'heure de début
export const applyEndTime = (
  startHours: string,
  startMinutes: string,
  endHours: string,
  endMinutes: string
): { endHours: string; endMinutes: string } => {
  if (
    parseInt(endHours) < parseInt(startHours) ||
    (parseInt(endHours) === parseInt(startHours) && parseInt(endMinutes) < parseInt(startMinutes))
  ) {
    return {
      endHours: startHours,
      endMinutes: startMinutes,
    };
  }
  return { endHours, endMinutes };
};

export const getDateOrDefault = (
  dateToValidate: Date | null | undefined,
  defaultDate: Date
): Date => {
  if (!dateToValidate || !(dateToValidate instanceof Date) || isNaN(dateToValidate.getTime())) {
    return new Date(defaultDate);
  }
  return dateToValidate;
};
// Formatage du titre de la vue calendrier (jour, semaine, mois)
export const formatViewTitle = (viewType: string, date: Date): string => {
  if (viewType === 'dayGridMonth') {
    // Format pour la vue mois: "Février 2025"
    return date.toLocaleDateString('fr-FR', {
      month: 'long',
      year: 'numeric',
    });
  } else if (viewType === 'timeGridWeek') {
    // Format pour la vue semaine: "10 mars – 16 mars 2025"
    // Trouver le lundi de la semaine
    const startDate = new Date(date);
    const day = startDate.getDay();
    const diff = startDate.getDate() - day + (day === 0 ? -6 : 1); // Ajuster lorsque jour = dimanche
    startDate.setDate(diff);

    // Le dimanche est le 6ème jour après le lundi
    const endDate = new Date(startDate);
    endDate.setDate(startDate.getDate() + 6);

    const startFormatted = startDate.getDate();
    const startMonth = startDate.toLocaleDateString('fr-FR', { month: 'short' });

    const endFormatted = endDate.getDate();
    const endMonth = endDate.toLocaleDateString('fr-FR', { month: 'long' });
    const year = endDate.getFullYear();

    return `${startFormatted} ${startMonth}. – ${endFormatted} ${endMonth} ${year}`;
  } else if (viewType === 'timeGridDay') {
    // Format pour la vue jour: "Mercredi 26 février 2025"
    return date.toLocaleDateString('fr-FR', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    });
  }

  // Format par défaut si le type de vue n'est pas reconnu
  return date.toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

export const startEndDatesForFormat = () => {
  const startDate = new Date();
  const endDate = new Date(startDate);

  return {
    startDate: startDate.toLocaleString('fr-FR', { timeZone: 'UTC' }),
    endDate: endDate.toLocaleString('fr-FR', { timeZone: 'UTC' }),
  };
};

export const formatDate = (date: Date): string => {
  let laDate = new Date(date);
  let year = laDate.getFullYear();
  let month = (laDate.getMonth() + 1).toString(); // Les mois commencent à 0
  let day = laDate.getDate().toString();

  return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
};

export const formatHours = (hours: number | null | undefined): string => {
  if (!hours || hours === 0) return '-';

  const isNegative = hours < 0;
  const absoluteHours = Math.abs(hours);
  const hoursInt = Math.floor(absoluteHours);
  const minutes = Math.round((absoluteHours - hoursInt) * 60);

  // Format « h00 », « h05 », « h15 », etc.
  return `${isNegative ? '-' : ''}${hoursInt}h${
    minutes === 0 ? '00' : minutes < 10 ? `0${minutes}` : minutes
  }`;
};

export const formatDateHoursWorked = (date) => {
  let dateObj;

  if (date instanceof Date) {
    dateObj = date;
  } else if (typeof date === 'string') {
    const parts = date.split('-');
    dateObj = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]), 12, 0, 0);
  }

  return new Intl.DateTimeFormat('fr-FR', {
    day: 'numeric',
    month: 'long',
  }).format(dateObj);
};

export const areDatesEqual = (dateStart, dateEnd) => {
  if (!dateStart || !dateEnd) return false;
  return dateStart == dateEnd;
};

export const getHoursFromRange = (activity) => {
  const start =
    activity.startDate instanceof Date ? activity.startDate : new Date(activity.startDate);
  const end = activity.endDate instanceof Date ? activity.endDate : new Date(activity.endDate);

  const diffMilliseconds = end.getTime() - start.getTime();
  const hours = diffMilliseconds / (1000 * 60 * 60);
  return hours;
};

export const isDateInCurrentWeek = (date: Date) => {
  return isThisWeek(date, { weekStartsOn: 1 });
};

export const roundTimeToNearest15Minutes = (dateStr: string | Date) => {
  const date = new Date(dateStr);
  const ms = 15 * 60 * 1000;
  return new Date(Math.round(date.getTime() / ms) * ms);
};