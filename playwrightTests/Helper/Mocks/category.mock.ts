import { MockConfig } from "../types";

export const categoryMocks = {
  getCategoriesByProjectSuccess: {
    url: "/project/*/categories",
    method: "GET",
    response: {
      status: 200,
      json: {
        categories: [
          {
            id: 1,
            name: "Test",
            description: "Nouveau format",
            billable: false,
            activities: [],
            createdAt: "2025-04-02T10:29:54-04:00",
            updatedAt: "2025-04-02T10:29:54-04:00",
            userId: 1,
            projectId: 1,
          },
        ],
      },
    },
  },
  getCategoriesByProjectDeletedSuccess: {
    url: "/project/*/categories",
    method: "GET",
    response: {
      status: 200,
      json: {
        categories: [],
      },
    },
  },
  deleteCategorySuccess: {
    url: "/category/1",
    method: "DELETE",
    response: {
      status: 200,
      json: {}
    },
  },
  deleteCategoryFailed: {
    url: "/category/1",
    method: "DELETE",
    response: {
      status: 400,
      json: {}
    },
  },
  renameCategoryFailed: {
    url: "/category",
    method: "PUT",
    response: {
      status: 400,
      json: {}
    },
  },
  renameCategorySuccess: {
    url: "/category",
    method: "PUT",
    response: {
      status: 200,
      json: {"updatedCategory":{"id":1,"name":"test","description":"Créé depuis l'application","billable":false,"activities":[],"createdAt":"0001-01-01T00:00:00Z","updatedAt":"2026-04-08T10:58:55.537-04:00","projectId":0}}
    }
  },
  getCategoryForProject1Success: {
    url: "/project/1/categories",
    method: "GET",
    response: {
      status: 200,
      json: {
        categories: [
          {
            id: 1,
            name: "Design",
            description: "Tâches de design",
            billable: true,
            activities: [],
            createdAt: "2025-04-02T10:29:54-04:00",
            updatedAt: "2025-04-02T10:29:54-04:00",
            userId: 1,
            projectId: 1,
          },
        ],
      },
    },
  },
  getCategoryForProject2Success: {
    url: "/project/2/categories",
    method: "GET",
    response: {
      status: 200,
      json: {
        categories: [
          {
            id: 2,
            name: "Développement",
            description: "Tâches de développement",
            billable: true,
            activities: [],
            createdAt: "2025-04-02T10:29:54-04:00",
            updatedAt: "2025-04-02T10:29:54-04:00",
            userId: 1,
            projectId: 2,
          },
        ],
      },
    },
  },
} satisfies Record<string, MockConfig>;
