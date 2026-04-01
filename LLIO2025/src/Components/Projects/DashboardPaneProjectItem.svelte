<script lang="ts">
  import { formatHours } from '../../utils/date';
  import { slide } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { ChevronDown } from 'lucide-svelte';
  import { getHoursColor } from '../../utils/displayUtils';
  import { calculateEmployeeTime, calculateRemainingTime } from '../../utils/CalculUtils';

  let { project } = $props();

 let open = $state(false); // état ouvert/fermé
 let selectedEmployee = $state(0);

  function toggle() {
    open = !open;
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
        {#each project.employees[selectedEmployee].categories as category}
          <tr class="bg-gray-50">
            <td class="py-2 pl-4 text-left">
              {category.name}
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