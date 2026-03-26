<script>
  let loading = false;

  async function handleDownload() {
    loading = true;
    try {
      const response = await fetch("http://localhost:8080/report/excel", {
        method: "GET",
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Erreur lors du téléchargement");
      }

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "activities.xlsx";
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error(error);
      alert("Impossible de télécharger le fichier Excel");
    } finally {
      loading = false;
    }
  }
</script>

<button on:click={handleDownload} disabled={loading}>
  {#if loading}Téléchargement...{:else}Exporter Excel{/if}
</button>