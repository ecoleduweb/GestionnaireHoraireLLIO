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

    await expect(page.locator('[data-testid="project-item"]').first()).toBeVisible();
  });

  test("filtre les projets avec les dates de début et de fin via l'API", async ({ page }) => {
    const testStartDate = "2025-01-01";
    const testEndDate = "2025-01-31";

    const dynamicMocker = new ApiMocker(page);
    await dynamicMocker
      .addMock(projectMocks.getDetailedProjectsFilteredSuccess)
      .apply();

    const mockFilteredProject =
      projectMocks.getDetailedProjectsFilteredSuccess.response.json.projects[0];

    await page.locator("#startDate").fill(testStartDate);
    await page.locator("#endDate").fill(testEndDate);

    const responsePromise = page.waitForResponse((response) => {
      const url = response.url();
      return (
        url.includes("projects/detailed") &&
        url.includes(`from=${testStartDate}`) &&
        url.includes(`to=${testEndDate}`) &&
        response.status() === 200
      );
    });

    await page.locator("#endDate").blur();
    await responsePromise;

    const visibleProjects = page.locator('[data-testid="project-item"]');

    await expect(visibleProjects).toHaveCount(1);
    await expect(visibleProjects.first()).toContainText(mockFilteredProject.name);
  });
});