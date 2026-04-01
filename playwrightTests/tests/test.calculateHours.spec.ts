import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";

test.describe("calcul d'heures", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([
        userMocks.userMeSuccess,
        activityMocks.getAllActivitiesDefaultWeekSuccess,
      ])
      .apply();

    await page.clock.install({
      time: new Date("2025-03-22T08:00:00-04:00"),
    });

    await page.goto("http://localhost:5002/calendar");
  });

  test("calcul d'heures success", async ({ page }) => {
    await page.waitForLoadState("networkidle");

    const totalHours = page.getByTestId("total-hours");
    const configButton = page.getByText("configurer");

    if (await totalHours.count()) {
      await totalHours.click();
    } else {
      await configButton.click();
    }

   
    await expect(
      page.getByText("Configuration des heures en banque")
    ).toBeVisible();

    
    await page.locator('input[name="startDate"]').fill("2025-03-01");
    await page.locator('input[name="hoursPerWeek"]').fill("35");
    await page.locator('input[name="offset"]').fill("5");

  
    await page.getByRole("button", { name: "Enregistrer" }).click();


  
    await expect(page.locator(".card")).toBeVisible();
  });
});