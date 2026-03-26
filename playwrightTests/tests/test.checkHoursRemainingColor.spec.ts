import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";

test.describe("checkRemainingHoursColor", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([
        userMocks.userMeSuccess,
        activityMocks.getAllActivitiesDefaultWeekSuccess,
        projectMocks.getProjectsListSuccess,
        projectMocks.getDetailedProjectsByUserSuccess,
      ])
      .apply();
    await page.clock.install({ time: new Date("2025-03-22T08:00:00") });
    await page.goto("http://localhost:5002/calendar");
    await page.waitForLoadState("networkidle");
  });

  test("negativeRemainingHours", async ({ page }) => {
    await page.getByRole("button", { name: "new-1 | new commut" }).click();
    await page.waitForSelector(".text-red-700", { state: "visible" });

    const redValues = page
      .locator(".text-red-700")
      .filter({ hasText: "-4h00" });
    await expect(redValues).toHaveCount(2);
    await expect(redValues.first()).toBeVisible();
  });
  test("positiveRemainingHours", async ({ page }) => {
    await page.getByRole("button", { name: "migr-2 | projet ! apre" }).click();
    await page.waitForSelector(".text-gray-700", { state: "visible" });
    const greyValues = page
      .locator(".text-gray-700")
      .filter({ hasText: "-1h00" });
    await expect(greyValues).toHaveCount(2);
    await expect(greyValues.first()).toBeVisible();
  });

  test("noEstimatedHours", async ({ page }) => {
    await page.getByRole("button", { name: "-mmm | le nom !" }).click();
    await page.waitForSelector(".text-gray-700", { state: "visible" });
    const greyValues = page
      .locator(".text-gray-700")
      .filter({ hasText: "-2h30" });
    await expect(greyValues).toHaveCount(2);
    await expect(greyValues.first()).toBeVisible();
  });

  test("noRemainingHasEstimate", async ({ page }) => {
    await page
      .getByRole("button", { name: "10htotal | 10 heures au total" })
      .click();
    await page.waitForSelector(".text-gray-700", { state: "visible" });
    const NoneValues = page.locator(".text-gray-400").filter({ hasText: "-" });
    await expect(NoneValues).toHaveCount(2);
    await expect(NoneValues.first()).toBeVisible();
  });
});
