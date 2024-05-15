import { Snapshot, columns } from "./columns";
import { DataTable } from "./data-table";

import React from "react";

async function getData(): Promise<Snapshot[]> {
  const res = await fetch("http://localhost:3000/api/marketdata", {
    cache: "no-store"
  });

  const data = await res.json();
  return data;
}

export default async function Page() {
  const data = await getData();

  return (
    <section className="py-24">
      <div className="container">
        <h1 className="text-3x1 font-bold">All Snapshots</h1>
        <DataTable columns={columns} data={data} />
      </div>
    </section>
  );
}
