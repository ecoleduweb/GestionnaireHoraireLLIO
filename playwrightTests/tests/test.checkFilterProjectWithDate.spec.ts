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

  test("filtre les projets avec les dates de début et de fin via l'API", async ({ page }) => {
    const testStartDate = "2025-01-01";
    const testEndDate = "2025-01-31";
    const mockFilteredProject = projectMocks.getDetailedProjectsSuccess.response.json.projects[0];

    await page.route('**/**/projects/detailed*', async (route) => {
      const url = new URL(route.request().url());
      const from = url.searchParams.get('from');
      const to = url.searchParams.get('to');

      if (from === testStartDate && to === testEndDate) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ projects: [mockFilteredProject] })
        });
      } else {
        await route.continue();
      }
    });

    await page.waitForSelector('[data-testid="project-item"]');

    await page.locator('#startDate').fill(testStartDate);
    await page.locator('#endDate').fill(testEndDate);
    
    // On prépare l'écouteur de réponse AVANT de déclencher l'action
    const responsePromise = page.waitForResponse(response => {
      const url = new URL(response.url());
      return url.pathname.includes('projects/detailed') && 
             url.searchParams.get('from') === testStartDate &&
             url.searchParams.get('to') === testEndDate &&
             response.status() === 200;
    });

    // Déclencher l'événement qui provoque l'appel réseau
    await page.locator('#endDate').dispatchEvent('change');

    // Attendre que la promesse soit résolue
    await responsePromise;

    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1);
    await expect(visibleProjects.first()).toContainText(mockFilteredProject.name);
  });
});