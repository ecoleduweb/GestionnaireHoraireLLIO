import type { Activity, Category, OutlookEvent } from '../Models';
import { initializeActivityDates, roundTimeToNearest15Minutes } from '../utils/date';

export const activityTemplate = {
  generate: (selectedEvent: OutlookEvent|null = null,
             selectedProject: {id: number, name: string} | null = null): Activity => {
    const { startDate, endDate } = initializeActivityDates();
    return {
      name: selectedEvent ? selectedEvent.subject : '',
      description: selectedEvent ? selectedEvent.body.content : '',
      userId: 1,
      projectId: selectedProject ? selectedProject.id : 0,
      categoryId: null,
      projectName: selectedProject ? selectedProject.name : '',
      startDate: selectedEvent ? roundTimeToNearest15Minutes(selectedEvent.start) : startDate,
      endDate: selectedEvent ? roundTimeToNearest15Minutes(selectedEvent.end) : endDate,
    };
  },

  time: {
    hours: Array.from({ length: 24 }, (_, i) => i.toString().padStart(2, '0')),
    minutes: ['00', '15', '30', '45'],
  },
} as const;
