import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";

test.describe("showConfirmationModal", () => {
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

  test("fermerModaleViaArrierePlanSansModification", async ({ page }) => {
    await page.getByText("Nouvelle activité").click();
    await expect(page.getByRole("heading", { name: "Nouvelle activité" })).toBeVisible();

    await page.mouse.click(600, 300);

    // Vérification
    await expect(page.getByRole("heading", { name: "Nouvelle activité" })).not.toBeVisible();
  });

  test("afficherAvertissementFermetureViaArrierePlanAvecModification", async ({ page }) => {
    await page.getByText("Nouvelle activité").click();
    await expect(page.getByRole("heading", { name: "Nouvelle activité" })).toBeVisible();

    await page.getByPlaceholder("Sélectionner un projet").click();
    await page.getByText(/^Projet sous-sol\b/).first().click();

    await page.waitForTimeout(500);

    await page.mouse.click(600, 300);

    const modalTitle = page.locator("h2").filter({ hasText: "Modifications non enregistrées" }).first();
    await expect(modalTitle).toBeVisible();

    await page.getByRole("button", { name: "Non, rester" }).first().click();
    await expect(modalTitle).not.toBeVisible();
    await expect(page.getByRole("heading", { name: "Nouvelle activité" })).toBeVisible();

    await page.waitForTimeout(300);
    await page.mouse.click(600, 300);
    await expect(modalTitle).toBeVisible();

    await page.getByRole("button", { name: "Oui, quitter" }).first().click();

    await expect(modalTitle).not.toBeVisible();
    await expect(page.getByRole("heading", { name: "Nouvelle activité" })).not.toBeVisible();
  });
});