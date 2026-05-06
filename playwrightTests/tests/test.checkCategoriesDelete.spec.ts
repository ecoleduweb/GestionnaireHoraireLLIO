import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";

test.describe("checkCategoriesDelete", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([
        projectMocks.getDetailedProjectsSuccess,
        activityMocks.getAllActivitySuccess,
        activityMocks.getAllActivitiesDefaultWeekSuccess,
        categoryMocks.getCategoriesByProjectSuccess,
        userMocks.userMeSuccess,
      ])
      .apply();
    await page.clock.install({ time: new Date("2025-03-22T08:00:00") });
    await page.goto("http://localhost:5002/calendar");
    await page.waitForLoadState("networkidle");
  });

  test("categoryDeleteSuccess", async ({ page }) => {
    // Configuration du mock pour la suppression réussie
    const mocker = new ApiMocker(page);
    await mocker
      .addMock(projectMocks.deleteProjectSuccess)
      .addMock(projectMocks.getDetailedProjectsSuccessAfterDelete)
      .addMock(userMocks.userMeSuccess)
      .addMock(categoryMocks.deleteCategorySuccess)
      .apply();

    await page.goto("http://localhost:5002/calendar");
    await page.waitForSelector(".fc-event", { state: "visible" });

    let activities = await page.locator(".fc-event").all();

    await activities[0].click();
    await page.waitForTimeout(2000);
    await page.getByRole("textbox", { name: "Catégorie *" }).click();

    page.on("dialog", async (dialog) => {
      expect(dialog.type()).toBe("confirm");
      expect(dialog.message()).toContain("Supprimer la catégorie Test ?");
      await dialog.accept();
    });

    await page.getByRole("button", { name: "Supprimer le projet" }).click();

    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([categoryMocks.getCategoriesByProjectDeletedSuccess])
      .apply();

    await page.getByRole("textbox", { name: "Catégorie *" }).click();

    await expect(
      page.getByRole("button", { name: "Supprimer le projet" }),
    ).toHaveCount(0, { timeout: 15000 });
  });

  test("deleteProjectError", async ({ page }) => {
    // Configuration du mock pour la suppression réussie
    const mocker = new ApiMocker(page);
    await mocker
      .addMock(projectMocks.deleteProjectSuccess)
      .addMock(projectMocks.getDetailedProjectsSuccessAfterDelete)
      .addMock(userMocks.userMeSuccess)
      .addMock(categoryMocks.deleteCategoryFailed)
      .apply();

    await page.goto("http://localhost:5002/calendar");
    await page.waitForSelector(".fc-event", { state: "visible" });

    let activities = await page.locator(".fc-event").all();

    await activities[0].click();
    await page.waitForTimeout(2000);
    await page.getByRole("textbox", { name: "Catégorie *" }).click();

    page.on("dialog", async (dialog) => {
      let type = dialog.type();

      if (type == "confirm") {
        expect(dialog.message()).toContain("Supprimer la catégorie Test ?");
        await dialog.accept();
      } else if (type == "alert") {
        expect(dialog.message()).toContain(
          "Erreur lors de la suppression de la catégorie",
        );
        await dialog.dismiss();
      }
    });

    await page.getByRole("button", { name: "Supprimer le projet" }).click();

    await page.getByRole("textbox", { name: "Catégorie *" }).click();

    await expect(
      page.getByRole("button", { name: "Supprimer le projet" }),
    ).toHaveCount(1, { timeout: 15000 });
  });
});
