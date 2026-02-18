import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";

test.describe("checkAddActivity", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([
        activityMocks.getByIdSuccess,
        activityMocks.getAllActivityEmpty,
        projectMocks.getProjectsListSuccess,
        projectMocks.getDetailedProjectsByUserSuccess,
        userMocks.userMeSuccess,
        categoryMocks.getCategoriesByProjectSuccess,
      ])
      .apply();
    await page.clock.install({ time: new Date("2025-03-22T08:00:00") });
    await page.goto("http://localhost:5002/calendar");
    await page.waitForLoadState("networkidle");
  });

  test("ajouterUneActivite", async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker.addMocks([activityMocks.addActivitySuccess]).apply();

    await page.getByText("Nouvelle activité").click();
    await page.waitForTimeout(2000);

    await page.getByPlaceholder("Nom de l'activité...").fill("asd");
    await page.locator("#activity-description").fill("asd");
    await page.getByPlaceholder('Sélectionner un projet').click();
    await page.getByText(/^Projet sous-sol\b/).first().click();
    await page.locator("#activity-category-search").first().click();
    await page.locator(".category-item").first().click();
    await page.getByText("Créer").click();
    await page.waitForSelector(".fc-event", { state: "visible" });
    await expect(
      page.locator(".fc-event").getByText("Projet sous-sol")
    ).toBeVisible();
  });

  test("ajouterUneActiviteSansNomDescription", async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([activityMocks.addActivitySuccessNoNameNoDescription])
      .apply();

    await page.getByText("Nouvelle activité").click();
    await page.waitForTimeout(2000);

    await page.getByPlaceholder('Sélectionner un projet').click();
    await page.getByText(/^Projet sous-sol\b/).first().click();
    await page.locator("#activity-category-search").first().click();
    await page.locator(".category-item").first().click();
    await page.getByText("Créer").click();

    await expect(page.locator(".fc-event-title-container")).toBeVisible();
    await expect(
      page.locator(".fc-event").getByText("Projet sous-sol")
    ).toBeVisible();
  });

  test("ajouterUneActiviteCasErreur", async ({page}) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
        .addMocks([activityMocks.addActivitySuccessNoNameNoDescription])
        .apply();

    await page.getByText("Nouvelle activité").click();
    await page.waitForTimeout(2000);
    await page.getByText("Créer").click();
    await expect(page.locator("span:has-text('Veuillez sélectionner un projet')")).toBeVisible();

    await page.getByPlaceholder('Sélectionner un projet').click();
    await page.getByText(/^Projet sous-sol\b/).first().click();
    await expect(page.locator("span:has-text('Veuillez sélectionner un projet')")).not.toBeVisible();
    await page.getByText("Créer").click();

    await expect(page.locator(".fc-event-title-container")).toBeVisible();
    await expect(
        page.locator(".fc-event").getByText("Projet sous-sol")
    ).toBeVisible();
  })
});
