export interface Money {
  amount: number;
  currency: string;
}

export interface Transaction {
  id: string;
  account_id: string;
  money: Money;
  direction: string;
  memo: string;
  state: string;
  error_reason?: string;
}

export interface Account {
  id: string;
  name: string;
  description: string;
  balance: Money;
}
