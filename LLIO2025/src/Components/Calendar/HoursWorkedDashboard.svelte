<script lang="ts">
import { UserApiService } from '../../services/UserApiService';
import { formatDateHoursWorked, areDatesEqual } from '../../utils/date';
import HoursWorkedConfigModal from './HoursWorkedConfigModal.svelte';

const { hoursTotal, dateStart, dateEnd, textHoursWorked } = $props();

let displayedHoursTotal = $state(hoursTotal);
let displayedTextHoursWorked = $state(textHoursWorked);

let showModal = $state(false);

$effect(() => {
  displayedHoursTotal = hoursTotal;
  displayedTextHoursWorked = textHoursWorked;
});

const openConfigModal = () => {
  showModal = true;
};

const closeConfigModal = () => {
  showModal = false;
};

const handleSave = async () => {
  try {
    const bank = await UserApiService.getTimeInBank();

    displayedHoursTotal = bank.timeInBank;
    displayedTextHoursWorked = bank.textHoursWorked;

  } catch (err) {
    alert("Erreur lors de la configuration de la banque d'heure");
  }
};
</script>

<div class="bilan-container">

<div class="header">
<h2>

{#if areDatesEqual(dateStart,dateEnd)}
Bilan du {formatDateHoursWorked(dateStart)}
{:else}
Bilan du {formatDateHoursWorked(dateStart)} au {formatDateHoursWorked(dateEnd)}
{/if}

</h2>
</div>

<span>
Vous avez travaillé <strong>{displayedHoursTotal}</strong> heures {displayedTextHoursWorked}.
</span>

<div class="config-section">

<button
type="button"
class="config-btn"
onclick={openConfigModal}
>
Configurer
</button>

</div>

</div>

{#if showModal}

<HoursWorkedConfigModal
onClose={closeConfigModal}
onSave={handleSave}
/>

{/if}

<style>

.bilan-container{
padding:1rem;
border:1px solid #ddd;
border-radius:5px;
}

.header{
margin-bottom:1rem;
}

.config-section{
margin-top:1rem;
padding-top:1rem;
border-top:1px solid #eee;
}

.config-btn{
background:#015e61;
color:white;
border:none;
cursor:pointer;
padding:0.75rem 1.5rem;
border-radius:4px;
}

</style>