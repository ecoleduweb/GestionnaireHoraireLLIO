import { MockConfig } from "../types";

export const projectMocks = {
  getDetailedProjectsSuccess: {
    url: "/projects/detailed",
    method: "GET",
    response: {
      status: 200,
      json: {
        projects: [
          {
            id: 1,
            uniqueId: "TEST-1",
            name: "AT-123",
            color: "blue",
            lead: "Katell Arnault de la Ménardière",
            isArchived: false,
            coLeads: ["Jean-François Jasmin"],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Développement", timeSpent: 30, timeEstimated: 25 },
                  { name: "Graphisme", timeSpent: 15, timeEstimated: 30 },
                ],
              },
              {
                name: "Jean-François Jasmin",
                categories: [
                  { name: "Développement", timeSpent: 20, timeEstimated: 50 },
                ],
              },
              {
                name: "Laure Desjardins",
                categories: [],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -45,
            totalTimeSpent: 65,
          },
          {
            id: 2,
            uniqueId: "TEST-2",
            name: "AT-456",
            color: "pink",
            lead: "Katell Arnault de la Ménardière",
            isArchived: false,
            coLeads: ["Marie Amélie Dubé", "Jimmy Paquet-Cormier"],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Développement", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 18, timeEstimated: 20 },
                ],
              },
              {
                name: "Ariane Dionne-Santerre",
                categories: [
                  { name: "Rédaction", timeSpent: 8, timeEstimated: 10 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 14, timeEstimated: 15 },
                ],
              },
            ],
            totalTimeEstimated: 200,
            totalTimeRemaining: 148,
            totalTimeSpent: 52,
          },
          {
            id: 3,
            uniqueId: "TEST-3",
            name: "FO-115",
            color: "yellow",
            lead: "Katell Arnault de la Ménardière",
            isArchived: false,
            coLeads: [],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 5, timeEstimated: 5 },
                ],
              },
              {
                name: "Annie Côté",
                categories: [
                  { name: "Graphisme", timeSpent: 22, timeEstimated: 22 },
                ],
              },
            ],
            totalTimeEstimated: 27,
            totalTimeRemaining: 0,
            totalTimeSpent: 27,
          },
          {
            id: 4,
            uniqueId: "TEST-4",
            name: "RA-224",
            color: "red",
            lead: "Steve Joncoux",
            isArchived: false,
            coLeads: [
              "Katell Arnault de la Ménardière",
              "Marjolaine Poirier",
              "Jimmy Paquet-Cormier",
            ],
            employees: [
              {
                name: "Steve Joncoux",
                categories: [
                  { name: "Développement", timeSpent: 16, timeEstimated: 20 },
                ],
              },
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 7, timeEstimated: 10 },
                ],
              },
              {
                name: "Marjolaine Poirier",
                categories: [
                  { name: "Rédaction", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 9, timeEstimated: 12 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 11, timeEstimated: 10 },
                ],
              },
            ],
            totalTimeEstimated: 0,
            totalTimeRemaining: -55,
            totalTimeSpent: 55,
          },
          {
            id: 5,
            uniqueId: "TEST-5",
            name: "AT-789",
            color: "blue",
            lead: "Katell Arnault de la Ménardière",
            isArchived: true,
            coLeads: ["Jean-François Jasmin"],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Développement", timeSpent: 30, timeEstimated: 25 },
                  { name: "Graphisme", timeSpent: 15, timeEstimated: 30 },
                ],
              },
              {
                name: "Jean-François Jasmin",
                categories: [
                  { name: "Développement", timeSpent: 20, timeEstimated: 50 },
                ],
              },
              {
                name: "Laure Desjardins",
                categories: [],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -45,
            totalTimeSpent: 65,
          },
          {
            id: 6,
            uniqueId: "TEST-6",
            name: "AT-987",
            color: "pink",
            lead: "Katell Arnault de la Ménardière",
            isArchived: true,
            coLeads: ["Marie Amélie Dubé", "Jimmy Paquet-Cormier"],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Développement", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 18, timeEstimated: 20 },
                ],
              },
              {
                name: "Ariane Dionne-Santerre",
                categories: [
                  { name: "Rédaction", timeSpent: 8, timeEstimated: 10 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 14, timeEstimated: 15 },
                ],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -32,
            totalTimeSpent: 52,
          },
          {
            id: 7,
            uniqueId: "TEST-7",
            name: "FO-789",
            color: "yellow",
            lead: "Katell Arnault de la Ménardière",
            isArchived: true,
            coLeads: [],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 5, timeEstimated: 8 },
                ],
              },
              {
                name: "Annie Côté",
                categories: [
                  { name: "Graphisme", timeSpent: 22, timeEstimated: 25 },
                ],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -7,
            totalTimeSpent: 27,
          },
          {
            id: 8,
            uniqueId: "TEST-8",
            name: "RA-987",
            color: "red",
            lead: "Steve Joncoux",
            isArchived: true,
            coLeads: [
              "Katell Arnault de la Ménardière",
              "Marjolaine Poirier",
              "Jimmy Paquet-Cormier",
            ],
            employees: [
              {
                name: "Steve Joncoux",
                categories: [
                  { name: "Développement", timeSpent: 16, timeEstimated: 20 },
                ],
              },
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 7, timeEstimated: 10 },
                ],
              },
              {
                name: "Marjolaine Poirier",
                categories: [
                  { name: "Rédaction", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 9, timeEstimated: 12 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 11, timeEstimated: 10 },
                ],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -35,
            totalTimeSpent: 55,
          },
        ],
      },
    },
  },
  getProjectsListSuccess: {
    url: "/projects",
    method: "GET",
    response: {
      status: 200,
      json: {
        projects: [
          {
            id: 1,
            managerId: 1,
            uniqueId: "Projet sous-sol",
            name: "Nommer le projet",
            status: 1,
            billable: false,
            activities: [],
            categories: [],
            createdAt: "2025-03-22T08:00:00",
            updatedAt: "2025-03-22T08:00:00",
            end_at: null,
          },
        ],
      },
    },
  },
  addProjectSuccess: {
    url: "/project",
    method: "POST",
    response: {
      status: 201,
      json: {
        project: {
          id: 6,
          managerId: 3,
          name: "Jérémie Lapointe",
          description: "das",
          status: 0,
          billable: true,
          activities: [],
          categories: [],
          createdAt: "2025-03-23T15:07:14.991-04:00",
          updatedAt: "2025-03-23T15:07:14.991-04:00",
          endAt: "0001-01-01T00:00:00Z",
        },
        response: "Le projet a bien été ajouté à la base de données",
      },
    },
  },
  deleteProjectSuccess: {
    url: '/project/*',
    method: 'DELETE',
    response: {
      status: 200,
      json: {"deleted":true}
    }
  },
  deleteProjectError: {
        url: '/project/*',
        method: 'DELETE',
        response: {
                status: 403,
                json: {"error":"Le projet a des activités associées, suppression impossible"}
        }
  },
  getDetailedProjectsSuccessAfterDelete: {
    url: "/projects/detailed",
    method: "GET",
    response: {
      status: 200,
      json: {
        projects: [
          {
            id: 2,
            uniqueId: "TEST-2",
            name: "AT-456",
            color: "pink",
            lead: "Katell Arnault de la Ménardière",
            isArchived: false,
            coLeads: ["Marie Amélie Dubé", "Jimmy Paquet-Cormier"],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Développement", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 18, timeEstimated: 20 },
                ],
              },
              {
                name: "Ariane Dionne-Santerre",
                categories: [
                  { name: "Rédaction", timeSpent: 8, timeEstimated: 10 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 14, timeEstimated: 15 },
                ],
              },
            ],
            totalTimeEstimated: 200,
            totalTimeRemaining: 148,
            totalTimeSpent: 52,
          },
          {
            id: 3,
            uniqueId: "TEST-3",
            name: "FO-115",
            color: "yellow",
            lead: "Katell Arnault de la Ménardière",
            isArchived: false,
            coLeads: [],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 5, timeEstimated: 5 },
                ],
              },
              {
                name: "Annie Côté",
                categories: [
                  { name: "Graphisme", timeSpent: 22, timeEstimated: 22 },
                ],
              },
            ],
            totalTimeEstimated: 27,
            totalTimeRemaining: 0,
            totalTimeSpent: 27,
          },
          {
            id: 4,
            uniqueId: "TEST-4",
            name: "RA-224",
            color: "red",
            lead: "Steve Joncoux",
            isArchived: false,
            coLeads: [
              "Katell Arnault de la Ménardière",
              "Marjolaine Poirier",
              "Jimmy Paquet-Cormier",
            ],
            employees: [
              {
                name: "Steve Joncoux",
                categories: [
                  { name: "Développement", timeSpent: 16, timeEstimated: 20 },
                ],
              },
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 7, timeEstimated: 10 },
                ],
              },
              {
                name: "Marjolaine Poirier",
                categories: [
                  { name: "Rédaction", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 9, timeEstimated: 12 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 11, timeEstimated: 10 },
                ],
              },
            ],
            totalTimeEstimated: 0,
            totalTimeRemaining: -55,
            totalTimeSpent: 55,
          },
          {
            id: 5,
            uniqueId: "TEST-5",
            name: "AT-789",
            color: "blue",
            lead: "Katell Arnault de la Ménardière",
            isArchived: true,
            coLeads: ["Jean-François Jasmin"],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Développement", timeSpent: 30, timeEstimated: 25 },
                  { name: "Graphisme", timeSpent: 15, timeEstimated: 30 },
                ],
              },
              {
                name: "Jean-François Jasmin",
                categories: [
                  { name: "Développement", timeSpent: 20, timeEstimated: 50 },
                ],
              },
              {
                name: "Laure Desjardins",
                categories: [],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -45,
            totalTimeSpent: 65,
          },
          {
            id: 6,
            uniqueId: "TEST-6",
            name: "AT-987",
            color: "pink",
            lead: "Katell Arnault de la Ménardière",
            isArchived: true,
            coLeads: ["Marie Amélie Dubé", "Jimmy Paquet-Cormier"],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Développement", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 18, timeEstimated: 20 },
                ],
              },
              {
                name: "Ariane Dionne-Santerre",
                categories: [
                  { name: "Rédaction", timeSpent: 8, timeEstimated: 10 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 14, timeEstimated: 15 },
                ],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -32,
            totalTimeSpent: 52,
          },
          {
            id: 7,
            uniqueId: "TEST-7",
            name: "FO-789",
            color: "yellow",
            lead: "Katell Arnault de la Ménardière",
            isArchived: true,
            coLeads: [],
            employees: [
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 5, timeEstimated: 8 },
                ],
              },
              {
                name: "Annie Côté",
                categories: [
                  { name: "Graphisme", timeSpent: 22, timeEstimated: 25 },
                ],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -7,
            totalTimeSpent: 27,
          },
          {
            id: 8,
            uniqueId: "TEST-8",
            name: "RA-987",
            color: "red",
            lead: "Steve Joncoux",
            isArchived: true,
            coLeads: [
              "Katell Arnault de la Ménardière",
              "Marjolaine Poirier",
              "Jimmy Paquet-Cormier",
            ],
            employees: [
              {
                name: "Steve Joncoux",
                categories: [
                  { name: "Développement", timeSpent: 16, timeEstimated: 20 },
                ],
              },
              {
                name: "Katell Arnault de la Ménardière",
                categories: [
                  { name: "Gestion", timeSpent: 7, timeEstimated: 10 },
                ],
              },
              {
                name: "Marjolaine Poirier",
                categories: [
                  { name: "Rédaction", timeSpent: 12, timeEstimated: 15 },
                ],
              },
              {
                name: "Jimmy Paquet-Cormier",
                categories: [
                  { name: "Développement", timeSpent: 9, timeEstimated: 12 },
                ],
              },
              {
                name: "Marie Amélie Dubé",
                categories: [
                  { name: "Graphisme", timeSpent: 11, timeEstimated: 10 },
                ],
              },
            ],
            totalTimeEstimated: 20,
            totalTimeRemaining: -35,
            totalTimeSpent: 55,
          },
        ],
      },
    },
  },
  getDetailedProjectsByUserSuccess:{
    url: "/projects/me/detailed",
    method: "GET",
    response: {
      status: 200,
      json: {
        projects: [
          {
            billable: false,
            coLeads: [],
            createdAt: "2025-06-23T16:43:12-04:00",
            employees: [
              {
                categories: [
                  {
                    name: "Par défaut",
                    timeEstimated: 0,
                    timeSpent: 24,
                  },
                ],
                name: " ",
              },
            ],
            id: 11,
            isArchived: false,
            lead: "Toumani Camara",
            managerId: 4,
            name: "new commut",
            totalTimeEstimated: 24,
            totalTimeRemaining: 0,
            totalTimeSpent: 24,
            uniqueId: "new-1",
            updatedAt: "2025-06-30T16:50:03-04:00",
          },
          {
            billable: true,
            coLeads: [],
            createdAt: "2025-06-23T16:27:39-04:00",
            employees: [
              {
                categories: [
                  {
                    name: "Par défaut",
                    timeEstimated: 0,
                    timeSpent: 1,
                  },
                ],
                name: " ",
              },
            ],
            id: 10,
            isArchived: false,
            lead: "Toumani Camara",
            managerId: 4,
            name: "projet ! apre smigration ",
            totalTimeEstimated: 0,
            totalTimeRemaining: -1,
            totalTimeSpent: 1,
            uniqueId: "migr-2",
            updatedAt: "2025-06-23T16:27:39-04:00",
          },
          {
            billable: false,
            coLeads: [],
            createdAt: "2025-06-29T17:24:21-04:00",
            employees: [
              {
                categories: [
                  {
                    name: "Par défaut",
                    timeEstimated: 0,
                    timeSpent: 2.5,
                  },
                ],
                name: " ",
              },
            ],
            id: 13,
            isArchived: false,
            lead: "Usager test",
            managerId: 1,
            name: "le nom !",
            totalTimeEstimated: 20,
            totalTimeRemaining: 17.5,
            totalTimeSpent: 2.5,
            uniqueId: "123-mmm",
            updatedAt: "2025-06-30T16:38:53-04:00",
          },
          {
            billable: false,
            coLeads: [],
            createdAt: "2025-07-01T11:43:49-04:00",
            employees: [
              {
                categories: [
                  {
                    name: "Par défaut",
                    timeEstimated: 0,
                    timeSpent: 15,
                  },
                ],
                name: " ",
              },
            ],
            id: 14,
            isArchived: false,
            lead: "Toumani Camara",
            managerId: 4,
            name: "10 heures au total go negatif",
            totalTimeEstimated: 10,
            totalTimeRemaining: -5,
            totalTimeSpent: 15,
            uniqueId: "10htotal",
            updatedAt: "2025-07-01T11:43:49-04:00",
          },
        ],
      },
    },
  },
  getTwoProjectsSuccess: {
    url: "/projects",
    method: "GET",
    response: {
      status: 200,
      json: {
        projects: [
          {
            id: 1,
            managerId: 1,
            uniqueId: "Projet sous-sol",
            name: "Nommer le projet",
            status: 1,
            billable: false,
            activities: [],
            categories: [],
            createdAt: "2025-03-22T08:00:00",
            updatedAt: "2025-03-22T08:00:00",
            end_at: null,
          },
          {
            id: 2,
            managerId: 1,
            uniqueId: "Projet grenier",
            name: "Nommer le projet",
            status: 1,
            billable: false,
            activities: [],
            categories: [],
            createdAt: "2025-03-22T08:00:00",
            updatedAt: "2025-03-22T08:00:00",
            end_at: null,
          },
        ],
      },
    },
  },
} satisfies Record<string, MockConfig>;
