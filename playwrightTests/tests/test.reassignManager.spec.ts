import { expect, test } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("Vérification de la réattribution du chargé de projet", () => {
  test.beforeEach(async ({ page }) => {
    await new ApiMocker(page)
      .addMocks([
        projectMocks.getDetailedProjectsSuccess,
        userMocks.getAllManagersSuccess,
        userMocks.userMeSuccess,
      ])
      .apply();

    await page.goto("http://localhost:5002/projects");
    await page.waitForLoadState("networkidle");
  });

  test("Réattribution du chargé de projet réussie", async ({ page }) => {
    const mocker = new ApiMocker(page);
    await mocker.addMock(projectMocks.reassignManagerSuccess).apply();

    const firstProject = page.getByTestId("project-item").first();
    await firstProject.getByText("Réattribuer").click();

    const heading = page.getByRole("heading", {
      name: "Réattribuer un chargé de projet",
    });
    await expect(heading).toBeVisible();

    const modal = heading
      .locator('xpath=ancestor::*[self::div][contains(@class,"shadow-xl")]')
      .first();

    const select = modal.locator("select#userId");
    await expect(select).toBeVisible();
    await expect(select.locator("option")).not.toHaveCount(0);
    await page.getByLabel("Utilisateur").selectOption("3");

    const confirmButton = modal.getByRole("button", { name: "Ajouter" });
    await expect(confirmButton).toBeEnabled();

    await mocker.clearMocks();
    await mocker
      .addMocks([
        projectMocks.getDetailedProjectsSuccess,
        projectMocks.reassignManagerSuccess,
        userMocks.getAllManagersSuccess,
        userMocks.userMeSuccess,
      ])
      .apply();

    await confirmButton.click();

    await expect(heading).not.toBeVisible();
  });
});
