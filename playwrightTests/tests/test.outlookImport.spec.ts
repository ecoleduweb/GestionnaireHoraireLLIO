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

  test("Import d'une journée sans évènements", async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker.addMocks([outlookMocks.getEventsNoEvent]).apply();
    // Load la page et fait la requête de base
    await page.goto("http://localhost:5002/calendar");
    await page.waitForLoadState("networkidle");

    // On va chercher le bouton d'import
    const button = page.locator(`.import-outlook-btn[data-date="2025-03-22"]`);

    // Vérification de l'état original du bouton
    await expect(button).toBeVisible();
    await expect(button).toHaveText("+ Outlook");
    await expect(button).not.toBeDisabled();

    // Clic
    await button.click();

    // Vérification du modal
    const modal = page.locator('.modal-overlay');
    await expect(modal).toBeVisible();
    await expect(modal.locator('.modal-title')).toHaveText('Alerte');
    await expect(modal).toContainText('Aucun évènement à importer pour la journée sélectionnée.');

    // Fermeture du modal
    await modal.locator('button', { hasText: 'Fermer' }).click();
    await expect(modal).not.toBeVisible();
  });

  test("outlookFailed", async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker.addMocks([outlookMocks.outlookFail]).apply();

    const dialogPromise = page.waitForEvent("dialog");

    page.getByLabel("22 mars").getByRole("button", { name: "Outlook" }).click();
    const dialog = await dialogPromise;

    expect(dialog.type()).toBe("alert");
    expect(dialog.message()).toContain(
      "Votre connexion à Outlook a expiré. Vous serez redirigés à la page de connexion pour vous reconnecter. Vous pourrez alors essayer d'importer vos évènements à nouveau.",
    );
    await dialog.dismiss();

    await page.waitForURL("http://localhost:5002/");
  });

  test("AddActivityFromOutlookWithProjetId", async ({ page }) => {
    const mocker = new ApiMocker(page);
    await mocker
      .addMocks([
        outlookMocks.getActivityWithProjectIdFromOutlook,
        userMocks.userMeSuccess,
        activityMocks.addActivitySuccess,
      ])
      .apply();

    await page.reload();
    await page.waitForLoadState("networkidle");
    await page
      .locator('[data-date="2025-03-22"] button:has-text("+ Outlook")')
      .click();

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
