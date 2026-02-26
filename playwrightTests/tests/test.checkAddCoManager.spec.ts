import {expect, test} from "@playwright/test";
import {ApiMocker} from "../Helper/mockApi";
import {projectMocks} from "../Helper/Mocks/project.mock";
import {userMocks} from "../Helper/Mocks/user.Mock";


test.describe('Vérification de l\'ajout de co-chargés de projet', () => {
    test.beforeEach(async ({page}) => {
        const apiMocker = new ApiMocker(page);
        await apiMocker.addMocks([
            projectMocks.getDetailedProjectsSuccess,
            userMocks.getAllManagersSuccess,
            userMocks.userMeSuccess,
        ])
            .apply();
        await page.clock.install({time: new Date('2025-03-22T08:00:00')});
        await page.goto('http://localhost:5002/projects');
        await page.waitForLoadState('networkidle');
    });

    test('Test de l\'ajout de co-chargés de projet réussi', async ({page}) => {
        const mocker = new ApiMocker(page);
        await mocker
            .addMock(projectMocks.createCoManagerSuccess)
            .apply();

        await expect(page.getByText('AT-123')).toHaveCount(2);
        const firstProject = page.getByTestId('project-item').first();
        await firstProject.locator('button').filter({ has: page.locator('svg') }).first().click();

        const heading = page.getByRole('heading', { name: 'Ajouter un co-chargé' });
        await expect(heading).toBeVisible();

        const modal = heading.locator('xpath=ancestor::*[self::div][contains(@class,"shadow-xl")]').first();

        const select = modal.locator('select#userId');
        await expect(select).toBeVisible();
        await expect(select.locator('option')).toHaveCount(2, { timeout: 10000 });

        await expect(select.locator('option[value="3"]')).toHaveCount(1);
        await select.selectOption('3');
        await expect(select).toHaveValue('3');

        await mocker.clearMocks();
        await mocker
            .addMocks([
                projectMocks.getDetailedProjectsSuccessWithCoLead3OnProject1,
                projectMocks.createCoManagerSuccess,
                userMocks.getAllManagersSuccess,
                userMocks.userMeSuccess,
            ])
            .apply()

        await modal.getByRole('button', { name: 'Ajouter' }).click();


        const coLeadLabel = firstProject.getByText('Co-chargé·e de projet', { exact: true });
        const coLeadNames = coLeadLabel.locator(
            'xpath=following-sibling::div[contains(@class,"text-sm")]'
        );
        await expect(coLeadNames.filter({ hasText: 'Marie Amélie Dubé' })).toHaveCount(1);
    })

    test('Test de l\'ajout de co-chargés de projet échoué', async ({page}) => {
        const mocker = new ApiMocker(page);
        await mocker
            .addMock(projectMocks.createCoManagerError)
            .apply();

        await expect(page.getByText('AT-123')).toHaveCount(2);
        const firstProject = page.getByTestId('project-item').first();
        await firstProject.locator('button').filter({ has: page.locator('svg') }).first().click();

        const heading = page.getByRole('heading', { name: 'Ajouter un co-chargé' });
        await expect(heading).toBeVisible();

        const modal = heading.locator('xpath=ancestor::*[self::div][contains(@class,"shadow-xl")]').first();

        const select = modal.locator('select#userId');
        await expect(select).toBeVisible();
        await expect(select.locator('option')).toHaveCount(2, { timeout: 10000 });

        await expect(select.locator('option[value="3"]')).toHaveCount(1);
        await select.selectOption('3');
        await expect(select).toHaveValue('3');

        await mocker.clearMocks();
        await mocker
            .addMocks([
                projectMocks.getDetailedProjectsByUserSuccess,
                projectMocks.createCoManagerError,
                userMocks.getAllManagersSuccess,
                userMocks.userMeSuccess,
            ])
            .apply()

        const dialogHandled = new Promise<void>((resolve) => {
            page.once('dialog', async (dialog) => {
                expect(dialog.type()).toBe('alert');
                expect(dialog.message()).toContain('Erreur');
                await dialog.accept();
                resolve();
            });
        });

        await modal.getByRole('button', { name: 'Ajouter' }).click();
        await dialogHandled;
    })

    test('Test de l\'ajout de co-chargés de projet validation', async ({page}) => {
        const mocker = new ApiMocker(page);
        await mocker.clearMocks();
        await mocker
            .addMocks([
                projectMocks.getDetailedProjectsSuccessWithCoLead3OnProject1,
                userMocks.getAllManagersSuccess,
                userMocks.userMeSuccess,
            ])
            .apply()

        await page.reload();
        await page.waitForLoadState('networkidle');

        await expect(page.getByText('AT-123')).toHaveCount(2);
        const firstProject = page.getByTestId('project-item').first();
        await firstProject.locator('button').filter({ has: page.locator('svg') }).first().click();

        const heading = page.getByRole('heading', { name: 'Ajouter un co-chargé' });
        await expect(heading).toBeVisible();

        const modal = heading.locator('xpath=ancestor::*[self::div][contains(@class,"shadow-xl")]').first();

        const select = modal.locator('select#userId');
        await expect(select).toBeVisible();
        await expect(select.locator('option')).toHaveCount(1, { timeout: 10000 });

        const onlyOption = select.locator('option').first();
        await expect(onlyOption).toHaveAttribute('value', '');
        await expect(onlyOption).toHaveText("Aucun utilisateur disponible à l'ajout");
        await expect(onlyOption).toBeDisabled();

        const confirmButton = page.getByRole('button', { name: 'Ajouter', exact: true });
        await expect(confirmButton).toBeDisabled();
    })
})