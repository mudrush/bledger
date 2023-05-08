import { Money } from "./interfaces";

export function gen_key(
  account_id: string,
  money: Money,
  direction: string
): string {
  return btoa(`${account_id}-${money.amount}-${money.currency}-${direction}`);
}
