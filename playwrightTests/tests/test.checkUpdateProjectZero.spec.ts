import { test, expect } from '@playwright/test';
import { ApiMocker } from '../Helper/mockApi';
import { projectMocks } from '../Helper/Mocks/project.mock';
import { userMocks } from '../Helper/Mocks/user.Mock';

test.describe('checkUpdateProjectZero', () => {

    test.beforeEach(async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            projectMocks.getDetailedProjectsSuccess,
            userMocks.userMeSuccess,
        ])
            .apply();
        await page.clock.install({ time: new Date('2025-03-22T08:00:00') });
        await page.goto('http://localhost:5002/projects');
        await page.waitForLoadState('networkidle');
    });

    test('UpdateProjectTimeEstimatedToZero', async ({ page }) => {
        // Cliquer sur le bouton modifier (aria-label) du premier projet
        await page.locator('button[aria-label="Modifier le projet"]').first().click();
        
        // Attendre le modal
        await page.waitForTimeout(500);
        
        // Trouver l'input du temps estimé
        const timeInput = page.locator('input[type="number"]').first();
        await expect(timeInput).toBeVisible();
        
        // Vider et entrer 0
        await timeInput.clear();
        await timeInput.fill('0');
        
        // Vérifier que la valeur est bien 0
        await expect(timeInput).toHaveValue('0');
    });
});
