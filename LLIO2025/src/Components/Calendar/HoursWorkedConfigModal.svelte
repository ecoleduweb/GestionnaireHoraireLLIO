<script lang="ts">
import { X } from 'lucide-svelte';
import { UserApiService } from '../../services/UserApiService';
import type { TimeBankConfig } from '../../Models/index';

let { onClose, onSave } = $props();

let config = $state<TimeBankConfig>({
  startDate: '',
  hoursPerWeek: 0,
  offset: 0
});

const updateStartDate = (event: Event) => {
  const input = event.currentTarget as HTMLInputElement;
  config.startDate = input.value;
};

const updateHoursPerWeek = (event: Event) => {
  const input = event.currentTarget as HTMLInputElement;
  const value = Number(input.value);

  
  config.hoursPerWeek = value < 0 ? 0 : value;
};

const updateOffset = (event: Event) => {
  const input = event.currentTarget as HTMLInputElement;
  const value = Number(input.value);

  config.offset = value < 0 ? 0 : value;
};

const handleSubmit = async () => {
  try {

    await UserApiService.saveTimeBankConfig(config);
    onSave();
    onClose();

  } catch (err) {

    console.error(err);
    alert("Erreur lors de la configuration de la banque d'heure");
  }
};
</script>

<div class="modal-overlay">

<div class="modal">

<div class="modal-header">

<h3>Configuration des heures en banque</h3>

<button class="close-btn" onclick={onClose}>
<X size={18}/>
</button>

</div>

<div class="modal-body">

<div class="form-group">

<label>Début de la période</label>

<input
type="date"
value={config.startDate}
oninput={updateStartDate}
class="form-input"
/>

</div>

<div class="form-group">

<label>Nombre d'heures par semaine</label>

<input
type="number"
value={config.hoursPerWeek}
min="0"
step="0.5"
oninput={updateHoursPerWeek}
class="form-input"
/>

</div>

<div class="form-group">

<label>Décalage (offset)</label>

<input
type="number"
value={config.offset}
min="0"
step="0.5"
oninput={updateOffset}
class="form-input"
/>

</div>

</div>

<div class="modal-footer">

<button class="btn-cancel" onclick={onClose}>
Annuler
</button>

<button class="btn-save" onclick={handleSubmit}>
Enregistrer
</button>

</div>

</div>

</div>

<style>

.modal-overlay{
position:fixed;
top:0;
left:0;
right:0;
bottom:0;
background-color:rgba(0,0,0,0.5);
display:flex;
align-items:center;
justify-content:center;
z-index:1000;
}

.modal{
background:white;
border-radius:8px;
box-shadow:0 4px 6px rgba(0,0,0,0.15);
max-width:400px;
width:90%;
}

.modal-header{
display:flex;
justify-content:space-between;
align-items:center;
padding:1.5rem;
border-bottom:1px solid #eee;
}

.modal-body{
padding:1.5rem;
}

.form-group{
margin-bottom:1.5rem;
}

label{
display:block;
margin-bottom:0.5rem;
font-weight:500;
font-size:0.9rem;
}

.form-input{
width:100%;
padding:0.75rem;
border:1px solid #ddd;
border-radius:4px;
font-size:1rem;
box-sizing:border-box;
}

.form-input:focus{
outline:none;
border-color:#015e61;
box-shadow:0 0 0 3px rgba(1,94,97,0.1);
}

.modal-footer{
display:flex;
gap:1rem;
padding:1.5rem;
border-top:1px solid #eee;
justify-content:flex-end;
}

.btn-cancel,
.btn-save{
padding:0.7rem 1.4rem;
border:none;
border-radius:4px;
font-weight:500;
cursor:pointer;
font-size:0.9rem;
}

.btn-cancel{
background:#f0f0f0;
color:#333;
}

.btn-save{
background:#015e61;
color:white;
}

</style>