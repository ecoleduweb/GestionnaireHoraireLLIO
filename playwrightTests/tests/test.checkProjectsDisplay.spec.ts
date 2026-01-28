import { test, expect } from '@playwright/test';
import { ApiMocker } from '../Helper/mockApi';
import { projectMocks } from '../Helper/Mocks/project.mock';
import { userMocks } from '../Helper/Mocks/user.Mock';

test.describe('checkProjectsDisplay', () => {

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

    test('projetsVisibles', async ({ page }) => {
        //Valide que les projets s'affichent sur la page principale ET dans le left pane
        await expect(page.getByText('AT-123')).toHaveCount(2);
        await expect(page.getByText('AT-456')).toHaveCount(2);
        await expect(page.getByText('FO-115')).toHaveCount(2);
        await expect(page.getByText('RA-224')).toHaveCount(2);        
        await expect(page.getByText('TEST-1')).toHaveCount(2);
        await expect(page.getByText('TEST-2')).toHaveCount(2);
        await expect(page.getByText('TEST-3')).toHaveCount(2);
        await expect(page.getByText('TEST-4')).toHaveCount(2);  
    });

    test('projetsArchivesVisibles', async ({ page }) => {
        await expect(page.getByText('Projets archivés (4)')).toBeVisible();
    })
});