import { test, expect } from '@playwright/test';
import { ApiMocker } from '../Helper/mockApi';
import { activityMocks } from '../Helper/Mocks/activity.mock';
import { userMocks } from '../Helper/Mocks/user.Mock';
import { projectMocks } from '../Helper/Mocks/project.mock';

test.describe('showActivityDetails', () => {

    test.beforeEach(async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            userMocks.userMeSuccess
        ])
        .apply();
        await page.clock.install({ time: new Date('2025-03-22T08:00:00-04:00') });         
    });

    test('showActivityDetail', async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
           activityMocks.getAllActivitiesDefaultWeekSuccess
        ]).apply();

        // Load la page et fait la requête de base 
        await page.goto('http://localhost:5002/calendar');
        await page.waitForSelector('.fc-event', { state: 'visible' });

        let activities = await page.locator('.fc-event').all();
        
        await activities[0].click();
        await page.waitForTimeout(2000);
        await expect(page.getByLabel('Nom')).toHaveValue('Veille de la fête ');
    
    });
});