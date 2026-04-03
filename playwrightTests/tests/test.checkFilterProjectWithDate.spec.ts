import { test, expect } from "@playwright/test";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("Filtrage par date isolé", () => {
  
  test("filtre les projets via API en utilisant les mocks", async ({ page }) => {
    const testStartDate = "2025-01-01";
    const testEndDate = "2025-01-31";
    
    // On extrait le tout premier projet de ton mock officiel pour simuler le résultat de la recherche
    const mockProject = projectMocks.getDetailedProjectsSuccess.response.json.projects[0];

    // 1. Mock de l'utilisateur (nécessaire pour accéder à la page)
    await page.route("**/user/me", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(userMocks.userMeSuccess.response.json),
      });
    });

    // 2. Mock de la requête avec les paramètres de date
    await page.route(`**/projects/detailed?from=${testStartDate}&to=${testEndDate}`, async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ projects: [mockProject] }),
      });
    });

    // 3. Mock du chargement initial (URL stricte sans paramètres)
    await page.route("**/projects/detailed", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(projectMocks.getDetailedProjectsSuccess.response.json),
      });
    });

    // 4. Exécution du test
    await page.goto("http://localhost:5002/projects");

    // Attente des 8 projets de base
    await page.waitForSelector('[data-testid="project-item"]');

    // Remplissage des champs de date
    await page.locator('#startDate').fill(testStartDate);
    await page.locator('#endDate').fill(testEndDate);
    await page.locator('#endDate').blur(); // Force le déclenchement de onchange

    // Vérification stricte
    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1, { timeout: 5000 });
    await expect(visibleProjects.first()).toContainText(mockProject.name);
  });

});