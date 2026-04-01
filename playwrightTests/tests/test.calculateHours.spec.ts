import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";

test.describe("calcul d'heures", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    let timeInBankCallCount = 0;

    await page.clock.install({
      time: new Date("2025-03-22T08:00:00-04:00"),
    });

    await page.route("**/user/time-bank", async (route) => {
      timeInBankCallCount += 1;

      await route.fulfill({
        status: 200,
        json: {
          isConfigured: timeInBankCallCount > 1,
          timeInBank: timeInBankCallCount > 1 ? 40 : 0,
          textHoursWorked: "cette semaine",
        },
      });
    });

    await page.route("**/user/time-bank/config", async (route) => {
      const method = route.request().method();

      if (method === "GET") {
        await route.fulfill({
          status: 200,
          json: {
            startDate: "2025-03-01",
            hoursPerWeek: 40,
            offset: 0,
          },
        });
        return;
      }

      await route.fulfill({
        status: 200,
        json: {
          startDate: "2025-03-01",
          hoursPerWeek: 40,
          offset: 0,
        },
      });
    });

    await apiMocker
      .addMocks([
        userMocks.userMeSuccess,
        activityMocks.getAllActivityEmpty,
        projectMocks.getProjectsListSuccess,
        projectMocks.getDetailedProjectsByUserSuccess,
        categoryMocks.getCategoriesByProjectSuccess,
      ])
      .apply();

    await page.goto("http://localhost:5002/calendar");
  });

  test("calcul d'heures success", async ({ page }) => {
    let savedConfig: Record<string, unknown> | null = null;

    await page.unroute("**/user/time-bank/config");
    await page.route("**/user/time-bank/config", async (route) => {
      const method = route.request().method();

      if (method === "PUT") {
        savedConfig = route.request().postDataJSON() as Record<string, unknown>;
      }

      await route.fulfill({
        status: 200,
        json: {
          startDate: "2025-03-01",
          hoursPerWeek: 40,
          offset: 0,
        },
      });
    });

    await page.waitForLoadState("networkidle");
    await expect(page.getByTestId("total-hours")).toHaveText("0");

    await page.getByTestId("total-hours").click();

    await page.fill('[name="startDate"]', "2025-03-01");
    await page.fill('[name="hoursPerWeek"]', "40");

    await page.getByRole("button", { name: "Enregistrer" }).click();

    await expect.poll(() => savedConfig).toEqual({
      startDate: "2025-03-01",
      hoursPerWeek: 40,
      offset: 0,
    });
    await expect(page.getByTestId("total-hours")).toHaveText("40");
  });
});
