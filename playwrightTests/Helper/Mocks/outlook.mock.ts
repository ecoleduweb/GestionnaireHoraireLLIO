import { MockConfig } from "../types";

export const outlookMocks = {
  getActivityWithoutProjectIdFromOutlook: {
    url: "/activities/me/outlook?date=2025-03-22",
    method: "GET",
    response: {
      status: 200,
      json: {
        date: "2025-03-22",
        events: [
          {
            id: "AAMkAGQ2MDNhYjVjLWFkMWMtNGNmZi05ZmVjLTdmZmQ4NjgwMjdkNwBGAAAAAABXx3nPK_GjQbVWwm_t6k5tBwAmsLVGJoOjQZrgmFxGm6c6AAAAAAENAAAmsLVGJoOjQZrgmFxGm6c6AADJTqHnAAA=",
            subject: "Test sans Id Projet",
            body: { content: "\r\n" },
            start: {
              dateTime: "2025-03-22T15:00:00.0000000",
              timeZone: "UTC",
            },
            end: {
              dateTime: "2025-03-22T16:00:00.0000000",
              timeZone: "UTC",
            },
            isAllDay: false,
          },
        ],
      },
    },
  },
  getEventsNoEvent: {
    url: "/activities/me/outlook?date=2025-03-22",
    method: "GET",
    response: {
      status: 200,
      json: { date: "2025-03-22", events: [] },
    },
  },
  outlookFail: {
    url: "/activities/me/outlook?date=*",
    method: "GET",
    response: {
      status: 401,
      json: {
        code: "GRAPH_EXPIRED",
        error:
          "Votre compte n'est pas connecté à votre compte Microsoft. Veuillez vous reconnecter et réessayer.",
      },
    },
  },
  getEventsFromOutlookSuccess: {
    url: "/activities/me/outlook?date=2025-03-22",
    method: "GET",
    response: {
      status: 200,
      json: {
        date: "2025-03-22",
        events: [
          {
            id: "AAMkAGQ2MDNhYjVjLWFkMWMtNGNmZi05ZmVjLTdmZmQ4NjgwMjdkNwBGAAAAAABXx3nPK_GjQbVWwm_t6k5tBwAmsLVGJoOjQZrgmFxGm6c6AAAAAAENAAAmsLVGJoOjQZrgmFxGm6c6AADJTqHnAAA=",
            subject: "AR145",
            body: { content: "\r\n" },
            start: { dateTime: "2025-03-22T12:00:00.0000000", timeZone: "UTC" },
            end: { dateTime: "2025-03-22T12:30:00.0000000", timeZone: "UTC" },
            isAllDay: false,
          },
          {
            id: "AAMkAGQ2MDNhYjVjLWFkMWMtNGNmZi05ZmVjLTdmZmQ4NjgwMjdkNwBGAAAAAABXx3nPK_GjQbVWwm_t6k5tBwAmsLVGJoOjQZrgmFxGm6c6AAAAAAENAAAmsLVGJoOjQZrgmFxGm6c6AADJTqHoAAA=",
            subject: "ahahaha",
            body: { content: "\r\n" },
            start: { dateTime: "2025-03-22T14:30:00.0000000", timeZone: "UTC" },
            end: { dateTime: "2025-03-22T17:30:00.0000000", timeZone: "UTC" },
            isAllDay: false,
          },
        ],
      },
    },
    getEventsNoEvent: {
      url: "/activities/me/outlook?date=2025-03-22",
      method: "GET",
      response: {
        status: 200,
        json: { date: "2025-03-22", events: [] },
      },
    },
    outlookFail: {
      url: "/activities/me/outlook?date=*",
      method: "GET",
      response: {
        status: 401,
        json: {
          code: "GRAPH_EXPIRED",
          error:
            "Votre compte n'est pas connecté à votre compte Microsoft. Veuillez vous reconnecter et réessayer.",
        },
      },
    },
  },
};
