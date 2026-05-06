import { expect, test } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { userMocks } from "../Helper/Mocks/user.Mock";

test.describe("Verification de la suppression de co-charges de projet", () => {
    test.beforeEach(async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker
            .addMocks([
                projectMocks.getDetailedProjectsSuccess,
                userMocks.getAllManagersSuccess,
                userMocks.userMeSuccess,
            ])
            .apply();
        await page.clock.install({ time: new Date("2025-03-22T08:00:00") });
        await page.goto("http://localhost:5002/projects");
        await page.waitForLoadState("networkidle");
    });

    test("Test de la suppression de co-charges de projet reussi", async ({ page }) => {
        const mocker = new ApiMocker(page);

        await expect(page.getByText("AT-123")).toHaveCount(2);

        const firstProject = page.getByTestId("project-item").first();
        const deleteCoManagerButton = firstProject.getByRole("button", {
            name: /Supprimer le co-charg.*Jean/i,
        });

        await expect(deleteCoManagerButton).toBeVisible();
        await deleteCoManagerButton.click();

        const heading = page.getByRole("heading", { name: /Supprimer un co-charg/i });
        await expect(heading).toBeVisible();

        await mocker.clearMocks();
        await mocker
            .addMocks([
                projectMocks.deleteCoManagerSuccess,
                projectMocks.getDetailedProjectsSuccessAfterCoManagerDelete,
                userMocks.getAllManagersSuccess,
                userMocks.userMeSuccess,
            ])
            .apply();

        await page.getByRole("button", { name: "Supprimer", exact: true }).click();

        const updatedFirstProject = page.getByTestId("project-item").first();
        await expect(
            updatedFirstProject.getByRole("button", { name: /Supprimer le co-charg.*Jean/i })
        ).toHaveCount(0);
    });
});
