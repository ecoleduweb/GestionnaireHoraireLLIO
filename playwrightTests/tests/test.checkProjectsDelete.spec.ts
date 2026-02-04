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
            .apply();

        await page.waitForTimeout(1000);

        await expect(page.getByText('AT-123')).toHaveCount(2);
        await page.getByRole('button', { name: 'Supprimer le projet' }).nth(0).click();
        await page.waitForTimeout(1000);
        await page.getByRole('button', { name: 'Supprimer', exact: true }).click();
        await mocker.clearMocks();
        await mocker
            .addMock(projectMocks.getDetailedProjectsSuccessAfterDelete)
            .addMock(userMocks.userMeSuccess)
            .apply();
        await expect(page.getByText('AT-123')).toHaveCount(0);
        await page.waitForLoadState('networkidle'); 
    });

    test('deleteProjectError', async ({ page }) => {
        // Configuration du mock pour la suppression non réussie
        const mocker = new ApiMocker(page);
        await mocker
            .addMock(projectMocks.deleteProjectError)
            .apply();

        await page.waitForTimeout(1000);

        await expect(page.getByText('AT-123')).toHaveCount(2);
        await page.waitForTimeout(1000);
        await page.getByRole('button', { name: 'Supprimer le projet' }).nth(0).click();
        page.on('dialog', async dialog => {
            await expect(dialog.type()).toBe('alert');

            await expect(dialog.message()).toContain('Erreur lors de la suppression du projet, il a soit une ou des activités liées à ce projet ou bien le projet est inexistant');
            
            await dialog.dismiss()
        });
        await page.getByRole('button', { name: 'Supprimer', exact: true }).click();
        
        await page.waitForTimeout(1000);

        await page.waitForLoadState('networkidle'); 
        
    });   
});