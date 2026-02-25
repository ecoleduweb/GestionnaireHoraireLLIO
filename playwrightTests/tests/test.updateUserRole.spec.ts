import { test, expect } from '@playwright/test';
import { ApiMocker } from '../Helper/mockApi';
import { userMocks } from '../Helper/Mocks/user.Mock';

test.describe('updateUserRole', () => {

    test.beforeEach(async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            userMocks.usersSuccess,
            userMocks.userMeSuccess,
        ])
            .apply();
        await page.clock.install({ time: new Date('2025-03-22T08:00:00') });
        await page.goto('http://localhost:5002/administrator');
        await page.waitForLoadState('networkidle');
    });

    test('updateRoleSuccess', async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            userMocks.updateUserRoleSuccess,
        ]).apply();
        
        // Configurer un écouteur pour intercepter l'alerte
        await page.once('dialog', async dialog => {
            expect(dialog.message()).toBe('Rôle mis à jour avec succès');
            await dialog.dismiss().catch(() => {});
        });
        await page.locator('#user-button').click();
        await page.locator('#userSelect').selectOption('JérémieTest Lapointe | jayboss@email.com');
        await page.locator('#roleSelect').selectOption('Chargé de projet');
        
        await page.getByText('Confirmer').click();
    });
    
    test('updateRoleError', async ({ page }) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            userMocks.updateUserRoleError,
        ]).apply();
        
        // Configurer un écouteur pour intercepter l'alerte
        await page.once('dialog', async dialog => {
            expect(dialog.message()).toBe('Erreur lors de la mise à jour du rôle');
            await dialog.dismiss().catch(() => {});
        });
        await page.locator('#user-button').click();
        await page.locator('#userSelect').selectOption('JérémieTest Lapointe | jayboss@email.com');
        await page.locator('#roleSelect').selectOption('Chargé de projet');
        
        await page.getByText('Confirmer').click();
    });
});