import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { outlookMocks } from "../Helper/Mocks/outlook.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";

test.describe("ImportOutlook", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([
        projectMocks.getDetailedProjectsByUserSuccess,
        categoryMocks.getCategoryForProject1Success,
        projectMocks.getProjectsListSuccess,
        userMocks.userMeSuccess,
        activityMocks.getAllActivityEmpty,
        userMocks.getTimeBankSuccess,
        userMocks.getTimeBankConfigSuccess,
      ])
      .apply();
    await page.clock.install({ time: new Date("2025-03-22T08:00:00-04:00") });
    await page.goto("http://localhost:5002/calendar");
    await page.waitForLoadState("networkidle");
  });
  test("AddActivityFromOutlookWithoutProjetId", async ({ page }) => {
    const mocker = new ApiMocker(page);
    await mocker
      .addMocks([
        outlookMocks.getActivityWithoutProjectIdFromOutlook,
        userMocks.userMeSuccess,
        activityMocks.addActivitySuccess,
      ])
      .apply();

    await page.reload();
    await page.waitForLoadState("networkidle");
    await page
      .locator('[data-date="2025-03-22"] button:has-text("+ Outlook")')
      .click();
    await expect(page.getByPlaceholder("Sélectionner un projet")).toBeVisible();
    await page.getByRole("textbox", { name: "Projet *" }).click();
    await page
      .getByText("Projet sous-sol | Nommer le projet", { exact: true })
      .click();

    await expect(
      page.getByRole("textbox", { name: "Nom (optionnel)" }),
    ).toHaveValue("Test sans Id Projet");

    const activityRequest = page.waitForRequest(
      (request) =>
        request.url().includes("/activity") && request.method() === "POST",
    );

    await page.getByRole("button", { name: "Importer", exact: true }).click();

    const request = await activityRequest;
    expect(request.method()).toBe("POST");

    await expect(
      page.getByText("Importation des évènements Outlook"),
    ).not.toBeVisible();
  });
});
