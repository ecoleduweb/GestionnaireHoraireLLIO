import { expect, test } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("Modification du taux horaire", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker
      .addMocks([
        projectMocks.getDetailedProjectsSuccess,
        userMocks.userMeSuccess,
      ])
      .apply();

    await page.goto("http://localhost:5002/projects");
    await page.waitForLoadState("domcontentloaded");
  });

  test("Affichage du taux horaire quand il est défini", async ({ page }) => {
    const firstProject = page.getByTestId("project-item").first();

    const employeeRow = firstProject
      .locator(".grid.grid-cols-6")
      .filter({ hasText: "Katell Arnault de la Ménardière" })
      .first();

    const rateCell = employeeRow.getByTestId("employee-hourly-rate");
    await expect(rateCell).toBeVisible();

    await expect(rateCell).toHaveText("85$/h");
  });

  test("Affichage d'un tiret quand le taux horaire est absent", async ({
    page,
  }) => {
    const firstProject = page.getByTestId("project-item").first();

    const employeeRow = firstProject
      .locator(".grid.grid-cols-6")
      .filter({ hasText: "Jean-François Jasmin" })
      .first();

    await employeeRow.click();

    await expect(employeeRow.getByTestId("employee-hourly-rate")).toHaveText(
      "-",
    );
  });

  test("Modification réussie du taux horaire", async ({ page }) => {
    const apiMocker = new ApiMocker(page);

    await apiMocker.addMock(projectMocks.updateEmployeeRateSuccess).apply();

    const firstProject = page.getByTestId("project-item").first();

    const firstEmployeeRow = firstProject
      .locator(".grid.grid-cols-6")
      .filter({ hasText: "Katell Arnault de la Ménardière" })
      .first();

    await firstEmployeeRow.click();

    const editButton = firstEmployeeRow.getByLabel("Modifier le taux horaire");
    await expect(editButton).toBeVisible();

    await editButton.click();

    const modal = page.locator("div.shadow-xl").first();

    const input = page.getByPlaceholder("Nouveau taux horaire");

    await expect(input).toBeVisible();
    await input.fill("75");
    await expect(input).toHaveValue("75");

    await page.getByRole("button", { name: "Enregistrer" }).click();

    await expect(modal).not.toBeVisible({ timeout: 5000 });
  });
  test("Affichage du total des taux horaires", async ({ page }) => {
    const firstProject = page.getByTestId("project-item").first();

    await expect(firstProject.getByTestId("total-hourly-rate")).toHaveText(
      "85$/h",
    );
  });
});
