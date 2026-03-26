import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";

test.describe("calcul d'heures", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await page.clock.install({
      time: new Date("2025-03-22T08:00:00-04:00"),
    });

    await apiMocker
      .addMocks([
        userMocks.userMeSuccess,
        activityMocks.getAllActivityEmpty,
        projectMocks.getProjectsListSuccess,
        categoryMocks.getCategoriesByProjectSuccess,
      ])
      .apply();

    await page.goto("http://localhost:5002/calendar");
  });
});

 test("calcul d'heures success", async ({ page }) => {
  const apiMocker = new ApiMocker(page);

  await apiMocker
    .addMocks([
      userMocks.saveTimeBankConfigSuccess,
      userMocks.getTimeBankSuccess
    ])
    .apply();

  await page.getByText("Configurer").click();

  await page.fill('[name="startDate"]', '2025-03-01');
  await page.fill('[name="hoursPerWeek"]', '40');

  await page.getByText("Enregistrer").click();

  await expect(page.getByTestId("total-hours")).toHaveText("40");
});
