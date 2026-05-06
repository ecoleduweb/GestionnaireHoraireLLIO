import { test, expect } from '@playwright/test';
import { ApiMocker } from '../Helper/mockApi';
import { projectMocks } from '../Helper/Mocks/project.mock';
import { userMocks } from '../Helper/Mocks/user.Mock';

test.describe('checkProjectsArchive', () => {

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

    test('projetArchiveSuccess', async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([projectMocks.archiveSuccess]).apply();
        //Valide que les projets s'affichent sur la page principale ET dans le left pane
        await expect(page.getByText('Projets archivés (4)')).toBeVisible();
        await page.getByRole('button', { name: 'Archiver le projet' }).first().click();
        page.once('dialog', dialog => {
            console.log(`Dialog message: ${dialog.message()}`);
            dialog.dismiss().catch(() => {});
        });
        await apiMocker.addMocks([projectMocks.getDetailedProjectsSuccessAfterArchive]).apply();

        await page.waitForTimeout(1000);

        await page.getByRole('button', { name: 'Confirmer' }).click();

        await page.waitForTimeout(1000);

        await expect(page.getByText('Projets archivés (5)')).toBeVisible();

    });

    test('projetArchiveFailed', async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([projectMocks.archiveFailed]).apply();

        await expect(page.getByText('Projets archivés (4)')).toBeVisible();
        await page.getByRole('button', { name: 'Archiver le projet' }).first().click();

        let dialogHandled = false;

        page.once('dialog', async dialog => {
            expect(dialog.type()).toBe('alert');
            expect(dialog.message()).toContain('Erreur lors de la archivation du projet.');
            dialogHandled = true;
            await dialog.dismiss();
        });

        await page.getByRole('button', { name: 'Confirmer' }).click();

        await expect.poll(() => dialogHandled).toBe(true);
        await expect(page.getByText('Projets archivés (4)')).toBeVisible();
    });

});