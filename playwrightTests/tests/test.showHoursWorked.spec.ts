import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("showHoursWorked", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([
        userMocks.userMeSuccess,
        userMocks.getTimeBankConfigNotConfigured,
        userMocks.getTimeBankNotConfiguredSuccess,
      ])
      .apply();

    await page.clock.install({ time: new Date("2025-03-22T08:00:00-04:00") });
  });

  test("showHoursWorkedDefault", async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([activityMocks.getAllActivitiesDefaultWeekSuccess])
      .apply();

    await page.goto("http://localhost:5002/calendar");

    await page.waitForSelector(".fc-event", { state: "visible" });

    
    await expect(page.getByRole("button", { name: "Configurer votre banque d'heures", exact: true })).toBeVisible();
  });

  test("showHoursWorkedMonth", async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([
        activityMocks.getAllActivitiesMonthSuccess,
        activityMocks.getAllActivitiesDefaultWeekSuccess,
        activityMocks.getAllActivitiesDaySuccess,
      ])
      .apply();

    await page.goto("http://localhost:5002/calendar");

    await page.waitForSelector(".fc-event", { state: "visible" });

    await page.getByRole("button", { name: "Mois", exact: true }).click();

   
    await expect(page.getByRole("button", { name: "Configurer votre banque d'heures", exact: true })).toBeVisible();
  });

  test("showHoursWorkedDay", async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([
        activityMocks.getAllActivitiesDaySuccess,
        activityMocks.getAllActivitiesDefaultWeekSuccess,
        activityMocks.getAllActivitiesMonthSuccess,
      ])
      .apply();

    await page.goto("http://localhost:5002/calendar");

    await page.waitForSelector(".fc-event", { state: "visible" });

    await page.getByRole("button", { name: "Jour", exact: true }).click();

    
    await expect(page.getByRole("button", { name: "Configurer votre banque d'heures", exact: true })).toBeVisible();
  });

  test("hoursWorkedNoActivities", async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([activityMocks.getAllActivitiesEmptyDefaultWeekSuccess])
      .apply();

    await page.goto("http://localhost:5002/calendar");

    await page.waitForSelector(".fc", { state: "visible" });

    
    await expect(page.getByRole("button", { name: "Configurer votre banque d'heures", exact: true })).toBeVisible();
  });
});
