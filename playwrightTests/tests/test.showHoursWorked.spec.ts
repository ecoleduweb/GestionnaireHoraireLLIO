import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("showHoursWorked", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker.addMocks([userMocks.userMeSuccess]).apply();
    await page.clock.install({ time: new Date("2025-03-22T08:00:00-04:00") });
  });

  test("showHoursWorkedDefault", async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([activityMocks.getAllActivitiesDefaultWeekSuccess])
      .apply();
    await page.goto("http://localhost:5002/calendar");
    await page.waitForSelector(".fc-event", { state: "visible" });
    await expect(page.locator(".bilan-container")).toContainText("8 heures cette semaine.");
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
    await expect(page.locator(".bilan-container")).toContainText("21 heures ce mois-ci.");
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
    await expect(page.locator(".bilan-container")).toContainText("7 heures aujourd'hui.");
  });

  test("hoursWorkedNoActivities", async ({ page }) => {
  const apiMocker = new ApiMocker(page);
  await apiMocker.addMocks([activityMocks.getAllActivityEmpty]).apply();

  await page.goto("http://localhost:5002/calendar");

  await page.waitForSelector(".fc", { state: "visible" });

  await expect(page.locator(".bilan-container")).toBeVisible();
  await expect(page.locator(".bilan-container")).toContainText("0 heures cette semaine.");
});
});
