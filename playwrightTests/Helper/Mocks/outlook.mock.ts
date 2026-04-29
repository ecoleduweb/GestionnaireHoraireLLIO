import { MockConfig } from "../types";

export const outlookMocks = {
  getEventsFromOutlookSuccess: {
    url: "/activities/me/outlook?date=2026-05-01",
    method: "GET",
    response: {
      status: 200,
      json: {
        date: "2026-05-01",
        events: [
          {
            id: "AAMkAGQ2MDNhYjVjLWFkMWMtNGNmZi05ZmVjLTdmZmQ4NjgwMjdkNwBGAAAAAABXx3nPK_GjQbVWwm_t6k5tBwAmsLVGJoOjQZrgmFxGm6c6AAAAAAENAAAmsLVGJoOjQZrgmFxGm6c6AADJTqHnAAA=",
            subject: "AR145",
            body: { content: "\r\n" },
            start: { dateTime: "2026-05-01T12:00:00.0000000", timeZone: "UTC" },
            end: { dateTime: "2026-05-01T12:30:00.0000000", timeZone: "UTC" },
            isAllDay: false,
          },
          {
            id: "AAMkAGQ2MDNhYjVjLWFkMWMtNGNmZi05ZmVjLTdmZmQ4NjgwMjdkNwBGAAAAAABXx3nPK_GjQbVWwm_t6k5tBwAmsLVGJoOjQZrgmFxGm6c6AAAAAAENAAAmsLVGJoOjQZrgmFxGm6c6AADJTqHoAAA=",
            subject: "ahahaha",
            body: { content: "\r\n" },
            start: { dateTime: "2026-05-01T14:30:00.0000000", timeZone: "UTC" },
            end: { dateTime: "2026-05-01T17:30:00.0000000", timeZone: "UTC" },
            isAllDay: false,
          },
        ],
      },
    },
  },
};
