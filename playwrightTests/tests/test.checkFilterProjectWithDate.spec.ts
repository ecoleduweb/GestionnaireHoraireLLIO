import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.beforeEach(async ({ page }) => {
  const apiMocker = new ApiMocker(page);
  // Charge les mocks par défaut pour l'état initial de la page
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
    // 1. Définir les dates de test
    const testStartDate = "2025-01-01";
    const testEndDate = "2025-01-31";

    // 2. Intercepter l'appel API spécifique généré par le changement de dates
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
              status: 1, // En cours
              managerId: 1,
              coLeads: []
            }
          ]
        })
      });
    });

    // 3. Attendre le chargement initial
    await page.waitForSelector('[data-testid="project-item"]');

    // 4. Remplir les champs de date (déclenche le onchange)
    await page.locator('#startDate').fill(testStartDate);
    await page.locator('#endDate').fill(testEndDate);

    // 5. Attendre que la requête API avec les dates soit interceptée et résolue
    await page.waitForResponse(response => 
      response.url().includes(`from=${testStartDate}`) && 
      response.url().includes(`to=${testEndDate}`) && 
      response.status() === 200
    );

    // 6. Vérifier que l'interface a bien été mise à jour avec la donnée mockée
    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1);
    await expect(visibleProjects.first()).toContainText("JAN-2025");
    await expect(visibleProjects.first()).toContainText("Projet de Janvier");
  });

});