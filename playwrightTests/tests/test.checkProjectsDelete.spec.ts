import { test, expect } from '@playwright/test';
import { ApiMocker } from '../Helper/mockApi';
import { projectMocks } from '../Helper/Mocks/project.mock';
import { userMocks } from '../Helper/Mocks/user.Mock';

test.describe('checkProjectsDelete', () => {

    test.beforeEach(async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            projectMocks.getDetailedProjectsSuccess,
            userMocks.userMeSuccess
        ])
            .apply();
        await page.clock.install({ time: new Date('2025-03-22T08:00:00') });
        await page.goto('http://localhost:5002/projects');
        await page.waitForLoadState('networkidle');
    });

    test('projectDeleteSuccess', async ({ page }) => {
        // Configuration du mock pour la suppression réussie
        const mocker = new ApiMocker(page);
        await mocker
            .addMock(projectMocks.deleteProjectSuccess)
            .addMock(projectMocks.getDetailedProjectsSuccessAfterDelete)
            .addMock(userMocks.userMeSuccess)
            .apply();

        await expect(page.getByText('AT-123')).toHaveCount(2);

        await page.getByRole('button', { name: 'Supprimer le projet' }).first().click();

        const confirmBtn = page.getByRole('button', { name: 'Confirmer', exact: true });
        await expect(confirmBtn).toBeVisible();
        await expect(confirmBtn).toBeEnabled();

        await confirmBtn.click();

        await expect(page.getByText('AT-123')).toHaveCount(0, { timeout: 15000 });

    });

    test('deleteProjectError', async ({ page }) => {
        // Configuration du mock pour la suppression non réussie
        const mocker = new ApiMocker(page);
        await mocker
            .addMock(projectMocks.deleteProjectError)
            .apply();

        await expect(page.getByText('AT-123')).toHaveCount(2);

        await page.getByRole('button', { name: 'Supprimer le projet' }).first().click();

        const confirmBtn = page.getByRole('button', { name: 'Confirmer', exact: true });
        await expect(confirmBtn).toBeVisible();
        await expect(confirmBtn).toBeEnabled();

        const dialogPromise = page.waitForEvent('dialog');
        await confirmBtn.click();

        const dialog = await dialogPromise;
        expect(dialog.type()).toBe('alert');
        expect(dialog.message()).toContain(
            'Erreur lors de la suppression du projet'
        );
        await dialog.dismiss();
    });   
});