<script lang="ts">
  import { Pencil, Plus } from "lucide-svelte";
  import type { Category } from "../../Models";
  import { formatHours } from '../../utils/date';
  import { getHoursColor } from '../../utils/displayUtils';
  import TextInputModal from "../TextInputModal.svelte";
  import DOMPurify from 'dompurify';

  let { categories, sendRenameCategory }: { categories: Category[], sendRenameCategory: (category: Category, newName: string) => Promise<boolean> } = $props(); //employee.categories
  let listCategories = $state(categories);

  let hoveredCategoryName: string = $state("");

  let selectedCategory: Category = $state();
  let enablePrompt = $state(false);

  function handleRenameCategory(category: Category) 
  {
    selectedCategory = category;
    enablePrompt = true;
  }

</script>

<table class="w-full">
    <tbody>
        {#each listCategories as category, categoryIndex}
            <tr
                class="border-b border-gray-200 {categoryIndex % 2 === 0
                ? 'bg-white'
                : 'bg-gray-50'}"
                
                onmouseenter={() => hoveredCategoryName = category.name }
                onmouseleave={() => hoveredCategoryName = "" }
            >
                <td class="py-2 text-left w-1/2 pl-4">{category.name}
                {#if category.name !== "Par défaut" && category.name == hoveredCategoryName}
                        <button
                            class="justify-end p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
                            onclick={() => handleRenameCategory(category)}
                            aria-label="Renommer la catégorie"
                        >
                        <Pencil size={16} />
                    </button>
                {/if}
                </td>
                <td class="py-2 text-right w-1/6">{formatHours(category.timeSpent)}</td>
                    <td class="py-2 text-right w-1/6">{formatHours(category.timeEstimated)}</td>
                    <td class="py-2 text-right w-1/6 {getHoursColor(category.timeSpent, category.timeEstimated)}">{formatHours(category.timeEstimated - category.timeSpent)}</td>
            </tr> 
        {/each}
        <tr>
            <td colspan="4" class="py-2 pl-4">
                <button
                    class="mt-2 inline-flex items-center bg-gray-100 border border-transparent rounded-4xl shadow-sm hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-grey-500 text-gray-700 text-xs"
                 >
                    <Plus class="w-3 h-3" />
                </button>
            </td>
        </tr>
    </tbody>
</table>

{#if enablePrompt}
    <TextInputModal
      modalTitle="Modification du nom de la catégorie"
      modalText={`Renommez la catégorie "<strong>` + DOMPurify.sanitize(selectedCategory.name) + `</strong>"`}
      errorText="Erreur lors de la suppression du projet, il a soit une ou des activités liées à ce projet ou bien le projet est inexistant"
      defaultTextInValue={DOMPurify.sanitize(selectedCategory.name)}
      onSuccess={async (val: string) => {
        if (await sendRenameCategory(selectedCategory, val)) {
            listCategories = listCategories.map((cat) => cat.id === selectedCategory.id ? {...cat, name: val} : cat);
        }
        enablePrompt = false;
      }}
      onClose={() => {
        enablePrompt = false;
      }}
    />
{/if}