<script lang="ts">
  import { formatHours } from '../../utils/date';
  import type { Project } from '../../Models/index.ts';
  import { slide } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { ChevronDown } from 'lucide-svelte';

  let { project } = $props();

 let open = $state(false); // état ouvert/fermé
 let selectedEmployee = $state(0);

  function toggle() {
    open = !open;
  }
const calculateEmployeeTime = (employee: any, type: 'spent' | 'estimated'): number => {
  return employee.categories.reduce(
    (sum: number, cat: any) =>
      sum + (type === 'spent' ? cat.timeSpent : cat.timeEstimated),
    0
  );
};const calculateRemainingTime = (timeSpent: number, timeEstimated: number): number => {
  return timeEstimated - timeSpent;
};
  function getProjectHoursColorBasedOnEstimatedHours(project: Project) {
    if (project.totalTimeEstimated <= 0) {
      return 'text-gray-400';
    } else if (project.totalTimeSpent > project.totalTimeEstimated) {
      return 'text-red-700';
    } else {
      return 'text-gray-700';
    }
  }
</script>

<div class="border-l-10 border-b" style="border-left-color: {project.color}">
  <div class="p-4">

    <!-- HEADER CLIQUABLE -->
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

      <ChevronDown
      class="w-4 h-4 shrink-0 transition-transform duration-200 {open ? 'rotate-90' : ''}"
      />
    </button>

   <div class="flex gap-2 mt-2">

</div>

    {#if open}
  <div
    class="px-2 pb-2 bg-white text-sm overflow-hidden"
    transition:slide={{ duration: 300, easing: quintOut }}
  >
    <table class="w-full border-separate border-spacing-y-1">

      <thead>
        <tr class="text-xs text-gray-400">
          <th class="text-left pl-4 font-normal w-1/4">Catégorie</th>
          <th class="text-right font-normal w-1/4">Temps passé</th>
          <th class="text-right font-normal w-1/4">Temps estimé</th>
          <th class="text-right pr-4 font-normal w-1/4">Temps restant</th>
        </tr>
      </thead>

      <tbody>
        {#each project.employees[selectedEmployee].categories as category}
          <tr class="bg-gray-50">
            <td class="py-2 pl-4 text-left">
              {category.name}
            </td>


            <td class="py-2 px-4 text-right">
              {formatHours(category.timeSpent)}
            </td>

  
            <td class="py-2 px-4 text-right text-gray-400">
              {category.timeEstimated
                ? formatHours(category.timeEstimated)
                : '-'}
            </td>

            <td
              class="py-2 pr-4  text-right"
              class:text-red-500={category.timeEstimated - category.timeSpent < 0}
              class:text-gray-400={category.timeEstimated - category.timeSpent >= 0}
            >
              {category.timeEstimated
                ? formatHours(category.timeEstimated - category.timeSpent)
                : '-'}
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

          <td
            class="pr-4 text-right"
            class:text-red-500={
              calculateRemainingTime(
                calculateEmployeeTime(project.employees[selectedEmployee], 'spent'),
                calculateEmployeeTime(project.employees[selectedEmployee], 'estimated')
              ) < 0
            }
          >
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