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

    const dynamicMocker = new ApiMocker(page);
    await dynamicMocker.addMock(projectMocks.getDetailedProjectsFilteredSuccess).apply();

    await page.waitForSelector('[data-testid="project-item"]');

    await page.locator('#startDate').fill(testStartDate);
    await page.locator('#endDate').fill(testEndDate);
    
    await page.locator('#endDate').dispatchEvent('change');

    await page.waitForResponse(
      response => 
        response.url().includes('projects/detailed') && 
        response.url().includes('from=') &&
        response.status() === 200
    );

    const mockFilteredProject = projectMocks.getDetailedProjectsFilteredSuccess.response.json.projects[0];
    const visibleProjects = page.locator('[data-testid="project-item"]');
    
    await expect(visibleProjects).toHaveCount(1);
    await expect(visibleProjects.first()).toContainText(mockFilteredProject.name);
  });
});