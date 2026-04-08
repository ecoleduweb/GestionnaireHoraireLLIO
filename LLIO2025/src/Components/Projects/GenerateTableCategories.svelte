<script lang="ts">
  import { Plus } from "lucide-svelte";
  import type { Category } from "../../Models";
  import { formatHours } from '../../utils/date';
  import { getHoursColor } from '../../utils/displayUtils';

  let { categories }: { categories: Category[] } = $props(); //employee.categories
  let listCategories = $state(categories);

</script>

<table class="w-full">
    <tbody>
        {#each listCategories as category, categoryIndex}
            <tr
                class="border-b border-gray-200 {categoryIndex % 2 === 0
                ? 'bg-white'
                : 'bg-gray-50'}"
            >
                <td class="py-2 text-left w-1/2 pl-4">{category.name}</td>
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