export const refreshTimeBankSignal= $state({ tick: 0 });
 
export function triggerTimeBankRefetch(): void {
  refreshTimeBankSignal.tick++;
}