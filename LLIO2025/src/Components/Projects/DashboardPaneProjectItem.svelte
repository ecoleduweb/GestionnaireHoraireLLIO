<script lang="ts">
  import { formatHours } from '../../utils/date';
  import { slide } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { ChevronDown, Pencil } from 'lucide-svelte';
  import { getHoursColor } from '../../utils/displayUtils';
  import { calculateEmployeeTime, calculateRemainingTime } from '../../utils/CalculUtils';
  import type { Category } from '../../Models';
  import TextInputModal from '../TextInputModal.svelte';
  import DOMPurify from 'dompurify';
  import { CategoryApiService } from '../../services/CategoryApiService';

  let { project } = $props();

 let open = $state(false); // état ouvert/fermé
 let selectedEmployee = $state(0);

  function toggle() {
    open = !open;
  }

  let categories: Category[] = $state(project.employees[selectedEmployee].categories);

  let hoveredCategoryName: string = $state("");

  let selectedCategory: Category = $state();
  let enablePrompt = $state(false);

  function handleRenameCategory(category: Category) 
  {
    selectedCategory = category;
    enablePrompt = true;
  }

  async function sendRenameCategory(category: Category, newName: string) {
    try {
      return await CategoryApiService.changeCategoryName(newName, category);
    }
    catch (error) {
      alert("Erreur - impossible de modifier le nom de la catégorie")
    }
    return false;
  }

</script>

<div class="border-l-10 border-b" style="border-left-color: {project.color}">
  <div class="p-4">

    <button
  type="button"
  class="flex justify-between items-center w-full text-left"
  onclick={toggle}
  aria-expanded={open}
>
  <div>
    <span class="text-black">{project.uniqueId}</span>
    <span class="text-gray-500 ml-2">|</span>
    <span class="text-gray-500 ml-2">{project.name}</span>
  </div>

  <div class="flex items-center gap-2">
   
    <ChevronDown
      class="w-4 h-4 shrink-0 transition-transform duration-200 {open ? 'rotate-90' : ''}"
    />
  </div>
</button>


    {#if open}
  <div
    class="px-2 pb-2 bg-white text-sm "
    transition:slide={{ duration: 300, easing: quintOut }}
  >
    <table class="w-full border-separate border-spacing-y-1">

      <thead>
        <tr class="text-xs text-gray-400">
          <th class="text-left pl-4 font-normal w-1/4">Catégorie</th>
          <th class="text-center font-normal w-1/4">Temps<br>passé</th>
          <th class="text-center font-normal w-1/4">Temps<br>estimé</th>
          <th class="text-right pr-4 font-normal w-1/4">Temps<br>restant</th>
        </tr>
      </thead>

      <tbody>
        {#each categories as category}
          <tr class="bg-gray-50"
            onmouseenter={() => hoveredCategoryName = category.name }
            onmouseleave={() => hoveredCategoryName = "" }
          >
            <td class="py-2 pl-4 text-left">
              {category.name}
              {#if category.name !== "Par défaut" && category.name == hoveredCategoryName}
                        <button
                            class="justify-end p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
                            onclick={() => handleRenameCategory(category)}
                            aria-label="Renommer la catégorie"
                        >
                        <Pencil size={10} />
                    </button>
                {/if}
            </td>


            <td class="py-2 px-4 text-right">
              {formatHours(category.timeSpent)}
            </td>

  
            <td class="py-2 px-4 text-right text-gray-400">
              {formatHours(category.timeEstimated)}
            </td>

           <td class="py-2 pr-4 text-right {getHoursColor(category.timeSpent, category.timeEstimated)}">
              {formatHours(category.timeEstimated - category.timeSpent)}
            </td>
          </tr>
        {/each}


        <tr class="bg-gray-100 font-medium">
          <td class="py-3 pl-4">Total</td>

          <td class="text-right px-4">
            {formatHours(calculateEmployeeTime(project.employees[selectedEmployee], 'spent'))}
          </td>

          <td class="text-right px-4 text-gray-400">
            {formatHours(calculateEmployeeTime(project.employees[selectedEmployee], 'estimated'))}
          </td>

          <td class="pr-4 text-right {getHoursColor(
              calculateEmployeeTime(project.employees[selectedEmployee], 'spent'),
              calculateEmployeeTime(project.employees[selectedEmployee], 'estimated')
            )}">
              {formatHours(
                calculateRemainingTime(
                calculateEmployeeTime(project.employees[selectedEmployee], 'spent'),
                calculateEmployeeTime(project.employees[selectedEmployee], 'estimated')
                )
            )}
          </td>
        </tr>

      </tbody>
    </table>
  </div>
{/if}

  </div>
</div>


{#if enablePrompt}
    <TextInputModal
      modalTitle="Modification du nom de la catégorie"
      modalText={`Renommez la catégorie "<strong>` + DOMPurify.sanitize(selectedCategory.name) + `</strong>"`}
      errorText="Erreur lors de la suppression du projet, il a soit une ou des activités liées à ce projet ou bien le projet est inexistant"
      defaultTextInValue={DOMPurify.sanitize(selectedCategory.name)}
      onSuccess={async (val: string) => {
        if (await sendRenameCategory(selectedCategory, val)) {
          categories = categories.map((cat) => cat.id === selectedCategory.id ? {...cat, name: val} : cat);
        }
        enablePrompt = false;
      }}
      onClose={() => {
        enablePrompt = false;
      }}
    />
{/if}