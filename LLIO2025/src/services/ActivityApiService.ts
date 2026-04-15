import type { Activity, ActivityUpdateResponse, OutlookEvent, OutlookEventRaw } from '../Models/index';
import { GET, POST, PUT, DELETE } from '../ts/server';

interface RawActivity extends Omit<Activity, 'startDate' | 'endDate'> {
  startDate: string;
  endDate: string;
}

const transformActivityStringToDates = (activity: RawActivity): Activity => ({
  ...activity,
  startDate: new Date(activity.startDate),
  endDate: new Date(activity.endDate),
});

const transformEventStringDates = (event: OutlookEventRaw): OutlookEvent => ({
  ...event,
  start: new Date(`${event.start.dateTime}${event.start.timeZone === "UTC" ? "Z" : ""}`),
  end: new Date(`${event.end.dateTime}${event.end.timeZone === "UTC" ? "Z" : ""}`),
});

const toStringDatesToString = (activity: Activity) => ({
  ...activity,
  startDate: activity.startDate.toISOString(),
  endDate: activity.endDate.toISOString(),
});

const createActivity = async (activity: Activity): Promise<Activity> => {
  // Préparation des données pour l'API
  const activityForApi = toStringDatesToString(activity);

  try {
    const response = await POST<typeof activityForApi, { activity: RawActivity }>(
      '/activity',
      activityForApi
    );
    return transformActivityStringToDates(response.activity);
  } catch (error) {
    console.error('Erreur lors de la création de tâche:', error);
    throw error;
  }
};

const updateActivity = async (activity: Activity): Promise<Activity> => {
  if (!activity) {
    console.error('La tâche est undefined');
    throw new Error('La tâche ou son ID est manquant');
  }

  if (!activity.id) {
    console.error("L'ID de la tâche est manquant", activity);
    throw new Error("L'ID de la tâche est manquant");
  }

  const activityForApi = toStringDatesToString(activity);

  try {
    const response = await PUT<typeof activityForApi, ActivityUpdateResponse>(
      '/activity',
      activityForApi
    );

    // Extraire la tâche mise à jour de l'objet de réponse
    const updatedactivityData = response.updated_activity;

    // Convertir les dates
    return {
      ...updatedactivityData,
      startDate: new Date(updatedactivityData.startDate),
      endDate: new Date(updatedactivityData.endDate),
    };
  } catch (error) {
    console.error('Erreur lors de la mise à jour de la tâche:', error);
    throw error;
  }
};

const deleteActivity = async (id: number): Promise<void> => {
  await DELETE(`/activity/${id}`);
};

const getAllActivitesFromRange = async (startDate: string, endDate: string) => {
  try {
    
    const response = await GET<{ activities: RawActivity[] }>(
      '/activities/me/detailed?startDate=' + startDate + '&endDate=' + endDate
    );

    if (response?.activities && Array.isArray(response.activities)) {
      return response.activities.map(transformActivityStringToDates);
    } else {
      console.error('Format de réponse inattendu :', response);
      alert('Erreur: Format de réponse inattendu lors de la récupération des activités.');
      return [];
    }
    
  } catch (error) {
    console.error('Erreur lors de la récupération des activités', error);
    alert('Erreur lors de la récupération des activités. Veuillez réessayer plus tard.');
    throw error;
  }

  ;
};

const getOutlookEvents = async (date: string): Promise<OutlookEvent[]> => {
  try {
    const response = await GET<{ events: OutlookEventRaw[] }>(
      '/activities/me/outlook?date=' + date,
      false
    );

    if (response?.events && Array.isArray(response.events)) {
      return response.events.map(transformEventStringDates);
    } else {
      console.error('Format de réponse inattendu :', response);
      alert('Erreur: Format de réponse inattendu lors de la récupération des évènements Outlook.');
      return [];
    }
  } catch (error) {
    const erreur = error as any;
    console.log(erreur.code);
    if (erreur?.code === 'GRAPH_EXPIRED') {
      console.error('Erreur lors de la récupération des événements Outlook : Votre compte Microsoft est déconnecté. Veuillez vous reconnecter au site web.', error);
    }
    console.error('Erreur lors de la récupération des événements Outlook', error);
    throw error;
  }
}

export const ActivityApiService = {
  createActivity,
  updateActivity,
  deleteActivity,
  getAllActivitesFromRange,
  getOutlookEvents
};
