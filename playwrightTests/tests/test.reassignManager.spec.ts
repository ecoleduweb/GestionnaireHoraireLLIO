import { expect, test } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("Réattribution du chargé de projet", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([
        projectMocks.getDetailedProjectsSuccess,
        projectMocks.getAvailableManagersSuccess,
        userMocks.userMeSuccess,
      ])
      .apply();

    await page.clock.install({ time: new Date("2025-03-22T08:00:00") });

    await page.goto("http://localhost:5002/projects");

    await page.waitForLoadState("domcontentloaded");
  });

  test("Réattribution réussie", async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker.addMock(projectMocks.reassignManagerSuccess).apply();

    const firstProject = page.getByTestId("project-item").first();
    await firstProject.getByRole("button", { name: "Réattribuer" }).click();
    const modal = page.locator("div.shadow-xl").first();

    const heading = modal.getByRole("heading", {
      name: "Réattribuer un chargé de projet",
    });

    await expect(heading).toBeVisible();

    const select = modal.locator("select#userId");
    await expect(select).toBeVisible();

    await expect
      .poll(async () => {
        return await select.locator("option").count();
      })
      .toBeGreaterThan(1);

    const options = select.locator("option");
    const value = await options.nth(1).getAttribute("value");

    if (!value) throw new Error("Option sans value");

    await select.selectOption(value);

    await expect(select).toHaveValue(value);

    const confirmButton = modal.getByRole("button", { name: "Ajouter" });
    await expect(confirmButton).toBeEnabled();
    await confirmButton.click();

    await expect(modal).not.toBeVisible({ timeout: 5000 });
  });
});
