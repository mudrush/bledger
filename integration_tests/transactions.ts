import { assert } from "console";
import fetch from "node-fetch";
import { testCreateAccount } from "./account";
import { gen_key } from "./helpers";
import { Account, Transaction } from "./interfaces";

export async function transactions() {
  await testFailedTransaction();
  await testPendingTransaction();
  await testImmediateTransaction();
  await testPendingTransactionExecution();
  await testInsufficientBalance();
  await testAccountBalanceChangedAfterTransaction();
}

export async function testAccountBalanceChangedAfterTransaction() {
  const acct = await testCreateAccount();

  const body = {
    money: {
      amount: 10,
      currency: "USD",
    },
    memo: "test tx",
    direction: "CREDIT",
    account_id: acct.id,
  };

  const idemKey = gen_key(acct.id, body.money, body.direction);

  const headers = {
    Accept: "*/*",
    "Idempotency-Key": idemKey,
    "Content-Type": "application/json",
  };

  const response = await fetch(
    "http://localhost:8080/v1/transactions/immediate",
    {
      method: "POST",
      body: JSON.stringify(body),
      headers: headers,
    }
  );

  const data = await response.text();
  const parsedData = JSON.parse(data) as Transaction;

  const response2 = await fetch(
    `http://localhost:8080/v1/accounts/${acct.id}`,
    {
      method: "GET",
      headers: headers,
    }
  );

  const data2 = await response2.text();
  const parsedData2 = JSON.parse(data2) as Account;

  assert(response.status === 201 || response.status === 200);
  assert(parsedData.account_id === body.account_id);
  assert(parsedData.money.amount === body.money.amount);
  assert(parsedData.money.currency === body.money.currency);
  assert(parsedData.direction === body.direction);
  assert(parsedData.state === "COMPLETED");
  assert(parsedData.memo === body.memo);
  assert(parsedData2.balance.amount === body.money.amount);
}

export async function testInsufficientBalance() {
  const acct = await testCreateAccount();

  const body = {
    money: {
      amount: 10,
      currency: "USD",
    },
    memo: "test tx",
    direction: "DEBIT",
    account_id: acct.id,
  };

  const idemKey = gen_key(acct.id, body.money, body.direction);

  const headers = {
    Accept: "*/*",
    "Idempotency-Key": idemKey,
    "Content-Type": "application/json",
  };

  const response = await fetch("http://localhost:8080/v1/transactions", {
    method: "POST",
    body: JSON.stringify(body),
    headers: headers,
  });

  const data = await response.text();
  const parsedData = JSON.parse(data) as Transaction;

  assert(response.status === 201 || response.status === 200);
  assert(parsedData.account_id === body.account_id);
  assert(parsedData.money.amount === body.money.amount);
  assert(parsedData.money.currency === body.money.currency);
  assert(parsedData.direction === body.direction);
  assert(parsedData.state === "FAILED");
  assert(parsedData.memo === body.memo);
  assert(parsedData.error_reason === "insufficient funds");
}

export async function testReverseTransaction() {
  const { account, transaction } = await testPendingTransaction();

  const headers = {
    Accept: "*/*",
    "Content-Type": "application/json",
  };

  const response = await fetch(
    `http://localhost:8080/v1/transactions/${transaction.id}`,
    {
      method: "DELETE",
      headers: headers,
    }
  );

  const data = await response.text();
  const parsedData = JSON.parse(data) as Transaction;

  assert(response.status === 201 || response.status === 200);
  assert(data !== undefined);
  assert(parsedData.account_id === account.id);
  assert(parsedData.money.amount === transaction.money.amount);
  assert(parsedData.money.currency === transaction.money.currency);
  assert(parsedData.direction === transaction.direction);
  assert(parsedData.state === "REVERSED");
}

export async function testPendingTransactionExecution() {
  const { account, transaction } = await testPendingTransaction();

  const headers = {
    Accept: "*/*",
    "Content-Type": "application/json",
  };

  const response = await fetch(
    `http://localhost:8080/v1/transactions/${transaction.id}`,
    {
      method: "PUT",
      headers: headers,
    }
  );

  const data = await response.text();
  const parsedData = JSON.parse(data) as Transaction;

  assert(response.status === 201 || response.status === 200);
  assert(data !== undefined);
  assert(parsedData.account_id === account.id);
  assert(parsedData.money.amount === transaction.money.amount);
  assert(parsedData.money.currency === transaction.money.currency);
  assert(parsedData.direction === transaction.direction);
  assert(parsedData.state === "COMPLETED");
}

export async function testPendingTransaction(): Promise<{
  account: Account;
  transaction: Transaction;
}> {
  const acct = await testCreateAccount();

  const body = {
    money: {
      amount: 10,
      currency: "USD",
    },
    memo: "test tx",
    direction: "CREDIT",
    account_id: acct.id,
  };

  const idemKey = gen_key(acct.id, body.money, body.direction);

  const headers = {
    Accept: "*/*",
    "Idempotency-Key": idemKey,
    "Content-Type": "application/json",
  };

  const response = await fetch("http://localhost:8080/v1/transactions", {
    method: "POST",
    body: JSON.stringify(body),
    headers: headers,
  });

  const data = await response.text();
  const parsedData = JSON.parse(data) as Transaction;

  assert(response.status === 201 || response.status === 200);
  assert(data !== undefined);
  assert(parsedData.account_id === body.account_id);
  assert(parsedData.money.amount === body.money.amount);
  assert(parsedData.money.currency === body.money.currency);
  assert(parsedData.direction === body.direction);
  assert(parsedData.state === "PENDING");

  return { account: acct, transaction: parsedData };
}

export async function testImmediateTransaction() {
  const acct = await testCreateAccount();

  const body = {
    money: {
      amount: 10,
      currency: "USD",
    },
    memo: "test tx",
    direction: "CREDIT",
    account_id: acct.id,
  };

  const idemKey = gen_key(acct.id, body.money, body.direction);

  const headers = {
    Accept: "*/*",
    "Idempotency-Key": idemKey,
    "Content-Type": "application/json",
  };

  const response = await fetch(
    "http://localhost:8080/v1/transactions/immediate",
    {
      method: "POST",
      body: JSON.stringify(body),
      headers: headers,
    }
  );

  const data = await response.text();
  const parsedData = JSON.parse(data) as Transaction;

  assert(response.status === 201 || response.status === 200);
  assert(data !== undefined);
  assert(parsedData.account_id === body.account_id);
  assert(parsedData.money.amount === body.money.amount);
  assert(parsedData.money.currency === body.money.currency);
  assert(parsedData.direction === body.direction);
  assert(parsedData.state === "COMPLETED");
}

export async function testFailedTransaction() {
  const body = {
    money: {
      amount: 10,
      currency: "USD",
    },
    memo: "test tx",
    direction: "CREDIT",
    // Create random ids
    account_id: String(Math.floor(Math.random() * 1000)),
  };

  const idemKey = "conflicting";

  const headers = {
    Accept: "*/*",
    "Idempotency-Key": idemKey,
    "Content-Type": "application/json",
  };

  const response = await fetch("http://localhost:8080/v1/transactions", {
    method: "POST",
    body: JSON.stringify(body),
    headers: headers,
  });

  assert(response.status === 400);
}
