import clientPromise from "../../lib/mongodb";
import { NextApiRequest, NextApiResponse } from "next";

export default async (req: NextApiRequest, res: NextApiResponse) => {
  try {
    const client = await clientPromise;
    const db = client.db("crypto_thp");
    const snapshots = await db
      .collection("ohlcv_snapshots")
      .find({})
      .sort({ request_timestamp: -1 })
      .toArray();
    res.json(snapshots);
  } catch (e) {
    console.error(e);
  }
};
