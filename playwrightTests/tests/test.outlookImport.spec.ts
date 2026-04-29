import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { outlookMocks } from "../Helper/Mocks/outlook.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";

test.describe("ImportOutlook", () => {
  test.beforeEach(async ({ page }) => {
    const apiMocker = new ApiMocker(page);
    await apiMocker
      .addMocks([
        projectMocks.getDetailedProjectsByUserSuccess,
        categoryMocks.getCategoriesByProjectSuccess,
        userMocks.userMeSuccess,
      ])
      .apply();
    await page.clock.install({ time: new Date("2026-05-01T06:00:00") });
    await page.goto("http://localhost:5002/projects");
    await page.waitForLoadState("networkidle");
  });
});
