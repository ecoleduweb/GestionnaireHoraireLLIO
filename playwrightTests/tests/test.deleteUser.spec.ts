
import { userMocks } from '../Helper/Mocks/user.Mock';
import { ApiMocker } from '../Helper/mockApi';
import { test, expect } from '@playwright/test';

test.describe('deleteUser', () => {
  
  test.beforeEach(async ({ page }) => {
    // Configuration des mocks de base pour charger la page
    const mocker = new ApiMocker(page);
    await mocker
      .addMock(userMocks.userMeSuccess)
      .addMock(userMocks.usersSuccess)
      .apply();

    // Naviguer vers la page de gestion des utilisateurs
    await page.goto('/administrator'); // Ajustez l'URL selon votre routing
  });

  test('deleteUserSuccess', async ({ page }) => {
    // Configuration du mock pour la suppression réussie
    const mocker = new ApiMocker(page);
    await mocker
      .addMock(userMocks.deleteUserSuccess)
      .apply();

    await page.waitForTimeout(1000); // Attendre que le sélecteur soit interactif 
    await page.locator('#user-button').click();
    await page.locator('#userSelect').click();
    await page.locator('#userSelect').selectOption({ label: 'Charle-ÉtienneTest Soucy | wong@email.com' });
    await page.getByText('Supprimer').click();
    await page.locator('#userSelect').click();
    await expect(page.locator('#userSelect')).not.toContainText('Charle-ÉtienneTest Soucy | wong@email.com');
    await page.waitForLoadState('networkidle'); 
    });

  test('deleteUserError', async ({ page }) => {
  
    // Configuration du mock pour la suppression avec erreur
    const mocker = new ApiMocker(page);
    await mocker
      .addMock(userMocks.deleteUserError)
      .apply();

    await page.waitForTimeout(1000); // Attendre que le sélecteur soit interactif
    await page.locator('#user-button').click();
    await page.locator('#userSelect').click();
    await page.locator('#userSelect').selectOption({ label: 'Charle-ÉtienneTest Soucy | wong@email.com' });
    await page.getByRole('button', { name: 'Supprimer' }).click();
    await page.locator('#userSelect').click();
    expect(page.locator('#userSelect')).toContainText('Charle-ÉtienneTest Soucy | wong@email.com');
    await page.waitForTimeout(2000);
    
    
  });   
});