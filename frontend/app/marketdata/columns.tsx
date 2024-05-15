"use client";

import { ColumnDef } from "@tanstack/react-table";

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Snapshot = {
  request_id: string;
  request_timestamp: number;
  venue_name: string;
  currency_base: string;
  currency_quote: string;
  open: number;
  high: number;
  close: number;
  volume: number;
};

export const columns: ColumnDef<Snapshot>[] = [
  {
    accessorKey: "request_timestamp",
    header: "Unix timestamp"
  },
  {
    accessorKey: "request_id",
    header: "ID"
  },
  {
    accessorKey: "venue_name",
    header: "Venue"
  },
  {
    accessorKey: "currency_base",
    header: "Base"
  },
  {
    accessorKey: "currency_quote",
    header: "Quote"
  },
  {
    accessorKey: "market_data.open",
    header: "Open"
  },
  {
    accessorKey: "market_data.close",
    header: "Close"
  },
  {
    accessorKey: "market_data.high",
    header: "High"
  },
  {
    accessorKey: "market_data.low",
    header: "Low"
  },
  {
    accessorKey: "market_data.volume",
    header: "Volume"
  }
];
