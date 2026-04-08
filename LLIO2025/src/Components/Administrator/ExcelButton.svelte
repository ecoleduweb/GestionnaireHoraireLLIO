<script>
 import { ReportApiService } from '../../services/ReportApiService';

 let isLoading = false;

    const handleButtonClick = async () => {
        isLoading = true
        try{
            await ReportApiService.getReportExcel()
        } catch (error) {  
            alert('Erreur lors du téléchargement du fichier Excel');
        } finally{
            isLoading = false
        }
    };
</script>
<button 
    onclick={() => handleButtonClick()}
    id="excel-button"
    type="button" 
    class="w-full mt-4 py-2 px-4 text-sm font-medium transition-colors bg-[#e6f0f0] text-[#005e61] rounded-md hover:bg-[#d0e6e6] flex items-center justify-center cursor-pointer"
    disabled={isLoading}
    >
     <svg class="w-4 h-4 mr-2" viewBox="0 0 192 192" xmlns="http://www.w3.org/2000/svg" fill="none">
      <g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
      <g id="SVGRepo_iconCarrier">
        <path d="M56 30c0-1.662 1.338-3 3-3h108c1.662 0 3 1.338 3 3v132c0 1.662-1.338 3-3 3H59c-1.662 0-3-1.338-3-3v-32m0-68V30" style="fill-opacity:.402658;stroke:#005e61;stroke-width:12;stroke-linecap:round;paint-order:stroke fill markers"></path>
        <rect width="68" height="68" x="-58.1" y="40.3" rx="3" style="fill:none;fill-opacity:.402658;stroke:#005e61;stroke-width:12;stroke-linecap:round;stroke-linejoin:miter;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" transform="translate(80.1 21.7)"></rect>
        <path d="M138.79 164.725V27.175M56.175 58.792H170M170 96H90.328M169 133.21H56.175M44.5 82l23 28m0-28-23 28" style="fill:none;stroke:#005e61;stroke-width:12;stroke-linecap:round;stroke-linejoin:round;stroke-dasharray:none;stroke-opacity:1"></path>
      </g>
    </svg> 
    {#if !isLoading}
        Générer un export Excel
    {:else}
        Génération du fichier en cours...
    {/if} 
</button>