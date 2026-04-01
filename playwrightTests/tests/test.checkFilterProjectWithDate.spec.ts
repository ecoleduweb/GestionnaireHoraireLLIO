import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.beforeEach(async ({ page }) => {
  const apiMocker = new ApiMocker(page);
  await apiMocker
    .addMocks([
      projectMocks.getDetailedProjectsSuccess,
      userMocks.userMeSuccess,
    ])
    .apply();

  await page.clock.install({ time: new Date("2025-03-22T08:00:00") });
  await page.goto("http://localhost:5002/projects");
  await page.waitForLoadState("networkidle");
});

test.describe("Recherche et filtrage de projets", () => {
  
  test("filtre les projets par texte en temps réel", async ({ page }) => {
    await page.waitForSelector('[data-testid="project-item"]');
    await expect(page.locator('[data-testid="project-item"]')).toHaveCount(8);

    const searchInput = page.getByTestId("project-search");
    await searchInput.fill("AT-123");

    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1);
    await expect(visibleProjects.first()).toContainText("AT-123");
  });

test("filtre les projets avec les dates de début et de fin via l'API", async ({ page }) => {
    const testStartDate = "2025-01-01";
    const testEndDate = "2025-01-31";

    // Étape 1 : On cible UNIQUEMENT la requête avec les deux paramètres.
    // Les requêtes intermédiaires seront gérées silencieusement par le ApiMocker du beforeEach.
    await page.route(`**/projects/detailed?from=${testStartDate}&to=${testEndDate}`, async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          projects: [
            {
              id: 999,
              uniqueId: "JAN-2025",
              name: "Projet de Janvier",
              status: 1,
              managerId: 1,
              coLeads: []
            }
          ]
        })
      });
    });

    // Étape 2 : Attendre que l'interface initiale (les 8 projets du mock par défaut) soit chargée
    await page.waitForSelector('[data-testid="project-item"]');

    // Étape 3 : Remplir les dates
    await page.locator('#startDate').fill(testStartDate);
    await page.locator('#endDate').fill(testEndDate);
    
    // Étape 4 : Forcer la perte de focus du dernier champ pour garantir le déclenchement du onchange Svelte
    await page.locator('#endDate').blur();

    // Étape 5 : Vérification (sans cliquer sur un bouton inutile)
    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1, { timeout: 10000 });
    await expect(visibleProjects.first()).toContainText("JAN-2025");
    await expect(visibleProjects.first()).toContainText("Projet de Janvier");
  });
});
