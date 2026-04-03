import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("Recherche et filtrage de projets", () => {
  
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([
        projectMocks.getDetailedProjectsSuccess,
        userMocks.userMeSuccess,
      ])
      .apply();

    await page.clock.install({ time: new Date("2025-03-22T08:00:00Z") });
    await page.goto("http://localhost:5002/projects");
    await page.waitForLoadState("networkidle");
  });

  test.afterEach(async ({ page }) => {
    // Nettoyer les routes interceptées pour éviter les interférences
    await page.unroute("**/*");
  });

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

    const mockFilteredProject = projectMocks.getDetailedProjectsSuccess.response.json.projects[0];

    // Créer un nouveau mocker pour cette route dynamique
    const dynamicMocker = new ApiMocker(page);
    await dynamicMocker.addMock({
      url: `/projects/detailed?from=${testStartDate}&to=${testEndDate}`,
      method: "GET",
      response: {
        status: 200,
        json: {
          projects: [mockFilteredProject]
        }
      }
    }).apply();

    await page.waitForSelector('[data-testid="project-item"]');

    // Remplir les dates et déclencher la recherche
    await page.locator('#startDate').fill(testStartDate);
    await page.locator('#endDate').fill(testEndDate);
    
    // Attendre la réponse API avant de vérifier le résultat
    await Promise.all([
      page.waitForResponse(response => 
        response.url().includes(`/projects/detailed?from=${testStartDate}&to=${testEndDate}`)
      ),
      page.locator('#endDate').blur()
    ]);

    // Vérifier que les projets filtrés s'affichent
    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1, { timeout: 10000 });
    await expect(visibleProjects.first()).toContainText(mockFilteredProject.name);
  });
});