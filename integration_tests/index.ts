import { accounts } from "./account";
import { transactions } from "./transactions";

async function main() {
  await transactions();
  await accounts();
}

main();
