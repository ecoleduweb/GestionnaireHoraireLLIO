import { MockConfig } from "../types";

export const outlookMocks = {
  getActivityWithoutProjectIdFromOutlook: {
    url: "/activities/me/outlook?date=2026-05-01",
    method: "GET",
    response: {
      status: 200,
      json: {
        date: "2026-05-01",
        events: [
          {
            id: "AAMkAGU2NjNjZjY5LWJiNTktNDk3YS1iOTA3LWRjZjJkYzg2YzI4MwBGAAAAAACobdNDVfcWRrhq17eN0GJQBwC9ZRNNh3HrTpUd-xS6ahyvAAAAAAENAAC9ZRNNh3HrTpUd-xS6ahyvAAG-SjUyAAA=",
            subject: "Test sans Id Projet",
            body: { content: "\r\n" },
            start: {
              dateTime: "2026-05-01T15:00:00.0000000",
              timeZone: "UTC",
            },
            end: {
              dateTime: "2026-05-01T16:00:00.0000000",
              timeZone: "UTC",
            },
            isAllDay: false,
          },
        ],
      },
    },
  },
};
