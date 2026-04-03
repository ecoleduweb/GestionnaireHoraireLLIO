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

test("filtre les projets avec les dates de début et de fin via l'API", async ({ page }) => {
    const testStartDate = "2025-01-01";
    const testEndDate = "2025-01-31";

    const mockFilteredProject = projectMocks.getDetailedProjectsSuccess.response.json.projects[0];

    // Intercept la requête AVANT de modifier les dates
    await page.route('**/projects/detailed**', (route) => {
        const url = new URL(route.request().url());
        const from = url.searchParams.get('from');
        const to = url.searchParams.get('to');

        // Si les dates correspondent, retourner le projet filtré
        if (from === testStartDate && to === testEndDate) {
            route.fulfill({
                status: 200,
                json: { projects: [mockFilteredProject] }
            });
        } else {
            // Sinon, continuer avec le mock par défaut
            route.continue();
        }
    });

    await page.waitForSelector('[data-testid="project-item"]');

    await page.locator('#startDate').fill(testStartDate);
    await page.locator('#endDate').fill(testEndDate);
    await page.locator('#endDate').blur();

    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1, { timeout: 10000 });
    await expect(visibleProjects.first()).toContainText(mockFilteredProject.name);
});
});