import { MongoClient } from "mongodb";
import dotenv from "dotenv";

dotenv.config({ path: "./../.env" });

const connectionString = process.env.MONGO_CONNECTION_STR;

const client = new MongoClient(connectionString);

let conn;
try {
  conn = await client.connect();
} catch (e) {
  console.error(e);
}

let db = conn.db("crypto_thp");

export default db;
