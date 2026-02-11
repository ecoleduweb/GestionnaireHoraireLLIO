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

test.describe("Recherche de projets", () => {
  test("filtre les projets en temps réel", async ({ page }) => {
    // 1️⃣ Attendre que les projets soient affichés
    await page.waitForSelector('[data-testid="project-item"]');

    // 2️⃣ Vérifier l’état initial
    await expect(page.locator('[data-testid="project-item"]')).toHaveCount(8);

    // 3️⃣ Taper dans la recherche
    const searchInput = page.getByTestId("project-search");
    await searchInput.fill("AT-123");

    // 4️⃣ Vérifier que seul le bon projet reste visible
    const visibleProjects = page.locator('[data-testid="project-item"]');
    await expect(visibleProjects).toHaveCount(1);
    await expect(visibleProjects.first()).toContainText("AT-123");
  });
});
