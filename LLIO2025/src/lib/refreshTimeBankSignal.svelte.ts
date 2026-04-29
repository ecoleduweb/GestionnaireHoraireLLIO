import { isDateInCurrentWeek } from '../utils/date';

export const refreshTimeBankSignal = $state({ tick: 0 });

export function triggerTimeBankRefetch(activityStartDate: Date): void {
  if (!isDateInCurrentWeek(activityStartDate)) refreshTimeBankSignal.tick++;
}
