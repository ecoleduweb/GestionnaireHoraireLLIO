<script lang="ts">
  import { Plus } from 'lucide-svelte';
  import { slide } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { formatHours } from '../../utils/date';

  let { project } = $props();
  let isDetailsVisible = $state([]);

  const calculateRemainingTime = (timeSpent: number, timeEstimated: number): number =>{
    return timeEstimated - timeSpent;
  }

  const calculateEmployeeTime = (employee: any, type: 'spent' | 'estimated'): number =>{
    return employee.categories.reduce(
      (sum: number, cat: any) => sum + (type === 'spent' ? cat.timeSpent : cat.timeEstimated),
      0
    );
  }
</script>

<style>
  .project {
    width: auto;
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin: 15px;
  }
</style>

<div class="project">
  <!-- Contenu du dashboard -->
  <div class="project-content">
    <div>
      <!-- Bordure gauche de couleur -->
      <div class="border-l-12 rounded" style="border-left-color: {project.color}">
        <div class="p-4">
          <div class="flex justify-between items-center">
            <div>
              <span class="text-xl">{project.uniqueId}</span>
              <span class="text-xl text-gray-500 ml-2">|</span>
              <span class="text-xl">{project.name}</span>
            </div>
          </div>
          <div class="flex mt-1">
            <!-- Section gauche de la ligne (20% width) -->
            <div class="w-1/5 flex-shrink-0">
              <div class="mt-1 text-xs text-gray-400">Chargé·e de projet</div>
              <div class="text-sm wrap-normal">{project.lead}</div>
              <button
                class="mt-2 inline-flex items-center px-3 py-1.5 bg-gray-100 border border-transparent rounded-4xl shadow-sm hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-grey-500 text-gray-700 text-xs font-medium"
              >
                <span class="text-xs">Réattribuer</span>
              </button>
              <hr class="mt-2 text-xs text-gray-400" />
              <div class="mt-1 text-xs text-gray-400">Co-chargé·e de projet</div>
              {#each project.coLeads as coLead}
                <div class="text-sm wrap-normal">{coLead}</div>
              {/each}
              <button
                class="mt-2 inline-flex items-center bg-gray-100 border border-transparent rounded-4xl shadow-sm hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-grey-500 text-gray-700 text-xs"
              >
                <Plus class="w-3 h-3" />
              </button>
            </div>

            <!-- <vr/> -->
            <div class="ml-4 border-l border-gray-300 h-auto"></div>

            <!-- Section droite de la ligne -->
            <div class="flex-1 pl-4">
              <div class="flex justify-end pr-2 mb-1">
                <div class="w-1/2"></div>
                <div class="w-1/6 text-right text-xs text-gray-500">Temps passé</div>
                <div class="w-1/6 text-right text-xs text-gray-500">Temps estimé</div>
                <div class="w-1/6 text-right text-xs text-gray-500">Temps restant</div>
              </div>
              {#each project.employees as employee, index}
                <button
                  class="w-full p-2 flex items-center justify-between hover:bg-gray-50 mt-1 cursor-pointer"
                  onclick={() => (isDetailsVisible[index] = !isDetailsVisible[index])}
                >
                  <span class="text-sm text-left w-1/2">{employee.name}</span>
                  <div class="flex w-1/2">
                    <div class="w-1/3 text-right text-sm">
                      {formatHours(
                        employee.categories.reduce((sum, cat) => sum + cat.timeSpent, 0)
                      )}
                    </div>
                    <div class="w-1/3 text-right text-sm">
                      {formatHours(
                        employee.categories.reduce((sum, cat) => sum + cat.timeEstimated, 0)
                      )}
                    </div>
                    <div
                      class="w-1/3 text-right text-sm"
                      class:text-red-500={calculateRemainingTime(
                        calculateEmployeeTime(employee, 'spent'),
                        calculateEmployeeTime(employee, 'estimated')
                      ) < 0}
                    >
                      {formatHours(
                        calculateRemainingTime(
                        calculateEmployeeTime(employee, 'spent'),
                        calculateEmployeeTime(employee, 'estimated')
                        )
                      )}
                    </div>
                  </div>
                  <svg
                    class="w-4 h-4 transform transition-transform ml-2 {isDetailsVisible[index]
                      ? 'rotate-180'
                      : ''}"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 9l-7 7-7-7"
                    ></path>
                  </svg>
                </button>
                {#if isDetailsVisible[index]}
                  <div
                    class="p-2 bg-white text-sm overflow-hidden"
                    transition:slide={{ duration: 300, easing: quintOut }}
                  >
                    <table class="w-full">
                      <tbody>
                        {#each employee.categories as category, categoryIndex}
                          <tr
                            class="border-b border-gray-200 {categoryIndex % 2 === 0
                              ? 'bg-white'
                              : 'bg-gray-50'}"
                          >
                            <td class="py-2 text-left w-1/2 pl-4">{category.name}</td>
                            <td class="py-2 text-right w-1/6">{formatHours(category.timeSpent)}</td>
                            <td class="py-2 text-right w-1/6"
                              >{formatHours(category.timeEstimated)}</td
                            >
                            <td
                              class="py-2 text-right w-1/6"
                              class:text-red-500={category.timeEstimated - category.timeSpent < 0}
                            >
                              {formatHours(category.timeEstimated - category.timeSpent)}
                            </td>
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
                  </div>
                {/if}
              {/each}

              <!-- Bouton pour ajouter un employé -->
              <div class="py-2 pl-4 mt-2">
                <button
                  class="inline-flex items-center bg-gray-100 border border-transparent rounded-4xl shadow-sm hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-grey-500 text-gray-700 text-xs"
                >
                  <Plus class="w-3 h-3" />
                </button>
              </div>

              <!-- Total du projet -->
              <div
                class="w-full p-3 flex items-center justify-between bg-gray-100 mt-3 font-medium"
              >
                <span class="text-2xl text-left w-1/2">Total</span>
                <div class="flex w-1/2">
                  <div class="w-1/3 text-right text-sm">
                    {formatHours(project.totalTimeSpent)}
                  </div>
                  <div class="w-1/3 text-right text-sm font-normal">
                    {formatHours(project.totalTimeEstimated)}
                  </div>
                  <div
                    class="w-1/3 text-right text-sm font-normal"
                    class:text-red-500={project.totalTimeRemaining < 0}
                  >
                    {formatHours(project.totalTimeRemaining)}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
