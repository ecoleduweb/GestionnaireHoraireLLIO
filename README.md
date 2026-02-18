Je n'ai pas réussi à fix le bug du calendrier.
Je pense que c'est un problème qui vient du component lui-même.

Constat:
    Le CSS n'est plus écouté quand qu'il change de taille, donc même mettre un max-width et max-height ne fonctionne pas
    La taille des jours est multiplier par 2.

Les essaies:
    J'ai individuellement mis en commentaire des bouts de codes de la page calendrier elle-même (+page.svelte) à CalendarService (calendar.service.ts) jusqu'au Component et rien n'a vraiment réparer le problème.
    J'ai essayer d'aussi changer le temps, pour voir s'il y a un changement et non.


(Enlever le message après problème résolut)