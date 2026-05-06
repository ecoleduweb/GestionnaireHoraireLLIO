import { test, expect } from '@playwright/test';
import { ApiMocker } from '../Helper/mockApi';
import { projectMocks } from '../Helper/Mocks/project.mock';
import { userMocks } from '../Helper/Mocks/user.Mock';
import { activityMocks } from '../Helper/Mocks/activity.mock';
import { categoryMocks } from '../Helper/Mocks/category.mock';

test.describe('checkCategoriesRename', () => {

    test.beforeEach(async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            projectMocks.getDetailedProjectsSuccess,
            activityMocks.getAllActivitySuccess,
            activityMocks.getAllActivitiesDefaultWeekSuccess,
            categoryMocks.getCategoriesByProjectSuccess,
            userMocks.userMeSuccess
        ])
            .apply();
        await page.clock.install({ time: new Date('2025-03-22T08:00:00') });
        await page.goto('http://localhost:5002/projects');
        await page.waitForLoadState('networkidle');
    });

    test('categoryRenameSuccess', async ({ page }) => {
        // Configuration du mock pour la suppression réussie
        const mocker = new ApiMocker(page);
        await mocker
            .addMock(projectMocks.deleteProjectSuccess)
            .addMock(projectMocks.getDetailedProjectsSuccessAfterDelete)
            .addMock(userMocks.userMeSuccess)
            .addMock(categoryMocks.renameCategorySuccess)
            .apply();

        await page.goto('http://localhost:5002/projects');
        
        await page.waitForTimeout(2000);

        await page.getByRole('button', { name: 'Katell Arnault de la Ménardière' }).first().click();
        await page.getByRole('cell', { name: 'Développement' }).hover();
        await page.getByRole('button', { name: 'Renommer la catégorie' }).click();
        await page.locator('form').getByRole('textbox').click();
        await page.locator('form').getByRole('textbox').fill('test');
        await page.getByRole('button', { name: 'Confirmer' }).click();
        await expect(page.getByRole('cell', { name: 'test' })).toHaveCount(1, { timeout: 15000 });
    });


    test('categoryRenameFail', async ({ page }) => {
        // Configuration du mock pour la suppression réussie
        const mocker = new ApiMocker(page);
        await mocker
            .addMock(projectMocks.deleteProjectSuccess)
            .addMock(projectMocks.getDetailedProjectsSuccessAfterDelete)
            .addMock(userMocks.userMeSuccess)
            .addMock(categoryMocks.renameCategoryFailed)
            .apply();

        await page.goto('http://localhost:5002/projects');
        
        await page.waitForTimeout(2000);

        await page.getByRole('button', { name: 'Katell Arnault de la Ménardière' }).first().click();
        await page.getByRole('cell', { name: 'Développement' }).hover();
        await page.getByRole('button', { name: 'Renommer la catégorie' }).click();
        await page.locator('form').getByRole('textbox').click();

        page.on('dialog', async dialog => {
            if (dialog.type() == 'alert') {
                expect(dialog.message()).toContain(
                    'Erreur - impossible de modifier le nom de la catégorie'
                );
                await dialog.dismiss();
            }

            
        });

        await page.locator('form').getByRole('textbox').fill('test');
        await page.getByRole('button', { name: 'Confirmer' }).click();
        await expect(page.getByRole('cell', { name: 'test' })).toHaveCount(0, { timeout: 15000 });
    });

});