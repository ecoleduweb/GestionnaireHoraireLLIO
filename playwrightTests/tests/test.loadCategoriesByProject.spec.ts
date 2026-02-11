import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";

test.describe("loadCategoriesByProject", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([
        userMocks.userMeSuccess,
        projectMocks.getTwoProjectsSuccess,
        projectMocks.getDetailedProjectsByUserSuccess,
        activityMocks.getAllActivityEmpty,
        activityMocks.getByIdSuccess,
        categoryMocks.getCategoriesByProjectSuccess,
        categoryMocks.getCategoryForProject1Success,
        categoryMocks.getCategoryForProject2Success,
      ])
      .apply();
    await page.clock.install({ time: new Date("2025-03-22T08:00:00-04:00") });
    await page.goto("http://localhost:5002/calendar");
    await page.waitForLoadState("networkidle");
  });

  test("shouldCreateActivityWithExistingCategoryAndVerifyProjectChangeBehavior", async ({
    page,
  }) => {
    const apiMocker = new ApiMocker(page);

    // Mock pour la création d'une activité
    await apiMocker.addMocks([activityMocks.addActivitySuccess]).apply();

    // === ÉTAPE 1: Créer une nouvelle activité ===
    // Cliquer sur le bouton "Nouvelle activité"
    await page.getByText("Nouvelle activité").click();

    // Attendre que le projet sélect se remplisse avec les options
    await page.waitForTimeout(500);

    // Vérifier que les options de projet sont disponibles
    const projectSelect = page.locator("#activity-project");
    await projectSelect.waitFor({ state: "visible" });

    // Sélectionner le premier projet (id=1)
    await projectSelect.selectOption("1");

    // Ouvrir le dropdown des catégories
    await page.locator("#activity-category-search").first().click();

    // Sélectionner la catégorie "Design" qui existe déjà
    await page.click(".category-item:has-text('Design')");

    // Remplir les détails de l'activité
    await page
      .getByPlaceholder("Nom de l'activité...")
      .fill("Ma première activité");
    await page.locator("#activity-description").fill("Description test");

    // Soumettre le formulaire
    await page.getByText("Créer").click();

    // === ÉTAPE 2: Cliquer sur l'activité pour l'éditer ===
    // Attendre que l'événement apparaisse sur le calendrier
    await page.waitForSelector(".fc-event", { state: "visible" });

    // Cliquer sur l'événement pour ouvrir la modale d'édition
    await page.click(".fc-event");

    // === ÉTAPE 3: Vérifier que "Design" est sélectionné ===
    const selectedCategory = page.locator(".select-value");
    await expect(selectedCategory).toContainText("Design");

    // === ÉTAPE 4: Changer de projet vers le deuxième ===
    const projectSelectEdit = page.locator("#activity-project");
    await projectSelectEdit.selectOption("2");

    // === ÉTAPE 5: Vérifier que "Design" n'existe pas dans le projet 2 ===
    // Ouvrir le dropdown des catégories
    await page.locator("#activity-category-search").first().click();

    // Vérifier que "Design" n'existe pas
    await expect(page.getByText("Design")).not.toBeVisible();

    // Vérifier que seul "Développement" est affiché
    await expect(page.getByText("Développement")).toBeVisible();
  });
});
