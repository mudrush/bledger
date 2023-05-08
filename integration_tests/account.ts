import { assert } from "console";
import fetch from "node-fetch";
import { Account } from "./interfaces";

export async function accounts() {
  // await testCreateAccount();
  // await testGetAccount();
}

export async function testGetAccount() {
  const body = {
    description: "test account",
    name: "test account",
    currency: "USD",
  };

  const headers = {
    Accept: "*/*",
    "Content-Type": "application/json",
  };

  const response = await fetch("http://localhost:8080/v1/accounts", {
    method: "POST",
    body: JSON.stringify(body),
    headers: headers,
  });

  const data = await response.text();
  const parsedData = JSON.parse(data) as Account;

  const response2 = await fetch(
    `http://localhost:8080/v1/accounts/${parsedData.id}`,
    {
      method: "GET",
      headers: headers,
    }
  );

  const data2 = await response2.text();
  const parsedData2 = JSON.parse(data2) as Account;

  assert(response2.status === 201 || response2.status === 200);
  assert(data !== undefined);
  assert(parsedData2.description === body.description);
  assert(parsedData2.name === body.name);
  assert(parsedData2.balance.amount === 0);
  assert(parsedData2.balance.currency === body.currency);
}

export async function testCreateAccount(): Promise<Account> {
  const body = {
    description: "test account",
    name: "test account",
    currency: "USD",
  };

  const headers = {
    Accept: "*/*",
    "Content-Type": "application/json",
  };

  const response = await fetch("http://localhost:8080/v1/accounts", {
    method: "POST",
    body: JSON.stringify(body),
    headers: headers,
  });

  const data = await response.text();
  const parsedData = JSON.parse(data) as Account;

  assert(response.status === 201 || response.status === 200);
  assert(data !== undefined);
  assert(parsedData.description === body.description);
  assert(parsedData.name === body.name);
  assert(parsedData.balance.amount === 0);
  assert(parsedData.balance.currency === body.currency);

  return parsedData;
}
