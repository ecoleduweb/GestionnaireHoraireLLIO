import { test, expect } from "@playwright/test";
import { ApiMocker } from "../Helper/mockApi";
import { userMocks } from "../Helper/Mocks/user.Mock";
import { outlookMocks } from "../Helper/Mocks/outlook.mock";
import { projectMocks } from "../Helper/Mocks/project.mock";
import { categoryMocks } from "../Helper/Mocks/category.mock";
import { activityMocks } from "../Helper/Mocks/activity.mock";

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
    await page.goto("http://localhost:5002/calendar");
    await page.waitForLoadState("networkidle");
  });

    test('Import d\'une journée sans évènements', async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            outlookMocks.getEventsNoEvent
        ]).apply();
        // Load la page et fait la requête de base
        await page.goto('http://localhost:5002/calendar');
        await page.waitForLoadState("networkidle");

        // On va chercher le bouton d'import
        const button = page.locator(
            `.import-outlook-btn[data-date="2026-05-01"]`
        );

        // Vérification de l'état original du bouton
        await expect(button).toBeVisible();
        await expect(button).toHaveText("+ Outlook");
        await expect(button).not.toBeDisabled();

        // Clic
        await button.click();

        // Vérification de l'état final
        await expect(button).toHaveText("✓ Aucun évènement à importer");
        await expect(button).toBeDisabled();
        await expect(button).toHaveAttribute("data-date", "2026-05-01");

    });

    test('outlookFailed', async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            outlookMocks.outlookFail,
        ])
            .apply();
        
        const dialogPromise = page.waitForEvent('dialog');
        
        page.getByLabel('21 mars').getByRole('button', { name: 'Outlook' }).click();
        const dialog = await dialogPromise;

        expect(dialog.type()).toBe('alert');
        expect(dialog.message()).toContain("Votre connexion à Outlook a expirée. Vous serez redirigés à la page de connexion pour vous reconnecter. Vous pourrez alors essayer d'importer vos évènements à nouveau.");
        await dialog.dismiss();

        await page.waitForURL('http://localhost:5002/');
    });
});