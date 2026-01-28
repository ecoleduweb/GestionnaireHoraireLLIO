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
